package main

import (
	"net/http"
	"io"
	"github.com/k0kubun/pp"
)

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

// from io.go
func TeeReadCloser(readIn io.Reader, readClose io.Closer, writeIn io.Writer, writeClose io.Closer) io.ReadCloser {
	return &teeReadCloser{
		readIn, readClose, writeIn, writeClose,
	}
}

type teeReadCloser struct {
	ri io.Reader
	rc io.Closer
	wi io.Writer
	wc io.Closer
}

func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.ri.Read(p)
	if n > 0 {
		if n, err := t.wi.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *teeReadCloser) Close() error {
	var rErr, wErr error
	if t.rc != nil {
		rErr = t.rc.Close()
	}
	if t.wc != nil {
		wErr = t.wc.Close()
	}
	if rErr != nil || wErr != nil {
		return pp.Errorf("TeeReadCloser.Close error:\n input reader error: %v \n output writer error %v", rErr, wErr)
	} else {
		return nil
	}
}

