package domain

type User struct {
	Id       int64
	Email    string
	Password []byte
	Ctime    int64
}
