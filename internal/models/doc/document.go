package doc

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"mdnav/internal/core"
	"mdnav/internal/pkg/markdown"
)

type Document struct {
	Name        string    `json:"name"`        // 文档标题
	Keywords    string    `json:"keywords"`    // 关键词，用于SEO和搜索
	Description string    `json:"description"` // 文档摘要，简短描述文档内容
	Published   bool      `json:"published"`   // 是否发布，false表示草稿
	Sort        int       `json:"sort"`        // 排序权重，数字越大优先级越高
	Icon        string    `json:"icon"`        // 文档图标URL
	Url         string    `json:"url"`         // 文档链接URL
	Slug        string    `json:"slug"`        // 文档唯一标识，用于URL路径
	CateSlug    string    `json:"cate_slug"`   // 分类唯一标识，用于URL路径
	Tags        []string  `json:"tags"`        // 文档标签列表
	Image       string    `json:"image"`       // 文档封面图片URL
	CreateTime  time.Time `json:"create_time"` // 创建时间
	Custom      any       `json:"custom"`      // 自定义数据
	UpdateTime  time.Time `json:"update_time"` // 修改时间，自动从文件属性获取
	Markdown    string    `json:"markdown"`    // Markdown原始内容
}

type DocumentsMap struct {
	documents map[string]Document
	tags      map[string][]string
	mx        sync.RWMutex
}

func New(ctx *core.Context) (docsMap *DocumentsMap, err error) {

	documentsMap := &DocumentsMap{
		documents: make(map[string]Document),
		tags:      make(map[string][]string),
	}

	documentsMap.documents, documentsMap.tags, err = getAllDocuments(ctx)
	if err != nil {
		return nil, err
	}

	return documentsMap, nil
}

// GetDocumentsMap 获取文档数组数据
func (d *DocumentsMap) GetDocumentsMap() map[string]Document {
	d.mx.RLock()
	defer d.mx.RUnlock()
	return d.documents
}

// GetDocumentsSlice 获取文档数组数据
func (d *DocumentsMap) GetDocumentsSlice() []Document {

	d.mx.RLock()
	defer d.mx.RUnlock()

	var docs []Document
	for _, doc := range d.documents {
		docs = append(docs, doc)
	}

	return docs
}

// GetDocumentBySlug 根据slug获取文档数据
func (d *DocumentsMap) GetDocumentBySlug(slug string) *Document {

	d.mx.RLock()
	defer d.mx.RUnlock()

	doc, ok := d.documents[slug]
	if ok {
		return &doc
	}

	return nil
}

// GetTags 获取tags map 数据
func (d *DocumentsMap) GetTags() map[string][]string {

	d.mx.RLock()
	defer d.mx.RUnlock()
	return d.tags
}

// GetDocumentsSlugByTag 根据tag获取文档slug
func (d *DocumentsMap) GetDocumentsSlugByTag(tagName string) []string {

	d.mx.RLock()
	defer d.mx.RUnlock()

	docsSlug, ok := d.tags[tagName]
	if ok {
		return docsSlug
	}

	return nil
}

func getAllDocuments(ctx *core.Context) (map[string]Document, map[string][]string, error) {

	documents := make(map[string]Document)
	tags := make(map[string][]string)

	dirPath := ctx.Conf.GetString("server.content_dir")

	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, nil, err
	}

	if !info.IsDir() {
		return nil, nil, errors.New("content_dir 不是目录")
	}

	walkDir := strings.TrimRight(info.Name(), "/") + "/"

	// 遍历目录
	err = filepath.WalkDir(walkDir, func(pathName string, d fs.DirEntry, err error) error {
		if err != nil {
			ctx.Log.Error(err.Error())
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		if d.Name() == "_index.md" {
			return nil
		}

		mdCont, err := markdown.Parser(pathName)
		if err != nil {
			ctx.Log.Error(err.Error())
			return nil // 继续处理其他文件
		}

		cateSlug := strings.TrimPrefix(path.Dir(pathName), walkDir)
		slug := strings.TrimSuffix(path.Join(cateSlug, d.Name()), ".md")
		sort.Strings(mdCont.Tags)
		document := Document{
			Name:        mdCont.Name,
			Keywords:    mdCont.Keywords,
			Description: mdCont.Description,
			Published:   mdCont.Published,
			Sort:        mdCont.Sort,
			Icon:        mdCont.Icon,
			Url:         mdCont.Url,
			Slug:        slug,
			CateSlug:    cateSlug,
			Tags:        mdCont.Tags,
			Image:       mdCont.Image,
			CreateTime:  mdCont.CreateTime,
			Custom:      mdCont.Custom,
			UpdateTime:  mdCont.UpdateTime,
			Markdown:    mdCont.Markdown,
		}

		for _, v := range document.Tags {
			tags[v] = append(tags[v], document.Slug)
		}

		documents[slug] = document

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return documents, tags, nil
}
