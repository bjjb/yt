package yt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// ContentTypeXWWWFormURLEncoded is the preferred MIME-type for video info
const ContentTypeXWWWFormURLEncoded = "application/x-www-form-urlencoded"

// InfoID is the regular expression which video IDs must match
var InfoID *regexp.Regexp

// InfoURL is the URL from which we fetch video info
var InfoURL *url.URL

// An InfoRepo is a repository from where YouTube video info can be obtained.
// It provides a single function (Get) which performs the search, optionally
// with caching.
type InfoRepo interface {
	Get(string) (*Info, error)
}

// An InfoCache provides methods to store and retrieve an Info.
type InfoCache interface {
	Get(string) *Info
	Put(string, *Info)
}

// An Info represents all the data that YouTube players can use to play media.
type Info struct {
	VideoDetails *struct {
		ID               string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelID        string   `json:"channelId"`
		ShortDescription string   `json:"shortDescription"`
		Thumbnail        *struct {
			Thumbnails []*struct {
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
	StreamingData *struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		Formats          []*struct {
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
		AdaptiveFormats []*struct {
			ITag      int    `json:"itag"`
			URL       string `json:"url"`
			MIMEType  string `json:"mimeType"`
			Bitrate   int    `json:"bitrate"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
			InitRange *struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"initRange"`
			IndexRange *struct {
				Start string `json:"start"`
				End   string `json:"end"`
			} `json:"indexRange"`
			LastModified     string `json:"lastModified"`
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

// An InfoClient can fetch info for a given video ID. A zero InfoClient uses
// defaults.
type InfoClient struct {
	InfoID *regexp.Regexp
	URL    *url.URL
	HTTP   *http.Client
}

// Get fetches the video info from it's URL (using it's http.Client).
func (i *InfoClient) Get(id string) (*Info, error) {
	m := i.InfoID
	if m == nil {
		m = InfoID
	}
	if !m.MatchString(id) {
		return nil, fmt.Errorf("invalid video ID %q", id)
	}

	u := i.URL
	if u == nil {
		u = InfoURL
	}
	u, err := url.Parse(u.String())
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("video_id", id)
	u.RawQuery = q.Encode()

	c := i.HTTP
	if c == nil {
		c = new(http.Client)
	}

	resp, err := c.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d (%s)", resp.StatusCode, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != ContentTypeXWWWFormURLEncoded {
		return nil, fmt.Errorf("unexpected response content type %q", contentType)
	}

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

	pr := values.Get("player_response")
	if pr == "" {
		return nil, errors.New("no player_response")
	}

	info := new(Info)
	err = json.NewDecoder(strings.NewReader(pr)).Decode(info)

	return info, err
}

// GetInfo gets video info for the video with the given ID
func GetInfo(id string) (*Info, error) {
	c := new(InfoClient)
	return c.Get(id)
}

func init() {
	InfoID = regexp.MustCompile("^[[:word:]]([[:word:]]|-){10}$")
	InfoURL, _ = url.Parse("https://youtube.com/get_video_info")
}
