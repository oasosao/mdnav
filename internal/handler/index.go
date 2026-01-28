package handler

import (
	"mdnav/internal/models/doc"
	"mdnav/internal/service"
	"mdnav/internal/utils/tpl"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Index(ctx *gin.Context) {

	data := service.GetCategoriesDocuments(doc.SortBySort, doc.Descending)
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
		// Tags:       service.GetAllTags(),
	}

	bytes, err := tpl.Render(h.TplDir, "index.html", result)
	if err != nil {
		h.Ctx.Log.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(bytes)
}
