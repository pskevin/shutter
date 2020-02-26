package node

import (
	"fmt"
	"sync"
)

type Node struct {
	nodeId: int
	balance: int
	inChannels: map[int](chan string)
	outChannels: map[int](chan string)

	sawToken: map[int](bool)
	canProceed: sync.RWMutex
	canRecv: sync.RWMutex

	noMarkerReceived: bool
	nodeState: int
	channelState: map[int]([]int)
	shouldRecordChannelState: map[int](bool)
}

func New(nodeId int, balance int) {
	node = &Node{
		nodeId:	nodeId,
		balance: balance,
		inChannels: make(map[int](chan string)),
		outChannels: make(map[int](chan string)),
		sawToken: make(map[int](bool))
	}

	return node
}


// This is assumed to be a command communicated by the master
func (self *Node) sendMessage(recvId int, amount int) {
	self.canProceed.RLock()
	self.canRecv.Lock()
	if amount > balance {
		fmt.Printf("ERR_SEND")
	} else {
		self.outChannels[recvId] <- amount
		self.balance -= amount
	}
	self.canRecv.Unlock()
	self.canProceed.RUnlock()
}

// TODO
func (self *Node) recvMessage(senderId int) {
	// TODO: make senderId optional in parameters
	select {
	case msg := <- self.inChannels[senderId]:
		if msg == -1 {
			self.canProceed.Lock()
			fmt.Printf("%d SnapshotToken -1", senderId)
			if self.noMarkerReceived {
				self.propagateSnapshot(senderId)
			} else if self.shouldRecordChannelState[senderId] == true {
				self.shouldRecordChannelState[senderId] = false
				doneSnapshot = true
				for _, stillRecording := range self.shouldRecordChannelState {
					doneSnapshot = doneSnapshot && !(stillRecording)
				}
				if doneSnapshot {
					self.firstMarkerReceived = false
				}
			} else {
				fmt.Println("Trying to take a new snapshot while another already going on")
			}
			self.canProceed.Unlock()
		} else {
			self.canProceed.RLock()
			self.canRecv.RLock()
			fmt.Printf("%d Transfer %d", senderId, msg)
			if self.shouldRecordChannelState[senderId] {
				channelState[senderId] = append(channelState[senderId], msg)
			}
			self.updateBalance(msg)
			self.canRecv.RUnlock()
			self.canProceed.RUnlock()
		}
	default:
		fmt.Println("Nothing to receive")
	}
}

func (self *Node) initiateSnapshot() {
	self.canProceed.Lock()
	self.nodeState = self.balance
	self.noMarkerReceived = false

	// Send marker message to all processes
	for _, channel := range self.outChannels {
		channel <- -1
	}

	// Start recording incoming channel messages
	for idx, _ := range self.shouldRecordChannelState {
		self.shouldRecordChannelState[idx] = true
	}
	self.canProceed.Unlock()
}

// Already locked at caller
func (self *Node) propagateSnapshot(senderId int) {
	self.nodeState = self.balance
	self.noMarkerReceived = false

	// Send marker message to all processes
	for _, channel := range self.outChannels {
		channel <- -1
	}

	// Start recording incoming channel messages
	// TODO: finish this
	for idx, _ := range self.shouldRecordChannelState {
		if idx == senderId:
			self.shouldRecordChannelState[idx] = true
	}
}

func (self *Node) updateBalance(amount int) {
	self.balance += amount
}