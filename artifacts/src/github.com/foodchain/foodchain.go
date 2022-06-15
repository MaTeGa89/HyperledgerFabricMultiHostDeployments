package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	// for IoT sensor

	"encoding/json"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SmartContract provides functions for managing an Asset
type FoodchainContract struct {
	contractapi.Contract
}

// Product Definition from the requirements analysis and the IoT sensors
type Product struct {
	ProductID string

	SensorID        string `json:"IoTSensorUnitID"`       //Number ID of IoT Sensor Unit installed into a critical point of the supply-chain
	RoomTemperature string `json:"RoomTemperature"`       //Room Temperature from DHT11 sensor
	Humidity        string `json:"Humidity"`              //Humidity from DHT11 sensor
	Temperature     string `json:"Temperature"`           //Liquid Temperature of Product (e.g. oil) from DS18B20 sensor
	Acidity         string `json:"pH"`                    //pH of Product (e.g. oil) from pH sensor
	Location        string `json:"Origin and Provenance"` //from GPS sensor
	NetQuantity     string `json:"NetQuantity"`           //from pressure sensor
	Time            string `json:"ReadingValueTime"`      // is the reading time of the measurement by the sensor

	// Partecipants definition
	FarmerID               string
	ManufactureID          string // Manufacter may be the Producer and/or Packer
	DistributorID          string
	Status                 string
	FarmerProcessDate      string
	ManufactureProcessDate string
	DistributorProcessDate string
	RegulatorID            string // Regulator may be an organization internal or external to the consortium of the supply-chain
}

func (t *FoodchainContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return setupFoodchainTracer(stub)
}

func (t *FoodchainContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "createProduct" {
		return t.createProduct(stub, args)
	} else if function == "manufactureProcessing" {
		return t.manufactureProcessing(stub, args)
	} else if function == "distributorProcessing" {
		return t.distributorProcessing(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}
	return shim.Error("Invalid function name")
}

func setupFoodchainTracer(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	ProductID := args[0]
	SensorID := args[1]
	FoodchainContract := Product{
		ProductID:       ProductID,
		SensorID:        SensorID,
		RoomTemperature: " ",
		Humidity:        " ",
		Temperature:     " ",
		Acidity:         " ",
		Location:        " ",
		NetQuantity:     " ",
		Time:            " ",
		FarmerID:        "",
	}

	ProductBytes, _ := json.Marshal(FoodchainContract)
	stub.PutState(FoodchainContract.ProductID, ProductBytes)

	return shim.Success(nil)
}

func (f *FoodchainContract) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	ProductID := args[0]
	ProductBytes, _ := stub.GetState(ProductID)
	pd := Product{}
	json.Unmarshal(ProductBytes, &pd)

	if pd.Status == "Product initiated" {
		pd.ProductID = "olives_1"
		currentts := time.Now()
		pd.FarmerProcessDate = currentts.Format("2022-06-10 18:10:10")
		pd.Status = "Olives harvested"
	} else {
		fmt.Printf("Product not initiated yet")
	}

	ProductBytes, _ = json.Marshal(pd)
	stub.PutState(ProductID, ProductBytes)

	return shim.Success(nil)
}

func (f *FoodchainContract) manufactureProcessing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	ProductID := args[0]
	ProductBytes, err := stub.GetState(ProductID)
	pd := Product{}
	err = json.Unmarshal(ProductBytes, &pd)
	if err != nil {
		return shim.Error(err.Error())

	}

	if pd.Status == "Product created" {
		pd.ManufactureID = "Manufacture_1"
		currentts := time.Now()
		pd.ManufactureProcessDate = currentts.Format("2022-06-10 18:10:10")
		pd.Status = "manufacture Process"
	} else {
		pd.Status = "Error"
		fmt.Printf("Product not initiated yet")
	}

	ProductBytes0, _ := json.Marshal(pd)
	err = stub.PutState(ProductID, ProductBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

func (f *FoodchainContract) distributorProcessing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	ProductID := args[0]
	ProductBytes, err := stub.GetState(ProductID)
	pd := Product{}
	err = json.Unmarshal(ProductBytes, &pd)
	if err != nil {
		return shim.Error(err.Error())
	}

	if pd.Status == "Distribution Process" {
		pd.DistributorID = "Distributor_1"
		currentts := time.Now()
		pd.DistributorProcessDate = currentts.Format("2022-06-10 18:10:10")
		pd.Status = "Distribution started"
	} else {
		pd.Status = "Error"
		fmt.Printf("Distribution not initiated yet")
	}

	ProductBytes0, _ := json.Marshal(pd)
	err = stub.PutState(ProductID, ProductBytes0)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (f *FoodchainContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ENIITY string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expected ENIITY Name")
	}

	ENIITY = args[0]
	Avalbytes, err := stub.GetState(ENIITY)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + ENIITY + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(Avalbytes)

}

func main() {

	err := shim.Start(new(FoodchainContract))
	if err != nil {
		fmt.Printf("Error creating new Foodchain Contract: %s", err)
	}

}
