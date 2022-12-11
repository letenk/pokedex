package web

import "github.com/letenk/pokedex/models/domain"

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Format for handle single response category
func FormatCategoryResponse(category domain.Category) CategoryResponse {
	formatter := CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
	return formatter
}

// Format for handle multiples response categorues
func FormatCategoriesResponse(category []domain.Category) []CategoryResponse {
	if len(category) == 0 {
		return []CategoryResponse{}
	}

	var formatters []CategoryResponse

	for _, data := range category {
		formatter := FormatCategoryResponse(data)
		formatters = append(formatters, formatter)
	}

	return formatters
}
