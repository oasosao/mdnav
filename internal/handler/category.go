package handler

import (
	"net/http"

	"mdnav/internal/store"
	"mdnav/internal/utils/tpl"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Category(ctx *gin.Context) {

	params := ctx.Param("slug")

	data, err := store.GetCategoryDocuments(h.Ctx, params, store.SortByUpdateTime, store.Ascending)
	if err != nil {
		h.Ctx.Logger.Error(err.Error())
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	bytes, err := tpl.Render("category.html", Index{
		Data: data,
		Site: store.GetSiteInfo(),
	})

	if err != nil {
		h.Ctx.Logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(bytes)
}
