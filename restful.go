package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type restful struct {
	Validator *Validator
	resources map[string]IController
}

func New() *restful {
	return &restful{
		Validator: &Validator{
			Validator: validator.New(),
		},
		resources: make(map[string]IController),
	}
}

// RegisterResource 注册资源
func (r *restful) RegisterResource(url string, ctrl IController) {
	ctrl.Init(ctrl, r)
	r.resources[url] = ctrl
}

// Mount 挂载全部controller
func (r *restful) Mount(router *gin.RouterGroup) {
	for url, ctrl := range r.resources {
		ctrl.Mount(router, url)
	}
}

// GetValidator 获取 common.IValidator
func (r *restful) GetValidator() IValidator {
	return r.Validator
}

// Validate 获取 *validator.Validate，用于注册自定义校验函数
func (r *restful) Validate() *validator.Validate {
	return r.Validator.Validator
}

func (r *restful) Print(prefix string) {
	for url, ctrl := range r.resources {
		ctrl.Print(prefix + url)
	}
}
