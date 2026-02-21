package domains

type CommentsDomain struct {
	Author    string `db:"author"    json:"author"`
	PostId    int    `db:"post_id"   json:"post_id"`
	UserId    int    `db:"user_id"   json:"user_id"`
	Content   string `db:"content"   json:"content"`
}
