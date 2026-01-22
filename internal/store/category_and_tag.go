package store

import (
	"encoding/json"
	"errors"
	"sort"

	"mdnav/internal/conf"
)

// Category 分类
type Category struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
}

// CategorySorter 分类排序器
type CategorySorter struct {
	categories []Category
}

func (s *CategorySorter) Len() int {
	return len(s.categories)
}

func (s *CategorySorter) Swap(i, j int) {
	s.categories[i], s.categories[j] = s.categories[j], s.categories[i]
}

func (s *CategorySorter) Less(i, j int) bool {
	return s.categories[i].Sort < s.categories[j].Sort
}

// SortCategories 分类排序方法
func SortCategories(cates []Category) []Category {
	sorter := &CategorySorter{
		categories: make([]Category, len(cates)),
	}
	copy(sorter.categories, cates)
	sort.Sort(sorter)
	return sorter.categories
}

func GetCategoryByName(name string) (Category, error) {

	cates := make(map[string]Category)
	for key, info := range getConfigCategories() {
		info.Slug = key
		cates[info.Name] = info
	}

	cate, ok := cates[name]
	if !ok {
		return Category{}, errors.New("not found cate name is " + name)
	}

	return cate, nil

}

func GetCategoryBySlug(slug string) (Category, error) {

	cates := make(map[string]Category)
	for key, info := range getConfigCategories() {
		info.Slug = key
		cates[info.Slug] = info
	}

	cate, ok := cates[slug]
	if !ok {
		return Category{}, errors.New("not found cate slug is " + slug)
	}

	return cate, nil

}

// GetConfigCategories 获取分类数据
func GetCategoriesSlice() ([]Category, error) {

	var categories []Category
	for key, cate := range getConfigCategories() {
		cate.Slug = key
		categories = append(categories, cate)
	}

	categories = SortCategories(categories)

	return categories, nil
}

// GetConfigTags 获取标签数据
func GetConfigTags() map[string]string {

	tagsMap := make(map[string]string)

	for k, v := range conf.Config().GetStringMapString("tags") {
		tagsMap[v] = k
	}

	return tagsMap
}

func getConfigCategories() map[string]Category {

	cateMap := conf.Config().GetStringMap("categories")

	cateByte, err := json.Marshal(cateMap)
	if err != nil {
		return nil
	}

	cates := make(map[string]Category)
	if err := json.Unmarshal(cateByte, &cates); err != nil {
		return nil
	}

	return cates
}
