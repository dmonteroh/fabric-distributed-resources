package pkg

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/dmonteroh/fabric-distributed-resources/internal"
)

func GetAllLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("latency").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToLatencyAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

func GetServersLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("latency").(*gateway.Contract)

	res, err := contract.EvaluateTransaction("GetServerAssets")
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, readRes)
}

func GetServersExceptIdLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	id := c.Param("id")
	readRes := getServersExcept(c, id)
	c.JSON(200, readRes)
}

func GetServersExceptSelfLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	id := c.ClientIP()
	readRes := getServersExcept(c, id)
	c.JSON(200, readRes)
}

func getServersExcept(c *gin.Context, id string) []internal.Asset {
	contract := c.MustGet("latency").(*gateway.Contract)
	res, err := contract.EvaluateTransaction("GetServerAssetsExceptId", id)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.JsonToAssetArray(string(res))
	if err != nil {
		panic(err.Error())
	}
	return readRes
}

func GetLatencyTargetsHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	id := c.ClientIP()
	targetAssets := getServersExcept(c, id)
	targets := internal.LatencyTargets{
		Source:  id,
		Targets: []internal.LatencyTarget{},
	}

	for _, asset := range targetAssets {
		target := internal.LatencyTargetFromMap(asset.Properties)
		targets.Targets = append(targets.Targets, target)
		//// Removed due to properties being added as a struct instead of a
		// if internal.KeysInStringMap(asset.Properties, targetKeys) {
		// 	target := internal.LatencyTargetFromMap(asset.Properties)
		// 	targets.Targets = append(targets.Targets, target)
		// }
	}

	c.JSON(200, targets)
}

func GetLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("latency").(*gateway.Contract)
	asset := c.Param("asset")

	res, err := contract.EvaluateTransaction("ReadAsset", asset)
	if err != nil {
		panic(err.Error())
	}
	readRes, err := internal.LatencyAssetJsonToStruct(string(res))
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, readRes)
}

func UpdateLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("latency").(*gateway.Contract)

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	inventory, err := internal.LatencyAssetJsonToStruct(string(jsonData))
	if err != nil {
		panic(err)
	}

	_, err = contract.SubmitTransaction("UpdateAsset", inventory.String())
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"key": inventory.ID})
}

func CreateLatencyHandler(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	contract := c.MustGet("latency").(*gateway.Contract)
	appType := c.MustGet("APP_TYPE").(string)

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	latencyResults, err := internal.LatencyResultsJsonToStruct(string(jsonData))
	latencyId := internal.CreateLatencyID(appType, latencyResults.Source, latencyResults.Timestamp)
	latencyAsset := internal.CreateLatencyAsset(latencyId, latencyResults)
	if err != nil {
		panic(err)
	}

	_, err = contract.SubmitTransaction("CreateAsset", latencyAsset.String())
	if err != nil {
		panic(err.Error())
	}
	c.JSON(200, gin.H{"key": latencyAsset.ID})
}
