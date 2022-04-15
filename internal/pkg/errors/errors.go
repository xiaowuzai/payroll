package errors

type MError struct {
	Code int   //http code
	Message string // 对外展示
	InnerMessage string // 内部定位问题
}

func (m *MError)Error() string{
	return m.Message
}

func New(code int, message, innerMessage string ) error{
	return &MError{
		Code: code,
		Message: message,
		InnerMessage: innerMessage,
	}
}

func Wrap(err error, message string) error {
	return nil
}