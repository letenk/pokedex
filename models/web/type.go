package web

import "github.com/letenk/pokedex/models/domain"

type TypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Format for handle single response type
func FormatTypeResponse(types domain.Type) TypeResponse {
	formatter := TypeResponse{
		ID:   types.ID,
		Name: types.Name,
	}
	return formatter
}

// Format for handle multiples response type
func FormatTypesResponse(types []domain.Type) []TypeResponse {
	if len(types) == 0 {
		return []TypeResponse{}
	}

	var formatters []TypeResponse

	for _, data := range types {
		formatter := FormatTypeResponse(data)
		formatters = append(formatters, formatter)
	}

	return formatters
}
