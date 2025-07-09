package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ElectronicDevice structure matching Go chaincode struct
type ElectronicDevice struct {
	DeviceID           string `json:"deviceId"`
	Brand              string `json:"brand"`
	DeviceType         string `json:"deviceType"`
	Color              string `json:"color"`
	ManufacturerName   string `json:"manufacturerName"`
	DateOfManufacture  string `json:"dateOfManufacture"`
}

func main() {
	router := gin.Default()

	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	// Home
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// POST /api/device — create electronic device
	router.POST("/api/device", func(ctx *gin.Context) {
		var req ElectronicDevice
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		fmt.Printf("Received device creation request: %+v\n", req)

		// Collect arguments to match CreateDevice(deviceID, brand, deviceType, color, manufacturerName, dateOfManufacture)
		args := []string{
			req.DeviceID,
			req.Brand,
			req.DeviceType,
			req.Color,
			req.ManufacturerName,
			req.DateOfManufacture,
		}

		fmt.Printf("Arguments for CreateDevice txn: %v\n", args)

		result := submitTxnFn(
			"org1", "autochannel", "Electronics-management",
			"ElectronicDeviceContract", "invoke", make(map[string][]byte),
			"CreateDevice",
			args...,
		)

		fmt.Printf("Transaction result: %v\n", result)
		ctx.JSON(http.StatusOK, gin.H{"result": result})
	})

	// GET /api/device/:id — read device by ID
	router.GET("/api/device/:id", func(ctx *gin.Context) {
		deviceId := ctx.Param("id")

		result := submitTxnFn(
			"org1", "autochannel", "Electronics-management",
			"ElectronicDeviceContract", "query", make(map[string][]byte),
			"ReadDevice", deviceId,
		)

		fmt.Printf("Result from chaincode: %v\n", result)
		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// GET /api/devices — get all devices
	router.GET("/api/devices", func(ctx *gin.Context) {
		result := submitTxnFn(
			"org1", "autochannel", "Electronics-management",
			"ElectronicDeviceContract", "query", make(map[string][]byte),
			"GetAllDevices",
		)

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.Run("localhost:3001")
}