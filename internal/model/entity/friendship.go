package entity

type FindAllFriendshipRequest struct {
	Limit      int
	Offset     int
	Search     string
	OrderBy    string
	SortBy     string
	OnlyFriend bool
	UserID     int64
}

type Friendship struct {
	ID        int64
	UserID    int64
	AddedBy   int64
	CreatedAt int64
	UpdatedAt int64
}
