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

	nodeState: int
	channelState: map[int](int)
}

func New(nodeId int, balance int) {
	node = &Node{
		nodeId:	nodeId,
		balance: balance,
		inChannels: make(map[int](chan string)),
		outChannels: make(map[int](chan string)),
		sawToken: make(map[int](bool))
	}

	// go node.start()

	return node
}


// func (self *Node) start() {
// 	for recvId, channel := range self.inChannels {
// 		select {
// 		case event := <
// 		}
// 	}
// }

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
			if self.inSnapshot {
				self.endSnapshot()
			}
			self.beginSnapshot(senderId)
		} else {
			self.updateBalance(msg)
		}

	}
}

func (self *Node) beginSnapshot(senderId int) {
	/* TODO
		record own state
		empty incoming channel from sender
		[]
	*/

	self.canProceed.Lock()
	self.snapshotState = self.balance
	for _, channel := range self.outChannels {
		channel <- -1
	}
	self.canProceed.Unlock()
}

func (self *Node) stopSnapshot(senderId int) {
	// Record channel state
}

func (self *Node) updateState(amount int) {
	self.canProceed.RLock()
	self.canRecv.RLock()
	self.balance += amount
}