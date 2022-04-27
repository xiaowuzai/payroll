package handler

type RequestId struct {
	Id string `json:"id" binding:"required"`
}

var ErrIdEmpty = "id 不能为空"
