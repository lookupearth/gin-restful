package restful

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"

	"github.com/lookupearth/restful/response"
)

type logIDSetter interface {
	SetLogID(string)
}

// Controller 定义 Restful 路由转发相关 api 结构体
type Controller struct {
	HaveDetail  bool
	urlHandlers map[string]map[HttpMethod]HandlerFunc

	// init阶段初始化
	instance interface{}
}

func NewController() *Controller {
	return &Controller{
		HaveDetail:  false,
		urlHandlers: make(map[string]map[HttpMethod]HandlerFunc),
	}
}

// Init Resource 初始化
func (ctrl *Controller) Init(instance interface{}, root IRoot) {
	ctrl.instance = instance
}

// Mount 将 Resource 方法注册到路由
func (ctrl *Controller) Mount(router *gin.RouterGroup, urlPath string) {
	instance := ctrl.instance
	var decorators []HandlerDecorator
	decorative, ok := instance.(IDecorator)
	if ok {
		decorators = decorative.GetDecorators()
	}

	getlist, ok := instance.(IList)
	if ok {
		if init, ok := instance.(IListInit); ok {
			init.InitList(instance.(IResource))
		}
		ctrl.RegisterMethod(ListMethod, HTTPMethodGet, "", getlist.List)
	}
	get, ok := instance.(IGet)
	if ok {
		if init, ok := instance.(IGetInit); ok {
			init.InitGet(instance.(IResource))
		}
		ctrl.RegisterMethod(DetailMethod, HTTPMethodGet, "", get.Get)
	}
	post, ok := instance.(IPost)
	if ok {
		if init, ok := instance.(IPostInit); ok {
			init.InitPost(instance.(IResource))
		}
		ctrl.RegisterMethod(ListMethod, HTTPMethodPost, "", post.Post)
	}
	put, ok := instance.(IPut)
	if ok {
		if init, ok := instance.(IPutInit); ok {
			init.InitPut(instance.(IResource))
		}
		ctrl.RegisterMethod(DetailMethod, HTTPMethodPut, "", put.Put)
	}
	patch, ok := instance.(IPatch)
	if ok {
		if init, ok := instance.(IPatchInit); ok {
			init.InitPatch(instance.(IResource))
		}
		ctrl.RegisterMethod(DetailMethod, HTTPMethodPatch, "", patch.Patch)
	}
	del, ok := instance.(IDelete)
	if ok {
		if init, ok := instance.(IDeleteInit); ok {
			init.InitDelete(instance.(IResource))
		}
		ctrl.RegisterMethod(DetailMethod, HTTPMethodDelete, "", del.Delete)
	}
	search, ok := instance.(ISearch)
	if ok {
		if init, ok := instance.(ISearchInit); ok {
			init.InitSearch(instance.(IResource))
		}
		ctrl.RegisterMethod(ListMethod, HTTPMethodPost, "_search", search.Search)
	}

	for path, methods := range ctrl.urlHandlers {
		// 安装装饰器，RegisterMethod阶段还没完成Init，只能在这里处理
		for method, handler := range methods {
			methods[method] = InstallDecorators(handler, decorators)
		}
		// 操作方法注册到路由
		router.Any(urlPath+path, func(c *gin.Context) {
			res := ctrl.httpProxy(methods)(c)
			if res != nil {
				res.Response(c)
			}
		})
	}
}

// RegisterMethod 操作方法和返回类型注册
func (ctrl *Controller) RegisterMethod(methodType MethodType, httpMethod HttpMethod, postfix string, handler HandlerFunc) {
	path := ""
	if methodType == DetailMethod {
		path += "/:id"
		ctrl.HaveDetail = true
	}

	if len(postfix) > 0 {
		if !strings.HasPrefix(postfix, "/") {
			postfix = "/" + postfix
		}
		path += postfix
	}

	if _, ok := ctrl.urlHandlers[path]; !ok {
		ctrl.urlHandlers[path] = make(map[HttpMethod]HandlerFunc)
	}

	methods := ctrl.urlHandlers[path]
	if _, ok := methods[httpMethod]; ok {
		panic(fmt.Sprintf("url[%s] method[%s] conflict", path, GetMethodName(httpMethod)))
	}

	// 等价于: ctrl.urlHandlers[path][method] = handler
	methods[httpMethod] = handler
}

// httpProxy HTTP Method 与操作方法的映射
func (ctrl *Controller) httpProxy(methods map[HttpMethod]HandlerFunc) HandlerFunc {
	return func(c *gin.Context) Response {
		if resource, ok := ctrl.instance.(IResource); ok {
			ContextWithResource(c, resource)
			switch c.Request.Method {
			case "POST":
				fallthrough
			case "PATCH":
				fallthrough
			case "PUT":
				var err error
				err = ContextWithRequestBody(c, c.Request)
				if err != nil {
					return response.NewError(500, err)
				}
			}
		}
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(*response.Error)
				if ok {
					err.Response(c)
				} else {
					panic(r)
				}
			}
		}()
		var res Response
		switch c.Request.Method {
		case "GET":
			if handler, ok := methods[HTTPMethodGet]; ok {
				res = handler(c)
			}
		case "POST":
			if handler, ok := methods[HTTPMethodPost]; ok {
				res = handler(c)
			}
		case "PUT":
			if handler, ok := methods[HTTPMethodPut]; ok {
				res = handler(c)
			}
		case "PATCH":
			if handler, ok := methods[HTTPMethodPatch]; ok {
				res = handler(c)
			}
		case "DELETE":
			if handler, ok := methods[HTTPMethodDelete]; ok {
				res = handler(c)
			}
		case "OPTIONS":
			allows := GetMethodsName(methods)
			allows = append(allows, "OPTIONS")
			res = &response.TextResponse{
				Status:  200,
				Text:    "",
				Headers: map[string]string{"Allow": strings.Join(allows, ",")},
			}
		}

		if res == nil {
			res = &response.TextResponse{
				Status: 405,
				Text:   "Method Not Allowed",
			}
		}
		if setter, ok := res.(logIDSetter); ok {
			logid := c.GetString("logid")
			setter.SetLogID(logid)
		}
		return res
	}
}

func (ctrl *Controller) Print(url string) {
	for path, methods := range ctrl.urlHandlers {
		for method, _ := range methods {
			fmt.Println(GetMethodName(method), url+path)
		}
	}
}
