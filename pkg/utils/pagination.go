package utils

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	defaultPage    = 1
	defaultPerPage = 10
	maxPerPage     = 100
)

// PaginationParams holds sanitized pagination and sorting inputs.
type PaginationParams struct {
	Page      int
	PerPage   int
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
	Search    string
}

// PaginationMeta is the standard meta block for paginated responses.
type PaginationMeta struct {
	CurrentPage        int   `json:"current_page"`
	PerPage            int   `json:"per_page"`
	LastPage           int   `json:"last_page"`
	Total              int64 `json:"total"`
	CurrentPageRecords int   `json:"current_page_record"`
}

// BuildPaginationParams sanitizes and validates pagination inputs and maps sort keys via whitelist.
func BuildPaginationParams(pageStr, perPageStr, sortBy, sortOrder, search string, sortWhitelist map[string]string) (PaginationParams, error) {
	page := parsePositiveInt(pageStr, defaultPage)
	per := parsePositiveInt(perPageStr, defaultPerPage)
	if per > maxPerPage {
		per = maxPerPage
	}

	sortByKey := strings.ToLower(strings.TrimSpace(sortBy))
	if sortByKey == "" {
		sortByKey = "created_at"
	}

	// Whitelist validation: if provided, validate sort_by against allowed fields
	var column string
	if sortWhitelist != nil {
		var ok bool
		column, ok = sortWhitelist[sortByKey]
		if !ok {
			return PaginationParams{}, errors.New("invalid sort_by")
		}
	} else {
		column = sortByKey
	}

	order := strings.ToLower(strings.TrimSpace(sortOrder))
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	search = strings.TrimSpace(search)

	offset := (page - 1) * per
	return PaginationParams{
		Page:      page,
		PerPage:   per,
		Offset:    offset,
		Limit:     per,
		SortBy:    column,
		SortOrder: order,
		Search:    search,
	}, nil
}

// BuildPaginationMeta computes the standard pagination meta.
func BuildPaginationMeta(total int64, params PaginationParams, currentCount int) PaginationMeta {
	lastPage := 0
	if params.PerPage > 0 {
		lastPage = int(math.Ceil(float64(total) / float64(params.PerPage)))
	}
	return PaginationMeta{
		CurrentPage:        params.Page,
		PerPage:            params.PerPage,
		LastPage:           lastPage,
		Total:              total,
		CurrentPageRecords: currentCount,
	}
}

func parsePositiveInt(val string, fallback int) int {
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
