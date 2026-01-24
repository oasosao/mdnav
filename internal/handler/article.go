package handler

import (
	"mdnav/internal/service"
	"mdnav/internal/utils/tpl"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Article(ctx *gin.Context) {

	params := strings.TrimPrefix(ctx.Param("slug"), "/")

	data := service.GetDocument(params)
	if data == nil {
		ctx.AbortWithStatus(404)
		return
	}

	// result := Response{
	// 	Status:  0,
	// 	Message: "success",
	// 	Result: Result{
	// 		Site: service.GetSiteInfo(h.Ctx),
	// 		Data: data,
	// 	},
	// }

	// ctx.JSON(200, result)

	result := Result{
		Site: service.GetSiteInfo(h.Ctx),
		Data: data,
	}

	bytes, err := tpl.Render(h.TplDir, "article.html", result)
	if err != nil {
		h.Ctx.Log.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(bytes)

}
