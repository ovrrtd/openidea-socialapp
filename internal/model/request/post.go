package request

type CreatePost struct {
	ContentHtml string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags        []string `json:"tags" validate:"required"`
	UserID      int64
}

type CreateComment struct {
	Content string `json:"comment" validate:"required,min=2,max=500"`
	PostID  string `json:"postId" validate:"required"`
	UserID  int64
}

type FindAllPost struct {
	Limit  int      `query:"limit"`
	Offset int      `query:"offset"`
	Search string   `query:"search"`
	Tags   []string `query:"searchTag"`
	UserID int64
}
