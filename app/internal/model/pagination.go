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

type UsersCursorBasedMetaData struct {
	Users      UsersData
	NextCursor int
}
