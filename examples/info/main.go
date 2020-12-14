package main

import (
	"fmt"
	"log"

	"github.com/bjjb/yt"
)

func main() {
	info, err := yt.GetInfo("aqz-KE-bpKQ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(info.VideoDetails.Title)
	fmt.Println(info.StreamingData.Formats[0].URL)
	fmt.Printf("Expires in %s seconds\n", info.StreamingData.ExpiresInSeconds)
}
