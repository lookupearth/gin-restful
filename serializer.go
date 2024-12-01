package restful

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"

	"github.com/lookupearth/restful/model"
	"github.com/lookupearth/restful/response"
)

type Serializer struct {
	model        *model.Model
	validator    IValidator
	partial      bool
	rawData      map[string]interface{}
	structData   interface{}
	structValue  reflect.Value
	withDefaults []string
}

// NewSerializer Serializer 实例化，设置partial=true后，不会解析默认值，不会校验未传入的字段
func NewSerializer(m *model.Model, v IValidator, partial bool) *Serializer {
	s := &Serializer{
		model:     m,
		validator: v,
		partial:   partial,
	}
	return s
}

// WithDefaults 固定解析部分字段的默认值
func (s *Serializer) WithDefaults(keys []string) ISerializer {
	s.withDefaults = keys
	return s
}

// setRawData 设置原始数据，过滤readonly部分
func (s *Serializer) setRawData(input map[string]interface{}) {
	rawData := make(map[string]interface{})
	for k, v := range input {
		if name, ok := s.model.Json2Name[k]; ok {
			field := s.model.Name2Field[name]
			if field.ReadOnly() {
				continue
			}
		}
		rawData[k] = v
	}
	s.rawData = rawData
}

// Parse 注意返回的是model的指针
func (s *Serializer) Parse(c *gin.Context, b []byte) error {
	data, err := s.model.Parse(b)
	if err != nil {
		return err
	}
	s.structData = data
	var rawData map[string]interface{}
	if err := json.Unmarshal(b, &rawData); err != nil {
		return err
	}
	s.setRawData(rawData)
	if s.partial == false {
		if err := s.model.ParseDefault(c, data, s.rawData); err != nil {
			return err
		}
	}
	if s.withDefaults != nil {
		if err := s.model.ParseDefaultWithKeys(c, data, s.withDefaults, s.rawData); err != nil {
			return err
		}
	}
	s.structValue = reflect.Indirect(reflect.ValueOf(s.structData))
	return nil
}

// ParseFromQuery 注意返回的是model的指针
func (s *Serializer) ParseFromQuery(c *gin.Context, query map[string]string) error {
	data, err := s.model.ParseFromQuery(query)
	if err != nil {
		return err
	}
	s.structData = data
	rawData := make(map[string]interface{})
	for k, v := range query {
		rawData[k] = v
	}
	s.setRawData(rawData)
	if s.partial == false {
		if err := s.model.ParseDefault(c, data, s.rawData); err != nil {
			return err
		}
	}
	if s.withDefaults != nil {
		if err := s.model.ParseDefaultWithKeys(c, data, s.withDefaults, s.rawData); err != nil {
			return err
		}
	}
	s.structValue = reflect.Indirect(reflect.ValueOf(s.structData))
	return nil
}

// ParseFromBody 从request解析，对Parse的封装
func (s *Serializer) ParseFromBody(c *gin.Context) error {
	requestBody := RequestBodyFromContext(c)
	body := requestBody.Get()
	if body == nil {
		return response.NewError(500, errors.New("body is nil"))
	}
	return s.Parse(c, body)
}

// Validate 根据struct tag校验输入数据
func (s *Serializer) Validate(c *gin.Context) *response.Error {
	if s.structData == nil {
		return response.NewErrorFromMsg(500, "Parse or ParseFromQuery should call before Validate")
	}
	if s.validator != nil {
		if s.partial == false {
			return s.validator.Validate(c, s.structData)
		} else {
			return s.validator.ValidatePartial(c, s.structData, s.model.FieldNames(s.rawData))
		}
	}
	return nil
}

// ValidateData 务必在 Validate 后调用，这里为了使用方便，未做错误检查
// 否则只能返回两个参数，不便于调用
// 该函数返回的map，key为数据库列，不一定是json/form的key
// 注意，该方法应仅在需要存储数据到db时使用，参数解析务必使用 JsonData 方法
func (s *Serializer) ValidateData() map[string]interface{} {
	validateData := make(map[string]interface{})
	if !s.structValue.IsValid() {
		return validateData
	}
	for name, _ := range s.model.Name2Field {
		column, ok := s.model.Name2Column[name]
		if !ok {
			continue
		}
		fv := s.structValue.FieldByName(name)
		if jsonKey, ok := s.model.Name2Json[name]; ok {
			if _, ok := s.rawData[jsonKey]; ok {
				validateData[column] = fv.Interface()
				continue
			}
		}
	}
	return validateData
}

// GetWithDefault 务必在 Validate 后调用，这里为了使用方便，未做错误检查
// 否则只能返回两个参数，不便于调用
// 该函数的key是json/form的key
func (s *Serializer) GetWithDefault(key string, value interface{}) interface{} {
	if !s.structValue.IsValid() {
		return value
	}
	if _, ok := s.rawData[key]; ok {
		if name, ok := s.model.Json2Name[key]; ok {
			fv := s.structValue.FieldByName(name)
			return fv.Interface()
		}
	}
	return value
}

// Get 务必在 Validate 后调用，这里为了使用方便，未做错误检查
// 该函数的key是json/form的key
// key不存在时，会返回error
func (s *Serializer) Get(key string) (interface{}, error) {
	if !s.structValue.IsValid() {
		return nil, errors.New("serializer is empty")
	}
	if _, ok := s.rawData[key]; ok {
		if name, ok := s.model.Json2Name[key]; ok {
			fv := s.structValue.FieldByName(name)
			return fv.Interface(), nil
		}
	}
	return nil, errors.New("key <" + key + "> not exists")
}

// JsonData 务必在 Validate 后调用，这里为了使用方便，未做错误检查
// 否则只能返回两个参数，不便于调用
// 该函数返回的map，key为数据库列，不一定是json/form的key
func (s *Serializer) JsonData() map[string]interface{} {
	jsonData := make(map[string]interface{})
	if !s.structValue.IsValid() {
		return jsonData
	}
	for name, _ := range s.model.Name2Field {
		jsonKey, ok := s.model.Name2Json[name]
		if !ok {
			continue
		}
		fv := s.structValue.FieldByName(name)
		if _, ok := s.rawData[jsonKey]; ok {
			jsonData[jsonKey] = fv.Interface()
			continue
		}
	}
	return jsonData
}

// StructData 务必在 Validate 后调用，这里为了使用方便，未做错误检查
// 该函数返回的类型与model相同，可以断言转回
func (s *Serializer) StructData() interface{} {
	return s.structData
}
