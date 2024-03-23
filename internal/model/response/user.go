package response

type User struct {
	ID          int64  `json:"userId"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Name        string `json:"name"`
	FriendCount int64  `json:"phoneCount"`
	ImageUrl    string `json:"imageUrl"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type Register struct {
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type Login struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
