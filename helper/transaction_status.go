package helper

import (
	"restoran/features/transaction/model"

	"github.com/midtrans/midtrans-go/coreapi"
)

func TransactionStatus(transactionStatusResp *coreapi.TransactionStatusResponse) model.Status {
	var status model.Status

	if transactionStatusResp.TransactionStatus == "capture" {
		if transactionStatusResp.FraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			status.Status = "challenge"
			status.Order = "unpaid"
		} else if transactionStatusResp.FraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			status.Status = "success"
			status.Order = "paid"
		}
	} else if transactionStatusResp.TransactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		status.Status = "success"
		status.Order = "paid"
	} else if transactionStatusResp.TransactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
	} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		status.Status = "failed"
		status.Order = "failed"
	} else if transactionStatusResp.TransactionStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		status.Status = "pending"
		status.Order = "unpaid"
	}

	return status
}
