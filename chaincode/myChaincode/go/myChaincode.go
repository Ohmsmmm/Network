package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {

}

type user struct {
	Name string `json:"name"`
	StdID string `json:"stdID"`
	Tel string `json:"tel"`
	Status string `json:"status"`
}

type wallet struct {
	WalletName string `json:"WalletName"`
	Money int `json:"Money"`
	Owner string `json:"Owner"`
}

// ============================================================
// initMarble - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//fmt.Println("abac Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createUser" {
		// Make payment of X units from A to B
		return t.createUser(stub, args)
	}else if function == "createWallet" {
		// Make payment of X units from A to B
		return t.createWallet(stub, args)
	}else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	justString := strings.Join(args,"")
	args = strings.Split(justString, "|")

	//   0        1       2       
	// "name", "stdID", "tel"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	
	fmt.Println("- start init marble")

	name := strings.ToLower(args[0])
	stdID := args[1]
	tel := args[2]
	

	userKey := "stdID|"+stdID

	// ==== Check if marble already exists ====
	userAsBytes, err := stub.GetState(userKey)
	if err != nil {
		return shim.Error("Failed to get stdID :" + err.Error())
	}else if userAsBytes != nil {
		return shim.Error("stdID already exist"+ userKey)
	}
	
	user := &user{
		Name : name,
		StdID : stdID,
		Tel : tel,
		Status : "true",
	}
	
	fmt.Println("- start init marble")

	// ==== Create marble object and marshal to JSON ====
	userJSONasBytes, err2 := json.Marshal(user)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	fmt.Println("- start init marble")

	err3 := stub.PutState(userKey, userJSONasBytes)//rewrite the user
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	fmt.Println("- start init marble")

	return shim.Success(nil)
}
func (t *SimpleChaincode) createWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	justString := strings.Join(args,"")
	args = strings.Split(justString, "|")
	
	// 	   0		  1			  2
	// 	WalletName 	 money		 owner
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	WalletName := strings.ToLower(args[0])
	// money := args[1]
	Owner := args[2]
	WalletKey := "Wallet|"+WalletName
	//check
	Walletbytes,err := stub.GetState(WalletKey)
	if Walletbytes != nil {
		Walletbytes := "Walletkey already exists : "+WalletKey
		return shim.Error(Walletbytes)
	}
	if err != nil {
		return shim.Error("Invalid transaction amount, excpting a integer value")
	}
	Money,err := strconv.Atoi(args[1])
	if err != nil{
		return shim.Error("money isn't Int" + err.Error())
	}
	wallet := &wallet {
		WalletName : WalletName,
		Money : Money,
		Owner : Owner,
	}
	WalletjJSONbytes,err := json.Marshal(wallet)
	if err != nil{
		return shim.Error("Marshal is Error" + err.Error())
	}
	err = stub.PutState(WalletKey,WalletjJSONbytes)
	if err != nil{
		return shim.Error("Putstate is error " + err.Error())
	}
	return shim.Success(nil)
}

	// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}