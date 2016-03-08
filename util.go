package main

import "net/http"

const (
	DEFAULT_REFERRER = "http://spapi.pixiv-app.net/"
	USER_AGENT = "PixivIOSApp/5.8.7"
)

var (
	DEFAULT_API_REQUEST_HEADERS = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer": DEFAULT_REFERRER,
		"User-Agent": USER_AGENT,
		"Accept-Encoding": "gzip, deflate",
	}

	DEFAULT_REQUEST_HEADERS = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer": DEFAULT_REFERRER,
		"User-Agent": USER_AGENT,
		"Accept-Encoding": "gzip, deflate",
	}
	COMMON_API_PARAMS = map[string]string{
		"include_stats": "true",
		"include_sanity_level": "true",
		"image_sizes": "small,px_128x128,px_480mw,large",
		"profile_image_sizes": "px_170x170,px_50x50",
	}
)

func updateRequest(req *http.Request, data map[string]string) {
	for k, v := range data {
		req.Header.Set(k, v)
	}
}

func setAPIRequestHeaders(req *http.Request) {
	updateRequest(req, DEFAULT_API_REQUEST_HEADERS)
}

func setRequestHeaders(req *http.Request) {
	updateRequest(req, DEFAULT_REQUEST_HEADERS)
}

func setCommonApiParams(params *map[string]string) {
	for k, v := range COMMON_API_PARAMS {
		(*params)[k] = v
	}
}
