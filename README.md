# golang-chat-example
Simple example to help me learn golang

A chat client (TODO) is intended to be used with the server, but a standard telnet connection can be used

Right now, the server will run on port ```5000``` but it will eventually accept port details from a config file.

To run the server

```
> git clone https://github.com/jhudson8/golang-chat-example.git
> cd golang-chat-example
> go run server.go
```

And, until the client has been created, in another terminal window

```
> telnet localhost 5000
> /user joe
> /message hello
> /leave
```

Note: you won't get any messages unless you connect with multiple clients (messages won't be echoed to sender).  You won't hear any activity if you are anonymous so you muse use ```/user {username}```.

Client messages are commands just like server messages.  You will see

When someone enters the chat server (and sets their username)
```
/enter [username]
```

When someone enters a message
```
/message [username] the message
```

When someone leaves the chat server
```
/leave [username]
```

When an unrecognized command is entered by the client
```
/unrecognized [username] commandName
```