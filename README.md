# golang-chat-example
Simple chat client and server example to help me learn golang


Chat Server
-----------

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


Chat Client
-----------
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

JSON Endpoint
----------
The JSON endpoint port can be configured using the ```JSONEndpointPort``` port (by default, 8080).  When the chat server is stated, the following endpoints are available

* ```/messages/all```: all messages
* ```/messages/search/{search term}```: example ```localhost:8080/messages/search/hello```
* ```/messages/user/{username}```: example ```localhost:8080/messages/user/joe```

The message query will only use the messages from the running server (previously logged messages will not be evaluated).


Chat Log
----------
Log files are in CSV format with the columns shown below.  You *must* set the ```LogFile``` config value to be the absolute file location or no logs will be created.

1. ***username***: the user that performed the action
2. ***action***: the action that was taken (```message```/```enter```/```leave```/```connect```/```disconnect```)
3. ***value***: the chat message or room that was entered or left
4. ***timestamp***: example ```Mar 12 2015 09.13.05 -0400 EDT```
5. ***ip***: example ```127.0.0.1:53594```
