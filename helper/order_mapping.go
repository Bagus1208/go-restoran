package helper

import "restoran/features/order/model"

func RequestToOrder(data model.OrderInput, total int) *model.Order {
	return &model.Order{
		NoTable: data.NoTable,
		Orders:  data.Orders,
		Total:   total,
	}
}
