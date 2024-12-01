package mixins

import (
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/response"
)

type IPostBefore interface {
	PostBefore(*gin.Context) error
}

type IPostAfter interface {
	PostAfter(*gin.Context, interface{}, interface{}) error
}

type PostMethod struct {
	Decorators []restful.HandlerDecorator

	handler  restful.HandlerFunc
	instance interface{}
}

func (c *PostMethod) InitPost(resource restful.IResource) {
	c.instance = resource
	c.handler = restful.InstallDecorators(c.post, c.Decorators)
}

// Post 添加数据（在新增数据时，未设置字段但有默认值时，会使用默认值）
func (c *PostMethod) post(ctx *gin.Context) restful.Response {
	resource := restful.ResourceFromContext(ctx)

	// before处理
	before, ok := c.instance.(IPostBefore)
	if ok {
		err := before.PostBefore(ctx)
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

	// GORM 实例化
	query := resource.QueryWithContext(ctx)
	query = query.Begin()
	defer func() {
		if r := recover(); r != nil {
			query.Rollback()
			panic(r)
		}
	}()
	// DB Create 操作
	validData := serializer.ValidateData()
	result := query.Create(validData)
	restful.CheckDBResult(result)

	// 获取新添加数据的ID
	var ret map[string]interface{}
	query.Select("last_insert_id() as id").Limit(1).Find(&ret)
	query.Commit()
	// DRDS环境，as未生效
	if v, ok := ret["last_insert_id()"]; ok {
		ret["id"], _ = model.ParsePrimaryKey(string(v.([]byte)))
	}

	// after处理
	after, ok := c.instance.(IPostAfter)
	if ok {
		err := after.PostAfter(ctx, ret["id"], validData)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	// 构造返回结果
	return &response.Response{
		Msg:    "",
		Status: 0,
		Data:   ret,
	}
}

// Post 添加数据（在新增数据时，未设置字段但有默认值时，会使用默认值）
func (c *PostMethod) Post(ctx *gin.Context) restful.Response {
	return c.handler(ctx)
}
