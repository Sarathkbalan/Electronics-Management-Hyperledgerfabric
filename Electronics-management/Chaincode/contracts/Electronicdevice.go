package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ElectronicDeviceContract contract for managing CRUD for ElectronicDevice
type ElectronicDeviceContract struct {
	contractapi.Contract
}

type ElectronicDevice struct {
	AssetType         string `json:"assetType"`
	DeviceID          string `json:"deviceId"`
	Color             string `json:"color"`
	DateOfManufacture string `json:"dateOfManufacture"`
	Brand             string `json:"brand"`
	DeviceType        string `json:"deviceType"`
	OwnedBy           string `json:"ownedBy"`
	Status            string `json:"status"`
}

type HistoryQueryResult struct {
	Record    *ElectronicDevice `json:"record"`
	TxId      string            `json:"txId"`
	Timestamp string            `json:"timestamp"`
	IsDelete  bool              `json:"isDelete"`
}

// DeviceExists returns true when asset with given ID exists in world state
func (c *ElectronicDeviceContract) DeviceExists(ctx contractapi.TransactionContextInterface, deviceID string) (bool, error) {
	data, err := ctx.GetStub().GetState(deviceID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

// CreateDevice creates a new instance of ElectronicDevice
func (c *ElectronicDeviceContract) CreateDevice(ctx contractapi.TransactionContextInterface, deviceID string, brand string, deviceType string, color string, manufacturerName string, dateOfManufacture string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "Org1MSP" {
		exists, err := c.DeviceExists(ctx, deviceID)
		if err != nil {
			return "", fmt.Errorf("%s", err)
		} else if exists {
			return "", fmt.Errorf("the device, %s already exists", deviceID)
		}

		device := ElectronicDevice{
			AssetType:         "electronicDevice",
			DeviceID:          deviceID,
			Color:             color,
			DateOfManufacture: dateOfManufacture,
			Brand:             brand,
			DeviceType:        deviceType,
			OwnedBy:           manufacturerName,
			Status:            "In Factory",
		}

		bytes, _ := json.Marshal(device)
		err = ctx.GetStub().PutState(deviceID, bytes)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Successfully added device %v", deviceID), nil
		}
	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// ReadDevice retrieves an instance of ElectronicDevice from the world state
func (c *ElectronicDeviceContract) ReadDevice(ctx contractapi.TransactionContextInterface, deviceID string) (*ElectronicDevice, error) {
	bytes, err := ctx.GetStub().GetState(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the device %s does not exist", deviceID)
	}

	var device ElectronicDevice
	err = json.Unmarshal(bytes, &device)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type ElectronicDevice")
	}

	return &device, nil
}

// DeleteDevice removes the instance from the world state
func (c *ElectronicDeviceContract) DeleteDevice(ctx contractapi.TransactionContextInterface, deviceID string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID == "Org1MSP" {
		exists, err := c.DeviceExists(ctx, deviceID)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return "", fmt.Errorf("the asset %s does not exist", deviceID)
		}

		err = ctx.GetStub().DelState(deviceID)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Device with id %v is deleted from world state.", deviceID), nil
		}
	} else {
		return "", fmt.Errorf("user under following MSP:%v cannot perform this action", clientOrgID)
	}
}

// GetAllDevices retrieves all the assets with assetType 'electronicDevice'
func (c *ElectronicDeviceContract) GetAllDevices(ctx contractapi.TransactionContextInterface) ([]*ElectronicDevice, error) {
	queryString := `{"selector":{"assetType":"electronicDevice"}}, "sort":[{ "deviceId": "desc"}]`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return deviceResultIteratorFunction(resultsIterator)
}

func deviceResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*ElectronicDevice, error) {
	var devices []*ElectronicDevice

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var device ElectronicDevice
		err = json.Unmarshal(queryResult.Value, &device)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}
	return devices, nil
}

// GetDevicesByRange gives a range of asset details based on a start key and end key
func (c *ElectronicDeviceContract) GetDevicesByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*ElectronicDevice, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return deviceResultIteratorFunction(resultsIterator)
}

// GetDeviceHistory returns the history of a device since creation
func (c *ElectronicDeviceContract) GetDeviceHistory(ctx contractapi.TransactionContextInterface, deviceID string) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(deviceID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var device ElectronicDevice
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &device)
			if err != nil {
				return nil, err
			}
		} else {
			device = ElectronicDevice{DeviceID: deviceID}
		}

		timestamp := response.Timestamp.AsTime()
		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &device,
			IsDelete:  response.IsDelete,
		}

		records = append(records, &record)
	}
	return records, nil
}