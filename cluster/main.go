package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
)

var (
	observerBegin   chan int
	observerCollect chan int
	observerPrint   chan int
)

func main() {
	observerBegin = make(chan int)
	observerCollect = make(chan int)
	observerPrint = make(chan int)
	go runObserver()

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
	go sender.SendMessage(rId, amount)
}

func Receive(args ...interface{}) {
	rId := toInt(args[0])
	receiver := nodes[rId]

	switch len(args) {
	case 1:
		go receiver.RecvMessage()
	case 2:
		sId := toInt(args[1])
		go receiver.RecvMessage(sId)
	}

}

// TODO: Once system is built, convert it to goroutine call
func ReceiveAll(args ...interface{}) {
	var wg sync.WaitGroup
	for _, nodeID := range initializedNodes {
		node := nodes[nodeID]
		wg.Add(1)
		go node.RecvNonBlocking(&wg)
	}
	wg.Wait()
}

// This command will be received from Observer, not Master
func ObserverBeginSnapshot(nodeID int) {
	node := nodes[nodeID]
	node.InitiateSnapshot()
}

// This function will send response to Observer, not Master
func ObserverCollectState() {
	for _, nodeID := range initializedNodes {
		node := nodes[nodeID]
		nodeState := node.GetState()
		nodeStates = append(nodeStates, nodeState)
	}
}

// Will be done in observer
func ObserverPrintSnapshot() {

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

func BeginSnapshot(args ...interface{}) {
	nodeId := toInt(args[0])
	observerBegin <- nodeId
}

func CollectState(args ...interface{}) {
	observerCollect <- 1
}

func PrintSnapshot(args ...interface{}) {
	observerPrint <- 1
}

func runObserver() {
	for {
		select {
		case nodeID := <-observerBegin:
			ObserverBeginSnapshot(nodeID)
		case <-observerCollect:
			ObserverCollectState()
		case <-observerPrint:
			ObserverPrintSnapshot()
		default:
		}
	}
}

func toInt(arg interface{}) int {
	s, ok := arg.(string)
	if ok {
		v, _ := strconv.Atoi(s)
		return v
	}
	return -1
}
