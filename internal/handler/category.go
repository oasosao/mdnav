package handler

import (
	"mdnav/internal/models/doc"
	"mdnav/internal/service"
	"mdnav/internal/utils/tpl"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Category(ctx *gin.Context) {

	params := ctx.Param("slug")

	data := service.GetCategoryDocumentsByCateSlug(params, doc.SortByUpdateTime, doc.Ascending)
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
		Site:       service.GetSiteInfo(h.Ctx),
		Data:       data,
		Categories: service.GetAllCategories(),
		Category:   service.GetCategoryBySlug(params),
	}

	bytes, err := tpl.Render(h.TplDir, "category.html", result)
	if err != nil {
		h.Ctx.Log.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(bytes)

}
