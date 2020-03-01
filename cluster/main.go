package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	// Master-Cluster Communications Initialization
	routeHandlerMap := RouteHandlerMap{
		"KillAll":       {handler: KillAll, block: false},
		"CreateNode":    {handler: CreateNode, block: true},
		"Send":          {handler: Send, block: false},
		"Receive":       {handler: Receive, block: false},
		"ReceiveAll":    {handler: ReceiveAll, block: false},
		"BeginSnapshot": {handler: BeginSnapshot, block: false},
		"CollectState":  {handler: CollectState, block: false},
		"PrintSnapshot": {handler: PrintSnapshot, block: false},
	}
	InitMasterClusterComm(routeHandlerMap)
}

/*

TODO: Shift the command handlers into different files instead of collecting all of them in one file

*/

// // KillAll executes 'KillAll' command on the cluster
// func KillAll(args ...interface{}) {
// 	fmt.Printf("\nAt function KillAll: %v", args)
// }

// // CreateNode executes 'CreateNode' command on the cluster
// func CreateNode(args ...interface{}) {
// 	time.Sleep(time.Millisecond)
// 	fmt.Printf("\nAt function CreateNode: %v", args)
// }

// // Send executes 'Send' command on the cluster
// func Send(args ...interface{}) {
// 	fmt.Printf("\nAt function Send: %v", args)
// }

// // Receive executes 'Receive' command on the cluster
// func Receive(args ...interface{}) {
// 	fmt.Printf("\nAt function Receive: %v", args)
// }

// // ReceiveAll executes 'ReceiveAll' command on the cluster
// func ReceiveAll(args ...interface{}) {
// 	fmt.Printf("\nAt function ReceiveAll: %v", args)
// }

// // BeginSnapshot executes 'BeginSnapshot' command on the cluster
// func BeginSnapshot(args ...interface{}) {
// 	fmt.Printf("\nAt function BeginSnapshot: %v", args)
// }

// // CollectState executes 'CollectState' command on the cluster
// func CollectState(args ...interface{}) {
// 	fmt.Printf("\nAt function CollectState: %v", args)
// }

// // PrintSnapshot blah
// func PrintSnapshot(args ...interface{}) {
// 	fmt.Printf("\nAt function PrintSnapshot: %v", args)
// }

var (
	MAX_NODES = 100000

	// Arrays are go-routine safe
	nodes [MAX_NODES]Node
)


func KillAll() {
	os.Exit(0)
}

// TODO iterate through nodes only when initialized
func CreateNode(id int, balance int) {
	new_node := New(id, balance)

	for _, n := range nodes {
		ch_in := make(chan int)
		new_node.CreateIncoming(n.nodeId, ch_in)
		n.CreateOutgoing(new_node.nodeId, ch_in)

		ch_out := make(chan int)
		new_node.CreateOutgoing(n.nodeId, ch_out)
		n.CreateIncoming(new_node.nodeId, ch_out)
	}

	nodes[new_node.nodeId] =  new_node
}

func Send(sId int, rId int, amount int) {
	sender := nodes[sId]
	sender.SendMessage(rId, amount)
}

func Receive(id ...int) {
	rId := id[0]
	receiver := nodes[rId]
	receiver.RecvMessage(id[1:])
}

func ReceiveAll() {
	for _, node := range nodes {
		go node.RecvNonBlocking()
	}
}

// This command will be received from Observer, not Master 
func BeginSnapshot(nodeId int) {
	node := nodes[nodeId]
	node.InitiateSnapshot()
}

// This function will send response to Observer, not Master
//TODO finish concurrent creations of states
func CollectState() {
	states := [MAX_NODES]State
	for _, node := range nodes {
		go node.GetState()
	}
}

// Will be done in observer
// TODO
func PrintSnapshot() {
}
