package query

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Apply is a convenience method that applies pagination, sorting, and filters in one go.
func Apply(query *gorm.DB, params *QueryParams) *gorm.DB {
	query = ApplyFilters(query, params)
	query = ApplyPagination(query, params)
	query = ApplySorting(query, params)
	return query
}

// ApplyPagination applies the pagination (limit/offset) to the query.
func ApplyPagination(query *gorm.DB, params *QueryParams) *gorm.DB {
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}
	return query
}

// ApplySorting applies sorting (ascending or descending) to the query.
func ApplySorting(query *gorm.DB, params *QueryParams) *gorm.DB {
	if params.Sort == "desc" {
		query = query.Order("created_at desc")
	} else {
		query = query.Order("created_at asc")
	}
	return query
}

// ApplyFilters applies the tags and search filter to the query.
func ApplyFilters(query *gorm.DB, params *QueryParams) *gorm.DB {
	// Apply tags filter (match any of the tags)
	if len(params.Tags) > 0 {
		query = query.Where("tags @> ?", pq.Array(params.Tags))
	}

	// Apply search filter
	if params.Search != "" {
		query = query.Where("tsv @@ plainto_tsquery(?)", params.Search)
	}

	// Apply date range filter (if provided)
	if params.Since != "" {
		query = query.Where("created_at >= ?", params.Since)
	}
	if params.Until != "" {
		query = query.Where("created_at <= ?", params.Until)
	}

	return query
}
