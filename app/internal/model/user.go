package model

type UserGenData struct {
	Name    string
	Surname string
}

type UserData struct {
	ID int
	UserGenData
}

type UsersData []UserData

type ResponseMeta struct {
	Error   string
	Success string
	Data    interface{}
}
