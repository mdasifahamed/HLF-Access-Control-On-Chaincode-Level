package main

import (
	"encoding/json"
	"fmt"
	"log"

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

func (contract *TestContract) Init_Assets(ctx contractapi.TransactionContextInterface) (string, error) {
	/*
		fucn Init_Assets() is method of the Asset Struct
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

		json_Asset, err := json.Marshal(asset) // Converting the asset struct to the json

		if err != nil {
			return "", fmt.Errorf("Failed To Convert The Asset To Json") // returns error if the it fails to convert i to the array
		}

		// it calls the shimpapi method over getstub() method to putstate()
		// to add the asset to ledger
		// PutState() method retquires two params
		// one is the key and its type must be a string
		// it is the key by which we can query over the ledger for an asset
		// another is the whole reprentation of the asset in json Format
		err = ctx.GetStub().PutState(asset.ID, json_Asset)

		if err != nil {
			return "", fmt.Errorf("Failed To Add Asset To The Ledger")

		}
	}

	return "All The Assets Are Added To The Ledger", nil

}

func (contract *TestContract) Create_Asset(ctx contractapi.TransactionContextInterface,
	_id string, _name string, _owner string) (string, error) {

	/* 	Transfer_Ownership() Create Asset (Note: Only Organization1 can use this fucntion)
	@params _id is the id of the asset
	@params _name  is the name of the asset
	@params _owner owner of the asset
	return: on successfull it will return successfull message else the error message
	*/

	proposer_client, err := ctx.GetClientIdentity().GetMSPID() // retriving the transaction prospe/initiater
	if err != nil {
		return "", fmt.Errorf("Failed to get Client Identity : %v", err)
	}

	if proposer_client != "Org1MSP" { // we only want that the org1 only can have access to create asset
		return "Client Is Not Permitted To Submit Transaction", nil

	}
	// Creating Asset Struct
	asset := Musics{
		ID:    _id,
		Name:  _name,
		Owner: _owner,
	}
	// Dumping The Asset Struct To Json
	json_Asset, err := json.Marshal(asset)
	if err != nil {

		return "", fmt.Errorf("Failed To Convert Asset To Json : %v", err)

	}

	err = ctx.GetStub().PutState(asset.ID, json_Asset) // Storing The Asset To the Ledger
	if err != nil {

		return "", fmt.Errorf("Failed To Add Asset To The Ledger: %v", err)

	}

	return "SuccessFully Added The Asset To The Ledger", nil

}

func (contract *TestContract) Transfer_Ownership(ctx contractapi.TransactionContextInterface, _asset_id string, _new_owner string) (string, error) {
	/* 	Transfer_Ownership() Transfer OwnerShip Of A Asset To Another One (Note: Only Organization2 can use this fucntion)
	@params _asset_id is the id of the asset which ownership will be changed
	@params _new_owner  new owner of the asset
	return: on successfull it will return successfull message else the error message
	*/
	proposer_client, err := ctx.GetClientIdentity().GetMSPID() // retriving the transaction prospe/initiater
	if err != nil {
		return "", fmt.Errorf("Failed to get Client Identity : %v", err)
	}

	if proposer_client != "Org2MSP" { // we only want that the org2 only can have access to transfer asset

		return "Client Is Not Permitted To Submit Transaction", nil

	}

	var asset Musics
	json_asset, err := ctx.GetStub().GetState(_asset_id) // retriving the asset

	if err != nil {
		return "", fmt.Errorf("Asset Not Found For The Id : %v", _asset_id)
	}

	err = json.Unmarshal(json_asset, &asset) // dumping the asset to go struct type from json

	if err != nil {
		return "", fmt.Errorf("Failed To Unmarshal Asset : %v", err)
	}

	asset.Owner = _new_owner // updating the asset owner

	json_Upadted_asset, err := json.Marshal(asset) // dumping the asset struct to json

	if err != nil {
		return "", fmt.Errorf("Failed To Convert The Asset To Json")
	}

	err = ctx.GetStub().PutState(_asset_id, json_Upadted_asset) // updating the asset to ledger

	if err != nil {
		return "", fmt.Errorf("Failed To Transfer Ownership")
	}

	return "Asset OwnerShip Transfered SuccessFully", nil

}

func (contract *TestContract) Get_Asset_History(ctx contractapi.TransactionContextInterface, _asset_id string) ([]*Musics, error) {

	/*
		Get_Asset_History() returns the history of changes have been done to a Asset
		@param : _asset_id Id of the asset whose changes history will be retrived

		return: array of Musics struct if successfull else error message
	*/

	history_iterator, err := ctx.GetStub().GetHistoryForKey(_asset_id) // return historyiterator struct

	if err != nil {
		return nil, fmt.Errorf("Failed Get History Set From The State :%v", err)
	}

	defer history_iterator.Close() // clonsig iterator after iteration but befroe retuning array.

	var assets []*Musics

	for history_iterator.HasNext() { //move to next key from the iterator

		asset_from_query, err := history_iterator.Next() //retrive query resposne
		if err != nil {
			return nil, fmt.Errorf("Failed Get Asset From The History Iterator :%v", err)
		}

		var asset Musics

		err = json.Unmarshal(asset_from_query.Value, &asset) //json data is on Value field
		if err != nil {

			return nil, fmt.Errorf("Failed To Unmarshal Asset : %v", err)

		}

		assets = append(assets, &asset) // appending the array
	}
	return assets, nil

}

func main() {
	// This Will Initiate The Chaincode On Low Level Api (ShimApi) To Interact
	test_Chaincode, err := contractapi.NewChaincode(&TestContract{})
	if err != nil {
		log.Panicf("Failed To Initiate Chainode %v", err)
	}
	// It Will Start The Chaincode For Further Interatction (Like Deploying Contract To The Ethereum)
	err = test_Chaincode.Start()
	if err != nil {
		log.Panicf("Failed To Start Chainode")
	}
}
