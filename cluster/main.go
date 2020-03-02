package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
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

// KillAll executes 'KillAll' command on the cluster
func KillAll(args ...interface{}) {
	os.Exit(0)
}

// CreateNode executes 'CreateNode' command on the cluster
func CreateNode(args ...interface{}) {
	id, initBalance := toInt(args[0]), toInt(args[1])
	newNode := New(id, initBalance)

	for _, presentID := range initializedNodes {
		n := nodes[presentID]
		chIn := make(chan int, MAX_NODES)
		newNode.CreateIncoming(n.nodeID, chIn)
		n.CreateOutgoing(newNode.nodeID, chIn)

		chOut := make(chan int, MAX_NODES)
		newNode.CreateOutgoing(n.nodeID, chOut)
		n.CreateIncoming(newNode.nodeID, chOut)
	}

	nodes[newNode.nodeID] = newNode
	initializedNodes = append(initializedNodes, id)
}

func Send(args ...interface{}) {
	sId, rId, amount := toInt(args[0]), toInt(args[1]), toInt(args[2])
	sender := nodes[sId]
	sender.SendMessage(rId, amount)
}

func Receive(args ...interface{}) {
	rId := toInt(args[0])
	receiver := nodes[rId]

	switch len(args) {
	case 1:
		receiver.RecvMessage()
	case 2:
		sId := toInt(args[1])
		receiver.RecvMessage(sId)
	}

}

// TODO: Once system is built, convert it to goroutine call
func ReceiveAll(args ...interface{}) {
	var wg sync.WaitGroup
	for _, nodeID := range initializedNodes {
		node := nodes[nodeID]
		wg.Add(1)
		node.RecvNonBlocking(&wg)
	}
	wg.Wait()
}

// This command will be received from Observer, not Master
func BeginSnapshot(args ...interface{}) {
	nodeId := toInt(args[0])
	node := nodes[nodeId]
	node.InitiateSnapshot()
}

// This function will send response to Observer, not Master
func CollectState(args ...interface{}) {
	for _, nodeID := range initializedNodes {
		node := nodes[nodeID]
		// TODO make this run concurrently
		nodeState := node.GetState()
		nodeStates = append(nodeStates, nodeState)
	}
}

// Will be done in observer
// TODO
func PrintSnapshot(args ...interface{}) {

	sort.Slice(nodeStates, func(i, j int) bool {
		return nodeStates[i].nodeID < nodeStates[j].nodeID
	})

	fmt.Println("---Node states")
	for _, state := range nodeStates {
		fmt.Printf("node %d = %d\n", state.nodeID, state.nodeState)
	}

	fmt.Println("---Channel states")
	for _, state := range nodeStates {
		for senderIdx, cStateArray := range state.channelState {
			channelAmount := 0
			for _, val := range cStateArray {
				channelAmount += val
			}
			fmt.Printf("channel (%d -> %d) = %d\n", senderIdx, state.nodeID, channelAmount)
		}
	}
}

// func main() {
// 	CreateNode(1, 1000)
// 	CreateNode(2, 500)
// 	Send(1, 2, 300)
// 	Send(2, 1, 100)
// 	Send(1, 2, 400)
// 	Receive(2, 1)
// 	BeginSnapshot(1)
// 	Receive(2)
// 	Receive(2)
// 	ReceiveAll()
// 	CollectState()
// 	PrintSnapshot()
// 	Send(1, 2, 5000)
// }

func toInt(arg interface{}) int {
	s, ok := arg.(string)
	if ok {
		v, _ := strconv.Atoi(s)
		return v
	}
	return -1
}
