package main

import "fmt"

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

// KillAll executes 'KillAll' command on the cluster
func KillAll(args ...interface{}) {
	fmt.Printf("\nAt function KillAll: %v", args)
}

// CreateNode executes 'CreateNode' command on the cluster
func CreateNode(args ...interface{}) {
	fmt.Printf("\nAt function CreateNode: %v", args)
}

// Send executes 'Send' command on the cluster
func Send(args ...interface{}) {
	fmt.Printf("\nAt function Send: %v", args)
}

// Receive executes 'Receive' command on the cluster
func Receive(args ...interface{}) {
	fmt.Printf("\nAt function Receive: %v", args)
}

// ReceiveAll executes 'ReceiveAll' command on the cluster
func ReceiveAll(args ...interface{}) {
	fmt.Printf("\nAt function ReceiveAll: %v", args)
}

// BeginSnapshot executes 'BeginSnapshot' command on the cluster
func BeginSnapshot(args ...interface{}) {
	fmt.Printf("\nAt function BeginSnapshot: %v", args)
}

// CollectState executes 'CollectState' command on the cluster
func CollectState(args ...interface{}) {
	fmt.Printf("\nAt function CollectState: %v", args)
}

// PrintSnapshot blah
func PrintSnapshot(args ...interface{}) {
	fmt.Printf("\nAt function PrintSnapshot: %v", args)
}
