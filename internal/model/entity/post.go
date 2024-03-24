package entity

type Post struct {
	ID          int64
	ContentHtml string
	Tags        string
	UserID      int64
	Creator     User
	Comments    []Comment
	CreatedAt   int64
	UpdatedAt   int64
}

type Comment struct {
	ID        int64
	Content   string
	Creator   User
	PostID    int64
	UserID    int64
	CreatedAt int64
	UpdatedAt int64
}

type FindAllPostRequest struct {
	Limit  int
	Offset int
	Search string
	Tags   []string
	UserID int64
}
