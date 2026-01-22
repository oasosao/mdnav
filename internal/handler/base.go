package handler

import (
	"net/http"

	"mdnav/internal/core"
	"mdnav/internal/pkg/logger"
	"mdnav/internal/utils/tpl"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Ctx *core.Context
}

type JsonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HtmlResponse struct {
	Site any
	Data any
}

type HttpError struct {
	Code int
	Msg  string
}

func Error(c *gin.Context, logger logger.Logger, status int, errMsg ...string) {

	eMsg := ""

	if len(errMsg) > 0 {
		eMsg = errMsg[0]
	} else {
		eMsg = http.StatusText(status)
	}

	httpError := struct {
		Code int
		Msg  string
	}{
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
