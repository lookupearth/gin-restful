package restful

import (
	"github.com/gin-gonic/gin"

	"github.com/lookupearth/restful/model"
	"github.com/lookupearth/restful/response"
	"gorm.io/gorm"
)

// Resource 定义 Restful Resource 结构体
type Resource struct {
	*Controller
	// 必选
	DB    *gorm.DB
	Model *model.Model

	// 方法设置
	model     interface{}
	instance  interface{}
	validator IValidator
	root      IRoot
}

func NewResource(m IModel) *Resource {
	return &Resource{
		Controller: NewController(),
		DB:         m.Database(),
		Model:      model.NewModel(m),
	}
}

// Init Resource 初始化
func (resource *Resource) Init(instance interface{}, root IRoot) {
	resource.Controller.Init(instance, root)
	resource.root = root
	resource.instance = instance
	resource.validator = root.GetValidator()
	resource.validator.Register(resource.Model.New())
	if resource.Controller.HaveDetail {
		if err := resource.Model.CheckPrimaryKey(); err != nil {
			panic("model need a PrimaryKey")
		}
	}
}

func (resource *Resource) GetDB() *gorm.DB {
	return resource.DB
}

func (resource *Resource) GetModel() *model.Model {
	return resource.Model
}

func (resource *Resource) GetSerializer(m *model.Model) ISerializer {
	return NewSerializer(m, resource.validator, false)
}

func (resource *Resource) GetPartialSerializer(m *model.Model) ISerializer {
	return NewSerializer(m, resource.validator, true)
}

func (resource *Resource) Query() *gorm.DB {
	return resource.DB.Model(resource.Model.New())
}

func (resource *Resource) QueryWithContext(ctx *gin.Context) *gorm.DB {
	return resource.DB.Model(resource.Model.New()).WithContext(ctx)
}

func (resource *Resource) QueryPrimaryKey(c *gin.Context) *gorm.DB {
	primaryKey := resource.GetPrimaryKey(c)
	return resource.QueryWithContext(c).Where(resource.Model.PrimaryKey+" = ?", primaryKey)
}

func (resource *Resource) GetPrimaryKey(c *gin.Context) interface{} {
	primaryKey, err := resource.Model.ParsePrimaryKey(c.Param(":id"))
	if err != nil {
		panic(response.NewError(404, err))
	}
	return primaryKey
}
