package helper

import "restoran/features/transaction/model"

func RequestToTransaction(data model.TransactionInput) *model.Transaction {
	return &model.Transaction{
		ID: data.OrderID,
	}
}

func TransactionToResponseInput(data *model.Transaction, token string, url string) *model.TransactionInputResponse {
	return &model.TransactionInputResponse{
		ID:          data.ID,
		OrderID:     data.OrderID,
		Status:      data.Status,
		Token:       token,
		RedirectURL: url,
	}
}

func TransactionToResponse(data *model.Transaction) *model.TransactionResponse {
	return &model.TransactionResponse{
		ID:            data.ID,
		OrderID:       data.OrderID,
		PaymentMethod: data.PaymentMethod,
		Status:        data.Status,
	}
}
