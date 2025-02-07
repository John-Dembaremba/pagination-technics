package model

type UserGenData struct {
	Name    string
	Surname string
}

type Pagination struct {
	CurrentPage int
	NextPage    int
	PrevPage    int
	TotalPages  int
}

type UserData struct {
	ID int
	UserGenData
}

type UsersData []UserData
