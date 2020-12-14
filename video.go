package yt

// A Video contains information about a YouTube video
type Video struct{}

// Type implements Result
func (v *Video) Type() string {
	return "video"
}
