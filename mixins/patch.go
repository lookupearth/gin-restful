package mixins

import (
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/response"
)

type IPatchBefore interface {
	PatchBefore(*gin.Context) error
}

type IPatchAfter interface {
	PatchAfter(*gin.Context, interface{}) error
}

type PatchMethod struct {
	Decorators   []restful.HandlerDecorator
	WithDefaults []string

	handler  restful.HandlerFunc
	instance interface{}
}

func (c *PatchMethod) InitPatch(resource restful.IResource) {
	c.instance = resource
	c.handler = restful.InstallDecorators(c.patch, c.Decorators)
}

// Patch 部分更新
//
//	仅更新用户提交的字段（遵循GORM结构体中读写限制的字段规则）
func (c *PatchMethod) patch(ctx *gin.Context) restful.Response {
	resource := restful.ResourceFromContext(ctx)

	// before处理
	before, ok := c.instance.(IPatchBefore)
	if ok {
		err := before.PatchBefore(ctx)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	model := resource.GetModel()
	serializer := resource.GetPartialSerializer(model).WithDefaults(c.WithDefaults)

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
	after, ok := c.instance.(IPatchAfter)
	if ok {
		var err error
		err = after.PatchAfter(ctx, updateData)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	return &response.Response{
		Msg:    "",
		Status: 0,
	}
}

// Patch 部分更新
func (c *PatchMethod) Patch(ctx *gin.Context) restful.Response {
	return c.handler(ctx)
}
