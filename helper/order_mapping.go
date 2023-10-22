package helper

import "restoran/features/order/model"

func RequestToOrder(data model.OrderInput, total int) *model.Order {
	return &model.Order{
		NoTable: data.NoTable,
		Orders:  data.Orders,
		Total:   total,
	}
}

func OrderToResponse(data *model.Order) model.OrderResponse {
	return model.OrderResponse{
		ID:      data.ID,
		NoTable: data.NoTable,
		Total:   data.Total,
		Status:  data.Status,
	}
}

func OrderDetailToResponse(data *model.OrderDetail) model.OrderDetailResponse {
	return model.OrderDetailResponse{
		ID:       data.ID,
		MenuID:   data.MenuID,
		Quantity: data.Quantity,
	}
}
