package pagination

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	// DefaultPageSize specifies the default page size
	DefaultPageSize = 50
	// MaxPageSize specifies the maximum page size
	MaxPageSize = 1000
	// PageVar specifies the query parameter name for page number
	PageVar = "page"
	// PageSizeVar specifies the query parameter name for page size
	PageSizeVar = "per_page"
)

// Pagination represents a paginated list of data items
type Pagination struct {
	// Meta is the pagination meta data
	Meta Meta `json:"meta"`

	// Data is the actual items array (slice)
	Data interface{} `json:"data"`
}

// Meta holds the pagination information
type Meta struct {
	// Page is the current page (index/number)
	Page int `json:"page"`

	// PerPage is the number of items on each page
	PerPage int `json:"per_page"`

	// PageCount says how many pages are there
	PageCount int `json:"page_count"`

	// TotalCount is the total number of data items
	TotalCount int `json:"total_count"`
}

// New creates the new Pagination data
func New(page, perPage, total int) *Pagination {
	if perPage <= 0 {
		perPage = DefaultPageSize
	}
	if perPage > MaxPageSize {
		perPage = MaxPageSize
	}
	pageCount := -1
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &Pagination{
		Meta: Meta{
			Page:       page,
			PerPage:    perPage,
			TotalCount: total,
			PageCount:  pageCount,
		},
	}
}

// NewFromRequest creates a Pages object using the query parameters found in the given HTTP request.
func NewFromRequest(r *http.Request) *Pagination {
	page := parseInt(r.URL.Query().Get(PageVar), 1)
	perPage := parseInt(r.URL.Query().Get(PageSizeVar), DefaultPageSize)
	return New(page, perPage, -1)
}

// parseInt parses a string into an integer
// if parsing is failed, defaultValue will be returned
func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

// Offset returns the OFFSET value that can be used in a SQL statement
func (p *Pagination) Offset() int {

	return (p.Meta.Page - 1) * p.Meta.PerPage
}

// Limit returns the LIMIT value that can be used in a SQL statement
func (p *Pagination) Limit() int {
	return p.Meta.PerPage
}

// BuildLinkHeader returns an HTTP header containing the links about the pagination.
func (p *Pagination) BuildLinkHeader(baseURL string, defaultPerPage int) string {
	links := p.BuildLinks(baseURL, defaultPerPage)
	header := ""
	if links[0] != "" {
		header += fmt.Sprintf("<%v>; rel=\"first\", ", links[0])
		header += fmt.Sprintf("<%v>; rel=\"prev\"", links[1])
	}
	if links[2] != "" {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf("<%v>; rel=\"next\"", links[2])
		if links[3] != "" {
			header += fmt.Sprintf(", <%v>; rel=\"last\"", links[3])
		}
	}
	return header
}

// BuildLinks returns the first, prev, next, and last links corresponding to the pagination
// A link could be an empty string if it is not needed
// e.g: if the pagination is at the first page, then both first and prev links will be empty
func (p *Pagination) BuildLinks(baseURL string, defaultPerPage int) [4]string {
	var links [4]string
	pageCount := p.Meta.PageCount
	page := p.Meta.Page
	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}
	if page > 1 {
		links[0] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, 1)
		links[1] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page-1)
	}
	if pageCount >= 0 && page < pageCount {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
		links[3] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, pageCount)
	} else if pageCount < 0 {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
	}
	if perPage := p.Meta.PerPage; perPage != defaultPerPage {
		for i := 0; i < 4; i++ {
			if links[i] != "" {
				links[i] += fmt.Sprintf("&%v=%v", PageSizeVar, perPage)
			}
		}
	}

	return links
}

// SetData sets the items data and their total db count
func (p *Pagination) SetData(items interface{}, totalItems int) {
	value := reflect.ValueOf(items)

	kind := value.Kind()
	if kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}

	if kind != reflect.Array && kind != reflect.Slice {
		panic(fmt.Sprintf("could not get len for a value of type %T", items))
	}

	p.Data = items

	p.Meta.TotalCount = totalItems
	p.Meta.PageCount = int(math.Ceil(float64(totalItems) / float64(p.Meta.PerPage)))
}
