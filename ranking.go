package pixiv

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type RankingCategory string

const (
	RANKING_CATEGORY_ALL    RankingCategory = "all"
	RANKING_CATEGORY_ILLUST RankingCategory = "illust"
	RANKING_CATEGORY_MANGA  RankingCategory = "manga"
	RANKING_CATEGORY_UGOIRA RankingCategory = "ugoira"
	RANKING_CATEGORY_NOVEL  RankingCategory = "novel"
)

type RankingMode string

const (
	RANKING_MODE_DAILY       RankingMode = "daily"
	RANKING_MODE_DAILY_R18   RankingMode = "daily_r18"
	RANKING_MODE_WEEKLY      RankingMode = "weekly"
	RANKING_MODE_WEEKLY_R18  RankingMode = "weekly_r18"
	RANKING_MODE_MONTHLY     RankingMode = "monthly"
	RANKING_MODE_MONTHLY_R18 RankingMode = "monthly_r18"

	// only allowed for 'all' categories
	RANKING_MODE_ORIGINAL   RankingMode = "original"
	RANKING_MODE_ROOKIE     RankingMode = "rookie"
	RANKING_MODE_MALE       RankingMode = "male"
	RANKING_MODE_MALE_R18   RankingMode = "male_r18"
	RANKING_MODE_FEMALE     RankingMode = "female"
	RANKING_MODE_FEMALE_R18 RankingMode = "female_r18"
	RANKING_MODE_R18G       RankingMode = "r18g"
)

var (
	allOnlyRankings map[RankingMode]bool = map[RankingMode]bool{
		RANKING_MODE_ORIGINAL:   true,
		RANKING_MODE_ROOKIE:     true,
		RANKING_MODE_MALE:       true,
		RANKING_MODE_MALE_R18:   true,
		RANKING_MODE_FEMALE:     true,
		RANKING_MODE_FEMALE_R18: true,
		RANKING_MODE_R18G:       true,
	}
)

type RankingItem struct {
	Work         Item `json:"work"`
	PreviousRank int  `json:"previous_rank"`
	Rank         int  `json:"rank"`
}

type RankingQuery struct {
	Category RankingCategory
	Mode     RankingMode
	Date     *time.Time
	PerPage  int
	APIEndpoint
}

type RankingResponse struct {
	Content string        `json:"content"`
	Mode    string        `json:"mode"`
	Date    string        `json:"date"`
	Works   []RankingItem `json:"works"`
}

func (px *Pixiv) RankingQuery(category RankingCategory, mode RankingMode, perPage int, date *time.Time) *RankingQuery {
	q := &RankingQuery{
		Category: category,
		Mode:     mode,
		Date:     date,
		PerPage:  perPage,
	}
	client, err := px.AuthClient()
	if err == nil && client != nil {
		q.DefaultClient(client)
	}
	return q
}

func (px *Pixiv) Ranking(category RankingCategory, mode RankingMode, perPage int, date *time.Time, page int) ([]RankingItem, error) {
	query := px.RankingQuery(category, mode, perPage, date)
	return query.Get(page)
}

func (rq *RankingQuery) Get(page int) ([]RankingItem, error) {
	if rq.client == nil {
		return nil, errors.New("Client is not activated!")
	}
	return rq.Fetch(rq.client, page)
}

func (rq *RankingQuery) Fetch(client *http.Client, page int) ([]RankingItem, error) {
	params := map[string]string{
		"mode":     string(rq.Mode),
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(rq.PerPage),
	}
	setCommonApiParams(&params)
	if rq.Date != nil {
	}
	var rankingResponse []RankingResponse
	err := rq.call(client, "v1/ranking/"+string(rq.Category), params, &rankingResponse)
	if err != nil {
		return nil, err
	}
	for _, rRes := range rankingResponse {
		for _, item := range rRes.Works {
			item.Work.SourceAPI = API_RANKING
		}
	}
	return rankingResponse[0].Works, nil
}
