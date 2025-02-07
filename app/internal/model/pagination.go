package model

type Pagination struct {
	CurrentPage int
	NextPage    int
	PrevPage    int
	TotalPages  int
}

type UsersPaginationMetaData struct {
	Users      UsersData
	Pagination Pagination
}
