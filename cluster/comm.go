package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// RouteHandlerMap registers function to execute on route
type RouteHandlerMap map[string]struct {
	handler func(...interface{})
	block   bool
}

var routeHandlerMap RouteHandlerMap

// InitMasterClusterComm setups communication between Master and the Cluster
func InitMasterClusterComm(_routeHandlerMap RouteHandlerMap) {
	routeHandlerMap = _routeHandlerMap

	// Setup router
	router := mux.NewRouter()
	router.HandleFunc("/KillAll", killAllRoute).Methods("GET")
	router.HandleFunc("/CreateNode", createNodeRoute).Methods("GET")
	router.HandleFunc("/Send", sendRoute).Methods("GET")
	router.HandleFunc("/Receive", receiveRoute).Methods("GET")
	router.HandleFunc("/ReceiveAll", receiveAllRoute).Methods("GET")
	router.HandleFunc("/BeginSnapshot", beginSnapshotRoute).Methods("GET")
	router.HandleFunc("/CollectState", collectStateRoute).Methods("GET")
	router.HandleFunc("/PrintSnapshot", printSnapshotRoute).Methods("GET")

	// Starting Cluster Server
	fmt.Println("\nStarting Cluster Server @ http://localhost:8118")
	if err := http.ListenAndServe(":8118", router); err != nil {
		fmt.Printf("\nmux server: %v\n", err)
	}
}

// killAllRoute handles 'KillAll' command from Master
func killAllRoute(w http.ResponseWriter, r *http.Request) {
	route := "KillAll"

	handleRoute(route)
}

// createNodeRoute handles 'CreateNode' command from Master
func createNodeRoute(w http.ResponseWriter, r *http.Request) {
	route := "CreateNode"

	query := r.URL.Query()
	nodeID, initAmount := query.Get("nodeID"), query.Get("initAmount")

	handleRoute(route, nodeID, initAmount)
}

// sendRoute handles 'Send' command from Master
func sendRoute(w http.ResponseWriter, r *http.Request) {
	route := "Send"

	query := r.URL.Query()
	senderID, receiverID, amount := query.Get("senderID"), query.Get("receiverID"), query.Get("amount")
	handleRoute(route, senderID, receiverID, amount)
}

// receiveRoute handles 'Receive' command from Master
func receiveRoute(w http.ResponseWriter, r *http.Request) {
	route := "Receive"

	query := r.URL.Query()
	receiverID, senderID := query.Get("receiverID"), query.Get("senderID")

	// Optional Sender ID
	if senderID != "" {
		handleRoute(route, receiverID, senderID)
	} else {
		handleRoute(route, receiverID, senderID)
	}
}

// receiveAllRoute handles 'ReceiveAll' command from Master
func receiveAllRoute(w http.ResponseWriter, r *http.Request) {
	route := "ReceiveAll"

	handleRoute(route)
}

// beginSnapshotRoute handles 'BeginSnapshot' command from Master
func beginSnapshotRoute(w http.ResponseWriter, r *http.Request) {
	route := "BeginSnapshot"

	query := r.URL.Query()
	nodeID := query.Get("nodeID")

	handleRoute(route, nodeID)
}

// collectStateRoute handles 'CollectState' command from Master
func collectStateRoute(w http.ResponseWriter, r *http.Request) {
	route := "CollectState"

	handleRoute(route)
}

// printSnapshotRoute handles 'PrintSnapshot' command from Master
func printSnapshotRoute(w http.ResponseWriter, r *http.Request) {
	route := "PrintSnapshot"

	handleRoute(route)
}

func handleRoute(route string, args ...interface{}) {
	fmt.Printf("\nAt /%s endpoint.", route)

	start := time.Now()
	if routeHandlerMap[route].block == true {
		routeHandlerMap[route].handler(args...)
	} else {
		go routeHandlerMap[route].handler(args...)
	}
	fmt.Printf("\n%s took %v", route, time.Since(start))
}
