package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//provide functions for managing an asset
type SmartContract struct {
      contractapi.Contract
    }

type Asset struct {
	ID        string `json:"ID"`        //the assetID
	AssetType string `json:"assetType"` //type of asset bus/nissan
	CarId     string `json:"carId"`     //deviceID
	Alias     string `json:"alias"`     //comment
	Hood      string `json:"hood"`      //from

}
type Iotdata struct {
	latitude	float64 `json:"latitude"`
	longitude	float64 `json:"longitude"`
	RainStatus	string `json:"RainStatus"`
}

	func (c *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", AssetType: "bus", CarId: "kca897f", Alias: "Woodini", Hood: "Kitengela"},
		{ID: "asset2", AssetType: "bus", CarId: "knu098e", Alias: "rockstar", Hood: "Embakasi"},
		{ID: "asset3", AssetType: "bus", CarId: "knu345g", Alias: "2pac", Hood: "Ronga"},
		{ID: "asset4", AssetType: "nissan", CarId: "kil345l", Alias: "sweper", Hood: "Kasarani"},
		{ID: "asset5", AssetType: "nissan", CarId: "kda908", Alias: "maddog", Hood: "Buru"},
		{ID: "asset6", AssetType: "nissan", CarId: "dji126", Alias: "einsten", Hood: "Kiambu"},
	}
	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put world state. %v", err)
		}
	}
	return nil
}

// CreateCar adds a new car to the world state with given details
	func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, assetType string, carId string, alias string, hood string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Asset %s already exits", id)
	}
	asset := Asset{
	ID:        id,
	AssetType: assetType,
	CarId:     carId,
	Alias:     alias,
	Hood:      hood,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return(err)
	}
	return ctx.GetStub().PutState(id, assetJSON)
}
	func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil

}
	func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, assetType string, carId string, alias string, hood string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:        id,
		AssetType: assetType,
		CarId:     carId,
		Alias:     alias,
		Hood:      hood,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)

}

//Deletes an asset
	func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}
	func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string)(bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil

}
	// Update changes the value with key in the world state
	func (s *SmartContract) CreateIotdata(ctx contractapi.TransactionContextInterface,iotdata Iotdata, key string, value string) error {
    	existing, err := ctx.GetStub().GetState(key)

    	if err != nil {
        	return fmt.Errorf("Unable to interact with world state")
    	}

    	if existing == nil {
        	return fmt.Errorf("Cannot update world state pair with key %s. Does not exist", key)
    	}

	iotdata = Iotdata{iotdata.latitude, iotdata.longitude, iotdata.RainStatus}
        iotdataAsBytes, err := json.Marshal(iotdata)

    	err = ctx.GetStub().PutState(key, []byte(iotdataAsBytes))
    	if err != nil {
        	return fmt.Errorf("Unable to interact with world state")
    	}

    	return nil
}

	func (s *SmartContract) getHistory(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil

}

	func main() {
	orgChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating Org-Asset basic chaincode: %v", err)
	}
	if err := orgChaincode.Start(); err != nil {
		log.Panicf("Error starting org-asset basic chaincode: %v", err)
	}
}
