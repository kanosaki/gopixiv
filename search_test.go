package pixiv

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	assert := assert.New(t)
	if _, found := os.LookupEnv("TEST_ALL"); !found {
		t.SkipNow()
	}
	pixiv, err := NewFromConfigFile("private/config.json")
	if err != nil {
		t.Error(err)
	}
	// a query for THE FIRST (oldest) illust in pixiv
	items, err := pixiv.Search("pixiv最古絵", 1)
	assert.Nil(err)
	assert.NotEmpty(items)
	kettle := items[len(items)-1]
	assert.Equal(20, kettle.ID)
	assert.Equal("2000年", kettle.Title)
	assert.Equal("馬骨", kettle.User.Name)
}
