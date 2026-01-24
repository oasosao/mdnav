package handler

import (
	"mdnav/internal/core"
)

// Handler HTTP请求处理器结构体，包含应用上下文
type Handler struct {
	Ctx    *core.Context // 应用上下文，包含日志记录器等核心组件
	TplDir string
}

// JsonResponse JSON响应结构体
type Response struct {
	Status  int    `json:"status"`  // 响应状态码，0表示成功，非0表示失败
	Message string `json:"message"` // 响应消息，描述请求结果
	Result  any    `json:"result"`  // 响应数据，根据请求返回对应的数据
}

type Result struct {
	Site any `json:"site"` // 站点信息，包含站点名称、关键词等配置
	Data any `json:"data"` // 页面数据，根据请求返回对应的数据
}
