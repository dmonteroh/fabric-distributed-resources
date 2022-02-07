/*
Copyright 2020 IBM All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	// ENVIROMENTAL VARIABLES
	//execMode := internal.GetEnv("EXEC_MODE", "DEBUG")
	//listenPort := internal.GetEnv("INTERNAL_PORT", "8080")

	// CONNECT TO THE FABRIC NETWORK
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")

	// FABRIC CALLS

	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: CreateAsset, creates new asset with ID, color, owner, size, and appraisedValue arguments")
	result, err = contract.SubmitTransaction("CreateAsset", "localhost", "{\"timestamp\":{\"timeLocal\":\"2022-02-07T11:54:21.222970965Z\",\"timeSeconds\":1644234861,\"timeNano\":1644234861222970965},\"host\":{\"hostname\":\"426ede137da2\",\"uptime\":16061,\"boottime\":1644218801,\"platform\":\"alpine\",\"virtualizationSystem\":\"docker\",\"virtualizationRole\":\"guest\",\"hostid\":\"a72ab14c-76bf-ea11-8105-842afd4cdfcb\"},\"cpuStats\":{\"modelName\":\"Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz\",\"vendorId\":\"GenuineIntel\",\"averageUsage\":9.746240601477155,\"coreUsage\":[5.2631578947179465,9.523809523840459,9.523809523757965,9.999999999954525,4.999999999977263,5.2631578947179465,23.809523809286638,23.809523809503187,5.000000000022737,5.000000000022737,10.000000000045475,4.7619047618789825]},\"memStats\":{\"total\":67268472832,\"available\":60582289408,\"used\":5795409920},\"diskStats\":[{\"device\":\"/dev/sda5\",\"path\":\"/app\",\"label\":\"\",\"fstype\":\"ext4\",\"total\":78693273600,\"used\":63548289024,\"usedPercent\":85.1262654153894}],\"procStats\":{\"totalProcs\":2182,\"createdProcs\":176236,\"runningProcs\":1,\"blockedProcs\":0},\"dockerStats\":[{\"containerID\":\"6440e6be7720f2fd842aea47b6a87de484414561a212acd71df7a6e04a915cee\",\"name\":\"/dev-peer0.org2.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981\",\"image\":\"dev-peer0.org2.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981-9e2ae745b02b13626fa3a2f2d71e307d9bfe37fc81485ce810802cb0859f0872\",\"status\":\"Up 5 minutes\",\"State\":\"running\"},{\"containerID\":\"f53321eee853444cf550909a339ed5182d74e8490a6dc7e64fa16f26fdde5154\",\"name\":\"/dev-peer0.org1.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981\",\"image\":\"dev-peer0.org1.example.com-basic_1.0-5f042bbcb3e3b1b4b6e6a25f30f746f263614a8b838865b6f72deb9cbd8ab981-1b5f8eb9971213e57fd41c10cac95d657f48b197f27f836e8cc606f0a4cf27fc\",\"status\":\"Up 5 minutes\",\"State\":\"running\"},{\"containerID\":\"670c610a8b5641cd060585be5d607fe1f0365974d4c6c02f915feb5ebc64b1a5\",\"name\":\"/cli\",\"image\":\"hyperledger/fabric-tools:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"8d933fb4377bb80f5eead71ccf160912e53abf54f62f96290a625cca89ea3800\",\"name\":\"/peer0.org2.example.com\",\"image\":\"hyperledger/fabric-peer:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"24069d7b1df3f28f556276af70be52d09dc3cf323581a327d7a7ad3084ca1396\",\"name\":\"/peer0.org1.example.com\",\"image\":\"hyperledger/fabric-peer:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"c9d549992b95170d7ae750c5dc6bf464b03c37188d33476dbb0e87126d52af6e\",\"name\":\"/orderer.example.com\",\"image\":\"hyperledger/fabric-orderer:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"f29776a0078c97076ca2ea5b4d3121bb4c4e307f676a6ba88981b7bb13195c26\",\"name\":\"/couchdb1\",\"image\":\"couchdb:3.1.1\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"3493a420a5df5b12efa116c279d57f46aecbbfaf0a820c0208a2462b7a068726\",\"name\":\"/couchdb0\",\"image\":\"couchdb:3.1.1\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"8ee4ad099a3e4446274444df66cfe3cb38aa257cfd5cf1b0012bec9ba71ed284\",\"name\":\"/ca_orderer\",\"image\":\"hyperledger/fabric-ca:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"d1b7e077111df8525e738388e70d99b5e75d8d611b3e6b1c8994a90f963d7832\",\"name\":\"/ca_org2\",\"image\":\"hyperledger/fabric-ca:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"674b38d89aa0d381c1b7476d773d0b346081a3c7202649c0aa81ad5b34661b7a\",\"name\":\"/ca_org1\",\"image\":\"hyperledger/fabric-ca:latest\",\"status\":\"Up 28 minutes\",\"State\":\"running\"},{\"containerID\":\"426ede137da2f30fa6aede67b757b6545a51d51e30a673e9cbf8cbebe4b4f0b2\",\"name\":\"/distributed-resource-collector\",\"image\":\"distributed-resource-collector\",\"status\":\"Up 3 seconds\",\"State\":\"running\"}]}")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadAsset, function returns an asset with a given assetID")
	result, err = contract.EvaluateTransaction("ReadAsset", "localhost")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: AssetExists, function returns 'true' if an asset with given assetID exist")
	result, err = contract.EvaluateTransaction("AssetExists", "localhost")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	// log.Println("--> Submit Transaction: TransferAsset asset1, transfer to new owner of Tom")
	// _, err = contract.SubmitTransaction("TransferAsset", "asset1", "Tom")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }

	// log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
	// result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
	// if err != nil {
	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
	// }
	// log.Println(string(result))
	log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
