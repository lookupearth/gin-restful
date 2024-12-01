package restful

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"

	"github.com/lookupearth/restful/response"
)

const (
	ctxResource    string = "resource"
	ctxRequestBody string = "requestBody"
)

// ContextWithResource 将 Resource 设置到 ctx 里去，之后可以使用 ResourceFromContext 读取到
func ContextWithResource(c *gin.Context, resource IResource) {
	c.Set(ctxResource, resource)
}

// ResourceFromContext 从 ctx 里读取 Resource
func ResourceFromContext(c *gin.Context) IResource {
	val, has := c.Get(ctxResource)
	if !has {
		return nil
	}
	return val.(IResource)
}

// RequestBody 请求body，为了支持修改专门设置
type RequestBody struct {
	Have  bool
	Value []byte
	Req   *http.Request
}

// NewRequestBody 创建RequestBody
func NewRequestBody(req *http.Request) (*RequestBody, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, response.NewError(500, err)
	}

	rb := &RequestBody{Req: req}
	if body != nil {
		rb.Have = true
		rb.Value = body
		rb.Req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	return rb, nil
}

func (rb *RequestBody) Set(b []byte) {
	if b != nil {
		rb.Have = true
		rb.Value = b
		rb.Req.Body = io.NopCloser(bytes.NewBuffer(b))
	} else {
		rb.Have = false
	}
}

func (rb *RequestBody) Get() []byte {
	if rb.Have {
		return rb.Value
	}
	return nil
}

// ContextWithRequestBody 将 RequestBody 设置到 ctx 里去，之后可以使用 RequestBodyFromContext 读取到
func ContextWithRequestBody(c *gin.Context, req *http.Request) error {
	rb, err := NewRequestBody(req)
	if err != nil {
		return err
	}
	c.Set(ctxRequestBody, rb)
	return nil
}

// RequestBodyFromContext 从 ctx 里读取 RequestBody
func RequestBodyFromContext(c *gin.Context) *RequestBody {
	val, has := c.Get(ctxRequestBody)
	if !has {
		return nil
	}
	return val.(*RequestBody)
}
