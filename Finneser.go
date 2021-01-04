package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	pb "github.com/hyperledger/fabric/protos/peer"
)

type AssetMgr struct {
}

//Define the asset
type OrgAsset struct {
	ID        string `json:"id"`        //the assetID
	AssetType string `json:"assettype"` //type of asset bus/nissan
	Status    string `json:"status"`    //status of asset
	Location  string `json:"location"`  //device location
	DeviceID  string `json:"deviceId"`  //deviceID
	Comment   string `json:"comment"`   //comment
	From      string `json:"from"`      //from
	To        string `json:"to"`        //to
}

//Define the Init and Invoke methods
func (c *AssetMgr) Init(stub shim.ChaincodeStubInterface) pb.Response {
	args := stub.GetStringArgs()
	if len(args) != 3 {
		return shim.Error("Incorect arguments expected key and a value")
		assetId := args[0]
		AssetType := args[1]
		deviceId := args[2]
	}
	//create Asset
	assetData := OrgAsset{
		Id:        assetId,
		AssetType: assetType,
		Status:    "FUL/HALF/EMPTY",
		Location:  "gpsReadings(args[2])",
		DeviceId:  deviceId,
		Comment:   "NICKNAME:WOODINI",
		From:      "Town",
		To:        "Destination:BURU"}
	assetBytes, _ := json.Marshal(assetData)
	assetErr := stub.PutState(assetId, assetBytes)
	if assetErr != nil {
		return shim.Error(fmt.Sprintf("failed to create asset:%s", args[0]))
	}
	return shim.Success(nil)

}
func (c *AssetMgr) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "Order":
		return c.Order(stub, args)
	case "Ship":
		return c.Ship(stub, args)
	case "Distribute":
		return c.Distribute(stub, args)
	case "Query":
		return c.query(stub, args)
	case "getHistory":
		return c.getHistory(stub, args)
	}
	return shim.Error("Invalid Function Name")
}

func (c *AssetMgr) Order(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return c.UpdateAsset(stub, args, "ORDER", "SCHOOL", "OEM")
}

//Update Asset Function
func (c *AssetMgr) UpdateAsset(stub shim.ChaincodeStubInterface, args []string, currentStatus string, from string, to string) pb.Response {
	assetId := args[0]
	comment := args[1]
	location := args[2]
	assetBytes, err := stub.GetState(assetId)
	orgAsset := OrgAsset{}

	if currentStatus == "ORDER" && OrgAsset.Status != "START" {
		return shim.Error(err.Error())
	} else if currentStatus == "SHIP" && orgAsset.Status != "ORDER" {
		return shim.Error(err.Error())
	} else if currentStatus == "Distribute" && orgAsset.Status != "SHIP" {
		return shim.Error(err.Error())
	}
	orgAsset.Comment = comment
	orgAsset.Status = currentStatus

	orgAsset0, _ := json.Marshal(orgAsset)
	err = stub.PutState(assetId, orgAsset0)

	return shim.Success(orgAsset0)

}
func (c *AssetMgr) Ship(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return c.UpdateAsset(stub, args, "SHIP", "SCHOOL", "OEM")
}
func (c *AssetMgr) Distribute(stub shim.ChiancodeStubInterface, args []string) pb.Response {
	return c.UpdateAsset(stub, args, "DISTRIBUTE", "SCHOOL", "OEM")
}
func (c *AssetMgr) getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId  string   `json:"txId"`
		Value OrgAsset `json:"Value"`
	}
	var history []AuditHistory
	var orgAsset OrgAsset
	assetId := args[0]
	//GEt HISTORY
	resultsIterator, err := stub.GetHistoryForKey(assetId)

	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		historyData, err := results.Iterator.Next()
		var tx AuditHistory
		tx.TxId = historyData.TxId
		json.Unmarshal(historyData.Value, &orgAsset)

		//orgAsset over

		history = append(history, tx) //add this tx to the list
	}

}

func main() {
	err := shim.Start(new(AssetMgr))
	if err != nil {
		fmt.Printf("Error creating new Asset Matatu :%s", err)
	}
}
