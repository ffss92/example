package posts

type Post struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Likes   int64  `json:"likes"`
	UserID  int64  `json:"-"`
}
