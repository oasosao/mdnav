package handler

import (
	"net/http"
	"path"

	"mdnav/internal/conf"
	"mdnav/internal/store"
	"mdnav/internal/utils/tpl"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Article(ctx *gin.Context) {

	params := ctx.Param("slug")

	ext := path.Ext(params)
	if ext != "" {
		dirPath := conf.Config().GetString("server.content_dir")
		fsPath := path.Join(dirPath, params)
		ctx.File(fsPath)
		return
	}

	data, err := store.GetDocument(params)
	if err != nil {
		h.Ctx.Logger.Error(err.Error())
		ctx.AbortWithStatus(404)
	} else {

		bytes, err := tpl.Render("article.html", HtmlResponse{
			Site: store.GetSiteInfo(),
			Data: data,
		})

		if err != nil {
			h.Ctx.Logger.Error(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(bytes)

	}

}
