package handler

import (
	"net/http"

	"mdnav/internal/store"
	"mdnav/internal/utils/tpl"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Tag(ctx *gin.Context) {

	params := ctx.Param("slug")
	data, err := store.GetTagDocuments(h.Ctx, params, store.SortByUpdateTime, store.Descending)
	if err != nil {
		h.Ctx.Logger.Error(err.Error())
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	bytes, err := tpl.Render("tag.html", HtmlResponse{
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

	// ctx.JSON(200, base.Response{Status: 0, Data: data, Message: "success"})
}
