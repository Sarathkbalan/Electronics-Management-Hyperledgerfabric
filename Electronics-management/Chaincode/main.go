package main

import (
	"log"

	"Electronics/contracts" // ✅ Make sure this matches your module path

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	deviceContract := new(contracts.ElectronicDeviceContract)
	orderContract := new(contracts.ElectronicsOrderContract)
	assignmentContract := new(contracts.ElectronicAssignmentContract)

	// Create multi-contract chaincode
	chaincode, err := contractapi.NewChaincode(deviceContract, orderContract ,assignmentContract)
	if err != nil {
		log.Panicf("Could not create chaincode: %v", err)
	}

	// Start chaincode
	if err := chaincode.Start(); err != nil {
		log.Panicf("Failed to start chaincode: %v", err)
	}
}