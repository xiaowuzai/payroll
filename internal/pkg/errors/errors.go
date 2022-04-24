package errors

import (
	"errors"
	"net/http"
)

type MyError struct {
	message  string `json:"message"`      // 对外提供
	inner    string `json:"interMessage"` // 对内定位使用
	httpCode int    `json:"httpCode"`     // HTTP Code
	code     string `json:"code"`
}

func (m *MyError) Error() string {
	return m.message
}

func (m *MyError) Message() string {
	return m.message
}

func (m *MyError) Inner() string {
	return m.inner
}

func (m *MyError) Code() string {
	return m.code
}

func (m *MyError) HTTPCode() int {
	return m.httpCode
}

// http code 500
func New(message string) error {
	return &MyError{
		code:     CodeError,
		message:  message,
		inner:    message,
		httpCode: http.StatusInternalServerError,
	}
}

// code is 500,
func NewInner(message, inner string) error {
	return &MyError{
		code:     CodeError,
		message:  message,
		inner:    inner,
		httpCode: http.StatusInternalServerError,
	}
}

func NewCode(message, inner string, httpCode int) error {
	return &MyError{
		message:  message,
		inner:    inner,
		httpCode: httpCode,
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Wrap() {
}

func UnWrap() {
}

// 数据未找到，返回400 并返回提示信息
func DataNotFound(inner string) error {
	return &MyError{
		message:  ErrNotFound,
		inner:    inner,
		httpCode: http.StatusBadRequest,
		code:     CodeNotFound,
	}
}

// 获取数据出错, 返回 500、
func ErrDataGet(inner string) error {
	return NewInner(ErrGet, inner)
}

// 数据插入出错， 返回500
func ErrDataInsert(inner string) error {
	return NewInner(ErrInsert, inner)
}

// 数据删除出错 返回500
func ErrDataDelete(inner string) error {
	return NewInner(ErrDelete, inner)
}

// 数据删除出错 返回500
func ErrDataUpdate(inner string) error {
	return NewInner(ErrUpdate, inner)
}
