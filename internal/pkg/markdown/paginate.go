package markdown

// PageResult 分页结果
type PageResult struct {
	Total      int        `json:"total"`       // 总文档数
	Page       int        `json:"page"`        // 当前页码
	PageSize   int        `json:"page_size"`   // 每页大小
	TotalPages int        `json:"total_pages"` // 总页数
	Documents  []Markdown `json:"list"`        // 当前页文档
}

// Paginate 分页函数
func Paginate(documents []Markdown, page, pageSize int) PageResult {
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
			Documents:  []Markdown{},
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
