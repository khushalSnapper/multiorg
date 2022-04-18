package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"log"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// type AssetPrivateDetailsWithBalance struct {
// 		TransactionID             string `json:"TransactionID"`
// 		BalanceQty                string `json:"BalanceQty"`
// 	        Status                    string `json:"Status"`
// 	        AssetName                 string `json:"AssetName"`
// }

type AssetPrivateDetails struct {
		TransactionID             string `json:"TransactionID"`
		MongoDBJson				  string `json:"MongoDBJson"`
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	TransactionID     string `json:"TransactionID"`
	TransactionEntity string `json:"TransactionEntity"`
	Branch            string `json:"Branch"`
	Module            string `json:"Module"`
	TransactionType   string `json:"TransactionType"`
	PartnerEntity     string `json:"PartnerEntity"`
	PartnerBranch     string `json:"PartnerBranch"`
	AssetID           string `json:"AssetID"`
	OrderID           string `json:"OrderID"`
	AssetLocation     string `json:"AssetLocation"`
	AssetType         string `json:"AssetType"`
	UOM               string `json:"UOM"`
	Qty               string `json:"Qty"`
	EffectiveDate     string `json:"EffectiveDate"`
	ExpiryDate        string `json:"ExpiryDate"`
	ReferenceAsset    string `json:"ReferenceAsset"`
	ReferenceOrder    string `json:"ReferenceOrder"`
	AllValues         string `json:"AllValues"`
	Acknowledgement   string `json:"Acknowledgement"`
	Transaction       string `json:"Transaction"`
	BalanceQty        string `json:"BalanceQty"`
	Status            string `json:"Status"`
	AssetName         string `json:"AssetName"`
}
 const assetCollection = "assetCollection"

// AssetPrivateDetails describes details that are private to owners
//type AssetPrivateDetails struct {
//	ID             string `json:"assetID"`
//	AppraisedValue int    `json:"appraisedValue"`
//}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) InsertAssetRecords(ctx contractapi.TransactionContextInterface, AssetData string, count string) (string, error) {
	
	if len(AssetData) == 0 {
		return "", fmt.Errorf("Please pass the correct Asset data")
	}
	i, err := strconv.Atoi(count)
	if err != nil {
		return "",fmt.Errorf("error when converting string to int: %v", err)
	}
	if i == 0 {
	// Get new asset from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "",fmt.Errorf("error getting transient: %v", err)
	}
	
	// Asset properties are private, therefore they get passed in transient field, instead of func args
	transientAssetJSON, ok := transientMap["asset_properties"]
	if !ok {
		//log error to stdout
		return "",fmt.Errorf("asset not found in the transient map input")
	}	
	
	// clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	// if(clientMSPID != "Org3MSP"){
	// type AssetTransientPrivateDetails struct {
	// 	TransactionID             string `json:"TransactionID"`
	// 	BalanceQty                string `json:"BalanceQty"`
	//         Status                    string `json:"Status"`
	//         AssetName                 string `json:"AssetName"`
	// }
	// var assetPrivateBalanceQty AssetTransientPrivateDetails
	// err = json.Unmarshal(transientAssetJSON, &assetPrivateBalanceQty)
	// if err != nil {
    //     	return "",fmt.Errorf("failed to unmarshal JSON: %v", err)
    // 	}
	
	// assetPrivateDetailsWithBalance := AssetPrivateDetailsWithBalance{
	// 	TransactionID:   assetPrivateBalanceQty.TransactionID,
	// 	BalanceQty:      assetPrivateBalanceQty.BalanceQty,
	// 	Status:          assetPrivateBalanceQty.Status,
	// 	AssetName:	  assetPrivateBalanceQty.AssetName,
	// }
	
	
	// assetPrivateDetailsWithBalanceAsBytes, err := json.Marshal(assetPrivateDetailsWithBalance)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal into JSON: %v", err)
	// }
	// err = ctx.GetStub().PutPrivateData(assetCollection, assetPrivateBalanceQty.TransactionID, assetPrivateDetailsWithBalanceAsBytes)
	// if err != nil {
	// 	return "",fmt.Errorf("failed to put asset private details: %v", err)
	// }
	
	// }
	// AssetPrivateDetails describes details that are private to owners
	type AssetTransientPrivateInput struct {
		TransactionID             string `json:"TransactionID"`
		MongoDBJson               string `json:"MongoDBJson"`
	}
	
	// Save order details to collection visible to owning organization
	var assetQtyDetails AssetTransientPrivateInput
	err = json.Unmarshal(transientAssetJSON, &assetQtyDetails)
	if err != nil {
        	return "",fmt.Errorf("failed to unmarshal JSON: %v", err)
    	}

    // Check if asset already exists
	//assetAsBytes, err := ctx.GetStub().GetPrivateData(assetQtyCollection, assetQtyDetails.TransactionID)
	//if err != nil {
	//	return fmt.Errorf("failed to get asset quantity details: %v", err)
	//} else if assetAsBytes != nil {
	//	fmt.Println("Asset quantity details already exists: " + assetQtyDetails.TransactionID)
	//	return fmt.Errorf("this asset quantity details already exists: " + assetQtyDetails.TransactionID)
	//}
	assetPrivateDetails := AssetPrivateDetails{
		TransactionID: assetQtyDetails.TransactionID,
		MongoDBJson:   assetQtyDetails.MongoDBJson,
	}
	
	assetPrivateDetailsAsBytes, err := json.Marshal(assetPrivateDetails)
	if err != nil {
		return "",fmt.Errorf("failed to marshal into JSON: %v", err)
	}
	
	// Get collection name for this organization.
	// orgCollection, err := getCollectionName(ctx)
	// if err != nil {
	// 	return "",fmt.Errorf("failed to infer private collection name for the org: %v", err)
	// }
	// Put asset appraised value into owners org specific private data collection
        err = ctx.GetStub().PutPrivateData(assetCollection, assetQtyDetails.TransactionID, assetPrivateDetailsAsBytes)
	if err != nil {
		return "",fmt.Errorf("failed to put asset private details: %v", err)
	}
	}
	var asset Asset
	err = json.Unmarshal([]byte(AssetData), &asset)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling Asset records %s", err.Error())
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling kyc records %s", err.Error())
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(asset.TransactionID, assetAsBytes)

}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) InsertOrderRecords(ctx contractapi.TransactionContextInterface, AssetData string, count string) (string, error) {
       	
	if len(AssetData) == 0 {
		return "", fmt.Errorf("Please pass the correct Asset data")
	}
	i , err := strconv.Atoi(count)
	if err != nil {
		return "",fmt.Errorf("error when converting string to int : %v", err)
	}
	if i == 0 {
	// Get new asset from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "",fmt.Errorf("error getting transient: %v", err)
	}
	
	// Asset properties are private, therefore they get passed in transient field, instead of func args
	transientAssetJSON, ok := transientMap["asset_properties"]
	if !ok {
		//log error to stdout
		return "",fmt.Errorf("asset not found in the transient map input")
	}	
	
	// clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	// if(clientMSPID != "Org3MSP"){
	// type AssetTransientPrivateDetails struct {
	// 	TransactionID             string `json:"TransactionID"`
	// 	BalanceQty                string `json:"BalanceQty"`
	//         Status                    string `json:"Status"`
	//         AssetName                 string `json:"AssetName"`
	// }
	// var assetPrivateBalanceQty AssetTransientPrivateDetails
	// err = json.Unmarshal(transientAssetJSON, &assetPrivateBalanceQty)
	// if err != nil {
    //     	return "",fmt.Errorf("failed to unmarshal JSON: %v", err)
    // 	}
	
	// assetPrivateDetailsWithBalance := AssetPrivateDetailsWithBalance{
	// 	TransactionID:   assetPrivateBalanceQty.TransactionID,
	// 	BalanceQty:      assetPrivateBalanceQty.BalanceQty,
	// 	Status:          assetPrivateBalanceQty.Status,
	// 	AssetName:	  assetPrivateBalanceQty.AssetName,
	// }
	
	
	// assetPrivateDetailsWithBalanceAsBytes, err := json.Marshal(assetPrivateDetailsWithBalance)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal into JSON: %v", err)
	// }
	// err = ctx.GetStub().PutPrivateData(assetCollection, assetPrivateBalanceQty.TransactionID, assetPrivateDetailsWithBalanceAsBytes)
	// if err != nil {
	// 	return "",fmt.Errorf("failed to put asset private details: %v", err)
	// }
	
	// }
	// AssetPrivateDetails describes details that are private to owners
	type AssetTransientPrivateInput struct {
		TransactionID             string `json:"TransactionID"`
		MongoDBJson               string `json:"MongoDBJson"`
	}
	
	// Save order details to collection visible to owning organization
	var assetQtyDetails AssetTransientPrivateInput
	err = json.Unmarshal(transientAssetJSON, &assetQtyDetails)
	if err != nil {
        	return "",fmt.Errorf("failed to unmarshal JSON: %v", err)
    	}
    	
    	// Check if asset already exists
	//assetAsBytes, err := ctx.GetStub().GetPrivateData(assetQtyCollection, assetQtyDetails.TransactionID)
	//if err != nil {
	//	return fmt.Errorf("failed to get asset quantity details: %v", err)
	//} else if assetAsBytes != nil {
	//	fmt.Println("Asset quantity details already exists: " + assetQtyDetails.TransactionID)
	//	return fmt.Errorf("this asset quantity details already exists: " + assetQtyDetails.TransactionID)
	//}
	assetPrivateDetails := AssetPrivateDetails{
		TransactionID: assetQtyDetails.TransactionID,
		MongoDBJson:   assetQtyDetails.MongoDBJson,
	}
	
	assetPrivateDetailsAsBytes, err := json.Marshal(assetPrivateDetails)
	if err != nil {
		return "",fmt.Errorf("failed to marshal into JSON: %v", err)
	}
	
	// Get collection name for this organization.
	// orgCollection, err := getCollectionName(ctx)
	// if err != nil {
	// 	return "",fmt.Errorf("failed to infer private collection name for the org: %v", err)
	// }
	// Put asset appraised value into owners org specific private data collection
        err = ctx.GetStub().PutPrivateData(assetCollection, assetQtyDetails.TransactionID, assetPrivateDetailsAsBytes)
	if err != nil {
		return "",fmt.Errorf("failed to put asset private details: %v", err)
	}
	}
	var asset Asset
	err = json.Unmarshal([]byte(AssetData), &asset)
	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling Asset records %s", err.Error())
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return "", fmt.Errorf("Failed while marshling kyc records %s", err.Error())
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(asset.TransactionID, assetAsBytes)
}


// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) GetAssetByTransactionID(ctx contractapi.TransactionContextInterface, TransactionID string) (*Asset, error) {
	if len(TransactionID) == 0 {
		return nil, fmt.Errorf("Please provide correct citizen Id")
	}

	assetAsBytes, err := ctx.GetStub().GetState(TransactionID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if assetAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", TransactionID)
	}

	asset := new(Asset)
	_ = json.Unmarshal(assetAsBytes, asset)

	return asset, nil

}

// ReadAssetPrivateDetails reads the asset private details in organization specific collection
func (s *SmartContract) GetAssetPrivateDetailsByTransactionID(ctx contractapi.TransactionContextInterface, collection string, TransactionID string) (*AssetPrivateDetails, error) {
	assetDetailsJSON, err := ctx.GetStub().GetPrivateData(collection, TransactionID) // Get the asset from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read asset details: %v", err)
	}
	if assetDetailsJSON == nil {
		return nil, nil
	}
	log.Printf("CreateAsset Put: collection %v, ID %v, owner %v", assetDetailsJSON, collection, TransactionID)
		var assetDetails *AssetPrivateDetails
		err = json.Unmarshal(assetDetailsJSON, &assetDetails)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
		}
	return assetDetails, nil

}

// ReadAssetPrivateDetails reads the asset private details in organization specific collection
// func (s *SmartContract) GetAssetBalancePrivateDetailsByTransactionID(ctx contractapi.TransactionContextInterface, collection string, TransactionID string) (*AssetPrivateDetailsWithBalance, error) {
// 	assetDetailsJSON, err := ctx.GetStub().GetPrivateData(collection, TransactionID) // Get the asset from chaincode state
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read asset details: %v", err)
// 	}
// 	if assetDetailsJSON == nil {
// 		return nil, nil
// 	}
// 	log.Printf("CreateAsset Put: collection %v, ID %v, owner %v", assetDetailsJSON, collection, TransactionID)

// 		var assetDetailsWithBalance *AssetPrivateDetailsWithBalance
// 		err = json.Unmarshal(assetDetailsJSON, &assetDetailsWithBalance)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
// 		}
// 	return assetDetailsWithBalance, nil

// }

func (s *SmartContract) GetAssetForQuery(ctx contractapi.TransactionContextInterface, queryString string) ([]Asset, error) {
	queryResults, err := s.getQueryResultForQueryString(ctx, queryString)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	return queryResults, nil

}
// getCollectionName is an internal helper function to get collection of submitting client identity.
func getCollectionName(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get the MSP ID of submitting client identity
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get verified MSPID: %v", err)
	}

	// Create the collection name
	orgCollection := clientMSPID + "PrivateCollection"

	return orgCollection, nil
}
func (s *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]Asset, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Asset{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		newAsset := new(Asset)

		err = json.Unmarshal(response.Value, newAsset)
		if err != nil {
			return nil, err
		}

		results = append(results, *newAsset)
	}
	return results, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, transactionID string, transactionentity string, branch string, module string, transactiontype string, partnerentity string, partnerbranch string, assetID string, orderID string, assetname string, assettype string, assetlocation string, balanceqty string, qty string, uom string, effectivedate string, expirydate string, referenceasset string, referenceorder string, status string, allvalues string, acknowledgement string, transaction string, count string) (string, error) {

	if len(transactionID) == 0 {

		return "", fmt.Errorf("Please pass the correct Asset data")
	}
	i , err := strconv.Atoi(count)
	if err != nil {
		return "",fmt.Errorf("error when converting string to int : %v", err)
	}
	if i == 0 {
		// Get new asset from transient map
		transientMap, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "",fmt.Errorf("error getting transient: %v", err)
		}
		
		// Asset properties are private, therefore they get passed in transient field, instead of func args
		transientAssetJSON, ok := transientMap["asset_properties"]
		if !ok {
			//log error to stdout
			return "",fmt.Errorf("asset not found in the transient map input")
		}	
		
		type AssetTransientPrivateInput struct {
			TransactionID             string `json:"TransactionID"`
			MongoDBJson               string `json:"MongoDBJson"`
		}
		
		// Save order details to collection visible to owning organization
		var assetQtyDetails AssetTransientPrivateInput
		err = json.Unmarshal(transientAssetJSON, &assetQtyDetails)
		if err != nil {
    	    	return "",fmt.Errorf("failed to unmarshal JSON: %v", err)
    		}
		
		assetPrivateDetails := AssetPrivateDetails{
			TransactionID: assetQtyDetails.TransactionID,
			MongoDBJson:   assetQtyDetails.MongoDBJson,
		}
		
		assetPrivateDetailsAsBytes, err := json.Marshal(assetPrivateDetails)
		if err != nil {
			return "",fmt.Errorf("failed to marshal into JSON: %v", err)
		}
		
		// Put asset appraised value into owners org specific private data collection
    	    err = ctx.GetStub().PutPrivateData(assetCollection, assetQtyDetails.TransactionID, assetPrivateDetailsAsBytes)
		if err != nil {
			return "",fmt.Errorf("failed to put asset private details: %v", err)
		}
	}

	assetAsBytes, err := ctx.GetStub().GetState(transactionID)

	if err != nil {
		return "", fmt.Errorf("Failed while unmarshling Asset records %s", err.Error())
	}

	if assetAsBytes == nil {

		return "", fmt.Errorf("the asset %s does not exist", transactionID)
	}

	// overwriting original asset with new asset
	asset := new(Asset)
	_ = json.Unmarshal(assetAsBytes, asset)

	asset.TransactionEntity = transactionentity
	asset.Branch = branch
	asset.Module = module
	asset.TransactionType = transactiontype
	asset.PartnerEntity = partnerentity
	asset.PartnerBranch = partnerbranch
	asset.AssetID = assetID
	asset.OrderID = orderID
	//asset.AssetName = assetname
	asset.AssetLocation = assetlocation
	asset.AssetType = assettype
	//asset.BalanceQty = balanceqty
	//asset.Qty = qty
	asset.UOM = uom
	asset.EffectiveDate = effectivedate
	asset.ExpiryDate = expirydate
	asset.ReferenceAsset = referenceasset
	asset.ReferenceOrder = referenceorder
	//asset.Status = status
	asset.AllValues = allvalues
	asset.Acknowledgement = acknowledgement
	asset.Transaction = transaction

	assetAsBytes, err = json.Marshal(asset)
	if err != nil {

		return "", fmt.Errorf("Failed marshal %s", err.Error())
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(asset.TransactionID, assetAsBytes)

}

func (s *SmartContract) DeleteAssetByTransactionId(ctx contractapi.TransactionContextInterface, TransactionID string, count string) (string, error) {
	if len(TransactionID) == 0 {
		return "", fmt.Errorf("Please provide correct contract Id")
	}

	i , err := strconv.Atoi(count)
	if err != nil {
		return "",fmt.Errorf("error when converting string to int : %v", err)
	}
	if i == 0 {
		// Get new asset from transient map
		transientMap, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "",fmt.Errorf("error getting transient: %v", err)
		}
		
		// Asset properties are private, therefore they get passed in transient field, instead of func args
		transientAssetJSON, ok := transientMap["asset_properties"]
		if !ok {
			//log error to stdout
			return "",fmt.Errorf("asset not found in the transient map input")
		}	
		
		type AssetTransientPrivateInput struct {
			TransactionID             string `json:"TransactionID"`
			MongoDBJson               string `json:"MongoDBJson"`
		}
		
		// Save order details to collection visible to owning organization
		var assetQtyDetails AssetTransientPrivateInput
		err = json.Unmarshal(transientAssetJSON, &assetQtyDetails)
		if err != nil {
    	    	return "",fmt.Errorf("failed to unmarshal JSON: %v", err)
    		}
		
		assetPrivateDetails := AssetPrivateDetails{
			TransactionID: assetQtyDetails.TransactionID,
			MongoDBJson:   assetQtyDetails.MongoDBJson,
		}
		
		assetPrivateDetailsAsBytes, err := json.Marshal(assetPrivateDetails)
		if err != nil {
			return "",fmt.Errorf("failed to marshal into JSON: %v", err)
		}
		
		// Put asset appraised value into owners org specific private data collection
    	    err = ctx.GetStub().PutPrivateData(assetCollection, assetQtyDetails.TransactionID, assetPrivateDetailsAsBytes)
		if err != nil {
			return "",fmt.Errorf("failed to put asset private details: %v", err)
		}
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(TransactionID)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, TransactionID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(TransactionID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) GetHistoryForAsset(ctx contractapi.TransactionContextInterface, carID string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(carID)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
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
		fmt.Printf("Error create Snapchain chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Snapchain chaincode: %s", err.Error())
	}

}