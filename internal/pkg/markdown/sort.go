package markdown

import (
	"sort"
)

// SortBy 排序类型
type SortBy string

const (
	SortByUpdateTime SortBy = "update_time" // 更新时间
	SortByCreateTime SortBy = "create_time" // 创建时间
	SortBySort       SortBy = "sort"        // 自定义排序
)

// SortOrder 排序顺序
type SortOrder string

const (
	Ascending  SortOrder = "asc"  // 升序
	Descending SortOrder = "desc" // 降序
)

// DocumentSorter 文档排序器
type DocumentSorter struct {
	documents []Markdown
	sortBy    SortBy
	order     SortOrder
}

func (s *DocumentSorter) Len() int {
	return len(s.documents)
}

func (s *DocumentSorter) Swap(i, j int) {
	s.documents[i], s.documents[j] = s.documents[j], s.documents[i]
}

func (s *DocumentSorter) Less(i, j int) bool {
	switch s.sortBy {
	case SortBySort:
		if s.order == Ascending {
			return s.documents[i].Sort < s.documents[j].Sort
		}
		return s.documents[i].Sort > s.documents[j].Sort
	case SortByCreateTime:
		if s.order == Ascending {
			return s.documents[i].CreateTime.Before(s.documents[j].CreateTime)
		}
		return s.documents[i].CreateTime.After(s.documents[j].CreateTime)

	case SortByUpdateTime:
		if s.order == Ascending {
			return s.documents[i].UpdateTime.Before(s.documents[j].UpdateTime)
		}
		return s.documents[i].UpdateTime.After(s.documents[j].UpdateTime)
	default:
		return s.documents[i].UpdateTime.After(s.documents[j].UpdateTime)
	}
}

// SortDocuments 对文档进行排序
func SortDocuments(docs []Markdown, sortBy SortBy, order SortOrder) []Markdown {

	sorter := &DocumentSorter{
		sortBy:    sortBy,
		order:     order,
		documents: docs,
	}

	sort.Sort(sorter)
	return sorter.documents
}
