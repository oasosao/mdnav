package cate

import "sort"

// CategorySorter 分类排序器
type categorySorter struct {
	categories []Category
}

func (s *categorySorter) Len() int {
	return len(s.categories)
}

func (s *categorySorter) Swap(i, j int) {
	s.categories[i], s.categories[j] = s.categories[j], s.categories[i]
}

func (s *categorySorter) Less(i, j int) bool {
	return s.categories[i].Sort < s.categories[j].Sort
}

// SortCategories 分类排序方法
func SortCategories(cates []Category) []Category {
	sorter := &categorySorter{
		categories: cates,
	}
	sort.Sort(sorter)
	return sorter.categories
}
