package tpl

import (
	"bytes"
	"html/template"
	"maps"
	"path"
	"time"

	"mdnav/internal/pkg/markdown"
	"mdnav/internal/utils"
)

var funcMaps = template.FuncMap{
	"md2html": func(md string) template.HTML {
		return markdown.ConvertMarkdownToHTML([]byte(md))
	},
	"timeFormat": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	},
}

func Render(tplDir, tplName string, data any, funcMap ...template.FuncMap) ([]byte, error) {

	tplFile := path.Join(tplDir, tplName)

	if !utils.PathExist(tplFile) {
		tplFile = path.Join(tplDir, "default.html")
	}

	if len(funcMap) > 0 {
		maps.Copy(funcMaps, funcMap[0])
	}

	tpl := template.New(tplName).Funcs(funcMaps)

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
