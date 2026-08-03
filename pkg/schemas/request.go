package schemas

type IdRequest struct {
	Id int64 `json:"id" form:"id"`
}

type NameRequest struct {
	Name string `json:"name" form:"name"`
}
