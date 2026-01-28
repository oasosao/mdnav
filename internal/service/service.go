package service

import (
	"mdnav/internal/core"
	"mdnav/internal/models"
	"mdnav/internal/models/cate"
	"mdnav/internal/models/doc"
	"sort"

	"go.uber.org/zap"
)

type CategoryDocuments struct {
	Category     cate.Category  `json:"category"`
	DocumentList []doc.Document `json:"document_list"`
}

type CategoryDocument struct {
	Category cate.Category `json:"category"`
	Document doc.Document  `json:"document"`
}

var categories *cate.CategoriesMap
var documents *doc.DocumentsMap
var cateDocsSlugMap *models.CateSlugDocsSlugMap

// LoadAllData 加载所有数据
func LoadAllData(ctx *core.Context) (err error) {

	categories = nil
	documents = nil
	cateDocsSlugMap = nil

	categories, err = cate.New(ctx)
	if err != nil {
		ctx.Log.Error("分类数据加载失败", zap.Error(err))
		return err
	}

	ctx.Log.Info("分类数据加载完成")

	documents, err = doc.New(ctx)
	if err != nil {
		ctx.Log.Error("文档数据加载失败", zap.Error(err))
		return err
	}
	ctx.Log.Info("文档数据加载完成")

	cateDocsSlugMap = models.GetCateDocsSlugMap(categories, documents)

	ctx.Log.Info("分类文档映射数据加载完成")

	return nil
}

// GetCategoriesDocuments 获取按分类文档归档好的数据
func GetCategoriesDocuments(sortBy doc.SortBy, order doc.SortOrder) []CategoryDocuments {

	var categoryDocuments []CategoryDocuments

	for _, category := range categories.GetCategoriesSlice() {

		docsSlug := cateDocsSlugMap.GetCateDocsSliceBySlug(category.Slug)
		var docs []doc.Document
		for _, docSlug := range docsSlug {
			d := documents.GetDocumentBySlug(docSlug)
			if d == nil {
				continue
			}
			docs = append(docs, *d)
		}

		docs = doc.SortDocuments(docs, sortBy, order)
		categoryDocuments = append(categoryDocuments, CategoryDocuments{Category: category, DocumentList: docs})
	}

	return categoryDocuments
}

// GetCategoryDocumentsByCateSlug 根据分类slug获取按分类文档归档好的数据
func GetCategoryDocumentsByCateSlug(cateSlug string, sortBy doc.SortBy, order doc.SortOrder) *CategoryDocuments {

	cateDoc := &CategoryDocuments{}

	category := categories.GetCategoriesBySlug(cateSlug)

	if category == nil {
		return nil
	}

	cateDoc.Category = *category
	var docs []doc.Document
	for _, docSlug := range cateDocsSlugMap.GetCateDocsSliceBySlug(cateSlug) {
		d := documents.GetDocumentBySlug(docSlug)
		if d == nil {
			continue
		}
		docs = append(docs, *d)
	}
	cateDoc.DocumentList = doc.SortDocuments(docs, sortBy, order)

	return cateDoc
}

// GetDocument 获取单个文档
func GetDocument(docSlug string) *CategoryDocument {

	categoryDocument := &CategoryDocument{}

	document := documents.GetDocumentBySlug(docSlug)
	if document == nil {
		return nil
	}

	categoryDocument.Document = *document
	categoryDocument.Category = *categories.GetCategoriesBySlug(document.CateSlug)

	return categoryDocument
}

// GetPageDocuments 获取分页后的文档数据
func GetPageDocuments(page, pageSize int, sortBy doc.SortBy, order doc.SortOrder) doc.PageResult {

	allDocuments := documents.GetDocumentsSlice()

	orderAllDocuments := doc.SortDocuments(allDocuments, sortBy, order)

	return doc.Paginate(orderAllDocuments, page, pageSize)
}

// GetAllCategoryMap 获取所有分类map数据 map[cateSlug]Category
func GetAllCategoryMap() map[string]cate.Category {
	return categories.GetCategoriesMap()
}

// GetTagDocuments 根据标签名获取文档数据
func GetTagDocuments(tagName string, sortBy doc.SortBy, order doc.SortOrder) []CategoryDocuments {

	var docs []doc.Document
	for _, docSlug := range documents.GetDocumentsSlugByTag(tagName) {
		d := documents.GetDocumentBySlug(docSlug)
		if d != nil {
			docs = append(docs, *d)
		}
	}

	// docs = doc.SortDocuments(docs, sortBy, order)

	cateSlugDocsMap := make(map[string][]doc.Document)

	for _, v := range docs {
		cateSlugDocsMap[v.CateSlug] = append(cateSlugDocsMap[v.CateSlug], v)
	}

	var tagDocuments []CategoryDocuments

	for cateSlug, docsVal := range cateSlugDocsMap {

		category := categories.GetCategoriesBySlug(cateSlug)
		if category != nil {
			categoryDocuments := CategoryDocuments{Category: *category, DocumentList: doc.SortDocuments(docsVal, sortBy, order)}
			tagDocuments = append(tagDocuments, categoryDocuments)
		}

	}

	return tagDocuments
}

// GetSiteInfo 获取网站信息
func GetSiteInfo(ctx *core.Context) map[string]any {
	return ctx.Conf.GetStringMap("site")
}

// GetAllCategories 获取所有分类数据
func GetAllCategories() []cate.Category {

	var cates []cate.Category

	for cateSlug, docsSlug := range cateDocsSlugMap.GetCateDocsSlugMap() {
		c := *categories.GetCategoriesBySlug(cateSlug)
		c.DocumentCount = len(docsSlug)
		cates = append(cates, c)
	}

	return cate.SortCategories(cates)
}

func GetCategoryBySlug(slug string) cate.Category {
	cates := categories.GetCategoriesBySlug(slug)
	return *cates
}

// GetAllTags 获取所有tag数据
func GetAllTags() []string {

	var tags []string
	for k := range documents.GetTags() {
		tags = append(tags, k)
	}
	sort.Strings(tags)
	return tags
}
