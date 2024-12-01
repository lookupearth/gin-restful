package restful

import (
	"errors"
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful/response"
	"gorm.io/gorm"
)

// GetQuery 获取get请求
func GetQuery(c *gin.Context) map[string]string {
	var params = make(map[string]string)

	for k, v := range c.Request.URL.Query() {
		params[k] = v[0]
	}

	return params
}

func CheckDBResult(result *gorm.DB) {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic(response.NewError(404, result.Error))
	}
	if result.Error != nil {
		panic(response.NewError(500, result.Error))
	}
}

func InstallDecorators(handler HandlerFunc, decorators []HandlerDecorator) HandlerFunc {
	if decorators != nil {
		for _, decorator := range decorators {
			handler = decorator(handler)
		}
	}
	return handler
}
