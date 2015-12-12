package model

import "time"

// Image corresponds to the data structure
// of image metadata pulled from PornHub.
// It references popular comments and tags,
// as well as general info about popularity.
type Image struct {
	Timestamp        time.Time `json:"timestamp,omitempty"`
	URI              string    `json:"uri,omitempty"`
	Votes            uint64    `json:"votes"`
	UpvotePercent    float64   `json:"upvote_percent,omitempty"`
	Views            uint64    `json:"views"`
	NumberOfComments uint64    `json:"num_comments"`
	PopularComments  []Comment `json:"popular_comments,omitempty"`
	Tags             []Tag     `json:"tags,omitempty"`
}

// Comment corresponds to the data structure
// for a PornHub image comment. It contains simple
// information about the comments. To not reveal
// too much info about casual users their actual
// profiles aren't crawled
type Comment struct {
	Username string `json:"username,omitempty"`
	//Timestamp  time.Time `json:"timestamp,omitempty"`
	NetUpvotes uint64 `json:"net_upvotes"`
	Text       string `json:"text"`
	//SubComments []Comment `json:"subcomments,omitempty"`
}

// Simple metadata about a user, including their
// username, number of subscribers, friends, views,
// etc.
type User struct {
	Username    string `json:"username"`
	Subscribers uint64 `json:"subscribers,omitempty"`
	Views       uint64 `json:"views"`
	Gender      string `json:"gender,omitempty"`
	Age         uint64 `json:"age,omitempty"`
}

// Tag corresponds to generic information about
// tags (just the name currently to save space)
type Tag string
