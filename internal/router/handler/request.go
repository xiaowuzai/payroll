package handler

type RequestId struct {
	Id string `json:"id" binding:"required"`
}

var idIsEmpty = "id 不能为空"
