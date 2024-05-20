package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TestContract struct {
	contractapi.Contract
}

type Asset struct {
	ID    string `json:ID`
	Name  string `json:Name`
	Owner string `json:Owner`
}

func (contract *Asset) InitAssets(ctx contractapi.TransactionContextInterface) (string, error) {

	assets := []Asset{
		{
			ID:    "1",
			Name:  "Land",
			Owner: "Jack",
		},

		{
			ID:    "2",
			Name:  "Car",
			Owner: "Patrick",
		},

		{
			ID:    "3",
			Name:  "Yacht",
			Owner: "Patrick",
		},
	}

	for _, asset := range assets {

		jsonAsset, err := json.Marshal(asset)

		if err != nil {
			return "", fmt.Errorf("Failed To Convert The Asset To Json")
		}

		err = ctx.GetStub().PutState(asset.ID, jsonAsset)

		if err != nil {
			return "", fmt.Errorf("Failed To Add Asset To The Ledger")

		}
	}

	return "All The Assets Are Added To The Ledger", nil

}

func main() {

}
