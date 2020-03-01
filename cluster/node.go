package main

import (
	"fmt"
	"os"
	"sync"
)

type Node struct {
	nodeId int
	balance int
	inChannels map[int](chan int)
	outChannels map[int](chan int)

	canProceed sync.RWMutex
	canRecv sync.RWMutex

	noMarkerReceived bool
	finishedSnapshot bool
	nodeState int
	channelState map[int]([]int)
	shouldRecordChannelState map[int](bool)
}

type State struct {
	nodeId int
	nodeState int
	channelState map[int]([] int)
}

func New(nodeId int, balance int) {
	node = &Node{
		nodeId:	nodeId,
		balance: balance,
		inChannels: make(map[int](chan int)),
		outChannels: make(map[int](chan int)),
	}

	return node
}

func NewState(nodeId int, nodeState int, channelState map[int]([] int)) {
	state = &State{
		nodeId: nodeId,
		nodeState: nodeState,
		channelState: channelState,
	}

	return state
}

func (self *Node) CreateIncoming(source int, channel chan) {
	self.inChannels[source] = channel
}

func (self *Node) CreateOutgoing(dest int, channel chan) {
	self.outChannels[dest] = channel
}


// This is assumed to be a command communicated by the master
func (self *Node) SendMessage(recvId int, amount int) {
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


// Here we assume that when senderId is given,
// the call is blocking. When senderId is not given,
// we try all inchannels to see if there is anything to
// receive. If we receive a message, we stop trying other
// channels. If not, we block on sender 0 (if we are not 0),
// 1 otherwise
// TODO: Reduce code repetition
func (self *Node) RecvMessage(sender ...int) {
	// TODO: make senderId optional in parameters
	senderSpecified := (len(sender)==1)
	if len(sender) > 1 {
		os.Exit(-1)
	}
	senderId := sender[0]

	if senderSpecified {
		select {
		case msg := <- self.inChannels[senderId]:
			if msg == -1 {
				self.canProceed.Lock()
				self.canRecv.RLock()
				fmt.Printf("%d SnapshotToken -1", senderId)
				if self.noMarkerReceived {
					self.propagateSnapshot(senderId)
				} else if self.shouldRecordChannelState[senderId] == true {
					self.shouldRecordChannelState[senderId] = false
					doneSnapshot := true
					for _, stillRecording := range self.shouldRecordChannelState {
						doneSnapshot = doneSnapshot && !(stillRecording)
					}
					if doneSnapshot {
						self.firstMarkerReceived = false
						self.finishedSnapshot = true
					}
				} else {
					fmt.Println("Trying to take a new snapshot while another already going on!")
				}
				self.canRecv.RUnlock()
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
		}
	} else {
		recvdFlag := false
		for idx, _ := range self.inChannels {
			select {
			case msg := <- self.inChannels[idx]:
				if msg == -1 {
					self.canProceed.Lock()
					self.canRecv.RLock()
					fmt.Printf("%d SnapshotToken -1", idx)
					if self.noMarkerReceived {
						self.propagateSnapshot(idx)
					} else if self.shouldRecordChannelState[idx] == true {
						self.shouldRecordChannelState[idx] = false
						doneSnapshot := true
						for _, stillRecording := range self.shouldRecordChannelState {
							doneSnapshot = doneSnapshot && !(stillRecording)
						}
						if doneSnapshot {
							self.firstMarkerReceived = false
							self.finishedSnapshot = true
						}
					} else {
						fmt.Println("Trying to take a new snapshot while another already going on!")
					}
					self.canRecv.RUnlock()
					self.canProceed.Unlock()
				} else {
					self.canProceed.RLock()
					self.canRecv.RLock()
					fmt.Printf("%d Transfer %d", idx, msg)
					if self.shouldRecordChannelState[idx] {
						channelState[idx] = append(channelState[idx], msg)
					}
					self.updateBalance(msg)
					self.canRecv.RUnlock()
					self.canProceed.RUnlock()
				}
				recvdFlag = true
				break
			default:
			}
		}

		if recvdFlag == false {
			if self.nodeId == 0 {
				senderId = 1
			} else {
				senderId = 0
			}

			select {
			case msg := <- self.inChannels[senderId]:
				if msg == -1 {
					self.canProceed.Lock()
					self.canRecv.RLock()
					fmt.Printf("%d SnapshotToken -1", senderId)
					if self.noMarkerReceived {
						self.propagateSnapshot(senderId)
					} else if self.shouldRecordChannelState[senderId] == true {
						self.shouldRecordChannelState[senderId] = false
						doneSnapshot := true
						for _, stillRecording := range self.shouldRecordChannelState {
							doneSnapshot = doneSnapshot && !(stillRecording)
						}
						if doneSnapshot {
							self.firstMarkerReceived = false
							self.finishedSnapshot = true
						}
					} else {
						fmt.Println("Trying to take a new snapshot while another already going on!")
					}
					self.canRecv.RUnlock()
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
			}
		}
	}
}

func (self *Node) InitiateSnapshot() {
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
func (self *Node) PropagateSnapshot(senderId int) {
	self.nodeState = self.balance
	self.noMarkerReceived = false

	// Send marker message to all processes
	for _, channel := range self.outChannels {
		channel <- -1
	}

	// Start recording incoming channel messages
	for idx, _ := range self.shouldRecordChannelState {
		if idx != senderId {
			self.shouldRecordChannelState[idx] = true
		}
		self.channelState[senderId] = nil
	}
}

func (self *Node) updateBalance(amount int) {
	self.balance += amount
}

func (self *Node) RecvNonBlocking() {
	for idx, _ := range self.inChannels {
		select {
		case msg := <- self.inChannels[idx]:
			if msg == -1 {
				self.canProceed.Lock()
				self.canRecv.RLock()
				// TODO - remove print below
				fmt.Printf("%d SnapshotToken -1", idx)
				if self.noMarkerReceived {
					self.propagateSnapshot(idx)
				} else if self.shouldRecordChannelState[idx] == true {
					self.shouldRecordChannelState[idx] = false
					doneSnapshot := true
					for _, stillRecording := range self.shouldRecordChannelState {
						doneSnapshot = doneSnapshot && !(stillRecording)
					}
					if doneSnapshot {
						self.firstMarkerReceived = false
						self.finishedSnapshot = true
					}
				} else {
					fmt.Println("Trying to take a new snapshot while another already going on!")
				}
				self.canRecv.RUnlock()
				self.canProceed.Unlock()
			} else {
				self.canProceed.RLock()
				self.canRecv.RLock()
				// TODO Remove print below
				fmt.Printf("%d Transfer %d", idx, msg)
				if self.shouldRecordChannelState[idx] {
					channelState[idx] = append(channelState[idx], msg)
				}
				self.updateBalance(msg)
				self.canRecv.RUnlock()
				self.canProceed.RUnlock()
			}
		default:
		}
	}
}

func (self *Node) GetState() {
	state = NewState(self.nodeId, self.nodeState, self.channelState)
	return state
}