package store

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"mdnav/internal/conf"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// Document Markdown Struct
type Document struct {
	Name        string            `yaml:"name" json:"name"`               // 文章标题
	Keywords    string            `yaml:"keywords" json:"keywords"`       // 关键词
	Description string            `yaml:"description" json:"description"` // 文章摘要
	Published   bool              `yaml:"published" json:"-"`             // 是否发布
	Sort        int               `yaml:"sort" json:"-"`                  // 文章排序
	Icon        string            `yaml:"icon" json:"icon"`               //icon
	Url         string            `yaml:"url" json:"url"`                 // 文章url
	Slug        string            `yaml:"slug" json:"slug"`               // 文章slug
	Category    string            `yaml:"category" json:"-"`              // 文章分类
	Tags        []string          `yaml:"tags" json:"-"`                  // 文章标签
	CreateTime  time.Time         `yaml:"create_time" json:"create_time"` // 发布日期
	UpdateTime  time.Time         `json:"update_time"`                    // 修改时间
	Markdown    string            `json:"markdown"`                       // Markdown内容
	TagsMap     map[string]string `json:"tags"`
	Image       string            `yaml:"image" json:"image"` // image
}

var htmlTagRegex = regexp.MustCompile("<[^>]*>")

// ParseFile markdown 文件解析方法
func ParseFile(filePath string) (doc Document, err error) {

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return doc, err
	}

	c, err := os.Open(filePath)
	if err != nil {
		return doc, err
	}
	defer c.Close()

	content, err := io.ReadAll(c)
	if err != nil {
		return doc, err
	}

	frontMatter, markdownContent, err := SplitFrontMatter(content)
	if err != nil {
		return doc, err
	}

	doc.Published = true
	doc.TagsMap = make(map[string]string)

	if frontMatter != nil {
		if err := yaml.Unmarshal(frontMatter, &doc); err != nil {
			return doc, err
		}
	}

	if !doc.Published {
		return Document{}, errors.New("doc published is fase")
	}

	doc.Markdown = markdownContent
	doc.UpdateTime = info.ModTime()

	doc.Slug = strings.TrimLeft(strings.TrimSuffix(filePath, ".md"), conf.Config().GetString("server.content_dir"))

	tagsMap := GetConfigTags()

	for _, flugTag := range doc.Tags {
		tag, ok := tagsMap[flugTag]
		if ok {
			doc.TagsMap[tag] = flugTag
		}
	}

	return doc, nil
}

// SplitFrontMatter
func SplitFrontMatter(content []byte) (frontCont []byte, mdContent string, err error) {

	contents := bytes.NewBuffer(content)
	if !strings.HasPrefix(contents.String(), "---") {
		return nil, contents.String(), nil
	}

	scanner := bufio.NewScanner(contents)

	var frontMatter bytes.Buffer
	var markdownContent bytes.Buffer

	state := "start"
	frontMatterLines := 0

	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case "start":
			if line == "---" {
				state = "frontmatter"
				frontMatterLines++
			}
		case "frontmatter":
			frontMatterLines++
			if line == "---" && frontMatterLines > 1 {
				state = "content"
				continue
			}
			frontMatter.WriteString(line + "\n")
		case "content":
			markdownContent.WriteString(line + "\n")
		}
	}

	return frontMatter.Bytes(), markdownContent.String(), scanner.Err()
}

// ConvertMarkdownToHTML 把markdown转换为html
func ConvertMarkdownToHTML(markdownContent []byte) template.HTML {

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(markdownContent, &buf); err != nil {
		return ""
	}

	return template.HTML(buf.String())
}
