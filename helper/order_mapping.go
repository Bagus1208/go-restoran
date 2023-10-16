package helper

import "restoran/features/order/model"

func RequestToOrder(data model.OrderInput, total int) *model.Order {
	return &model.Order{
		Orders: data.Orders,
		Total:  total,
	}
}
