package markdown

import (
	"bufio"
	"bytes"
	"html/template"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// Document Markdown文档结构体，用于表示单个Markdown文档的元数据和内容
type Markdown struct {
	Name        string    `yaml:"name"`        // 文档标题
	Keywords    string    `yaml:"keywords"`    // 关键词，用于SEO和搜索
	Description string    `yaml:"description"` // 文档摘要，简短描述文档内容
	Published   bool      `yaml:"published"`   // 是否发布，false表示草稿
	Sort        int       `yaml:"sort"`        // 排序权重，数字越大优先级越高
	Icon        string    `yaml:"icon"`        // 文档图标URL
	Url         string    `yaml:"url"`         // 文档链接URL
	Tags        []string  `yaml:"tags"`        // 文档标签列表
	Image       string    `yaml:"image"`       // 文档封面图片URL
	CreateTime  time.Time `yaml:"create_time"` // 创建时间
	Custom      any       `yaml:"custom"`      // 自定义数据
	Slug        string    `yaml:"slug"`        // 文档唯一标识，用于URL路径
	Category    string    `yaml:"category"`    // 文档所属分类名
	UpdateTime  time.Time // 修改时间，自动从文件属性获取
	Markdown    string    // Markdown原始内容
}

var htmlTagRegex = regexp.MustCompile("<[^>]*>")

// ParseFile markdown 文件解析方法
func Parser(filePath string) (markdownDoc Markdown, err error) {

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return markdownDoc, err
	}

	c, err := os.Open(filePath)
	if err != nil {
		return markdownDoc, err
	}
	defer c.Close()

	content, err := io.ReadAll(c)
	if err != nil {
		return markdownDoc, err
	}

	frontMatter, markdownContent, err := splitFrontMatter(content)
	if err != nil {
		return markdownDoc, err
	}

	markdownDoc.Published = true

	if frontMatter != nil {
		if err := yaml.Unmarshal(frontMatter, &markdownDoc); err != nil {
			return markdownDoc, err
		}
	}

	// if !markdownDoc.Published {
	// 	return Markdown{}, errors.New("文档未发布")
	// }

	markdownDoc.Markdown = markdownContent
	markdownDoc.UpdateTime = info.ModTime()

	return markdownDoc, nil
}

// SplitFrontMatter
func splitFrontMatter(content []byte) (frontCont []byte, mdContent string, err error) {

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
