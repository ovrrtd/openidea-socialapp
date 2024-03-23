package request

type Register struct {
	CredentialType  string `json:"credentialType" validate:"required"`
	CredentialValue string `json:"credentialValue" validate:"required"`
	Name            string `json:"name" validate:"required,min=5,max=50"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type Login struct {
	CredentialType  string `json:"credentialType" validate:"required"`
	CredentialValue string `json:"credentialValue" validate:"required"`
	Password        string `json:"password" validate:"required,min=5,max=15"`
}

type UpdateAccount struct {
	ID       int64  `validate:"required"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	ImageUrl string `json:"imageUrl" validate:"required"`
}

type LinkPhone struct {
	ID    int64  `validate:"required"`
	Phone string `json:"phone" validate:"required"`
}

type LinkEmail struct {
	ID    int64  `validate:"required"`
	Email string `json:"email" validate:"required"`
}
