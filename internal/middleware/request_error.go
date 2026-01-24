package middleware

import (
	"mdnav/internal/conf"
	"mdnav/internal/core"
	"mdnav/internal/utils/tpl"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestError(ctx *core.Context) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Next()

		status := c.Writer.Status()
		if status >= 300 {

			eMsg := http.StatusText(status)

			httpError := struct {
				Code int
				Msg  string
			}{
				Code: status,
				Msg:  eMsg,
			}

			tplDir := conf.Get().GetString("template.dir")
			bytes, err := tpl.Render(tplDir, "error.html", httpError)
			if err != nil {
				ctx.Log.Error(err.Error())
				c.AbortWithStatus(500)
				return
			}

			c.Writer.WriteHeader(status)
			_, err = c.Writer.Write(bytes)
			if err != nil {
				ctx.Log.Error(err.Error())
			}
		}
	}

}
