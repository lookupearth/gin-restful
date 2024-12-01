package restful

import (
	"context"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
	"github.com/lookupearth/restful/model"
	"github.com/lookupearth/restful/response"
	"gorm.io/gorm"
)

type Response interface {
	Response(*gin.Context)
}

type HandlerFunc func(*gin.Context) Response

type IModel interface {
	Database() *gorm.DB
}

type IRoot interface {
	RegisterResource(string, IController)
	Mount(*gin.RouterGroup)
	GetValidator() IValidator
	Print(string)
	Validate() *validator.Validate
}

type ISerializer interface {
	WithDefaults([]string) ISerializer
	Parse(*gin.Context, []byte) error
	ParseFromQuery(*gin.Context, map[string]string) error
	ParseFromBody(*gin.Context) error
	Validate(*gin.Context) *response.Error
	ValidateData() map[string]interface{}
	GetWithDefault(string, interface{}) interface{}
	Get(string) (interface{}, error)
	JsonData() map[string]interface{}
	StructData() interface{}
}

// IController 接口定义，启动时可以调用
type IController interface {
	Init(interface{}, IRoot)
	Mount(*gin.RouterGroup, string)
	RegisterMethod(MethodType, HttpMethod, string, HandlerFunc)
	Print(string)
}

// IResource 接口定义，运行时可以调用
type IResource interface {
	// Query 获取查询句柄，推荐优先用QueryWithContext/QueryPrimaryKey
	Query() *gorm.DB
	// QueryWithContext 获取已绑定context的查询句柄
	QueryWithContext(*gin.Context) *gorm.DB
	// QueryPrimaryKey 获取添加PrimaryKey条件的查询句柄
	QueryPrimaryKey(*gin.Context) *gorm.DB

	// GetDB 获取gorm DB实例
	GetDB() *gorm.DB
	// GetModel 获取资源的Model封装
	GetModel() *model.Model
	// GetSerializer 获取具体Model的序列化实例
	GetSerializer(*model.Model) ISerializer
	// GetPartialSerializer 获取具体Model的Partial序列化实例
	GetPartialSerializer(*model.Model) ISerializer

	// GetPrimaryKey 获取PrimaryKey
	GetPrimaryKey(*gin.Context) interface{}
}

type IList interface {
	List(*gin.Context) Response
}

type IListInit interface {
	InitList(IResource)
}

type IGet interface {
	Get(*gin.Context) Response
}

type IGetInit interface {
	InitGet(IResource)
}

type IPost interface {
	Post(*gin.Context) Response
}

type IPostInit interface {
	InitPost(IResource)
}

type IPut interface {
	Put(*gin.Context) Response
}

type IPutInit interface {
	InitPut(IResource)
}

type IPatch interface {
	Patch(*gin.Context) Response
}

type IPatchInit interface {
	InitPatch(IResource)
}

type IDelete interface {
	Delete(*gin.Context) Response
}

type IDeleteInit interface {
	InitDelete(IResource)
}

type ISearch interface {
	Search(*gin.Context) Response
}

type ISearchInit interface {
	InitSearch(IResource)
}

type IDecorator interface {
	GetDecorators() []HandlerDecorator
}

type IValidator interface {
	Register(interface{})
	Validate(context.Context, interface{}) *response.Error
	ValidatePartial(context.Context, interface{}, []string) *response.Error
}
