## `hub-db`
#### A PornHub images meta-data dataset

`hub-db` is a dataset of information about albums in the adult website [PornHub](pornhub.com). This application crawls through the 'most viewed' search pages (which pages are defined in the [`config.json`](config.json) file) and recursively crawls all albums on those pages, and the images from those albums. No images are saved, but links to the images as well as tag metadata, upload timestamp, comments, etc. are saved.

This repository includes both the code to crawl PornHub to get this data as well as the dataset itself (when crawling is done).

### JSON Data Structure

Below is the structure for the json, which is dumped by BasePage into separate files. These files can be then combined to form whatever smaller and more specific dataset you want.

Example JSON dumps can be found in [`data`](data/), along with the actual dataset fully compiled.

```go
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

// Image corresponds to the data structure
// of image metadata pulled from PornHub.
// It references popular comments and tags,
// as well as general info about popularity.
type Image struct {
	Timestamp        time.Time `json:"timestamp,omitempty"`
	URI              string    `json:"uri,omitempty"`
	Votes            uint64    `json:"votes"`
	UpvotePercent    float64   `json:"upvote_percent,omitempty"`
	Views            uint64    `json:"vites"`
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
	Username   string `json:"username,omitempty"`
	NetUpvotes uint64 `json:"net_upvotes"`
	Text       string `json:"text"`
}

// Tag corresponds to generic information about
// tags (just the name currently to save space)
type Tag string
```

## Endpoints

#### `GET /`
##### Status Check

Endpoint returns some generic info about how the crawling is going.

Example Response:

```json
HTTP/1.1 200 OK
Content-Length: 104
Content-Type: application/json
Date: Sat, 12 Dec 2015 01:07:34 GMT

{
    "errors": [],
    "pages_crawled": 2,
    "saving_json_to": "./data",
    "status": "running",
    "uptime": "16m15.553127188s"
}
```

# LICENSE (data = PornHub's ; code = MIT)

All data remains within the rights reserved by [PornHub](http://www.pornhub.com/information#terms).

For the code (not the data):

```
The MIT License

Copyright (c) 2015 Conner DiPaolo @cdipaolo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
