package tpl

import (
	"bytes"
	"html/template"
	"path"
	"time"

	"mdnav/internal/conf"
	"mdnav/internal/store"
	"mdnav/internal/utils"
)

var funcMap = template.FuncMap{
	"md2html": func(md string) template.HTML {
		return store.ConvertMarkdownToHTML([]byte(md))
	},
	"timeFormat": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
}

func Render(tplName string, data any) ([]byte, error) {

	tplDir := conf.Config().GetString("template.dir")

	tplFile := path.Join(tplDir, tplName)

	if !utils.PathExist(tplFile) {
		tplFile = path.Join(tplDir, conf.Config().GetString("template.default"))
	}

	tpl := template.New(tplName).Funcs(funcMap)

	tpl, err := tpl.ParseFiles(tplFile)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
