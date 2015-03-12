# golang-chat-example
Simple chat client and server example to help me learn golang

To run the server

Clone the repo
```
> git clone https://github.com/jhudson8/golang-chat-example.git
> cd golang-chat-example
```

Edit the server/client configuration as you need (```config.json```)
```
{
  "Port": "5555",
  "Hostname": "localhost",
  "HasEnteredTheRoomMessage": "[%s] has entered the room \"%s\"",
  "HasLeftTheRoomMessage": "[%s] has left the room \"%s\"",
  "HasEnteredTheLobbyMessage": "[%s] has entered the lobby",
  "HasLeftTheLobbyMessage": "[%s] has left the lobby",
  "ReceivedAMessage": "[%s] says: %s",
  "LogFile": ""
}
```

Start the server
```
> go run server.go
```

In another terminal window, create as many clients as you wish
```
> go run client.go {username}       // for example: "go run client.go joe"
hello, this message will be sent  // all clients connected to the lobby will receive this message
/enter SomeRoom                   // enter the private room called "SomeRoom" - only other clients in this room will see messages
this will only be in SomeRoom     // any messages when in a room will only be visible bo others in the same room
/leave SomeRoom                   // go back to the lobby
/disconnect                       // disconnect from the chat server
```
