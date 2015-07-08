package gopixiv
import (
	"net/http"
	"net/url"
	"encoding/json"
	"io"
	"io/ioutil"
	"compress/gzip"
)

var DefaultRequestHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"User-Agent": "PixivIOSApp/5.6.0",
	//"Accept": "*/*",
	"Accept-Encoding": "gzip, deflate",
	"Referer": "http://spapi.pixiv.net/", // non-auth (image servers...)
	// "Referer": "http://www.pixiv.net", // auth (api server ...)
	//"Accept-Language": "ja-jp",
}

var DefaultBaseUrl = "https://public-api.secure.pixiv.net/v1/"

type APIEndpoint struct {
	BaseUrl        string
	RequestHeaders map[string]string
}


type APIResponse struct {
	Status   string `json:"status"`
	Response *json.RawMessage `json:"response"`
}

func ReadAPIResponse(input io.ReadCloser) (*APIResponse, error) {
	result := new(APIResponse)
	uncomp, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(uncomp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InitAPIEndpoint() *APIEndpoint {
	return &APIEndpoint{
		BaseUrl: DefaultBaseUrl,
		RequestHeaders: DefaultRequestHeaders,
	}
}

func (ep *APIEndpoint) Get(path string, params map[string]string) {

}

func (ep *APIEndpoint) Url(path string, params map[string]string) (string, error) {
	u, err := url.Parse(ep.BaseUrl)
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
	for k, v := range ep.RequestHeaders {
		req.Header.Set(k, v)
	}
	return nil
}



