./master CreateNode 1 1000
./master CreateNode 2 500
./master Send 1 2 300
./master Send 2 1 100
./master Send 1 2 400
./master Receive 2 1
./master BeginSnapshot 1
./master Receive 2
./master Receive 2
./master ReceiveAll
./master CollectState
./master PrintSnapshot
./master Send 1 2 5000
./master KillAll