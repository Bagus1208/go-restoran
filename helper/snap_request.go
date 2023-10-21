package helper

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func CreateSnapRequest(snapClient snap.Client, orderID string, totalPrice int64) (*snap.Response, error) {
	var snapRequest = &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: totalPrice,
		},
		EnabledPayments: snap.AllSnapPaymentType,
	}

	snapResponse, err := snapClient.CreateTransaction(snapRequest)
	if err != nil {
		return nil, err
	}

	return snapResponse, nil
}
