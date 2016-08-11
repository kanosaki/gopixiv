package pixiv

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/k0kubun/pp"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	DUMPFILE_FORMAT = "20060102030405.000"
)

var DefaultRequestHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"User-Agent":   USER_AGENT,
	//"Accept": "*/*",
	"Accept-Encoding": "gzip, deflate",
	"Referer":         DEFAULT_REFERRER, // non-auth (image servers...)
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
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
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

// returns res.Body
// if DEBUG and DUMP_COMMUNICATION_DIR is set
// returned ReadCloser will be a TeeReadCloser and it will dump request and response to
// the specified directory.
// DUMP_COMMUNICATION_DIR will be created if not exists.
// Dump files format is hard coded (for now) and it like
// 20060102030405.000.req.txt or 20060102030405.000.res_body.json.gz
func (ep *APIEndpoint) openResponse(res *http.Response) (io.ReadCloser, error) {
	ret := res.Body
	if *DEBUG && len(*DUMP_COMMUNICATION_DIR) > 0 {
		dumpDir := *DUMP_COMMUNICATION_DIR
		if _, err := os.Stat(dumpDir); err == nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(dumpDir, 0755); err != nil {
					logrus.Fatalf("Cannot create dump dir %v", err)
				}
			} else if err != nil {
				logrus.Fatalf("Error during prepareing dump dir %v", err)
			}
		}
		// prepare file dump
		timestamp := time.Now().Format(DUMPFILE_FORMAT)
		dumpBody := path.Join(dumpDir, fmt.Sprintf("%s.res_body.json", timestamp))
		bodyFp, err := os.Create(dumpBody)
		if err != nil {
			return nil, err
		}
		ret = TeeReadCloser(ret, ret, bodyFp, bodyFp)

		dumpReq := path.Join(dumpDir, fmt.Sprintf("%s.req.txt", timestamp))
		reqFp, err := os.Create(dumpReq)
		defer reqFp.Close()
		if err != nil {
			return nil, err
		}
		if err = res.Request.Write(reqFp); err != nil {
			logrus.Error(err)
		}
	}
	return ret, nil
}

func ReadAPIResponse(input io.Reader) (*APIResponse, error) {
	result := new(APIResponse)
	dec := json.NewDecoder(input)
	err := dec.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ep *APIEndpoint) readAndParse(resData io.Reader, ret interface{}) error {
	apiResponse, err := ReadAPIResponse(resData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(apiResponse.Response, ret)
	if err != nil {
		return pp.Errorf("JSON Parse failed: %v", err)
	}
	if apiResponse.Status != "success" {
		return pp.Errorf("Request failed: %v", apiResponse)
	}
	return nil
}

// call this api with given client, path and params. And stores its response to ret
// through json.Unmarshal
func (ep *APIEndpoint) call(client *http.Client, path string, params map[string]string, ret interface{}) error {
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
	defer func() {
		// Drain and close the body to let the Transport reuse the connection
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	body, err := ep.openResponse(res)
	if err != nil {
		return err
	}
	decBody, err := gzip.NewReader(body)
	if err != nil {
		return err
	}
	ep.readAndParse(decBody, ret)
	return nil
}
