package main

import (
	"os"
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

const MAX_NODES = 100000

var (
	// Arrays are go-routine safe
	nodes            [MAX_NODES]*Node
	initializedNodes []int
	nodeStates       []State
)

// KillAll ...
func KillAll() {
	os.Exit(0)
}

// CreateNode ...
// TODO iterate through nodes only when initialized
func CreateNode(id int, balance int) {
	newNode := New(id, balance)

	for _, presentID := range initializedNodes {
		n := nodes[presentID]
		chIn := make(chan int)
		newNode.CreateIncoming(n.nodeID, chIn)
		n.CreateOutgoing(newNode.nodeID, chIn)

		chOut := make(chan int)
		newNode.CreateOutgoing(n.nodeID, chOut)
		n.CreateIncoming(newNode.nodeID, chOut)
	}

	nodes[newNode.nodeID] = newNode
	initializedNodes = append(initializedNodes, id)
}

// Send ...
func Send(sID int, rID int, amount int) {
	sender := nodes[sID]
	sender.SendMessage(rID, amount)
}

// Receive ...
func Receive(id ...int) {
	rID := id[0]
	receiver := nodes[rID]
	recvParams := id[1:]
	receiver.RecvMessage(recvParams)
}

func ReceiveAll() {
	for _, node := range nodes {
		go node.RecvNonBlocking()
	}
}

// This command will be received from Observer, not Master
func BeginSnapshot(nodeID int) {
	node := nodes[nodeID]
	node.InitiateSnapshot()
}

// This function will send response to Observer, not Master
func CollectState() {
	var states [MAX_NODES]State
	for _, node := range nodes {
		// TODO make this run concurrently
		nodeState := node.GetState()
		nodeStates = append(nodeStates, nodeState)
	}
}

// Will be done in observer
// TODO
func PrintSnapshot() {

}
