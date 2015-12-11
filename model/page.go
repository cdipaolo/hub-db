package model

// BasePage corresponds to the most viewed
// photos page at a given page number
type BasePage struct {
	URI    string  `json:"uri"`
	Page   uint64  `json:"page"`
	Albums []Album `json:"albums"`
}

// Album corresponds to a image album
type Album struct {
	Title          string  `json:"title"`
	URI            string  `json:"URI"`
	Votes          uint64  `json:"votes"`
	UpvotePercent  float64 `json:"upvote_percent,omitempty"`
	Views          uint64  `json:"views"`
	Segment        string  `json:"segment"`
	Tags           []Tag   `json:"tags"`
	NumberOfImages uint64  `json:"num_images"`
	Images         []Image `json:"images"`
}
