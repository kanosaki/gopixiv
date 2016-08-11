package pixiv

import (
	"fmt"
	"github.com/k0kubun/pp"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type ItemType string

const (
	ITEM_TYPE_ILLUSTRATION ItemType = "illustration"
	ITEM_TYPE_MANGA        ItemType = "manga"
	ITEM_TYPE_UGOIRA       ItemType = "ugoira"
)

// where this item data came from?
type ItemSourceAPI int

const (
	API_RANKING ItemSourceAPI = iota
	API_SEARCH
	API_DETAIL
)
const (
	itemTimeStampFormat = "2006-01-02 03:04:05"
)

type User struct {
	IsFriend         bool          `json:"is_friend"`
	Stats            string        `json:"stats"`
	Name             string        `json:"name"`
	IsFollower       bool          `json:"is_follower"`
	ProfileImageUrls ImageUrlsPack `json:"profile_image_urls"`
	Account          string        `json:"account"`
	Profile          string        `json:"profile"`
	IsFollowing      bool          `json:"is_following"`
	IsPremium        bool          `json:"is_premium"`
	ID               int           `json:"id"`
}

type ItemStats struct {
	CommentedCount int           `json:"commented_count"`
	FavoritedCount ItemFavorites `json:"favorited_count"`
	Score          int           `json:"score"`
	ScoredCount    int           `json:"scored_count"`
	ViewsCount     int           `json:"views_count"`
}

type ItemFavorites struct {
	Public  int `json:"public"`
	Private int `json:"private"`
}

type Item struct {
	SourceAPI ItemSourceAPI `json:"-"`
	Caption   string        `json:"caption"`
	// "right_to_left" or "none"
	BookStyle string        `json:"book_style"`
	IsManga   bool          `json:"is_manga"`
	ImageUrls ImageUrlsPack `json:"image_urls"`
	Width     int           `json:"width"`
	// ugoira
	ItemType           ItemType     `json:"type"`
	CreatedTimeExpr    string       `json:"created_time"`
	FavoriteID         int          `json:"favorite_id"`
	ContentType        string       `json:"content_type"`
	PageCount          int          `json:"page_count"`
	Tags               []string     `json:"tags"`
	Tools              []string     `json:"tools"`
	User               User         `json:"user"`
	SanityLevel        string       `json:"sanity_level"`
	ReUploadedTimeExpr string       `json:"reuploaded_time"`
	Stats              ItemStats    `json:"stats"`
	IsLiked            bool         `json:"is_liked"`
	Metadata           ItemMetadata `json:"metadata"`
	Publicity          int          `json:"publicity"`
	Title              string       `json:"title"`
	ID                 int          `json:"id"`
	Height             int          `json:"height"`
	AgeLimit           string       `json:"age_limit"`
}

func (self *Item) Extension() string {
	for _, imgUrl := range self.ImageUrls {
		u, err := url.Parse(imgUrl)
		if err != nil {
			continue
		} else {
			return path.Ext(u.Path)
		}
	}
	return ""
}

func (self *Item) IsFilled() bool {
	return self.SourceAPI == API_DETAIL
}

func (self *Item) CreatedTime() time.Time {
	if self.CreatedTimeExpr == "" {
		return time.Time{}
	}

	ret, err := time.ParseInLocation(itemTimeStampFormat, self.CreatedTimeExpr, time.Local)
	if err != nil {
		panic(err)
	}
	return ret
}

func (self *Item) ReUploadedTime() time.Time {
	if self.ReUploadedTimeExpr == "" {
		return time.Time{}
	}

	ret, err := time.ParseInLocation(itemTimeStampFormat, self.ReUploadedTimeExpr, time.Local)
	if err != nil {
		panic(err)
	}
	return ret
}

// Fill this item details through given pixiv client
func (self *Item) Fill(px *Pixiv) (*Item, error) {
	return fetchItemDetail(px, self)
}

func (self *Item) emulateImageUrlOf(size ImageSize, page int) (string, error) {
	baseUrl, ok := self.ImageUrls[size]
	if !ok {
		return "", pp.Errorf("Image size %v is not available", size)
	}
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	dir, originalName := path.Split(u.Path)
	extension := path.Ext(originalName)
	var name string
	switch size {
	case SIZE_128x128:
		name = fmt.Sprintf("%d_p%d_square1200.%s", self.ID, page, extension)
	case SIZE_480x960:
		fallthrough
	case SIZE_SMALL:
		fallthrough
	case SIZE_MEDIUM:
		name = fmt.Sprintf("%d_p%d_master1200.%s", self.ID, page, extension)
	case SIZE_LARGE:
		name = fmt.Sprintf("%d_p%d.%s", self.ID, page, extension)
	default:
		return "", pp.Errorf("Unsupported image size: %s", size)
	}
	u.Path = path.Join(dir, name)
	return u.String(), nil
}

func (self *Item) OpenImage(px *Pixiv, size ImageSize, pace int) (io.ReadCloser, error) {
	url, ok := self.ImageUrls[size]
	if !ok {
		return nil, pp.Errorf("Image size %v is not available", size)
	}
	client, err := px.PlainClient()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setRequestHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func fetchItemDetail(px *Pixiv, origin *Item) (*Item, error) {
	if origin.IsFilled() {
		return origin, nil
	} else {
		api := &WorkDetail{}
		client, err := px.AuthClient()
		if err != nil {
			return nil, err
		} else {
			return api.Fetch(client, origin)
		}
	}
}

type WorkDetail struct {
	APIEndpoint
}

func (self *WorkDetail) Fetch(client *http.Client, origin *Item) (*Item, error) {
	params := map[string]string{}
	setCommonApiParams(&params)
	var resp []Item
	err := self.call(client, fmt.Sprintf("v1/works/%d.json", origin.ID), params, &resp)
	if err != nil {
		return nil, err
	}
	for _, item := range resp {
		item.SourceAPI = API_DETAIL
	}
	return &resp[0], nil
}
