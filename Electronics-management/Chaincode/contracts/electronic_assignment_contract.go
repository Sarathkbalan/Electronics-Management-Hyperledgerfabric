package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ElectronicAssignmentContract handles assignments of devices to retailers (Org3MSP only)
type ElectronicAssignmentContract struct {
	contractapi.Contract
}

// DeviceAssignment represents the assignment record
type DeviceAssignment struct {
	AssetType     string `json:"assetType"`
	DeviceID      string `json:"deviceID"`
	RetailerName  string `json:"retailerName"`
	Quantity      string `json:"quantity"`
}

// DeviceExists checks whether a device ID exists in world state
func (e *ElectronicAssignmentContract) DeviceExists(ctx contractapi.TransactionContextInterface, deviceID string) (bool, error) {
	bytes, err := ctx.GetStub().GetState(deviceID)
	if err != nil {
		return false, err
	}
	return bytes != nil, nil
}

// AssignDeviceToRetailer assigns a device to a retailer (Org3MSP only)
func (e *ElectronicAssignmentContract) AssignDeviceToRetailer(ctx contractapi.TransactionContextInterface, deviceID string, retailerName string, quantity string) (string, error) {
	// Check MSP
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID != "Org3MSP" {
		return "", fmt.Errorf("MSP %v is not authorized to assign devices to retailers", clientOrgID)
	}

	// Check device exists
	exists, err := e.DeviceExists(ctx, deviceID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("device %v does not exist", deviceID)
	}

	// Create assignment
	assignment := DeviceAssignment{
		AssetType:    "DeviceAssignment",
		DeviceID:     deviceID,
		RetailerName: retailerName,
		Quantity:     quantity,
	}

	// Store in world state (keyed by deviceID)
	bytes, _ := json.Marshal(assignment)
	err = ctx.GetStub().PutState(deviceID, bytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Device %v assigned to retailer %v with quantity %v", deviceID, retailerName, quantity), nil
}

// ReadDeviceAssignment retrieves an assignment by device ID
func (e *ElectronicAssignmentContract) ReadDeviceAssignment(ctx contractapi.TransactionContextInterface, deviceID string) (*DeviceAssignment, error) {
	// Check device exists
	exists, err := e.DeviceExists(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("device %v does not exist", deviceID)
	}

	// Get assignment
	bytes, err := ctx.GetStub().GetState(deviceID)
	if err != nil {
		return nil, err
	}
	if bytes == nil {
		return nil, fmt.Errorf("assignment for device %v does not exist", deviceID)
	}

	// Unmarshal
	var assignment DeviceAssignment
	err = json.Unmarshal(bytes, &assignment)
	if err != nil {
		return nil, err
	}

	return &assignment, nil
}