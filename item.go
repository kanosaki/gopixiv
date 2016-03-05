package main

type User struct {
	IsFriend         bool `json:"is_friend"`
	Stats            string `json:"stats"`
	Name             string `json:"name"`
	IsFollower       bool `json:"is_follower"`
	ProfileImageUrls map[string]string `json:"profile_image_urls"`
	Account          string `json:"account"`
	Profile          string `json:"profile"`
	IsFollowing      bool `json:"is_following"`
	IsPremium        bool `json:"is_premium"`
	ID               int `json:"id"`
}

type ItemStats struct {
	CommentedCount int `json:"commented_count"`
	FavoritedCount ItemFavorites `json:"favorited_count"`
	Score          int `json:"score"`
	ScoredCount    int `json:"scored_count"`
	ViewsCount     int `json:"views_count"`
}

type ItemFavorites struct {
	Public  int `json:"public"`
	Private int `json:"private"`
}

type Item struct {
	Caption            string `json:"caption"`
	// "right_to_left" or "none"
	BookStyle          string `json:"book_style"`
	IsManga            bool `json:"is_manga"`
	ImageUrls          map[string]string `json:"image_urls"`
	Width              int `json:"width"`
	// ugoira
	ItemType           string `json:"type"`
	CreatedTimeExpr    string `json:"created_time"`
	FavoriteID         int `json:"favorite_id"`
	ContentType        string `json:"content_type`
	PageCount          int `json:"page_count"`
	Tags               []string `json:"tags"`
	Tools              []string `json:"tools"`
	User               User `json:"user"`
	SanityLevel        string `json:"sanity_level"`
	ReUploadedTimeExpr string `json:"reuploaded_time"`
	Stats              ItemStats `json:"stats"`
	IsLiked            bool `json:"is_liked"`
	Metadata           string `json:"metadata"`
	Publicity          int `json:"publicity"`
	Title              string `json:"title"`
	ID                 int `json:"id"`
	Height             int `json:"height"`
	AgeLimit           string `json:"age_limit"`
}
