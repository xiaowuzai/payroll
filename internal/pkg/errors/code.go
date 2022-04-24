package errors

// 错误码格式
// 服务编号 + 错误类型 + 错误码编号
/*
	服务编号：01

	错误类型：
		数据库错误: 01
*/

var (
	ErrNotFound = "数据不存在"
	ErrInsert   = "数据插入失败"
	ErrGet      = "获取数据失败"
	ErrUpdate   = "修改数据失败"
	ErrDelete   = "数据删除失败"
)

var (
	CodeSuccess  = "0101"
	CodeError    = "0102"
	CodeDefault  = "0103"
	CodeNotFound = "0104"
)
