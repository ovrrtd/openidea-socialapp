package response

type GetPosts struct {
	PostID string `json:"postId"`
	Post   struct {
		PostInHtml string   `json:"postInHtml"`
		Tags       []string `json:"tags"`
		CreatedAt  string   `json:"createdAt"`
	} `json:"post"`
	Comments []struct {
		Comment string `json:"comment"`
		Creator struct {
			UserID      string `json:"userId"`
			Name        string `json:"name"`
			ImageURL    string `json:"imageUrl"`
			FriendCount int    `json:"friendCount"`
		} `json:"creator"`
		CreatedAt string `json:"createdAt"`
	} `json:"comments"`
	Creator struct {
		UserID      string `json:"userId"`
		Name        string `json:"name"`
		ImageURL    string `json:"imageUrl"`
		FriendCount int    `json:"friendCount"`
		CreatedAt   string `json:"createdAt"`
	} `json:"creator"`
}
