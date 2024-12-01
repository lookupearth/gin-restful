package mixins

import (
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/response"
)

type IPutBefore interface {
	PutBefore(*gin.Context) error
}

type IPutAfter interface {
	PutAfter(*gin.Context, interface{}) error
}

type PutMethod struct {
	Decorators []restful.HandlerDecorator

	handler  restful.HandlerFunc
	instance interface{}
}

func (c *PutMethod) InitPut(resource restful.IResource) {
	c.instance = resource
	c.handler = restful.InstallDecorators(c.put, c.Decorators)
}

// Put 全量更新（在更新数据时，未设置字段但有默认值时，会使用默认值）
func (c *PutMethod) put(ctx *gin.Context) restful.Response {
	resource := restful.ResourceFromContext(ctx)

	// before处理
	before, ok := c.instance.(IPutBefore)
	if ok {
		err := before.PutBefore(ctx)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	model := resource.GetModel()
	serializer := resource.GetSerializer(model)

	if err := serializer.ParseFromBody(ctx); err != nil {
		return response.NewError(400, err)
	}
	if err := serializer.Validate(ctx); err != nil {
		return err
	}
	updateData := serializer.ValidateData()

	// GORM 实例化
	query := resource.QueryPrimaryKey(ctx)
	// DB Update 操作
	result := query.Updates(updateData)
	restful.CheckDBResult(result)

	// after处理
	after, ok := c.instance.(IPutAfter)
	if ok {
		err := after.PutAfter(ctx, updateData)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	// 构造返回结果
	return &response.Response{
		Msg:    "",
		Status: 0,
	}
}

// Put 全量更新（在更新数据时，未设置字段但有默认值时，会使用默认值）
func (c *PutMethod) Put(ctx *gin.Context) restful.Response {
	return c.handler(ctx)
}
