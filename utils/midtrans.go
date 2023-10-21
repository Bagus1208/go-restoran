package utils

import (
	"restoran/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

func MidtransSnapClient(config config.Config) snap.Client {
	// midtrans.ServerKey = config.MT_Server_Key
	// midtrans.Environment = midtrans.Sandbox

	var snapClient snap.Client
	snapClient.New(config.MT_Server_Key, midtrans.Sandbox)

	return snapClient
}

func MidtransCoreAPIClient(config config.Config) coreapi.Client {
	var coreAPIClient coreapi.Client
	coreAPIClient.New(config.MT_Server_Key, midtrans.Sandbox)

	return coreAPIClient
}
