package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("CLDChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const Initiator = "user_type1_0"
const Host = "user_type2_0"
const Home = "user_type3_0"
const CFA = "user_type4_0"

//==============================================================================================================================
//	 Status types - Claim Approval lifecycle is broken down into 5 statuses, this is part of the business logic to determine what can
//					be done to the vehicle at points in it's lifecycle
//==============================================================================================================================
const STATE_INITIATE = 0
const STATE_HOST = 1
const STATE_HOME = 2
const STATE_CFA = 3

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Claim - Defines the structure for a Claim object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Claim struct {
	ClaimID        string `json:"claimId"`
	ServiceDate    string `json:"serviceDate"`
	ProviderID     string `json:"providerId"`
	MemberID       string `json:"memberId"`
	SubscriberID   string `json:"subscriberId"`
	ProcedureCode  string `json:"procedureCode"`
	ChargedAmount  string `json:"chargedAmount"`
	ApprovedAmount string `json:"approvedAmount"`
	LocalPlanCode  string `json:"localPlanCode"`
	RemotePlanCode string `json:"remotePlanCode"`
	CostShare      string `json:"costShare"`
	AdjustmentFlag string `json:"adjustmentFlag"`
	Owner          string `json:"owner"`
}

//==============================================================================================================================
//	 Init - Initialize the process by creating one record in system validating owner and then storing the information
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A, B, C, D, E, F, G string // Entities

	var err error
	var callerName string

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}

	//get caller name
	callerName, err = t.get_caller_data(stub)
	//if callerName != Initiator { // Only the Provider can create a new claim

	//	return nil, fmt.Errorf("Permission Denied. User is not authorized to create record%s==%s", callerName, Initiator)

	//}
	//if err != nil {
	//	return nil, fmt.Errorf("Not got the user details from back end")
	//}
	// Initialize the chaincode
	A = args[0]
	B = args[1]
	C = args[2]
	D = args[3]
	E = args[4]
	F = args[5]
	G = args[6]

	_, err = t.create_claim(stub, callerName, A, B, C, D, E, F, G)
	if err != nil {
		return nil, fmt.Errorf("Not able to create claim")
	}
	return nil, nil
}

//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================
func (t *SimpleChaincode) get_caller_data(stub shim.ChaincodeStubInterface) (string, error) {

	user, err := t.get_username(stub)

	if err != nil {
		return "", err
	}

	return user, nil
}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

	username, err := stub.ReadCertAttribute("username")
	if err != nil {
		return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error())
	}
	return string(username), nil
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Vehicle - Creates the initial JSON for the vehcile and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_claim(stub shim.ChaincodeStubInterface, caller string, arg0 string, arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string) ([]byte, error) {
	var c Claim
	var err error
	claimID := "\"claimId\":\"" + arg0 + "\", "
	ServiceDate := "\"serviceDate\":\"" + arg1 + "\", "
	ProviderID := "\"providerId\":\"" + arg2 + "\", "
	MemberID := "\"memberId\":\"" + arg3 + "\", "
	SubscriberID := "\"subscriberId\":\"" + arg4 + "\", "
	ProcedureCode := "\"procedureCode\":\"" + arg5 + "\", "
	ChargedAmount := "\"chargedAmount\":\"" + arg6 + "\", "
	ApprovedAmount := "\"approvedAmount\":\"UNDEFINED\", "
	LocalPlanCode := "\"localPlanCode\":\"UNDEFINED\", "
	RemotePlanCode := "\"remotePlanCode\":\"UNDEFINED\", "
	CostShare := "\"costShare\":\"UNDEFINED\", "
	AdjustmentFlag := "\"adjustmentFlag\":\"UNDEFINED\", "
	Owner := "\"owner\":\"" + caller + "\" "

	claim_json := "{" + claimID + ServiceDate + ProviderID + MemberID + SubscriberID + ProcedureCode + ChargedAmount + ApprovedAmount + LocalPlanCode + RemotePlanCode + CostShare + AdjustmentFlag + Owner + "}" // Concatenates the variables to create the total JSON object

	err = json.Unmarshal([]byte(claim_json), &c) // Convert the JSON defined above into a Claim object for go

	if err != nil {
		return nil, errors.New("Invalid JSON object")
	}

	record, err := stub.GetState(c.ClaimID)
	// If not an error then a record exists so cant create a new claim with this claimID as it must be unique

	if record != nil {
		return nil, errors.New("Claim already exists")
	}

	_, err = t.save_changes(stub, c)

	if err != nil {
		fmt.Printf("CREATE_CLAIM: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	return nil, nil

}

//==============================================================================================================================
// save_changes - Writes to the ledger the Claim struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, c Claim) (bool, error) {

	bytes, err := json.Marshal(c)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error converting Claim record: %s", err)
		return false, errors.New("Error converting Claim record")
	}

	err = stub.PutState(c.ClaimID, bytes)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error storing Claim record: %s", err)
		return false, errors.New("Error storing Claim record")
	}

	return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil

}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))

	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}
