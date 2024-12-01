package restful

const (
	TimeFormat string = "2006-01-02 15:04:05"
)

// HTTP METHOD 方法定义
const (
	HTTPMethodGet HttpMethod = iota + 1
	HTTPMethodPost
	HTTPMethodPut
	HTTPMethodPatch
	HTTPMethodDelete
)

// Restful请求方法类型，针对列表和详情页有单独的处理
const (
	ListMethod   MethodType = 1
	DetailMethod MethodType = 2
)

// GetMethodName 根据 HTTP Method 返回对应的操作值
func GetMethodName(method HttpMethod) string {
	switch method {
	case HTTPMethodGet:
		return "GET"
	case HTTPMethodPost:
		return "POST"
	case HTTPMethodPut:
		return "PUT"
	case HTTPMethodPatch:
		return "PATCH"
	case HTTPMethodDelete:
		return "DELETE"
	}
	return ""
}

// GetMethodsName 获取 HTTP 请求方法列表
func GetMethodsName(methods map[HttpMethod]HandlerFunc) []string {
	allows := make([]string, 0)
	for k := range methods {
		switch k {
		case HTTPMethodGet:
			allows = append(allows, GetMethodName(HTTPMethodGet))
		case HTTPMethodPost:
			allows = append(allows, GetMethodName(HTTPMethodPost))
		case HTTPMethodPut:
			allows = append(allows, GetMethodName(HTTPMethodPut))
		case HTTPMethodPatch:
			allows = append(allows, GetMethodName(HTTPMethodPatch))
		case HTTPMethodDelete:
			allows = append(allows, GetMethodName(HTTPMethodDelete))
		}
	}
	return allows
}
