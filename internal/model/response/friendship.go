package response

type FindAllFriendships struct {
	ID          string `json:"userId"`
	Name        string `json:"name"`
	ImageUrl    string `json:"imageUrl"`
	FriendCount int64  `json:"friendCount"`
	CreatedAt   string `json:"createdAt"`
}
