package main

import (
	"net/http"
	"strconv"
	"strings"
	"errors"
	"fmt"
)

type SearchQuery struct {
	Query   string
	Mode    SearchMode
	Types   []ItemType
	Sort    SearchResultSort
	Order   SearchResultOrders
	Period  SearchPeriods
	PerPage int
	APIEndpoint
}

type SearchMode string

const (
	SEARCH_MODE_CAPTION SearchMode = "caption"
	SEARCH_MODE_TAG SearchMode = "tag"
)

type SearchResultOrders string

const (
	SEARCH_ORDER_ASCENDING = "asc"
	SEARCH_ORDER_DESCENDING = "desc"
)

type SearchResultSort string

const (
	SEARCH_SORT_DATE = "date"
// NOTE: popularity sort is only available in paid account
	SEARCH_SORT_POPULARITY = "popular"
)

type SearchPeriods string

const (
	SEARCH_PERIOD_ALL = "all"
	SEARCH_PERIOD_DAY = "day"
	SEARCH_PERIOD_WEEK = "week"
	SEARCH_PERIOD_MONTH = "month"
)

func (px *Pixiv) SearchQuery(queryString string) *SearchQuery {
	q := &SearchQuery{
		Query: queryString,
		Mode: SEARCH_MODE_TAG,
		Types: []ItemType{ITEM_TYPE_ILLUSTRATION, ITEM_TYPE_MANGA, ITEM_TYPE_UGOIRA},
		Sort: SEARCH_SORT_DATE,
		Order: SEARCH_ORDER_DESCENDING,
		Period: SEARCH_PERIOD_ALL,
		PerPage: 50,
	}
	client, err := px.AuthClient()
	if err == nil && client != nil {
		q.DefaultClient(client)
	}
	return q
}

func (px *Pixiv) Search(query string, page int) ([]Item, error) {
	q := px.SearchQuery(query)
	return q.Get(page)
}

func (rq *SearchQuery) Get(page int) ([]Item, error) {
	if (rq.client == nil) {
		return nil, errors.New("Client is not activated!")
	}
	return rq.Fetch(rq.client, page)
}

func (rq *SearchQuery) Fetch(client *http.Client, page int) ([]Item, error) {
	if len(rq.Query) == 0 || len(strings.TrimSpace(rq.Query)) == 0 {
		return nil, errors.New(fmt.Sprintf("SerchQuery: %s is invalid! Empty or whitespace query is not allowed", rq.Query))
	}
	var searchTypes []string
	for _, item := range rq.Types {
		searchTypes = append(searchTypes, string(item))
	}
	params := map[string]string{
		"q": rq.Query,
		"mode": string(rq.Mode),
		"types": strings.Join(searchTypes, ","),
		"order": string(rq.Order),
		"sort": string(rq.Sort),
		"period": string(rq.Period),
		"page":  strconv.Itoa(page),
		"per_page": strconv.Itoa(rq.PerPage),
	}
	setCommonApiParams(&params)
	var searchResponse []Item
	err := rq.execGet(client, "v1/search/works", params, &searchResponse)
	if err != nil {
		return nil, err
	}
	for _, item := range searchResponse {
		item.SourceAPI = API_SEARCH
	}
	return searchResponse, nil
}
