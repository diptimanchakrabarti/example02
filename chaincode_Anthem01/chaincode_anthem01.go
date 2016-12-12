/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A, B, C, D, E, F, G string          // Entities
	var Aval, Cval int                      // Asset holdings
	var Bval, Dval, Eval, Fval, Gval string //Asset holdings
	var err error

	if len(args) != 14 {
		return nil, errors.New("Incorrect number of arguments. Expecting 14")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for first asset holding")
	}
	B = args[2]
	Bval = args[3]

	C = args[4]
	Cval, err = strconv.Atoi(args[5])
	if err != nil {
		return nil, errors.New("Expecting integer value for third asset holding")
	}

	D = args[6]
	Dval = args[7]

	E = args[8]
	Eval = args[9]

	F = args[10]
	Fval = args[11]

	G = args[12]
	Gval = args[13]

	fmt.Printf("Aval = %d, Bval = %s, Cval = %d, Dval = %s, Eval = %s, Fval = %s, Gval = %s\n", Aval, Bval, Cval, Dval, Eval, Fval, Gval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(Bval))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(C, []byte(strconv.Itoa(Cval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(D, []byte(Dval))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(E, []byte(Eval))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(F, []byte(Fval))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(G, []byte(Gval))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	var G string // Entities
	// Asset holdings
	var Gval string //Asset holdings
	var X string    // Transaction value
	var err error

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	G = args[6]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Gvalbytes, err := stub.GetState(G)
	if err != nil {
		return nil, errors.New("Failed to get state of status")
	}
	if Gvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Gval = (string(Gvalbytes))

	// Perform the execution
	X = args[7]
	if err != nil {
		return nil, errors.New("Invalid value or no value")
	}

	Gval = X
	//Bval = Bval - X
	//Eval = Eval + X
	fmt.Printf("G = %s, X = %s, Gval = %s\n", G, X, Gval)

	// Write the state back to the ledger
	err = stub.PutState(G, []byte(Gval))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var G string // Entities
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Need claim amount and subscriberid")
	}

	G = args[4]

	// Get the state from the ledger
	Gvalbytes, err := stub.GetState(G)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + G + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Gvalbytes == nil {
		jsonResp := "{\"Error\":\"No status for " + G + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Status\":\"" + G + "\",\"Value\":\"" + string(Gvalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Gvalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
