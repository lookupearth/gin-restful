package mixins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lookupearth/restful/field"
	"strings"

	"github.com/lookupearth/restful"
	"github.com/lookupearth/restful/model"
	"github.com/lookupearth/restful/response"
	"gorm.io/gorm"
)

type IListBefore interface {
	ListBefore(*gin.Context, map[string]string) error
}

type IListAfter interface {
	// ListAfter 后置操作，对返回列表元素中数据进行加工/数据结构修改，不适合对列表元素进行增删
	ListAfter(*gin.Context, interface{}) (interface{}, error)
}

type ListParams struct {
	Echo    field.ExInt64  `json:"echo"`
	Page    field.ExInt64  `json:"page"`
	Size    field.ExInt64  `json:"size"`
	Offset  field.ExInt64  `json:"offset"`
	Limit   field.ExInt64  `json:"limit"`
	OrderBy field.ExString `json:"orderBy"`
	Search  field.ExString `json:"search"`
}

type ListMethod struct {
	Offset       int
	Limit        int
	OrderBy      []string
	SearchFields []string

	SearchParams interface{}
	Decorators   []restful.HandlerDecorator

	ListModel   *model.Model
	SearchModel *model.Model
	handler     restful.HandlerFunc
	instance    interface{}
}

func (c *ListMethod) InitList(resource restful.IResource) {
	c.instance = resource
	c.handler = restful.InstallDecorators(c.list, c.Decorators)
	c.ListModel = model.NewModel(&ListParams{})
	if c.SearchParams != nil {
		c.SearchModel = model.NewModel(c.SearchParams)
	}
}

func (c *ListMethod) SearchQuery(query *gorm.DB, search string) *gorm.DB {
	search = strings.TrimSpace(search)
	if len(c.SearchFields) > 0 && len(search) > 0 {
		terms := strings.Split(search, " ")
		for _, term := range terms {
			for _, field := range c.SearchFields {
				query = query.Or(fmt.Sprintf("`%s` LIKE ?", field), "%"+term+"%")
			}
		}
		return query
	}
	return nil
}

func (c *ListMethod) ParseOrderBy(orderStr string) []string {
	ret := make([]string, 0)
	orders := strings.Split(orderStr, ",")
	for _, order := range orders {
		order = strings.TrimSpace(order)
		if len(order) > 0 {
			ret = append(ret, order)
		}
	}
	return ret
}

func (c *ListMethod) Paginate(query *gorm.DB, listData *ListParams) *gorm.DB {
	page := int(listData.Page)
	size := int(listData.Size)
	offset := int(listData.Offset)
	limit := int(listData.Limit)
	if page != 0 || size != 0 {
		if page == 0 {
			page = 1
		}
		if size == 0 {
			size = c.Limit
		}
		offset = (page - 1) * size
		limit = size
	} else {
		if offset == 0 {
			offset = c.Offset
		}
		if limit == 0 {
			limit = c.Limit
		}
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	return query
}

// List 查询数据列表，遵循 Restful 查询规范
func (c *ListMethod) list(ctx *gin.Context) restful.Response {
	resource := restful.ResourceFromContext(ctx)

	params := restful.GetQuery(ctx)
	// before处理
	before, ok := c.instance.(IListBefore)
	if ok {
		err := before.ListBefore(ctx, params)
		if err != nil {
			return response.NewError(500, err)
		}
	}
	// GET 请求参数校验+提取
	m := resource.GetModel()
	query := resource.QueryWithContext(ctx)
	if c.SearchModel != nil {
		searchSerializer := resource.GetSerializer(c.SearchModel)
		if err := searchSerializer.ParseFromQuery(ctx, params); err != nil {
			return response.NewError(400, err)
		}
		if err := searchSerializer.Validate(ctx); err != nil {
			return err
		}
		searchData := searchSerializer.JsonData()
		for key, value := range searchData {
			query = c.SearchModel.Where(query, key, value)
		}
	}
	listSerializer := resource.GetSerializer(c.ListModel)
	if err := listSerializer.ParseFromQuery(ctx, params); err != nil {
		return response.NewError(400, err)
	}
	if err := listSerializer.Validate(ctx); err != nil {
		return err
	}
	listData := listSerializer.StructData().(*ListParams)
	// like检索逻辑，实际使用注意性能
	subQuery := c.SearchQuery(resource.Query(), string(listData.Search))
	if subQuery != nil {
		query = query.Where(subQuery)
	}
	// 获取数量
	var total int64
	query.Count(&total)

	// 排序
	orders := c.OrderBy
	if len(listData.OrderBy) > 0 {
		orders = c.ParseOrderBy(string(listData.OrderBy))
	}
	if len(orders) > 0 {
		for _, order := range orders {
			query = query.Order(order)
		}
	}
	// 分页
	query = c.Paginate(query, listData)

	results := m.NewSlice()
	result := query.Find(results)
	restful.CheckDBResult(result)

	echo := int(listData.Echo)

	// after处理
	after, ok := c.instance.(IListAfter)
	if ok {
		var err error
		results, err = after.ListAfter(ctx, results)
		if err != nil {
			return response.NewError(500, err)
		}
	}

	return &response.Response{
		Msg:    "",
		Status: 0,
		Data:   results,
		Total:  &total,
		Echo:   echo,
	}
}

// List 查询数据列表，遵循 Restful 查询规范
func (c *ListMethod) List(ctx *gin.Context) restful.Response {
	return c.handler(ctx)
}
