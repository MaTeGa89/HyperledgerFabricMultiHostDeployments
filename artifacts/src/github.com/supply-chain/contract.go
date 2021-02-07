package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type VaccineBatch struct {
	ID                      string                    `json:"id"`
	TempSensorID            string                    `json:"tempSensorId"`
	ManufacturingDate       uint64                    `json:"manufacturingDate"`
	ExpiryDate              uint64                    `json:"expiryDate"`
	ItemCount               string                    `json:"itemCount"`
	AddedAt                 uint64                    `json:"addedAt"`
	DocType                 string                    `json:"docType"`
	Owner                   string                    `json:"owner"`
	Description             string                    `json:"description"`
	TemperatureLocationData []TemperatureLocationData `json:"temperatureLocationData"`
}

type TemperatureLocationData struct {
	Temperature         string `json:"temperature"`
	TimeStamp           string `json:"timestamp"`
	TemperatureSensorID string `'json:"temperatureSensorId"`
	Longitude           string `json:"longitude"`
	Latitude            string `json:"latitude"`
}

func (s *SmartContract) CreateVaccineBatch(ctx contractapi.TransactionContextInterface, vaccineBatchData string) error {

	if len(vaccineBatchData) == 0 {
		return fmt.Errorf("Please pass the correct contract data")
	}

	var vaccineBatch VaccineBatch
	err := json.Unmarshal([]byte(vaccineBatchData), &vaccineBatch)
	if err != nil {
		return fmt.Errorf("Failed while unmarshling contract. %s", err.Error())
	}

	vaccineBatchAsBytes, err := json.Marshal(vaccineBatch)
	if err != nil {
		return fmt.Errorf("Failed while marshling contract. %s", err.Error())
	}

	return ctx.GetStub().PutState(vaccineBatch.ID, vaccineBatchAsBytes)
}

func (s *SmartContract) GetVaccineBatchById(ctx contractapi.TransactionContextInterface, vaccineBatchId string) (*VaccineBatch, error) {
	if len(vaccineBatchId) == 0 {
		return nil, fmt.Errorf("Please provide correct contract Id")
		// return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	VaccineBatchAsBytes, err := ctx.GetStub().GetState(vaccineBatchId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if VaccineBatchAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", vaccineBatchId)
	}

	vaccineBatch := new(VaccineBatch)
	_ = json.Unmarshal(VaccineBatchAsBytes, vaccineBatch)

	return vaccineBatch, nil

}

func (s *SmartContract) UpdateVaccineBatch(ctx contractapi.TransactionContextInterface, vaccineBatchID string, temperatureLocationData string) (string, error) {

	if len(vaccineBatchID) == 0 {
		return "", fmt.Errorf("Please pass the correct visitor id")
	}

	vaccineBatchAsBytes, err := ctx.GetStub().GetState(vaccineBatchID)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if vaccineBatchAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", vaccineBatchID)
	}

	vaccineBatch := new(VaccineBatch)
	_ = json.Unmarshal(vaccineBatchAsBytes, vaccineBatch)

	var temperatureLocationDataLocal TemperatureLocationData
	err = json.Unmarshal([]byte(temperatureLocationData), &temperatureLocationDataLocal)

	vaccineBatch.TemperatureLocationData = append(vaccineBatch.TemperatureLocationData, temperatureLocationDataLocal)

	vaccineBatchAsBytes, err = json.Marshal(vaccineBatch)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling invoice. %s", err.Error())
	}

	//  txId := ctx.GetStub().GetTxID()

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(vaccineBatch.ID, vaccineBatchAsBytes)

}

func (s *SmartContract) GetHistoryForAsset(ctx contractapi.TransactionContextInterface, vaccineBatchId string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(vaccineBatchId)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return string(buffer.Bytes()), nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
