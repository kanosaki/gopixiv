package main

import (
	"testing"
	"github.com/k0kubun/pp"
)

func TestFetchRanking(t *testing.T) {
	pixiv, err := NewFromConfigFile("private/config.json")
	if err != nil {
		t.Error(err)
	}
	ranking, err := pixiv.Ranking("all", "daily", 50, nil, 1)
	if err != nil {
		t.Error(err)
	}
	if len(ranking) == 0 {
		t.Error("Empty result!")
	}
}
