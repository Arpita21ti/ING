package utils

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// PaginationParams holds pagination information
type PaginationParams struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

// PagedResult represents a paginated result set
type PagedResult struct {
	Items      interface{}      `json:"items"`
	Pagination PaginationParams `json:"pagination"`
}

// Pagination constants
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ExtractPaginationParams extracts pagination parameters from the request query
func ExtractPaginationParams(r *http.Request) PaginationParams {
	query := r.URL.Query()

	page := parseIntParam(query, "page", DefaultPage)
	if page < 1 {
		page = DefaultPage
	}

	pageSize := parseIntParam(query, "pageSize", DefaultPageSize)
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// parseIntParam parses an integer parameter from query parameters
func parseIntParam(query url.Values, key string, defaultValue int) int {
	strValue := query.Get(key)
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}

	return value
}

// CalculatePagination calculates pagination values
func CalculatePagination(totalItems int64, params PaginationParams) PaginationParams {
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.PageSize)))

	return PaginationParams{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// ApplyPagination applies pagination to SQL queries
func ApplyPagination(query string, params PaginationParams) (string, int, int) {
	offset := (params.Page - 1) * params.PageSize
	limit := params.PageSize

	return query + " LIMIT ? OFFSET ?", limit, offset
}

// NewPagedResult creates a new paged result
func NewPagedResult(items interface{}, totalItems int64, params PaginationParams) PagedResult {
	pagination := CalculatePagination(totalItems, params)

	return PagedResult{
		Items:      items,
		Pagination: pagination,
	}
}

// BuildPaginationLinks builds pagination links for HATEOAS
func BuildPaginationLinks(baseURL string, pagination PaginationParams) map[string]string {
	links := make(map[string]string)

	// Parse base URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return links
	}

	// Get query parameters
	query := u.Query()

	// Self link
	query.Set("page", strconv.Itoa(pagination.Page))
	query.Set("pageSize", strconv.Itoa(pagination.PageSize))
	u.RawQuery = query.Encode()
	links["self"] = u.String()

	// First page link
	query.Set("page", "1")
	u.RawQuery = query.Encode()
	links["first"] = u.String()

	// Last page link
	query.Set("page", strconv.Itoa(pagination.TotalPages))
	u.RawQuery = query.Encode()
	links["last"] = u.String()

	// Previous page link (if not on first page)
	if pagination.Page > 1 {
		query.Set("page", strconv.Itoa(pagination.Page-1))
		u.RawQuery = query.Encode()
		links["prev"] = u.String()
	}

	// Next page link (if not on last page)
	if pagination.Page < pagination.TotalPages {
		query.Set("page", strconv.Itoa(pagination.Page+1))
		u.RawQuery = query.Encode()
		links["next"] = u.String()
	}

	return links
}

// HasNextPage checks if there is a next page
func HasNextPage(pagination PaginationParams) bool {
	return pagination.Page < pagination.TotalPages
}

// HasPrevPage checks if there is a previous page
func HasPrevPage(pagination PaginationParams) bool {
	return pagination.Page > 1
}

// GetSkip calculates the number of items to skip for MongoDB/NoSQL
func GetSkip(pagination PaginationParams) int64 {
	return int64((pagination.Page - 1) * pagination.PageSize)
}

// GetLimit returns the page size as limit
func GetLimit(pagination PaginationParams) int64 {
	return int64(pagination.PageSize)
}
