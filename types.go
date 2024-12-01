package restful

// MethodType 操作方法类型定义
type MethodType int32

// HttpMethod 原生方法
type HttpMethod int32

// HandlerDecorator 装饰器/中间件方法
type HandlerDecorator func(HandlerFunc) HandlerFunc
