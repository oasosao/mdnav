package handler

import (
	"net/http"

	"mdnav/internal/store"
	"mdnav/internal/utils/tpl"

	"github.com/gin-gonic/gin"
)

type Index struct {
	Data any
	Site any
}

func (h *Handler) Index(ctx *gin.Context) {
	data, err := store.GetCategoryDocuments(h.Ctx, "", store.SortByUpdateTime, store.Ascending)
	if err != nil {
		h.Ctx.Logger.Error(err.Error())
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	bytes, err := tpl.Render("index.html", Index{
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
