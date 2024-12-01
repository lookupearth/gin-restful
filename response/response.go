package response

import (
	"github.com/gin-gonic/gin"
)

// Response 响应内容
// swagger:response Response
type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
	Echo   int         `json:"echo,omitempty"`
	Total  *int64      `json:"total,omitempty"`
	From   string      `json:"from,omitempty"`
	LogID  string      `json:"logid,omitempty"`
}

func (response *Response) Response(c *gin.Context) {
	var httpCode int
	if response.Status == 0 {
		httpCode = 200
	} else if response.Status >= 400 && response.Status <= 599 {
		httpCode = response.Status
	} else {
		httpCode = 500
	}
	c.JSON(httpCode, response)
}

// SetLogID 为 Response 结构体设置 LogID 字段
func (response *Response) SetLogID(logid string) {
	response.LogID = logid
}

type TextResponse struct {
	Status  int
	Text    string
	Headers map[string]string
}

func (response *TextResponse) Response(c *gin.Context) {
	var httpCode int
	if response.Status == 0 {
		httpCode = 200
	} else if response.Status >= 400 && response.Status <= 599 {
		httpCode = response.Status
	} else {
		httpCode = 500
	}
	for key, value := range response.Headers {
		c.Header(key, value)
	}
	c.String(httpCode, response.Text)
}
