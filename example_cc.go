package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//BlockChainTaskChainCode
type BlockChainTaskChainCode struct {
}

//contacts structure
type vehicles struct {
	ObjectType        string `json:"docType"`
    VehicleID         string  `json:"vehicleID"`
	Ownership            string  `json:"ownership"`
    Status       string  `json:"status"`
    CreatedDate       string  `json:"createdDate"`
    LastModifiedDate string `json:"lastModifiedDate"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(BlockChainTaskChainCode))
	if err != nil {
		fmt.Printf("Error starting BlockChain Task chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *BlockChainTaskChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *BlockChainTaskChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke is running " + function)

	// Handle different functions
	if function == "createVehicles" {
		return t.createVehicles(stub, args)
	} else if function == "queryCreatedVehicles" {
	  return t.queryCreatedVehicles(stub, args)
	} else if function == "querySpecificVehicle" {
	  return t.querySpecificVehicle(stub, args)
	} else if function == "transferOwnership" {
	  return t.transferOwnership(stub, args)
	} else if function == "queryCraetedVehicleByManufacturer" {
	  return t.queryCraetedVehicleByManufacturer(stub, args)
	} else if function == "getVehicleHistory" {
	  return t.getVehicleHistory(stub, args)
	} 
	
	
	    
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
//  save Vehicles Data into Blockchain
// ============================================================
func (t *BlockChainTaskChainCode) createVehicles(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- start createVehicles")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//Data Fields
	vehicleID := args[0]
	ownership := "Manufacturer"
	status := "New"
	createdDate := args[1]
	lastModifiedDate := ""
	
	// ==== Check if vehicleID already exists ====
	createVehiclessAsBytes, err := stub.GetState(vehicleID)
	if err != nil {
		return shim.Error("Error checking vehicle id " + err.Error())
	} else if createVehiclessAsBytes != nil {
		fmt.Println("vehicle ID already exists: " + vehicleID)
		return shim.Error("vehicle Id already exists: " + vehicleID)
	}

	// ==== Create Vehicles and marshal to JSON ====
	objectType := "vehicles"
	vehicle := &vehicles{objectType, vehicleID, ownership, status, createdDate, lastModifiedDate}
	createVehiclessAsBytes, err = json.Marshal(vehicle)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save Vehicles to state ===
	err = stub.PutState(vehicleID, createVehiclessAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== createVehicles saved and indexed. Return success ====
	fmt.Println("- end createVehicles")

	return shim.Success(nil)
}

// ============================================================
//  Query Created Vehicles Data from Blockchain
// ============================================================
func (t *BlockChainTaskChainCode) queryCreatedVehicles(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vehicles\",\"status\":\"New\"}}")
	
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ============================================================
//  Query Specific Vehicle Data from Blockchain
// ============================================================
func (t *BlockChainTaskChainCode) querySpecificVehicle(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	vehicleID := args[0]
	
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vehicles\",\"vehicleID\":\"%s\"}}", vehicleID)
	
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ============================================================
// Transfer Ownership Data into Blockchain
// ============================================================
func (t *BlockChainTaskChainCode) transferOwnership(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- start transferOwnership")
	
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	//Data Fields
	vehicleID := args[0]
	ownership := args[1]
	status := args[2]
	createdDate := args[3]
	lastModifiedDate := args[4]
	
	// ==== Check if vehicleID already exists ====
	updateVehiclessAsBytes, err := stub.GetState(vehicleID)
	if err != nil {
		return shim.Error("Failed to get vehicle details for the given vehicle id" + err.Error())
	} else if updateVehiclessAsBytes == nil {
		fmt.Println("vehicle ID does not exists: " + vehicleID)
		return shim.Error("vehicle Id does not exists: " + vehicleID)
	}

	// ==== update Vehilcle and marshal to JSON ====
	objectType := "vehicles"
	vehicle := &vehicles{objectType, vehicleID, ownership, status, createdDate, lastModifiedDate}
	updateVehiclessAsBytes, err = json.Marshal(vehicle)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save New Vehicle Owner to state ===
	err = stub.PutState(vehicleID, updateVehiclessAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== transferOwnership saved and indexed. Return success ====
	fmt.Println("- end transferOwnership")

	return shim.Success(nil)
}

//  Query All Vehicles Created by Maufacturer Data from Blockchain
// ============================================================
func (t *BlockChainTaskChainCode) queryCraetedVehicleByManufacturer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vehicles\",\"ownership\":\"Manufacturer\"}}")
	
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//============================================================
//Query to get the Vehicle history
//============================================================
func (t *BlockChainTaskChainCode) getVehicleHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	vehicleID := args[0]

	fmt.Printf("- start getVehicleHistory: %s\n", vehicleID)

	resultsIterator, err := stub.GetHistoryForKey(vehicleID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the trip history
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON trip history)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getVehicleHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}