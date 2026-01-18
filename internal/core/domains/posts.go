package domains

type PostsDomain struct {
	Author      string `db:"author"      json:"author"`
	UserId      int    `db:"user_id"     json:"user_id"`
	Title       string `db:"title"       json:"title"`
	Description string `db:"description" json:"description"`
}
