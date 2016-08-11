package pixiv

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

// read local test file and assert its contents
func TestJsonParsing(t *testing.T) {
	assert := assert.New(t)
	dat, err := os.Open("test/detail1.json")
	if err != nil {
		t.Fatalf("Cannot open test data! %v", err)
	}
	api := &WorkDetail{}
	rets := []Item{}
	err = api.readAndParse(dat, &rets)
	assert.Nil(err)
	assert.True(len(rets) > 0)
	ret := rets[0]

	assert.Equal(12345678, ret.ID)
	assert.Equal("Captions Here", ret.Caption)
	assert.Equal("MangaSampleFile1", ret.Title)
	assert.Equal([]string{"オリジナル", "tag2", "tag3"}, ret.Tags)
	assert.Equal([]string{"a", "b"}, ret.Tools)
	assert.Equal(1000, ret.Width)
	assert.Equal(500, ret.Height)
	assert.Equal(1234, ret.Stats.ScoredCount)
	assert.Equal(4321, ret.Stats.Score)
	assert.Equal(10000, ret.Stats.ViewsCount)
	assert.Equal(123, ret.Stats.FavoritedCount.Public)
	assert.Equal(543, ret.Stats.FavoritedCount.Private)
	assert.Equal(32, ret.Stats.CommentedCount)
	assert.Equal("all-age", ret.AgeLimit)
	assert.Equal(time.Date(2016, 3, 7, 12, 9, 50, 0, time.Local), ret.CreatedTime())
	assert.Equal(time.Date(2016, 3, 8, 11, 9, 30, 0, time.Local), ret.ReUploadedTime())
	assert.Equal(987654, ret.User.ID)
	assert.Equal("hogehoge", ret.User.Account)
	assert.Equal("ほげほげ", ret.User.Name)
	assert.Equal(true, ret.IsManga)
	assert.Equal(false, ret.IsLiked)
	assert.Equal(0, ret.FavoriteID)
	assert.Equal(8, ret.PageCount)
	assert.Equal("none", ret.BookStyle)
	assert.Equal(ITEM_TYPE_ILLUSTRATION, ret.ItemType)
}
