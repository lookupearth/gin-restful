package mixins

import (
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/response"
)

type IGetBefore interface {
	GetBefore(*gin.Context) error
}

type IGetAfter interface {
	// GetAfter 后置操作，加工/替换返回值
	GetAfter(*gin.Context, interface{}) (interface{}, error)
}

type GetMethod struct {
	Decorators []restful.HandlerDecorator

	handler  restful.HandlerFunc
	instance interface{}
}

func (c *GetMethod) InitGet(resource restful.IResource) {
	c.instance = resource
	c.handler = restful.InstallDecorators(c.get, c.Decorators)
}

func (c *GetMethod) get(ctx *gin.Context) restful.Response {
	resource := restful.ResourceFromContext(ctx)

	// before处理
	before, ok := c.instance.(IGetBefore)
	if ok {
		err := before.GetBefore(ctx)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	model := resource.GetModel()
	// GORM 实例化
	data := model.New()
	query := resource.QueryPrimaryKey(ctx)
	// DB Query 操作
	result := query.First(data)
	restful.CheckDBResult(result)

	// after处理
	after, ok := c.instance.(IGetAfter)
	if ok {
		var err error
		data, err = after.GetAfter(ctx, data)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	return &response.Response{
		Msg:    "",
		Status: 0,
		Data:   data,
	}
}

// Get 查询单条数据
func (c *GetMethod) Get(ctx *gin.Context) restful.Response {
	return c.handler(ctx)
}
