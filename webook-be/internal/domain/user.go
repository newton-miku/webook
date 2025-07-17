package domain

type User struct {
	Id       int64
	Email    string
	Password []byte
	Ctime    int64
}

type UserProfile struct {
	Id       int64
	UID      int64
	Email    string
	Birthday string
	Nickname string
	Phone    string
	Summary  string `json:"AboutMe"`
}
