package main

import (
	"net/http"
	"net/url"
	"encoding/json"
	"io"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/k0kubun/pp"
)

var DefaultRequestHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"User-Agent": USER_AGENT,
	//"Accept": "*/*",
	"Accept-Encoding": "gzip, deflate",
	"Referer": DEFAULT_REFERRER, // non-auth (image servers...)
	// "Referer": "http://www.pixiv.net", // auth (api server ...)
	//"Accept-Language": "ja-jp",
}

var DefaultBaseUrl = "https://public-api.secure.pixiv.net/v1/"

type APIEndpoint struct {
	BaseUrl        string
	RequestHeaders map[string]string
	client         *http.Client
}

type APIResponse struct {
	Status   string `json:"status"`
	Response *json.RawMessage `json:"response"`
}

func ReadAPIResponse(input io.Reader) (*APIResponse, error) {
	result := new(APIResponse)
	dec, err := DecodeResponse(input)
	if err != nil {
		return nil, err
	}
	err = dec.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DecodeResponse(input io.Reader) (*json.Decoder, error) {
	uncomp, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}
	return json.NewDecoder(uncomp), nil
}

func (ep *APIEndpoint) Url(path string, params map[string]string) (string, error) {
	u, err := url.Parse(ep.BaseUrl)
	if err != nil {
		return "", err
	}
	if len(ep.BaseUrl) == 0 {
		u, err = url.Parse(DefaultBaseUrl)
	}
	if ep.RequestHeaders == nil {
		ep.RequestHeaders = DefaultRequestHeaders
	}
	if err != nil {
		return "", err
	}
	baseQuery := u.Query()
	for k, v := range params {
		baseQuery.Set(k, v)
	}
	u.RawQuery = baseQuery.Encode()
	u.Path = path
	return u.String(), nil
}

func (ep *APIEndpoint) RequestGET(path string, params map[string]string) (*http.Request, error) {
	urlString, err := ep.Url(path, params)
	if err != nil {
		return nil, err
	}
	fmt.Println(urlString)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	err = ep.PrepareRequest(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (ep *APIEndpoint) PrepareRequest(req *http.Request) error {
	headers := ep.RequestHeaders
	if headers == nil {
		headers = DefaultRequestHeaders
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return nil
}

func (ep *APIEndpoint) DefaultClient(cl *http.Client) error {
	if cl == nil {
		return errors.New("Nil client is not allowed")
	}
	ep.client = cl
	return nil
}

func (ep *APIEndpoint) execGet(client *http.Client, path string, params map[string]string, ret interface{}) error {
	req, err := ep.RequestGET(path, params)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return pp.Errorf("%v - \n %v", res.Status, res.Request)
	}
	defer res.Body.Close()
	apiResponse, err := ReadAPIResponse(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(*apiResponse.Response, &ret)
	if err != nil {
		return err
	}
	if apiResponse.Status != "success" {
		return pp.Errorf("Request filed: %v", apiResponse)
	}
	return nil
}
