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
	latencyContract := internal.GetEnv("LATENCY_SC", "latency-sc")

	// MAP VARIABLES INTO MAP
	variables := map[string]string{
		"APP_TYPE":  appType,
		"EXEC_MODE": execMode,
	}

	// CONNECT TO THE FABRIC NETWORK
	network := initFabric()
	// GET CONTRACTS
	resourcesSC := network.GetContract(resourcesContract)
	inventorySC := network.GetContract(inventoryContract)
	latencySC := network.GetContract(latencyContract)

	// INIT IS FOR DEBUGGING PURPOSES
	// _, err := inentorySC.SubmitTransaction("InitLedger")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }

	// INITIALIZE HTTP SERVER AND ADD MIDDLEWARE
	r := gin.Default()
	r.Use(internal.EnviromentMiddleware(variables))
	r.Use(internal.ContractMiddleware("resources", resourcesSC))
	r.Use(internal.ContractMiddleware("inventory", inventorySC))
	r.Use(internal.ContractMiddleware("latency", latencySC))

	// --- APP HTTP ROUTES
	// ASSETS
	r.GET("/resources", pkg.GetAllResourcesHandler)
	r.GET("/resources/:asset", pkg.GetResourceHandler)
	r.PUT("/resources", pkg.UpdateResourceHandler)
	r.POST("/resources", pkg.UpsertResourceHandler)
	// INVENTORY
	r.GET("/inventory", pkg.GetAllInventoryHandler)
	r.GET("/inventory/servers", pkg.GetServersInventoryHandler)
	r.GET("/inventory/:asset", pkg.GetInventoryHandler)
	r.PUT("/inventory", pkg.UpdateInventoryHandler)
	r.POST("/inventory", pkg.CreateInventoryHandler)
	// LATENCY
	r.GET("/latency", pkg.GetAllLatencyHandler)
	r.GET("/latency/servers", pkg.GetServersLatencyHandler)
	r.GET("/latency/servers/except/self", pkg.GetServersExceptSelfLatencyHandler)
	r.GET("/latency/servers/except/:id", pkg.GetServersExceptIdLatencyHandler)
	r.GET("/latency/servers/targets", pkg.GetLatencyTargetsHandler)
	r.GET("/latency/:asset", pkg.GetLatencyHandler)
	r.PUT("/latency", pkg.UpdateLatencyHandler)
	r.POST("/latency", pkg.CreateLatencyHandler)
	// -- COLLECTOR
	r.POST("/collector", pkg.UpsertResourceHandler)
	r.POST("/measurement", pkg.CreateLatencyHandler)

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
