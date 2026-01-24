package models

import (
	"mdnav/internal/models/cate"
	"mdnav/internal/models/doc"
	"sync"
)

type CateSlugDocsSlugMap struct {
	categoryDocuments map[string][]string
	mx                sync.RWMutex
}

func GetCateDocsSlugMap(catesMap *cate.CategoriesMap, docsMap *doc.DocumentsMap) *CateSlugDocsSlugMap {

	cateDocsSlugMap := &CateSlugDocsSlugMap{
		categoryDocuments: make(map[string][]string),
	}

	catesSlice := catesMap.GetCategoriesSlice()
	documentsSlice := docsMap.GetDocumentsSlice()

	for _, cate := range catesSlice {

		var docsSlug []string
		for _, doc := range documentsSlice {
			if cate.Slug == doc.CateSlug {
				docsSlug = append(docsSlug, doc.Slug)
			}
		}

		cateDocsSlugMap.categoryDocuments[cate.Slug] = docsSlug
	}

	return cateDocsSlugMap
}

func (c *CateSlugDocsSlugMap) GetCateDocsSlugMap() map[string][]string {

	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.categoryDocuments
}

func (c *CateSlugDocsSlugMap) GetCateDocsSliceBySlug(cateSlug string) []string {

	c.mx.RLock()
	defer c.mx.RUnlock()

	docsSlugSlice, ok := c.categoryDocuments[cateSlug]
	if ok {
		return docsSlugSlice
	}

	return nil
}
