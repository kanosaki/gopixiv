package gopixiv

import (
	"time"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/k0kubun/pp"
)

type RankingItem struct {
	Work         Illust `json:"work"`
	PreviousRank int  `json:"previous_rank"`
	Rank         int  `json:"rank"`
}

type RankingQuery struct {
	Category string
	Mode     string
	Date     *time.Time
	PerPage  int
	*APIEndpoint
}

type RankingResponse struct {
	Content string `json:"content"`
	Mode    string `json:"mode"`
	Date    string `json:"date"`
	Works   []RankingItem `json:"works"`
}

func (px *Pixiv) RankingQuery(category string, mode string, perPage int, date *time.Time) *RankingQuery {
	query := &RankingQuery{
		Category: category,
		Mode: mode,
		Date: date,
		PerPage: perPage,
		APIEndpoint: InitAPIEndpoint(),
	}
	return query
}

func (px *Pixiv) Ranking(category string, mode string, perPage int, date *time.Time, page int) ([]RankingItem, error) {
	query := px.RankingQuery(category, mode, perPage, date)
	client, err := px.AuthClient()
	if err != nil {
		return nil, err
	}
	return query.Get(client, page)
}

func (rq *RankingQuery) Get(client *http.Client, page int) ([]RankingItem, error) {
	params := map[string]string{
		"mode": rq.Mode,
		"include_stats": "true",
		"include_sanity_level": "true",
		"image_sizes": "small,px_128x128,px_480mw,large",
		// "profile_image_sizes": "px_170x170,px_50x50",
		"page":  strconv.Itoa(page),
		"per_page": strconv.Itoa(rq.PerPage),
	}
	if rq.Date != nil {
	}
	req, err := rq.RequestGET("v1/ranking/" + rq.Category, params)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, pp.Errorf("%v - \n %v", res.Status, res.Request)
	}
	defer res.Body.Close()
	apiResponse, err := ReadAPIResponse(res.Body)
	if err != nil {
		return nil, err
	}
	var rankingResponse []RankingResponse
	err = json.Unmarshal(*apiResponse.Response, &rankingResponse)
	if err != nil {
		return nil, err
	}
	return rankingResponse[0].Works, nil
}