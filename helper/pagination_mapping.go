package helper

import "restoran/features/menu/model"

func QueryParamToPagination(data model.QueryParam, totalItem int64) *model.Pagination {
	return &model.Pagination{
		Page:       data.Page,
		PageSize:   data.PageSize,
		TotalItems: totalItem,
	}
}
