package ui

// Pagination gerencia a paginação dos resultados
type Pagination struct {
	currentPage  int
	itemsPerPage int
	totalItems   int
}

// NewPagination cria uma nova instância de paginação
func NewPagination(itemsPerPage int) *Pagination {
	return &Pagination{
		currentPage:  0,
		itemsPerPage: itemsPerPage,
		totalItems:   0,
	}
}

// SetTotalItems define o total de itens
func (p *Pagination) SetTotalItems(total int) {
	p.totalItems = total
}

// GetCurrentPage retorna a página atual (0-indexed)
func (p *Pagination) GetCurrentPage() int {
	return p.currentPage
}

// GetTotalPages retorna o número total de páginas
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

// NextPage avança para a próxima página
func (p *Pagination) NextPage() bool {
	if p.currentPage < p.GetTotalPages()-1 {
		p.currentPage++
		return true
	}
	return false
}

// PrevPage volta para a página anterior
func (p *Pagination) PrevPage() bool {
	if p.currentPage > 0 {
		p.currentPage--
		return true
	}
	return false
}

// GetPageItems retorna os índices dos itens da página atual
func (p *Pagination) GetPageItems() (start, end int) {
	start = p.currentPage * p.itemsPerPage
	end = start + p.itemsPerPage
	
	if end > p.totalItems {
		end = p.totalItems
	}
	
	return start, end
}

// Reset reseta para a primeira página
func (p *Pagination) Reset() {
	p.currentPage = 0
}
