package model

// DB structure
type Users struct {
	UserID       int    `json:"id" db:"id"`
	NickName     string `json:"nickname" db:"nickname"`
	Name         string `json:"name" db:"name"`
	LastName     string `json:"lastname" db:"lastname"`
	EmailAddress string `json:"email" db:"email"`
	Password     string `json:"password" db:"password"`
	Status       string `json:"status" db:"status"`
	CreatedDate  string `json:"created_at" db:"created_at"`
	UpdatedDate  string `json:"updated_at" db:"updated_at"`
	DeleteDate   string `json:"delete_at" db:"delete_at"`
	Likes        int    `json:"likes" db:"likes"`
}

// User's request structures
type UserSignUp struct {
	NickName string `json:"nickname"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password"`
}

// REST API response structures
type ApiResponse struct {
	UserID       int    `json:"id" db:"id"`
	NickName     string `json:"nickname" db:"nickname"`
	Name         string `json:"name" db:"name"`
	LastName     string `json:"lastname" db:"lastname"`
	EmailAddress string `json:"email" db:"email"`
	Status       string `json:"status" db:"status"`
	CreatedDate  string `json:"created_at" db:"created_at"`
	UpdatedDate  string `json:"updated_at" db:"updated_at"`
	DeleteDate   string `json:"delete_at" db:"delete_at"`
	Likes        int    `json:"likes" db:"likes"`
}

// structure which API expect for user info update
type UpdateUser struct {
	UserID       int    `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"lastname"`
	EmailAddress string `json:"email"`
	Likes        int    `json:"likes"`
}
