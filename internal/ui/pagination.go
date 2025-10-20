package ui

type Pagination struct {
	currentPage  int
	itemsPerPage int
	totalItems   int
}

func NewPagination(itemsPerPage int) *Pagination {
	return &Pagination{
		currentPage:  0,
		itemsPerPage: itemsPerPage,
		totalItems:   0,
	}
}

func (p *Pagination) SetTotalItems(total int) {
	p.totalItems = total
}

func (p *Pagination) GetCurrentPage() int {
	return p.currentPage
}

func (p *Pagination) GetTotalPages() int {
	if p.totalItems == 0 {
		return 0
	}
	pages := p.totalItems / p.itemsPerPage
	if p.totalItems%p.itemsPerPage > 0 {
		pages++
	}
	return pages
}

func (p *Pagination) NextPage() bool {
	if p.currentPage < p.GetTotalPages()-1 {
		p.currentPage++
		return true
	}
	return false
}

func (p *Pagination) PrevPage() bool {
	if p.currentPage > 0 {
		p.currentPage--
		return true
	}
	return false
}

func (p *Pagination) GetPageItems() (start, end int) {
	start = p.currentPage * p.itemsPerPage
	end = start + p.itemsPerPage

	if end > p.totalItems {
		end = p.totalItems
	}

	return start, end
}

func (p *Pagination) Reset() {
	p.currentPage = 0
}
