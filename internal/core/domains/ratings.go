package domains

type Vote struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	PostID   int    `json:"post_id"`
	VoteType string `json:"vote_type"` // "like", "dislike"
}

type VoteAction string // like | dislike

type RatingStats struct {
	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
	NetScore int `json:"net_score"` // likes - dislikes
}

type RatingStatsWithUserVote struct {
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	NetScore int    `json:"net_score"` // likes - dislikes
	UserVote string `json:"user_vote"` // "like", "dislike", or empty if not voted
}
