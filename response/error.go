package response

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Msg    string
	Status int
	Data   interface{}
}

func NewError(status int, err error) *Error {
	err2, ok := err.(*Error)
	if ok && err2 != nil {
		return err2
	}
	return NewErrorFromMsg(status, err.Error())
}

func NewErrorFromMsg(status int, msg string) *Error {
	if status == 0 {
		status = 1
	}
	return &Error{
		Msg:    msg,
		Status: status,
	}
}

func (e *Error) Error() string {
	return e.Msg
}

func (e *Error) GetStatus() int {
	return e.Status
}

func (e *Error) GetData() interface{} {
	return e.Data
}

func (e *Error) Response(c *gin.Context) {
	res := &Response{
		Status: e.GetStatus(),
		Msg:    e.Error(),
		Data:   e.GetData(),
	}
	res.Response(c)
}
