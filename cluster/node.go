package main

import (
	"fmt"
	"os"
	"sync"
)

// Node ...
type Node struct {
	nodeID      int
	balance     int
	inChannels  map[int](chan int)
	outChannels map[int](chan int)

	canProceed sync.RWMutex
	canRecv    sync.RWMutex

	noMarkerReceived         bool
	finishedSnapshot         bool
	nodeState                int
	channelState             map[int]([]int)
	shouldRecordChannelState map[int](bool)
}

// State ...
type State struct {
	nodeID       int
	nodeState    int
	channelState map[int]([]int)
}

// New ...
func New(nodeID int, balance int) {
	node = &Node{
		nodeID:      nodeID,
		balance:     balance,
		inChannels:  make(map[int](chan int)),
		outChannels: make(map[int](chan int)),
		// canProceed: true,
		// canRecv: true,

		noMarkerReceived:         true,
		finishedSnapshot:         false,
		channelState:             make(map[int]([]int)),
		shouldRecordChannelState: make(map[int]bool),
	}

	return node
}

// NewState ...
func NewState(nodeID int, nodeState int, channelState map[int]([]int)) {
	state = &State{
		nodeID:       nodeID,
		nodeState:    nodeState,
		channelState: channelState,
	}

	return state
}

// CreateIncoming ...
func (this_node *Node) CreateIncoming(source int, channel chan int) {
	this_node.inChannels[source] = channel
}

// CreateOutgoing ...
func (this_node *Node) CreateOutgoing(dest int, channel chan int) {
	this_node.outChannels[dest] = channel
}

// SendMessage ...
// This is assumed to be a command communicated by the master
func (this_node *Node) SendMessage(recvID int, amount int) {
	this_node.canProceed.RLock()
	this_node.canRecv.Lock()
	if amount > balance {
		fmt.Printf("ERR_SEND")
	} else {
		this_node.outChannels[recvID] <- amount
		this_node.balance -= amount
	}
	this_node.canRecv.Unlock()
	this_node.canProceed.RUnlock()
}

// RecvMessage ...
// Here we assume that when senderID is given,
// the call is blocking. When senderID is not given,
// we try all inchannels to see if there is anything to
// receive. If we receive a message, we stop trying other
// channels. If not, we block on sender 0 (if we are not 0),
// 1 otherwise
// TODO: Reduce code repetition
func (this_node *Node) RecvMessage(sender ...int) {
	// TODO: make senderID optional in parameters
	senderSpecified := (len(sender) == 1)
	if len(sender) > 1 {
		os.Exit(-1)
	}
	if senderSpecified {
		senderID := sender[0]
	}

	if senderSpecified {
		select {
		case msg := <-this_node.inChannels[senderID]:
			if msg == -1 {
				this_node.canProceed.Lock()
				this_node.canRecv.RLock()
				fmt.Printf("%d SnapshotToken -1", senderID)
				if this_node.noMarkerReceived {
					this_node.propagateSnapshot(senderID)
				} else if this_node.shouldRecordChannelState[senderID] == true {
					this_node.shouldRecordChannelState[senderID] = false
					doneSnapshot := true
					for _, stillRecording := range this_node.shouldRecordChannelState {
						doneSnapshot = doneSnapshot && !(stillRecording)
					}
					if doneSnapshot {
						this_node.firstMarkerReceived = false
						this_node.finishedSnapshot = true
					}
				} else {
					fmt.Println("Trying to take a new snapshot while another already going on!")
				}
				this_node.canRecv.RUnlock()
				this_node.canProceed.Unlock()
			} else {
				this_node.canProceed.RLock()
				this_node.canRecv.RLock()
				fmt.Printf("%d Transfer %d", senderID, msg)
				if this_node.shouldRecordChannelState[senderID] {
					channelState[senderID] = append(channelState[senderID], msg)
				}
				this_node.updateBalance(msg)
				this_node.canRecv.RUnlock()
				this_node.canProceed.RUnlock()
			}
		}
	} else {
		recvdFlag := false
		for idx := range this_node.inChannels {
			select {
			case msg := <-this_node.inChannels[idx]:
				if msg == -1 {
					this_node.canProceed.Lock()
					this_node.canRecv.RLock()
					fmt.Printf("%d SnapshotToken -1", idx)
					if this_node.noMarkerReceived {
						this_node.propagateSnapshot(idx)
					} else if this_node.shouldRecordChannelState[idx] == true {
						this_node.shouldRecordChannelState[idx] = false
						doneSnapshot := true
						for _, stillRecording := range this_node.shouldRecordChannelState {
							doneSnapshot = doneSnapshot && !(stillRecording)
						}
						if doneSnapshot {
							this_node.firstMarkerReceived = false
							this_node.finishedSnapshot = true
						}
					} else {
						fmt.Println("Trying to take a new snapshot while another already going on!")
					}
					this_node.canRecv.RUnlock()
					this_node.canProceed.Unlock()
				} else {
					this_node.canProceed.RLock()
					this_node.canRecv.RLock()
					fmt.Printf("%d Transfer %d", idx, msg)
					if this_node.shouldRecordChannelState[idx] {
						channelState[idx] = append(channelState[idx], msg)
					}
					this_node.updateBalance(msg)
					this_node.canRecv.RUnlock()
					this_node.canProceed.RUnlock()
				}
				recvdFlag = true
				break
			default:
			}
		}

		if recvdFlag == false {
			if this_node.nodeID == 0 {
				senderID = 1
			} else {
				senderID = 0
			}

			select {
			case msg := <-this_node.inChannels[senderID]:
				if msg == -1 {
					this_node.canProceed.Lock()
					this_node.canRecv.RLock()
					fmt.Printf("%d SnapshotToken -1", senderID)
					if this_node.noMarkerReceived {
						this_node.propagateSnapshot(senderID)
					} else if this_node.shouldRecordChannelState[senderID] == true {
						this_node.shouldRecordChannelState[senderID] = false
						doneSnapshot := true
						for _, stillRecording := range this_node.shouldRecordChannelState {
							doneSnapshot = doneSnapshot && !(stillRecording)
						}
						if doneSnapshot {
							this_node.firstMarkerReceived = false
							this_node.finishedSnapshot = true
						}
					} else {
						fmt.Println("Trying to take a new snapshot while another already going on!")
					}
					this_node.canRecv.RUnlock()
					this_node.canProceed.Unlock()
				} else {
					this_node.canProceed.RLock()
					this_node.canRecv.RLock()
					fmt.Printf("%d Transfer %d", senderID, msg)
					if this_node.shouldRecordChannelState[senderID] {
						channelState[senderID] = append(channelState[senderID], msg)
					}
					this_node.updateBalance(msg)
					this_node.canRecv.RUnlock()
					this_node.canProceed.RUnlock()
				}
			}
		}
	}
}

// InitiateSnapshot ...
func (this_node *Node) InitiateSnapshot() {
	this_node.canProceed.Lock()
	this_node.nodeState = this_node.balance
	this_node.noMarkerReceived = false

	// Send marker message to all processes
	for _, channel := range this_node.outChannels {
		channel <- -1
	}

	// Start recording incoming channel messages
	for idx := range this_node.shouldRecordChannelState {
		this_node.shouldRecordChannelState[idx] = true
	}
	this_node.canProceed.Unlock()
}

// PropagateSnapshot ...
// Already locked at caller
func (this_node *Node) PropagateSnapshot(senderID int) {
	this_node.nodeState = this_node.balance
	this_node.noMarkerReceived = false

	// Send marker message to all processes
	for _, channel := range this_node.outChannels {
		channel <- -1
	}

	// Start recording incoming channel messages
	for idx := range this_node.shouldRecordChannelState {
		if idx != senderID {
			this_node.shouldRecordChannelState[idx] = true
		}
		this_node.channelState[senderID] = nil
	}
}

func (this_node *Node) updateBalance(amount int) {
	this_node.balance += amount
}

func (this_node *Node) recvNonBlocking() {
	for idx := range this_node.inChannels {
		select {
		case msg := <-this_node.inChannels[idx]:
			if msg == -1 {
				this_node.canProceed.Lock()
				this_node.canRecv.RLock()
				// TODO - remove print below
				fmt.Printf("%d SnapshotToken -1", idx)
				if this_node.noMarkerReceived {
					this_node.propagateSnapshot(idx)
				} else if this_node.shouldRecordChannelState[idx] == true {
					this_node.shouldRecordChannelState[idx] = false
					doneSnapshot := true
					for _, stillRecording := range this_node.shouldRecordChannelState {
						doneSnapshot = doneSnapshot && !(stillRecording)
					}
					if doneSnapshot {
						this_node.firstMarkerReceived = false
						this_node.finishedSnapshot = true
					}
				} else {
					fmt.Println("Trying to take a new snapshot while another already going on!")
				}
				this_node.canRecv.RUnlock()
				this_node.canProceed.Unlock()
			} else {
				this_node.canProceed.RLock()
				this_node.canRecv.RLock()
				// TODO Remove print below
				fmt.Printf("%d Transfer %d", idx, msg)
				if this_node.shouldRecordChannelState[idx] {
					channelState[idx] = append(channelState[idx], msg)
				}
				this_node.updateBalance(msg)
				this_node.canRecv.RUnlock()
				this_node.canProceed.RUnlock()
			}
		default:
		}
	}
}

// GetState ...
func (this_node *Node) GetState() {
	state = NewState(this_node.nodeID, this_node.nodeState, this_node.channelState)
	return state
}
