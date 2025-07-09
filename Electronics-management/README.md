# Electronics-Management

## Bring up the Test network

```bash
cd fabric-samples/test-network
```
```bash
./network.sh up createChannel -c autochannel -ca -s couchdb
```
## Addding org3
```bash
cd addOrg3
./addOrg3.sh up -c autochannel -ca -s couchdb
```
```bash
cd ..
```
## Deploy Chaincode
```bash
./network.sh deployCC -ccn Electronics-management -ccp ../../Electronics-management/Chaincode/ -ccl go -c autochannel -ccv 1.0 -ccs 1 -cccg ../../Electronics-management/Chaincode/collections.json
```

## Deploy chaincode with updation
```bash
./network.sh deployCC -ccn Electronics-management -ccp ../../Electronics-management/Chaincode/ -ccl go -c autochannel -ccv 2.0 -ccs 2 -cccg ../../Electronics-management/Chaincode/collections.json
```

## General Environment variables
```bash
export FABRIC_CFG_PATH=$PWD/../config/

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_TLS_ENABLED=true
```

## Environment variables for Org1
```bash
export CORE_PEER_LOCALMSPID=Org1MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

export CORE_PEER_ADDRESS=localhost:7051
```
## Invoke-Create Electronic device

```bash
 peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C autochannel \
  -n Electronics-management \
  --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
  -c '{"function":"CreateDevice","Args":["DEV-001","Samsung","Phone","Black","Samsung Inc","2025-01-01"]}'



```
## Query- Get Electronics by ID
```bash
peer chaincode query   -C autochannel   -n Electronics-management   -c '{"Args":["ElectronicDeviceContract:ReadDevice","DEV-002"]}'
```
## Query -Get all Electronics
```bash
peer chaincode query \
  -C autochannel \
  -n Electronics-management \
  -c '{ "Args": ["ElectronicDeviceContract:GetAllDevices"]}'
```

## Environment variables for Org2
```bash
export CORE_PEER_LOCALMSPID=Org2MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

export CORE_PEER_ADDRESS=localhost:9051
```
## CreateOrder Only Org2
```bash
echo -n "{\"brand\":\"LG\",\"deviceType\":\"Television\",\"color\":\"Black\",\"dealerName\":\"ElectroWorld\"}" | base64

```
```bash
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C autochannel \
  -n Electronics-management \
  --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
  --transient '{"brand":"U2Ftc3VuZw==","deviceType":"UGhvbmU=","color":"QmxhY2s=","dealerName":"RWxlY3Ryb01hcnQ="}' \
  -c '{"function":"ElectronicsOrderContract:CreateOrder","Args":["ORD-001"]}'

```
## Read Order by ID
```bash
peer chaincode query   -C autochannel   -n Electronics-management   -c '{"function":"ElectronicsOrderContract:ReadOrder","Args":["ORD-001"]}'
```
## Read all orders
```bash
peer chaincode query   -C autochannel   -n Electronics-management   -c '{"function":"ElectronicsOrderContract:GetAllOrders","Args":[]}'

```
## Query Command for GetOrdersByRange
```bash
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C autochannel \
  -n Electronics-management \
  --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
  --transient '{"brand":"U2Ftc3VuZw==","deviceType":"UGhvbmU=","color":"QmxhY2s=","dealerName":"RWxlY3Ryb01hcnQ="}' \
  -c '{"function":"ElectronicsOrderContract:GetOrdersByRange","Args":["ORD-001","ORD-002"]}'
```

## Environment variables for Org3:
```bash
export CORE_PEER_LOCALMSPID=Org3MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp

export CORE_PEER_ADDRESS=localhost:11051
```
## Invoke Command — Assign Electronic Device to Buyer (Register Device)
```bash
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile $ORDERER_CA \
  -C autochannel \
  -n Electronics-management \
  --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
   -c '{"function":"ElectronicAssignmentContract:AssignDeviceToRetailer","Args":["DEV-001","BestElectro","30"]}'

```
## Query Command — Read Electronic Assignment
```bash
peer chaincode query \
  -C autochannel \
  -n Electronics-management \
  --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
  -c '{"function":"ElectronicAssignmentContract:ReadDeviceAssignment","Args":["DEV-001"]}'
```

## Query for gell all medicine by ORG-3
```bash
 peer chaincode query \
  -C autochannel \
  -n Electronics-management \
  -c '{"Args":["ElectronicDeviceContract:GetAllDevices"]}' \
  --peerAddresses localhost:11051 \
  --tlsRootCertFiles $ORG3_PEER_TLSROOTCERT
```
