package cmd

import (
	"fmt"
	"log"
	"net/http"
)

// BeginSnapshot a
func BeginSnapshot(args []string) {
	command, nodeID := "BeginSnapshot", args[0]
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s?nodeID=%s", ClusterURL, command, nodeID))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// CollectState s
func CollectState(args []string) {
	command := "CollectState"
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s", ClusterURL, command))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// CreateNode s
func CreateNode(args []string) {
	command, nodeID, initAmount := "CreateNode", args[0], args[1]
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s?nodeID=%s&initAmount=%s", ClusterURL, command, nodeID, initAmount))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// KillAll s
func KillAll(args []string) {
	command := "KillAll"
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s", ClusterURL, command))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// PrintSnapshot s
func PrintSnapshot(args []string) {
	command := "PrintSnapshot"
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s", ClusterURL, command))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// Receive s
func Receive(args []string) {
	command, receiverID := "Receive", args[0]
	req := fmt.Sprintf("%s/%s?receiverID=%s", ClusterURL, command, receiverID)
	if len(args) == 2 {
		senderID := args[1]
		req = fmt.Sprintf("%s&senderID=%s", req, senderID)
	}

	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// ReceiveAll s
func ReceiveAll(args []string) {
	command := "ReceiveAll"
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s", ClusterURL, command))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}

// Send s
func Send(args []string) {
	command, senderID, receiverID, amount := "Send", args[0], args[1], args[2]
	fmt.Printf("\n%s called\n", command)
	res, err := http.Get(fmt.Sprintf("%s/%s?senderID=%s&receiverID=%s&amount=%s", ClusterURL, command, senderID, receiverID, amount))
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}
