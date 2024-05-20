package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Chaincode Struct Which Is The Main Chaincode On Which All The Functinalties Will Be Applied
type TestContract struct {
	contractapi.Contract
}

type Musics struct {
	ID    string `json:ID`
	Name  string `json:Name`
	Owner string `json:Owner`
}

func (contract *Musics) InitAssets(ctx contractapi.TransactionContextInterface) (string, error) {
	/*
		fucn InitAssets() is method of the Asset Struct
		Usage: It Adds Dummy  Assets To The Legder
		@param: 'ctx'  is type of 'contractapi.TransactionContextInterface' Which Provide all the
		ShimAPI Methods To Interact Within The Networks

		returns: A String If The Transaction is SuccessFull or it will be if the the function fails it will return an error message


	*/

	// dummy asset
	assets := []Musics{
		{
			ID:    "1",
			Name:  "Shiloutes",
			Owner: "Avicii",
		},

		{
			ID:    "2",
			Name:  "Levels",
			Owner: "Avicii",
		},

		{
			ID:    "3",
			Name:  "Sunshine",
			Owner: "Avicii",
		},
	}
	// Iteration over the assets array
	// to retrive every asset
	for _, asset := range assets {

		jsonAsset, err := json.Marshal(asset) // Converting the asset struct to the json

		if err != nil {
			return "", fmt.Errorf("Failed To Convert The Asset To Json") // returns error if the it fails to convert i to the array
		}

		// it calls the shimpapi method over getstub() method to putstate()
		// to add the asset to ledger
		// PutState() method retquires two params
		// one is the key and its type must be a string
		// it is the key by which we can query over the ledger for an asset
		// another is the whole reprentation of the asset in json Format
		err = ctx.GetStub().PutState(asset.ID, jsonAsset)

		if err != nil {
			return "", fmt.Errorf("Failed To Add Asset To The Ledger")

		}
	}

	return "All The Assets Are Added To The Ledger", nil

}

func main() {

}
