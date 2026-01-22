package store

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"mdnav/internal/conf"
	"mdnav/internal/core"
)

// AllDocumentsList
type AllDocumentsList struct {
	Data map[string]Document
	Lock sync.RWMutex
}

type ResultMap struct {
	Data map[string][]string
	Lock sync.RWMutex
}

// 按分类归档所有文档
type CategoryDocuments struct {
	CateInfo  Category   `json:"cate"`
	Documents []Document `json:"list"`
}

// 定义单个文档
type CategoryDocument struct {
	CateInfo Category `json:"cate"`
	DocInfo  Document `json:"doc"`
}

type ResultType int

const (
	CateType ResultType = iota
	TagType
)

// 所有文章数据
var allDocumentsList AllDocumentsList

// LoadAllDocuments 加载所有数据
func LoadAllDocuments(ctx *core.Context) error {

	allDocumentsList.Lock.Lock()
	defer allDocumentsList.Lock.Unlock()

	allDocumentsList.Data = make(map[string]Document)

	dirPath := conf.Config().GetString("server.content_dir")

	info, err := os.Stat(dirPath)
	if err != nil {
		ctx.Logger.Error("目录不存在或无法访问: " + err.Error())
		return err
	}

	if !info.IsDir() {
		errMsg := "路径不是目录: " + dirPath
		ctx.Logger.Error(errMsg)
		return errors.New(errMsg)
	}

	// 遍历目录
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			ctx.Logger.Error("遍历目录: " + err.Error())
			return err
		}

		// 如果是根目录自身，跳过
		if path == dirPath {
			return nil
		}

		// 跳过非Markdown文件和目录
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		mdDocument, err := ParseFile(path)
		if err == nil {
			allDocumentsList.Data[mdDocument.Slug] = mdDocument
		}

		// allData := fmt.Sprintf("%+v", allDocumentsList.Data)

		// ctx.Logger.Info(allData)

		return nil
	})

	if err != nil {
		ctx.Logger.Error("遍历目录失败: " + err.Error())
		return err
	}

	return nil
}

// GetCategoryDocuments 根据分类获取文档 如果分类名(cateName)为空返回所有分类和文档
func GetCategoriesDocumnets(ctx *core.Context, sortBy SortBy, sortOrder SortOrder) ([]CategoryDocuments, error) {

	allDocument := getAllResultTypeDocuments(CateType)

	var categoryDocuments []CategoryDocuments

	cates, err := GetCategoriesSlice()
	if err != nil {
		return nil, err
	}

	for _, info := range cates {
		docs, ok := allDocument[info.Name]
		if ok {
			docs = SortDocuments(docs, sortBy, sortOrder)
			categoryDocuments = append(categoryDocuments, CategoryDocuments{CateInfo: info, Documents: docs})
		}
	}

	return categoryDocuments, nil
}

func GetCategoryDocuments(ctx *core.Context, cateSlug string, sortBy SortBy, sortOrder SortOrder) (CategoryDocuments, error) {

	allDocument := getAllResultTypeDocuments(CateType)

	cate, err := GetCategoryBySlug(cateSlug)
	if err != nil {
		return CategoryDocuments{}, err
	}

	docs, ok := allDocument[cate.Name]
	if ok {
		docs = SortDocuments(docs, sortBy, sortOrder)
	}

	return CategoryDocuments{CateInfo: cate, Documents: docs}, nil
}

// GetTagDocuments 根据标签获取文档
func GetTagDocuments(ctx *core.Context, tagSlug string, sortBy SortBy, sortOrder SortOrder) ([]CategoryDocuments, error) {

	var resultData []CategoryDocuments

	data, err := GetCategoriesDocumnets(ctx, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	for _, cateData := range data {

		var docs []Document
		for _, doc := range cateData.Documents {
			for slug := range doc.TagsMap {
				if tagSlug == slug {
					docs = append(docs, doc)
				}
			}
		}

		if len(docs) > 0 {
			resultData = append(resultData, CategoryDocuments{CateInfo: cateData.CateInfo, Documents: docs})
		}
	}

	return resultData, nil
}

// GetDocument 获取单个文档
func GetDocument(docSlug string) (CategoryDocument, error) {
	allDocumentsList.Lock.Lock()
	defer allDocumentsList.Lock.Unlock()

	doc, ok := allDocumentsList.Data[docSlug]
	if !ok {
		return CategoryDocument{}, errors.New("not found document")
	}

	cate, err := GetCategoryByName(doc.Category)
	if err != nil {
		return CategoryDocument{}, err
	}

	return CategoryDocument{CateInfo: cate, DocInfo: doc}, nil
}

func GetSiteInfo() map[string]any {
	return conf.Config().GetStringMap("site")
}

func loadResultMap(typeName ResultType) map[string][]string {

	resultMap := make(map[string][]string)

	allDocumentsList.Lock.Lock()
	defer allDocumentsList.Lock.Unlock()

	for index, val := range allDocumentsList.Data {

		switch typeName {
		case CateType:
			resultMap[val.Category] = append(resultMap[val.Category], index)
		case TagType:
			for _, tag := range val.Tags {
				resultMap[tag] = append(resultMap[tag], index)
			}
		}
	}

	return resultMap
}

func getAllResultTypeDocuments(resultType ResultType) (allDocment map[string][]Document) {

	allDocment = make(map[string][]Document)

	for name, docsIndex := range loadResultMap(resultType) {
		var docs []Document
		for _, index := range docsIndex {
			doc, ok := allDocumentsList.Data[index]
			if ok {
				docs = append(docs, doc)
			}
		}
		allDocment[name] = docs
	}

	return allDocment
}
