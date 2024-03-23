package entity

type User struct {
	ID          int64
	Name        string
	Phone       string
	Email       string
	FriendCount int64
	ImageUrl    string
	Password    string
	CreatedAt   int64
	UpdatedAt   int64
}
