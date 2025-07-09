package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ElectronicsOrderContract contract for managing CRUD for ElectronicsOrder
type ElectronicsOrderContract struct {
	contractapi.Contract
}

type ElectronicsOrder struct {
	AssetType  string `json:"assetType"`
	Color      string `json:"color"`
	DealerName string `json:"dealerName"`
	Brand      string `json:"brand"`       // Previously Make
	DeviceType string `json:"deviceType"`  // Previously Model
	OrderID    string `json:"orderID"`
}

func getCollectionName() string {
	return "ElectronicsOrderCollection"
}

// OrderExists returns true when asset with given ID exists in private data collection
func (o *ElectronicsOrderContract) OrderExists(ctx contractapi.TransactionContextInterface, orderID string) (bool, error) {
	collectionName := getCollectionName()

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, orderID)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

// CreateOrder creates a new instance of ElectronicsOrder
func (o *ElectronicsOrderContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderID string) (string, error) {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "Org2MSP" {
		exists, err := o.OrderExists(ctx, orderID)
		if err != nil {
			return "", fmt.Errorf("could not read from world state. %s", err)
		} else if exists {
			return "", fmt.Errorf("the asset %s already exists", orderID)
		}

		order := new(ElectronicsOrder)

		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", err
		}

		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide the private data of brand, deviceType, color, dealerName")
		}

		brand, exists := transientData["brand"]
		if !exists {
			return "", fmt.Errorf("the brand was not specified in transient data. Please try again")
		}
		order.Brand = string(brand)

		deviceType, exists := transientData["deviceType"]
		if !exists {
			return "", fmt.Errorf("the deviceType was not specified in transient data. Please try again")
		}
		order.DeviceType = string(deviceType)

		color, exists := transientData["color"]
		if !exists {
			return "", fmt.Errorf("the color was not specified in transient data. Please try again")
		}
		order.Color = string(color)

		dealerName, exists := transientData["dealerName"]
		if !exists {
			return "", fmt.Errorf("the dealer was not specified in transient data. Please try again")
		}
		order.DealerName = string(dealerName)

		order.AssetType = "electronicsOrder"
		order.OrderID = orderID

		bytes, _ := json.Marshal(order)
		collectionName := getCollectionName()

		return fmt.Sprintf("Order with ID %v added successfully", orderID), ctx.GetStub().PutPrivateData(collectionName, orderID, bytes)
	} else {
		return fmt.Sprintf("Order cannot be created by organisation with MSPID %v", clientOrgID), nil
	}
}

// ReadOrder retrieves an instance of ElectronicsOrder from the private data collection
func (o *ElectronicsOrderContract) ReadOrder(ctx contractapi.TransactionContextInterface, orderID string) (*ElectronicsOrder, error) {
	exists, err := o.OrderExists(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", orderID)
	}

	return ReadPrivateState(ctx, orderID)
}

func ReadPrivateState(ctx contractapi.TransactionContextInterface, orderID string) (*ElectronicsOrder, error) {
	collectionName := getCollectionName()

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, orderID)
	if err != nil {
		return nil, err
	}

	order := new(ElectronicsOrder)
	err = json.Unmarshal(bytes, order)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private data to type ElectronicsOrder")
	}

	return order, nil
}

// DeleteOrder deletes an instance of ElectronicsOrder from the private data collection
func (o *ElectronicsOrderContract) DeleteOrder(ctx contractapi.TransactionContextInterface, orderID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}

	if clientOrgID == "Org1MSP" || clientOrgID == "Org2MSP" {
		exists, err := o.OrderExists(ctx, orderID)
		if err != nil {
			return fmt.Errorf("could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("the asset %s does not exist", orderID)
		}

		collectionName := getCollectionName()
		return ctx.GetStub().DelPrivateData(collectionName, orderID)
	} else {
		return fmt.Errorf("organisation with MSPID %v cannot delete the order", clientOrgID)
	}
}

// GetAllOrders retrieves all the assets with assetType 'electronicsOrder'
func (o *ElectronicsOrderContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*ElectronicsOrder, error) {
	collectionName := getCollectionName()
	queryString := `{"selector":{"assetType":"electronicsOrder"}}`

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return orderResultIteratorFunction(resultsIterator)
}

// GetOrdersByRange gives a range of order details based on a start key and an end key
func (o *ElectronicsOrderContract) GetOrdersByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*ElectronicsOrder, error) {
	collectionName := getCollectionName()
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return orderResultIteratorFunction(resultsIterator)
}

// Helper iterator function
func orderResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*ElectronicsOrder, error) {
	var orders []*ElectronicsOrder

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var order ElectronicsOrder
		err = json.Unmarshal(queryResult.Value, &order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}