package crawl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cdipaolo/hub-db/model"

	"github.com/PuerkitoBio/goquery"
)

// Constants hold URLS that are commonly needed
// within the crawler
const (
	ImageDateFormat             = "January 2, 2006"
	BaseURI                     = "http://pornhub.com"
	BasePageURIAllCategoriesFmt = "http://www.pornhub.com/albums/female-gay-male-misc-shemale-straight-uncategorized?o=mv&t=a&page=%d"
)

var (
	// AlbumViewsRegexp will match ">632 views"
	// but not ">632abc views"
	AlbumViewsRegexp = regexp.MustCompile(">[0-9]+ views")

	// AlbumPercentRegexp will match ">76%<" but
	// not ">1<" or ">1006%<"
	AlbumPercentRegexp = regexp.MustCompile(">[0-9]{1,3}%<")

	// AlbumVotesRegexp matches "(1235 votes" but not
	// "(a1431 votes"
	AlbumVotesRegexp = regexp.MustCompile(`\([0-9]+ votes`)

	ImageTimestampRegex = regexp.MustCompile(`Uploaded on [A-Z][a-z]+ [0-9]{1,2}, [0-9]{4}`)

	// CutTabsAndNewlines matches any tabs or newlines
	CutTabsAndNewlines = strings.NewReplacer("\n", "", "\t", "")

	// DropCommas takes the commas out of a string
	DropCommas = strings.NewReplacer(",", "")
)

func GetBaseURIAllCategoriesForPage(page uint64) string {
	return "http://www.pornhub.com/albums/female-gay-male-misc-shemale-straight-uncategorized?o=mv&t=a&page=" + strconv.FormatUint(page, 10)
}

// GetBasePageByPageNumber takes in a page number
// and returns a BasePage struct (only descenced
// down into the album struct, not individual
// images) corresponding to that page number. Note
// that this doesn't filter by segment of the site.
func GetBasePageByPageNumber(page uint64) (*model.BasePage, error) {
	p := &model.BasePage{
		URI:  GetBaseURIAllCategoriesForPage(page),
		Page: page,
	}

	doc, err := goquery.NewDocument(p.URI)
	if err != nil {
		return nil, err
	}

	// now get some albums
	albums := []model.Album{}
	doc.Find(".photoAlbumListBlock > a").Each(func(i int, s *goquery.Selection) {
		// Now get the some mets info about the albums
		link, exists := s.Attr("href")
		if !exists {
			fmt.Printf("Cannot find album link in < %v >\n", p.URI)
			return
		}
		album, err := GetAlbumFromURI(BaseURI + link)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		albums = append(albums, *album)
	})

	p.Albums = albums
	return p, nil
}

// GetAlbumFromURI takes in an album page URL
//     ie. http://www.pornhub.com/album/66212
// and returns a corresponding *model.Album for
// that page
func GetAlbumFromURI(uri string) (*model.Album, error) {
	album := &model.Album{
		URI: uri,
	}
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		return nil, err
	}

	// get album metadata
	album.Title = CutTabsAndNewlines.Replace(doc.Find(".photoAlbumTitleV2").Text())

	upvotePercent := doc.Find("#ratingAlbumInfo > div > span").Text()
	if upvotePercent == "" {
		return nil, fmt.Errorf("Couldn't find the album rating percentage in < %v >", uri)
	}
	percent, err := strconv.Atoi(upvotePercent[:len(upvotePercent)-1])
	if err != nil {
		return nil, err
	}
	// cut the "%" and go from 98% to 0.98
	album.UpvotePercent = float64(percent) / 100

	votes := AlbumVotesRegexp.FindString(doc.Find("#ratingAlbumInfo > div").Text())
	if votes == "" {
		return nil, fmt.Errorf("Couldn't find the album votes in < %v >", uri)
	}
	// chop the "(x votes" to "x" and convert
	votesInt, err := strconv.Atoi(votes[1 : len(votes)-6])
	if err != nil {
		return nil, err
	}
	album.Votes = uint64(votesInt)

	views := doc.Find("#photoAlbumRatingBox > .photoBoxContContainer > #likeBlockContent > #ratingAlbumInfo > div#viewsPhotAlbumCounter").Text()
	if views == "" {
		return nil, fmt.Errorf("Couldn't find the album views in < %v >", uri)
	}
	// chop the "x views" to "x" and convert
	viewsInt, err := strconv.Atoi(views[:len(views)-6])
	if err != nil {
		return nil, err
	}
	album.Views = uint64(viewsInt)

	// this is a little hackey but it'll let me get only
	// the text not in any child nodes, delete any tabs
	// and newlines, and remove the leading space
	album.Segment = CutTabsAndNewlines.Replace(doc.Find("#segmentCont").Children().Remove().End().Text())[1:]
	tags := []model.Tag{}
	doc.Find(".tagLabel").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, model.Tag(s.Text()))
	})
	album.Tags = tags

	images := []model.Image{}
	doc.Find(".photoAlbumListBlock > a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			fmt.Println("Couldn't find image link in < %v >", uri)
			return
		}
		image, err := GetImageFromURI(BaseURI + link)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		images = append(images, *image)
	})
	album.Images = images
	album.NumberOfImages = uint64(len(images))

	return album, nil
}

// GetImageFromURI takes in a URI for an image
// page
//     ie. http://www.pornhub.com/photo/540124
// and returns a corresponding &model.Image
// struct for that page.
func GetImageFromURI(uri string) (*model.Image, error) {
	image := &model.Image{}
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		return nil, err
	}

	var exists bool
	image.URI, exists = doc.Find("#photoImageSection > div > a > img").Attr("src")
	if !exists {
		return nil, fmt.Errorf("Couldn't find image URI on < %v >", uri)
	}

	upvotePercent := doc.Find("#votePercentageNumber").Text()
	if upvotePercent == "" {
		return nil, fmt.Errorf("Couldn't find the image rating percentage in < %v >", uri)
	}
	percent, err := strconv.Atoi(upvotePercent)
	if err != nil {
		return nil, err
	}
	// go from 98 to 0.98
	image.UpvotePercent = float64(percent) / 100

	votesInt, err := strconv.Atoi(doc.Find("#voteCountNumber").Text())
	if err != nil {
		return nil, err
	}
	image.Votes = uint64(votesInt)

	viewsInt, err := strconv.Atoi(DropCommas.Replace(doc.Find("#photoInfoSection > div > strong").Text()))
	if err != nil {
		return nil, err
	}
	image.Views = uint64(viewsInt)

	commentsString := doc.Find("#cmtWrapper > h2 > span").Text()
	if commentsString == "" {
		return nil, fmt.Errorf("Couldn't find the number of comments in < %v >", uri)
	}
	// cut "(x)" into "x" and convert
	commentsInt, err := strconv.Atoi(commentsString[1 : len(commentsString)-1])
	if err != nil {
		return nil, err
	}
	image.NumberOfComments = uint64(commentsInt)

	tags := []model.Tag{}
	doc.Find(".tagText").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, model.Tag(s.Text()))
	})
	image.Tags = tags

	date := ImageTimestampRegex.FindString(doc.Find("#photoWrapper > .photoColumnRight > #userInformation > ul").Text())[12:]
	image.Timestamp, err = time.Parse(ImageDateFormat, date)
	if err != nil {
		return nil, err
	}

	// get comments
	image.PopularComments = []model.Comment{}
	doc.Find("#cmtWrapper > #cmtContent .commentBlock").Each(func(i int, s *goquery.Selection) {
		c := getCommentFromSelection(s)
		if c == nil {
			return
		}
		image.PopularComments = append(image.PopularComments, *c)
	})

	return image, nil
}

// getCommentFromSelection takes in a goquery selection and
// returns a corresponding model.Comment struct filled
// with information about the comment. Because this is
// always embedded within a goquery .Each() block and can't
// throw errors, no error is returned.
func getCommentFromSelection(s *goquery.Selection) *model.Comment {
	c := &model.Comment{}

	c.Text = CutTabsAndNewlines.Replace(s.Find(".commentMessage").Text())

	upvotes, err := strconv.Atoi(s.Find(".voteTotal").Text())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	c.NetUpvotes = uint64(upvotes)

	c.Username = s.Find(".usernameLink").Text()

	return c
}
