package crawl

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cdipaolo/hub-db/model"

	"github.com/stretchr/testify/assert"
)

// hold some test pages
const (
	PageNumber1 = 1
	PageNumber2 = 356
	AlbumURI1   = "http://www.pornhub.com/album/66212"
	AlbumURI2   = "http://www.pornhub.com/album/169283"

	ImageURI1 = "http://www.pornhub.com/photo/93867162"
	ImageURI2 = "http://www.pornhub.com/photo/3136936"
)

func TestGetImageFromURIShouldPass1(t *testing.T) {
	image, err := GetImageFromURI(ImageURI1)
	assert.Nil(t, err, "Error should not be nil")

	assert.NotNil(t, image.Timestamp, "Timestamp shouldn't be nil")
	assert.NotEqual(t, image.URI, "", "CDN image should not be nil")
	assert.True(t, image.Votes > 65, "Votes should be greater than 65")
	assert.True(t, image.UpvotePercent >= 0 && image.UpvotePercent <= 1, "Percent should be in [0,1]")
	assert.True(t, image.Views > 5000, "Image views should be greater than 5000")
	assert.True(t, image.NumberOfComments > 10, "There should be more than 10 comments")
	assert.True(t, len(image.PopularComments) > 4, "There should be more than 4 popular comments saved")
	assert.True(t, len(image.Tags) > 5, "There should be more than 5 tags")

	j, _ := json.Marshal(image)
	fmt.Println(string(j))
}

func TestGetImageFromURIShouldPass2(t *testing.T) {
	image, err := GetImageFromURI(ImageURI2)
	assert.Nil(t, err, "Error should not be nil")

	assert.NotNil(t, image.Timestamp, "Timestamp shouldn't be nil")
	assert.NotEqual(t, image.URI, "", "CDN image should not be nil")
	assert.True(t, image.Votes > 0, "Votes should be greater than 0")
	assert.True(t, image.UpvotePercent >= 0 && image.UpvotePercent <= 1, "Percent should be in [0,1]")
	assert.True(t, image.Views > 50, "Image views should be greater than 50")
	assert.NotNil(t, image.PopularComments, "Popular comments shouldn't be nil")
	assert.True(t, len(image.Tags) > 1, "There should be more than 1 tag")
}

func TestGetAlbumFromURIShouldPass1(t *testing.T) {
	album, err := GetAlbumFromURI(AlbumURI1)
	assert.Nil(t, err, "Error should not be nil")

	assert.Equal(t, album.Title, "Amateur Beach Sex - Public Sex is HOT!!!", "Title should match")
	assert.Equal(t, album.URI, AlbumURI1, "URI should match")
	assert.True(t, album.Views > 95000, "There should be more than 95,000 views")
	assert.True(t, album.Votes > 300, "There should be more than 300 votes")
	assert.True(t, album.UpvotePercent >= 0 && album.UpvotePercent <= 1, "The upvote percent should be in [0,1]")
	assert.Equal(t, album.Segment, "Straight Sex", "Segment should match")
	assert.Equal(t, album.Tags, []model.Tag{"beach", "bikini", "blowjob", "cum", "exhibitionist", "flasher", "fuck", "girlfriend", "outdoor", "public", "pussy", "sex", "shot", "teen", "wife"}, "Tags should match ground truth")

	assert.NotNil(t, album.Images, "Album images should not be nil")
	assert.Len(t, album.Images, 5, "Album should have 5 images")
}

func TestGetAlbumFromURIShouldPass2(t *testing.T) {
	album, err := GetAlbumFromURI(AlbumURI2)
	assert.Nil(t, err, "Error should not be nil")

	assert.Equal(t, album.Title, "This is me 15 years ago.", "Title should match")
	assert.Equal(t, album.URI, AlbumURI2, "URI should match")
	assert.True(t, album.Views > 69000, "There should be more than 69,000 views")
	assert.True(t, album.Votes > 850, "There should be more than 850 votes")
	assert.True(t, album.UpvotePercent >= 0 && album.UpvotePercent <= 1, "The upvote percent should be in [0,1]")
	assert.Equal(t, album.Segment, "Solo Female", "Segment should match")
	assert.Equal(t, album.Tags, []model.Tag{"lol"}, "Tags should match ground truth")

	assert.NotNil(t, album.Images, "Album images should not be nil")
	assert.Len(t, album.Images, 20, "Album should have 20 images")
}

/*
func TestGetBasePageByPageNumberShouldPass1(t *testing.T) {
	p, err := GetBasePageByPageNumber(PageNumber1)
	assert.Nil(t, err, "Error should be nil")

	assert.NotNil(t, p, "Page should not be nil")
	assert.True(t, len(p.Albums) > 10, "Should be more than 10 albums in a page.")

	j, _ := json.Marshal(p)
	f, err := ioutil.TempFile("", "base_page_1_hub_db.json")
	assert.Nil(t, err, "Error should be nil creating temp file")
	fmt.Printf("\n\nSaving test base page output to < %v >", f.Name())

	_, err = f.Write(j)
	assert.Nil(t, err, "File write error should be nil")
}*/
