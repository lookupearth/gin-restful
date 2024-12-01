package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/field"
	"github.com/lookupearth/restful/mixins"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Demo struct {
	*restful.Resource
	*mixins.GetMethod
	*mixins.ListMethod
	*mixins.PostMethod
	*mixins.PatchMethod
	*mixins.PutMethod
	*mixins.DeleteMethod
}

type searchParams struct {
	Name   string          `gorm:"column:name" json:"name"`
	Status int32           `gorm:"column:status" json:"status"`
	Start  field.Timestamp `gorm:"column:create_time;" json:"start" operate:">="`
	End    field.Timestamp `gorm:"column:create_time;" json:"end" operate:"<"`
}

const TableName = "demo"

type Model struct {
	ID      int64            `gorm:"column:id;primaryKey;->" json:"id"`
	Name    string           `gorm:"column:name" json:"name" validate:"required"`
	Status  field.ExInt64    `gorm:"column:status" json:"status" default:"2"`
	Content field.JSONObject `gorm:"column:content" json:"content" default:"{}"`

	CreateTime field.Timestamp `gorm:"column:create_time" json:"create_time"`
	UpdateTime field.Timestamp `gorm:"column:update_time;->" json:"update_time"`
}

func (*Model) TableName() string {
	return TableName
}

func (m *Model) Database() *gorm.DB {
	return db
}

func NewDemo() *Demo {
	demo := &Demo{
		Resource:  restful.NewResource(&Model{}),
		GetMethod: &mixins.GetMethod{},
		ListMethod: &mixins.ListMethod{
			Offset:       0,
			Limit:        10,
			OrderBy:      []string{"id desc"},
			SearchParams: &searchParams{},
		},
		PostMethod:   &mixins.PostMethod{},
		PatchMethod:  &mixins.PatchMethod{},
		PutMethod:    &mixins.PutMethod{},
		DeleteMethod: &mixins.DeleteMethod{},
	}
	return demo
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	r := gin.Default()
	api := r.Group("/api")

	root := restful.New()
	root.RegisterResource("/demo", NewDemo())
	root.Mount(api)
	root.Print("/api")
	r.Run("localhost:8080")
}
