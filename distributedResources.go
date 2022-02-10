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

	"github.com/dmonteroh/fabric-distributed-resources/internal"
	"github.com/dmonteroh/fabric-distributed-resources/pkg"
	"github.com/gin-gonic/gin"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	// ENVIROMENTAL VARIABLES
	appType := internal.GetEnv("APP_TYPE", "single_insert")
	execMode := internal.GetEnv("EXEC_MODE", "DEBUG")
	listenPort := internal.GetEnv("INTERNAL_PORT", "8080")
	resourcesContract := internal.GetEnv("RESROUCES_SC", "resources-sc")
	inventoryContract := internal.GetEnv("INVENTORY_SC", "inventory-sc")

	// MAP VARIABLES INTO MAP
	variables := map[string]string{
		"APP_TYPE":  appType,
		"EXEC_MODE": execMode,
	}

	// CONNECT TO THE FABRIC NETWORK
	network := initFabric()
	// GET CONTRACTS
	resourcesSC := network.GetContract(resourcesContract)
	inentorySC := network.GetContract(inventoryContract)

	// INIT IS FOR DEBUGGING PURPOSES
	// _, err := inentorySC.SubmitTransaction("InitLedger")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }

	// INITIALIZE HTTP SERVER AND ADD MIDDLEWARE
	r := gin.Default()
	r.Use(internal.EnviromentMiddleware(variables))
	r.Use(internal.ContractMiddleware(resourcesSC))
	r.Use(internal.ContractMiddleware(inentorySC))

	// SAVE VARIABLES INSIDE GIN CONTEXT

	// APP HTTP ROUTES
	r.GET("/assets", pkg.GetAllAssetsHandler)
	r.GET("/assets/:asset", pkg.GetAssetHandler)
	r.POST("/assets", pkg.UpsertAssetHandler)
	r.PUT("/assets", pkg.UpdateAssetHandler)
	r.POST("/collector", pkg.UpsertAssetHandler)

	// START HTTP SERVER
	r.Run(":" + listenPort)
}

func initFabric() *gateway.Network {
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

	return network
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
