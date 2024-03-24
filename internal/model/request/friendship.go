package request

type FindAllFriendships struct {
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
	Search     string `query:"search"`
	OrderBy    string `query:"orderBy"`
	SortBy     string `query:"sortBy"`
	OnlyFriend bool   `query:"onlyFriend"`
	UserID     int64
}

type CreateFriendship struct {
	UserID  string `json:"userId" validate:"required"`
	AddedBy int64
}

type DeleteFriendship struct {
	Friend1 string `json:"userId" validate:"required"`
	Friend2 int64
}
