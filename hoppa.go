package main

func createDoc(stub shim.ChaincodeStubInterface, args string) pb.Response {

    // Get Transient map

    transientMap, err := stub.GetTransient()



    // Do some error checks

    if err != nil {

        return shim.Error(“Error getting transient: ” + err.Error())

    }

    // Note that the chaincode has to access the same key: ‘asset’

    if _, ok := transientMap[“asset”]; !ok {

        return shim.Error(“asset must be a key in the transient map”)

    }

    if len(transientMap[“asset”]) == 0 {

        return shim.Error(“asset in transient map must be a non-empty JSON string”)

    }



    // Unmarshal the transient object to a struct

    var newObject MyStruct

    err := json.Unmarshal([]byte(transientMap[“asset”], &newObject)

    if err != nil {

        return shim.Error(err.Error())

    }

    // Unmarshal the args object to a struct

    var oldObject MyStruct

    err := json.Unmarshal([]byte(args, &oldObject)

    if err != nil {

        return shim.Error(err.Error())

    }



    // Write to Private & Channel Ledger

    // First to the private

    bytesPvt, err := json.Marshal(newObject)

    if err != nil {

        LOGGER.Error(“Error converting to bytes:”, err)

        return nil, err

    }



    err = stub.PutPrivateData(newObject.PDCCollectionName, newObject.keyField, bytesPvt)

    if err != nil {

        LOGGER.Error(“Error invoking on chaincode:”, err)

        return nil, err

    }



    // Second to the channel

    bytesChn, err := json.Marshal(oldObject)

    if err != nil {

        LOGGER.Error(“Error converting to bytes:”, err)

        return nil, err

    }



    err = stub.PutState(oldObject.keyField, bytesChn)

    if err != nil {

        LOGGER.Error(“Error invoking on chaincode:”, err)

        return nil, err

    }

     return bytesChn, nil

}

func getPvtDocByTxnId(stub shim.ChaincodeStubInterface, id string, pdcName string) (pb.Response, error) {

    LOGGER.Info(“Entering getPvtDocByTxnId “)

    var jsonResp string



    valAsbytes, err := stub.GetPrivateData(pdcName, id)



    if err != nil {

        jsonResp = “{\”Error\”:\”Failed to get state for ” + id + “\”}”

        return shim.Error(jsonResp), errors.New(jsonResp)

    } else if valAsbytes == nil {

        jsonResp = “{\”Error\”:\”Txn does not exist: ” + id + “\”}”

        return shim.Error(jsonResp), errors.New(jsonResp)

    }



    return shim.Success(valAsbytes), nil

}
