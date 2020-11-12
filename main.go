package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const rex = "^[[:word:]]([[:word:]]|-){10}$"

var regex *regexp.Regexp
var fs http.Handler

// https://youtube.com/get_video_info?video_id=
const videoInfoURL = "https://youtube.com/get_video_info?video_id="

func main() {
	var videoID, webRoot, tmpl, addr string
	var daemon bool
	var t *template.Template
	flag.StringVar(&videoID, "i", "", "get video info")
	flag.StringVar(&webRoot, "D", "", "serve files from a static directory")
	flag.BoolVar(&daemon, "d", false, "start a server")
	flag.StringVar(&addr, "b", ":8080", "HTTP bind address")
	flag.StringVar(&tmpl, "f", "{{.}}", "format output")
	flag.Parse()

	t, err := template.New("cli").Parse(tmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid format specifier (%s)", err)
		os.Exit(2)
	}

	if daemon {
		if webRoot != "" {
			fs = http.FileServer(http.Dir(webRoot))
			fmt.Fprintf(os.Stdout, "Serving files from %s, ", webRoot)
		}
		fmt.Fprintf(os.Stdout, "listening on %s...", addr)
		err := http.ListenAndServe(addr, http.HandlerFunc(serveHTTP))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(9)
		}
		os.Exit(0)
	}

	if videoID != "" {
		info, err := getVideoInfo(videoID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(3)
		}
		if err = t.Execute(os.Stdout, info); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(4)
		}
	}
}

type videoInfo struct {
	VideoDetails struct {
		ID               string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelID        string   `json:"channelId"`
		ShortDescription string   `json:"shortDescription"`
		Thumbnail        struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`

		AverageRating float64 `json:"averageRating"`
		AllowRatings  bool    `json:"allowRatings"`
		ViewCount     string  `json:"viewCount"`
		Author        string  `json:"author"`
		IsPrivate     bool    `json:"isPrivate"`
		IsLiveContent bool    `json:"isLiveContent"`
	} `json:"videoDetails"`
	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		Formats          []struct {
			ITag             int    `json:"itag"`
			URL              string `json:"url"`
			MIMEType         string `json:"mimeType"`
			Bitrate          int    `json:"bitrate"`
			Width            int    `json:"width"`
			Height           int    `json:"height"`
			Quality          string `json:"quality"`
			QualityLabel     string `json:"qualityLabel"`
			LastModified     string `json:"lastModified"`
			ContentLength    string `json:"contentLength"`
			FPS              int    `json:"fps"`
			ApproxDurationMS string `json:"approxDurationMs"`
		} `json:"formats"`
		AdaptiveFormats []struct {
			ITag      int    `json:"itag"`
			URL       string `json:"url"`
			MIMEType  string `json:"mimeType"`
			Bitrate   int    `json:"bitrate"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
			InitRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"initRange"`
			IndexRange struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"indexRange"`
			LastModified     string `json:lastModified"`
			ContentLength    string `json:"contentLength"`
			Quality          string `json:"quality"`
			FPS              int    `json:"fps"`
			QualityLabel     string `json:"qualityLabel"`
			ProjectionType   string `json:"projectionType"`
			AverageBitrate   int    `json:"averageBitrate"`
			ApproxDurationMS string `json:"approxDurationMs"`
		} `json:"adaptiveFormats"`
	} `json:"streamingData"`
}

func getVideoInfo(id string) (*videoInfo, error) {
	if regexp.MustCompile("^[[:word:]][[[:word]]-]{10}$").MatchString(id) {
		return nil, fmt.Errorf("invalid video ID %s", id)
	}
	u, err := url.Parse(videoInfoURL + id)
	if err != nil {
		return nil, err
	}
	u.Query().Set("video_id", id)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d (%s)", resp.StatusCode, resp.Status)
	}
	if resp.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return nil, fmt.Errorf("unexpected response content %s", resp.Header.Get("Content-Type"))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, err
	}
	if values.Get("status") == "fail" {
		return nil, fmt.Errorf("errorcode %s (%s)", values.Get("errorcode"), values.Get("reason"))
	}
	if values.Get("player_response") == "" {
		return nil, errors.New("no player_response")
	}
	info := new(videoInfo)
	err = json.NewDecoder(strings.NewReader(values.Get("player_response"))).Decode(info)
	return info, err
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Methods", http.MethodGet)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Accept")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if regex.MatchString(r.URL.Path[1:]) {
		info, err := getVideoInfo(r.URL.Path[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	fs.ServeHTTP(w, r)
}

func init() {
	regex = regexp.MustCompile(rex)
	fs = http.HandlerFunc(http.NotFound)
}
