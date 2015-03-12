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

In other terminal windows, create as many clients as you wish
```
> go run client.go
```

You can send commands or messages.  Commands begin with ```/``` and messages are anything else.
The commands are available

* ```enter```: enter a private room (only messages from others in the same private room will be visible).  No need to explicitely create the room and you can only be in a single room at a time.
* ```leave```: leave a private room to go back to the main lobby
* ```disconnect```: disconnect from the chat server

A sample client session is below
```
> go run client.go joe
hello everyone, I am now in the lobby
/enter SomeRoom
now, this message will only be seen by others in "SomeRoom"
/leave
hi everyone, I'm back in the lobby now
/disconnect
```

Log files are in the following format each value should be url decoded when examining the files
```
username: {username}, action: {action}, value: {action specific value}, timestamp: {timestamp}, ip: {client IP address}
```

* username: example ```joe```
* action: ```message```/```enter```/```leave```/```connect```/```disconnect```
* value: the chat message or room that was entered or left
* timestamp: example ```Mar 12 2015 09.13.05 -0400 EDT```
* ip: example ```127.0.0.1%3A53594```; remember, it is url encoded so the %3A is a ":"
