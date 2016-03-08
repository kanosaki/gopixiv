package main

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

const (
// multi page urls
// array of ImageUrlsPack
// this field might be appear when pages is greater than 1
	ITEM_METADATA_PAGES = "pages"
// ugoira zip urls
// this field appears when "type" is "ugoira"
// like:
//"zip_urls": {
//    "ugoira600x600": "http://i4.pixiv.net/img-zip-ugoira/img/2016/03/04/00/01/28/00000000_ugoira600x600.zip"
//  },
	ITEM_METADATA_ZIP_URLS = "zip_urls"
// ugoira frames
// like:
//  "frames": [
//    {
//    "delay_msec": 80
//    },
//    {
//    "delay_msec": 80
//    },
//    {
//    "delay_msec": 80
//    },
//    ..
	ITEM_METADATA_FRAMES = "frames"
)

// TODO: add support for other metadata types (frames, zip_urls, ...)
type ItemMetadata map[string]interface{}

func (self ItemMetadata) UnmarshalJSON(b []byte) error {
	rawData := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &rawData); err != nil {
		return err
	}
	for key, data := range rawData {
		switch key {
		case ITEM_METADATA_PAGES:
		case ITEM_METADATA_ZIP_URLS:
		case ITEM_METADATA_FRAMES:
		default:
			logrus.Warningf("Unknown metadata! Key: %s, Value: %v", key, data)
		}
	}
	return nil
}

type ImageSize string

const (
// for user profile image, not available for Item images
	SIZE_50x50 ImageSize = "px_50x50"
	SIZE_170x170 ImageSize = "px_170x170"
	SIZE_128x128 ImageSize = "px_128x128"
// might be 150x150
	SIZE_SMALL ImageSize = "small"
// might be 1200x1200
	SIZE_MEDIUM ImageSize = "medium"
	SIZE_480x960 ImageSize = "px_480mw"
// might be raw image?
	SIZE_LARGE ImageSize = "large"
)

//var (
// resizedImagesPattern = regexp.MustCompile(`(?P<id>\d+)_p(?P<page>\d+)_[^.]+\.(?<ext>.+)`)
// filenamePatterns = map[ImageSize]*regexp.Regexp{
//	SIZE_50x50: regexp.MustCompile(`(?P<id>\d+)_s\.(?<ext>.+)`),
//	SIZE_170x170: regexp.MustCompile(`(?P<id>\d+)\.(?<ext>.+)`),
//	SIZE_SMALL: resizedImagesPattern,
//	SIZE_MEDIUM: resizedImagesPattern,
//	SIZE_128x128: resizedImagesPattern,
//	SIZE_480x960: resizedImagesPattern,
//	SIZE_LARGE: regexp.MustCompile(`(?P<id>\d+)_p(?P<page>\d+)\.(?<ext>.+)`),
//}
//)

type ImageUrlsPack map[ImageSize]string

func (self ImageUrlsPack) Get(size ImageSize) (ret string, ok bool) {
	ret, ok = self[size]
	return
}
