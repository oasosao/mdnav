package store

import (
	"sort"
)

// PageResult 分页结果
type PageResult struct {
	Total      int        `json:"total"`       // 总文档数
	Page       int        `json:"page"`        // 当前页码
	PageSize   int        `json:"page_size"`   // 每页大小
	TotalPages int        `json:"total_pages"` // 总页数
	Documents  []Document `json:"list"`        // 当前页文档
}

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
	documents []Document
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
func SortDocuments(docs []Document, sortBy SortBy, order SortOrder) []Document {
	sorter := &DocumentSorter{
		documents: make([]Document, len(docs)),
		sortBy:    sortBy,
		order:     order,
	}
	copy(sorter.documents, docs)
	sort.Sort(sorter)
	return sorter.documents
}

// Paginate 分页函数
func Paginate(documents []Document, page, pageSize int) PageResult {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	total := len(documents)
	totalPages := (total + pageSize - 1) / pageSize // 向上取整

	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return PageResult{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			Documents:  []Document{},
		}
	}

	if end > total {
		end = total
	}

	return PageResult{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Documents:  documents[start:end],
	}
}
