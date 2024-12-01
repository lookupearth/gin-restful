// Package restful
package restful

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/lookupearth/restful/response"
)

const TableNameDemo = "demo"

// Gorm 说明文档：
//   [This] V2: https://gorm.io/docs/models.html#Creating-Updating-Time-Unix-Milli-Nano-Seconds-Tracking
// 	        V1： https://v1.gorm.io/docs/models.html
// 更新数据时，只会兼容Gorm中的部分属性
//
// 结构体中定义的属性由两部分程序接管：
// 	1. column/primaryKey/default 由自由程序接管处理
//  2. <-:create/->/<- 由GORM程序接管处理(在执行GORM中方法的时会自动识别，无需程序关注)

// Activity mapped from table <activity>
type DemoTable struct {
	ID       int64  `gorm:"column:id;primaryKey;->" json:"id"`              // 自增id
	Name     string `gorm:"column:name" json:"name"`                        // 活动名称
	Status   int32  `gorm:"column:status" json:"status" default:"2"`        // 状态
	ISDelete int32  `gorm:"column:is_delete;deleteKey;<-" json:"is_delete"` // 删除标记字段
}

func (d *DemoTable) ValidateCtx(ctx context.Context, sl validator.StructLevel) {
	fmt.Println("DemoTable ValidateCtx")
}

func (d *DemoTable) Database() *gorm.DB {
	return &gorm.DB{}
}

// TableName Activity's table name
func (*DemoTable) TableName() string {
	return TableNameDemo
}

type DemoResource struct {
	*Resource
}

// Get 查询单条数据
func (d *DemoResource) publish(c *gin.Context) Response {
	return &response.Response{
		Msg:    "",
		Status: 0,
		Data:   "publish test",
	}
}

func NewDemo() *DemoResource {
	demo := &DemoResource{
		Resource: NewResource(&DemoTable{}),
	}

	demo.RegisterMethod(DetailMethod, HTTPMethodGet, "publish", demo.publish)
	return demo
}

func testValidateCtx(ctx context.Context, fl validator.FieldLevel) bool {
	fmt.Println("testValidateCtx")
	return true
}

func TestRestful(t *testing.T) {
	app := gin.New()
	apiRouter := app.Group("/api")
	root := New()
	_ = root.Validate().RegisterValidationCtx("test", testValidateCtx)
	root.RegisterResource("/demo/demo", NewDemo())
	root.Mount(apiRouter)
	root.Print("")
}
