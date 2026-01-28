package cate

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"mdnav/internal/core"
	"mdnav/internal/pkg/markdown"
)

// Category 分类
type Category struct {
	Name          string    `json:"name"`
	Keywords      string    `json:"keywords"`
	Description   string    `json:"description"`
	Slug          string    `json:"slug"`
	Icon          string    `json:"icon"`
	Markdown      string    `json:"markdown"`
	Image         string    `json:"image"`
	Sort          int       `json:"sort"`
	Custom        any       `json:"custom"`
	Published     bool      `json:"published"`
	DocumentCount int       `json:"document_count"`
	CreateTIme    time.Time `json:"create_time"`
	UpdateTIme    time.Time `json:"update_time"`
}

type CategoriesMap struct {
	categories map[string]Category
	mx         sync.RWMutex
}

func New(ctx *core.Context) (cates *CategoriesMap, err error) {

	catesMap := &CategoriesMap{
		categories: make(map[string]Category),
	}

	catesMap.categories, err = getAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	return catesMap, nil
}

// GetConfigCategories 获取分类数组数据
func (c *CategoriesMap) GetCategoriesSlice() []Category {

	c.mx.RLock()
	defer c.mx.RUnlock()

	var cates []Category
	for _, cate := range c.categories {
		cates = append(cates, cate)
	}

	return SortCategories(cates)
}

func (c *CategoriesMap) GetCategoriesMap() map[string]Category {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.categories
}

// GetCategoriesBySlug 根据slug获取分类数据
func (c *CategoriesMap) GetCategoriesBySlug(slug string) *Category {

	c.mx.RLock()
	defer c.mx.RUnlock()

	cate, ok := c.categories[slug]
	if ok {
		return &cate
	}

	return nil
}

func getAllCategories(ctx *core.Context) (map[string]Category, error) {

	categories := make(map[string]Category)

	dirPath := ctx.Conf.GetString("server.content_dir")

	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("is not dir")
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

			mdCont, err := markdown.Parser(pathName)
			if err != nil {
				ctx.Log.Error(err.Error())
				return nil // 继续处理其他文件
			}

			cateSlug := strings.TrimPrefix(path.Dir(pathName), walkDir)
			categories[cateSlug] = Category{
				Name:        mdCont.Name,
				Keywords:    mdCont.Keywords,
				Description: mdCont.Description,
				Slug:        cateSlug,
				Icon:        mdCont.Icon,
				Markdown:    mdCont.Markdown,
				Image:       mdCont.Image,
				Sort:        mdCont.Sort,
				Custom:      mdCont.Custom,
				Published:   mdCont.Published,
				CreateTIme:  mdCont.CreateTime,
				UpdateTIme:  mdCont.UpdateTime,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return categories, nil

}
