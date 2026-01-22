package base

import (
	"net/http"

	"mdnav/internal/pkg/logger"
	"mdnav/internal/utils/tpl"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HttpError struct {
	Code int
	Msg  string
}

var ErrorCodeMsg map[int]string = map[int]string{
	404: "找不到访问页面",
	500: "服务器错误",
}

func Error(c *gin.Context, logger logger.Logger, status int, errMsg ...string) {

	eMsg := ""

	if len(errMsg) > 0 {
		eMsg = errMsg[0]
	} else {
		eMsg = http.StatusText(status)
	}

	httpError := HttpError{
		Code: status,
		Msg:  eMsg,
	}

	bytes, err := tpl.Render("error.html", httpError)
	if err != nil {
		logger.Error(err.Error())
		c.Writer.WriteHeader(500)
		return
	}

	c.Writer.WriteHeader(status)

	_, err = c.Writer.Write(bytes)
	if err != nil {
		logger.Error(err.Error())
	}

}
