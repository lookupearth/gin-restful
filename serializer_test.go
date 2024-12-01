// Package restful
package restful

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lookupearth/restful/field"
	"net/http/httptest"
	"reflect"
	"testing"

	valid "github.com/go-playground/validator/v10"
	"github.com/lookupearth/restful/model"
)

type Activity struct {
	ID      int64           `gorm:"column:cid;primaryKey;->" json:"id" validate:"required"` // 自增id
	Status  *field.ExInt64  `gorm:"column:-" json:"status" validate:"required"`             // 状态
	Status2 field.ExInt64   `gorm:"column:cstatus2" json:"status2" default:"1" `            // 状态
	Name    string          `gorm:"column:cname" json:"name" validate:"required"`           // 名称
	Name2   *field.ExString `default:""`
	Name3   string          `gorm:"column:name3" json:"name3,readonly" default:"123"` // 名称
}

func newSerializer(data interface{}, partial bool) *Serializer {
	m := model.NewModel(data)
	va := &Validator{
		Validator: valid.New(),
	}
	return NewSerializer(m, va, partial)
}

func TestSerializerParse(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	s := newSerializer(&Activity{}, false)
	b := []byte(`{"id":"1","status":456,"name3":"789"}`)
	err := s.Parse(ctx, b)
	if err == nil {
		t.Errorf("Serializer.Parse fail, error=%v", err)
	}
	fmt.Println("TestSerializerParse", err)
}

func TestSerializerValid(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	s := newSerializer(&Activity{}, false)
	b := []byte(`{"id":1,"status":456}`)
	if err := s.Parse(ctx, b); err != nil {
		t.Errorf("Serializer.Parse fail, error=%v", err)
	}
	err := s.Validate(ctx)
	if err == nil {
		t.Errorf("Serializer.Validate fail, error=%v", err)
	}
	fmt.Println("TestSerializerValid", err)
}

func TestSerializer(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	s := newSerializer(&Activity{}, false)
	b := []byte(`{"id":1,"status":456,"name":"aaa","name3":"789"}`)
	err := s.Parse(ctx, b)
	if err != nil {
		t.Errorf("Serializer.Parse fail, error=%v", err)
	}
	if err := s.Validate(ctx); err != nil {
		t.Errorf("Serializer.Validate fail, error=%v", err)
	}
	name2 := field.ExString("")
	status := field.ExInt64(456)
	validateData := map[string]interface{}{
		"cid":      int64(1),
		"cstatus2": field.ExInt64(1),
		"cname":    "aaa",
		"name2":    &name2,
		"name3":    "123",
	}
	jsonData := map[string]interface{}{
		"id":      int64(1),
		"status":  &status,
		"status2": field.ExInt64(1),
		"name":    "aaa",
		"Name2":   &name2,
		"name3":   "123",
	}
	if !reflect.DeepEqual(validateData, s.ValidateData()) {
		t.Errorf("Serializer.ValidateData fail, expect=%v got=%v", validateData, s.ValidateData())
	}
	if !reflect.DeepEqual(jsonData, s.JsonData()) {
		t.Errorf("Serializer.ValidateData fail, expect=%v got=%v", validateData, s.ValidateData())
	}
}

func TestPartial(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	s := newSerializer(&Activity{}, true).WithDefaults([]string{"Name3"})
	data := map[string]string{
		"id":    "1",
		"name":  "ddd",
		"name3": "888",
	}
	if err := s.ParseFromQuery(ctx, data); err != nil {
		t.Errorf("Serializer.Parse fail, error=%v", err)
	}
	if err := s.Validate(ctx); err != nil {
		t.Errorf("Serializer.Validate fail, error=%v", err)
	}
	validateData := map[string]interface{}{
		"cid":   int64(1),
		"cname": "ddd",
		"name3": "123",
	}
	jsonData := map[string]interface{}{
		"id":    int64(1),
		"name":  "ddd",
		"name3": "123",
	}
	fmt.Println(s.ValidateData())
	fmt.Println(s.JsonData())
	if !reflect.DeepEqual(validateData, s.ValidateData()) {
		t.Errorf("Serializer.ValidateData fail, expect=%v got=%v", validateData, s.ValidateData())
	}
	if !reflect.DeepEqual(jsonData, s.JsonData()) {
		t.Errorf("Serializer.ValidateData fail, expect=%v got=%v", validateData, s.ValidateData())
	}
	v, err := s.Get("id")
	if err != nil {
		t.Errorf("Serializer.Get fail, error=%v", err)
	}
	if v != int64(1) {
		t.Errorf("Serializer.ValidateData fail, expect=%v got=%v", int64(1), v)
	}
	if _, err := s.Get("id2"); err == nil {
		t.Errorf("Serializer.Get fail, error=%v", err)
	}
	v2 := s.GetWithDefault("name", "aaa")
	v2s, ok := v2.(string)
	if !ok {
		t.Errorf("Serializer.GetWithDefault type fail, expect=string got=%v", reflect.TypeOf(v2).Kind())
	}
	if v2s != "ddd" {
		t.Errorf("Serializer.GetWithDefault fail, expect=%v got=%v", "ddd", v2s)
	}
	v3 := s.GetWithDefault("status2", field.ExInt64(2))
	v3i, ok := v3.(field.ExInt64)
	if !ok {
		t.Errorf("Serializer.GetWithDefault type fail, expect=field.ExInt64 got=%v", reflect.TypeOf(v3).Kind())
	}
	if v3i != 2 {
		t.Errorf("Serializer.GetWithDefault fail, expect=%v got=%v", 2, v3i)
	}
}
