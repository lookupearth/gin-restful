package mixins

import (
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/response"
)

type IDeleteBefore interface {
	DeleteBefore(*gin.Context) error
}

type IDeleteAfter interface {
	// DeleteAfter 后置操作，加工/替换返回值
	DeleteAfter(*gin.Context, interface{}) error
}

type DeleteMethod struct {
	Decorators []restful.HandlerDecorator

	handler  restful.HandlerFunc
	instance interface{}
}

func (c *DeleteMethod) InitDelete(resource restful.IResource) {
	c.instance = resource
	c.handler = restful.InstallDecorators(c.delete, c.Decorators)
}

func (c *DeleteMethod) delete(ctx *gin.Context) restful.Response {
	resource := restful.ResourceFromContext(ctx)

	// before处理
	before, ok := c.instance.(IDeleteBefore)
	if ok {
		err := before.DeleteBefore(ctx)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	model := resource.GetModel()

	// GORM 实例化
	query := resource.QueryPrimaryKey(ctx)

	data := model.New()
	result := query.Delete(data)
	restful.CheckDBResult(result)

	// after处理
	after, ok := c.instance.(IDeleteAfter)
	if ok {
		var err error
		err = after.DeleteAfter(ctx, data)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	return &response.Response{
		Msg:    "",
		Status: 0,
	}
}

// Delete 标记删除数据，不支持批量
func (c *DeleteMethod) Delete(ctx *gin.Context) restful.Response {
	return c.handler(ctx)
}
