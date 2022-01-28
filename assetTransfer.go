/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
//now I want to add the subscriber to this application and just keep the create asset function
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	////////////////////subscriber////////////
	"context"
	"strconv"
	"sync"
	"time"

	//////////////////////////////////////////
	geometry_msgs "github.com/TIERS/rclgo-msgs/geometry_msgs/msg"
	std_msgs "github.com/TIERS/rclgo-msgs/std_msgs/msg"
	"github.com/TIERS/rclgo/pkg/rclgo"
)

var contract *gateway.Contract

var message_counter = 0
var face_counter = 0
var counter = "drone"

var drone_x_UWB, drone_y_UWB float64
var count string
var pos_x []string
var pos_y []string
var drone_x, drone_y string
var drone_x_position, drone_y_position string

/////////////////////UWB/////////////////////
func Handler_UWB(s *rclgo.Subscription) {

	msg_position := geometry_msgs.PoseStampedTypeSupport.New()
	_, err := s.TakeMessage(msg_position)
	if err != nil {
		fmt.Println("failed to take message:", err)
		return
	}
	// fmt.Println(s.TopicName)
	// fmt.Println(s.Ros2MsgType.TypeSupport())
	drone_x_UWB = msg_position.(*geometry_msgs.PoseStamped).Pose.Position.X
	drone_y_UWB = msg_position.(*geometry_msgs.PoseStamped).Pose.Position.Y
	drone_x_position = strconv.FormatFloat(drone_x_UWB, 'f', 4, 64)
	drone_y_position = strconv.FormatFloat(drone_y_UWB, 'f', 4, 64)
}

/////////////////////face detector///////////////
func Handler_face_detector(s *rclgo.Subscription) {
	msg_face_found := std_msgs.StringTypeSupport.New()
	_, err2 := s.TakeMessage(msg_face_found)
	if err2 != nil {
		fmt.Println("failed to take message:", err2)
		return
	}
	// fmt.Println(s.TopicName)
	// fmt.Println(s.Ros2MsgType.TypeSupport())
	fmt.Println(msg_face_found.(*std_msgs.String).Data)
	face_counter += 1
	count = fmt.Sprint("XFace ", face_counter, " - ", time.Now().Format(time.RFC850))
	// salma := fmt.Sprint(drone_x_UWB[len(drone_x_UWB)-1], drone_y_UWB[len(drone_y_UWB)-1])
	log.Println("--> Submit Transaction: CreateLanding, creates new landing with ID, color, owner, size, and appraisedValue arguments")
	result, err := contract.SubmitTransaction("CreateLanding", count, drone_x_position, drone_y_position)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: GetAllLandings, function returns all the current landings on the ledger")
	result, err = contract.EvaluateTransaction("GetAllLandings")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

}

///////////////////////simple publisher/////////////////////
// func Handler(s *rclgo.Subscription) {
// 	msg_object_found := std_msgs.StringTypeSupport.New()
// 	_, err2 := s.TakeMessage(msg_object_found)
// 	if err2 != nil {
// 		fmt.Println("failed to take message:", err2)
// 		return
// 	}
// 	fmt.Println(s.TopicName)
// 	fmt.Println(s.Ros2MsgType.TypeSupport())
// 	fmt.Println(msg_object_found.(*std_msgs.String).Data)
// 	msg_data := msg_object_found.(*std_msgs.String).Data
// 	data_list := strings.Split(msg_data, ",")
// 	// drone_name := data_list[0]
// 	drone_x = data_list[1]
// 	drone_y = data_list[2]
// 	// ////////////////////////////////////////////
// 	// log.Println("--> Submit Transaction: CreateLanding, creates new landing with ID, color, owner, size, and appraisedValue arguments")
// 	// result, err := contract.SubmitTransaction("CreateLanding", drone_name, drone_x, drone_y)
// 	// //result, err := contract.SubmitTransaction("CreateLanding", msg.(*std_msgs.String).Data, "10", "20")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to Submit transaction: %v", err)
// 	// }
// 	// log.Println(string(result))

// 	// log.Println("--> Evaluate Transaction: GetAllLandings, function returns all the current landings on the ledger")
// 	// result, err = contract.EvaluateTransaction("GetAllLandings")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to evaluate transaction: %v", err)
// 	// }
// 	// log.Println(string(result))

// }

func main() {
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

	contract = network.GetContract("basic")

	///////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////

	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of landings on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	//////////////////////////////////////////////////////////////////

	log.Println("--> Evaluate Transaction: GetAllLandings, function returns all the current landings on the ledger")
	result, err = contract.EvaluateTransaction("GetAllLandings")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
	//////////////////////////////////////////////////////////////////

	log.Println("============ ROS node starts ============")
	var doneChannel = make(chan bool)
	var wg sync.WaitGroup

	ctx, quitFunc := context.WithCancel(context.Background())

	// Setup ROS nodes
	rclArgs, rclErr := rclgo.NewRCLArgs("")
	if rclErr != nil {
		log.Fatal(rclErr)
	}

	rclContext, rclErr := rclgo.NewContext(&wg, 0, rclArgs)
	if rclErr != nil {
		log.Fatal(rclErr)
	}
	defer rclContext.Close()

	rclNode, rclErr := rclContext.NewNode("communicate", "publisher_test")
	if rclErr != nil {
		log.Fatal(rclErr)
	}
	/////////////////////UWB/////////////////////
	sub, _ := rclNode.NewSubscription("/position", geometry_msgs.PoseStampedTypeSupport, Handler_UWB)
	go func() {
		err := sub.Spin(ctx, 1*time.Second)
		log.Printf("Subscription failed: %v", err)
	}()
	////////////////face detection///////////////
	sub2, _ := rclNode.NewSubscription("/face_found", std_msgs.StringTypeSupport, Handler_face_detector)
	go func() {
		err := sub2.Spin(ctx, 1*time.Second)
		log.Printf("Subscription failed: %v", err)
	}()
	///////////////object detection//////////////
	// sub3, _ := rclNode.NewSubscription("/object_found", std_msgs.StringTypeSupport, Handler)
	// go func() {
	// 	err := sub3.Spin(ctx, 1*time.Second)
	// 	log.Printf("Subscription failed: %v", err)
	// }()
	////////////////////////////////////////////////////////////////////
	log.Println("--> Evaluate Transaction: GetAllLandings, function returns all the current landings on the ledger")
	result, err = contract.EvaluateTransaction("GetAllLandings")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
	////////////////////////////////////////////////////////////////////
	<-doneChannel
	quitFunc()
	wg.Wait()
	log.Printf("Signing off - BYE")
	////////////////////////////////////////////////////////////////////
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
