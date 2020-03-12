# shutter

## Protocol Description
The node chosen to initiate the Snapshot records its node state and sends snapshot token (-1 in our implementation) to all of its outgoing channels. The node starts recording messages from all incoming channels.

Nodes which receive a snapshot token for the first time, record their node state, set the channel state for their incoming channel from the sender as empty, send a snapshot token to all of their outgoing channels and start recording the channel states for all their incoming channels except for the one coming from the sender of their first snapshot token. This all happens atomically, so that no other application-based state change can happen in middle.

Nodes which receive a snapshot token after their first time stop recording the channel state at the incoming channel of the token. Once all channels are stopped from recording, the node is ready to send its partial snapshot results to the observer.

The observer collects the node and channel states from each node.

## Implementation
We implement the cluster in Golang. We consider the master to be split into two processes for a cleaner implementation. The first process (called request-master for now onward) serves as an interface to the master. It parses the commands file and issues http request to another process (cluster-master). Cluster-master stores an array of references to each node object in the cluster. Running in the same process, as a go-routine (Go-runtime managed thread), there is also an observer. When a method of a node object is called, we spawn a new go-routine to handle the call (which provides some similarity to an actual distributed system)

Communication:

1. The request-master and cluster-master communicate via http requests for the sake of simplicity. The request-master itself is modelled as a cli.

2. The cluster-master and observer communicate using Golang's message passing primitive between go-routines, called channels. The channel is FIFO ordered by design.

3. Communication between the cluster-master and nodes (Send, Receive, ReceiveAll triggers) and between observer and nodes (CollectSnapshot, PrintSnapshot) happens directly by calling a method of the specified node.

4. Communication between nodes (Sending amounts, tokens) happens via go channels. Each node has an incoming and outgoing channel associated with every other node. Only when receive is called, a message is taken out from the head of the channel

Our node abstraction keeps track of its state using state variables like shouldRecordChannelState, noMarkerReceived, and finishedSnapshot.

Go channels provide the abstraction of FIFO channels used in Chandy-Lamport algorithm.

## Instructions
The code requires setup of Golang.
Once Golang is installed, run the ./build.sh to build Go binary. This builds binaries specific to multiple platforms (MacOS, Linux and Windows). Depending on the platform you are working on, run binary with it's spec.

For example, if we ran on MacOS,
In one terminal, start the cluster:  `./cluster@darwin-amd64`
In another terminal, invoke the master to run parse test file: `./master@darwin-amd64 ParseFile master/examples/run_1.txt`

## Example Test Cases
We tried the provided test cases, and also tried creating more nodes after some amount of Sends and Receives have happened. Moreover, we have attached a video showing how our Final E2E run works for reference.
Example:
```
CreateNode 1 1000
CreateNode 2 500
CreateNode 3 300
Send 1 2 300
Send 2 1 100
CreateNode 4 600
Send 1 3 100
Send 1 2 400
BeginSnapshot 1
Send 4 2 100
Receive 2 1
Send 4 1 700
Receive 3
Send 2 4 50
ReceiveAll
CollectState
PrintSnapshot
KillAll
```
```
Started by Node 1
1 Transfer 300
ERR_SEND
1 Transfer 100
---Node states
node 1 = 200
node 2 = 1050
node 3 = 400
node 4 = 500
---Channel states
channel (1 -> 2) = 0
channel (1 -> 3) = 0
channel (1 -> 4) = 0
channel (2 -> 1) = 100
channel (2 -> 3) = 0
channel (2 -> 4) = 50
channel (3 -> 1) = 0
channel (3 -> 2) = 0
channel (3 -> 4) = 0
channel (4 -> 1) = 0
channel (4 -> 2) = 100
channel (4 -> 3) = 0
```

## Group Members
Aashaka Shah, UTEID: as88752, UTCS ID: aashaka

Kevin Sijo Puthusseri, UTEID: ksp2236 , UTCS ID: pskevin
