# golang-chat-example
Simple example to help me learn golang

A chat client (TODO) is intended to be used with the server, but a standard telnet connection can be used

Right now, the server will run on port ```5000``` but it will eventually accept port details from a config file.

To run the server

```
> git clone https://github.com/jhudson8/golang-chat-example.git
> cd golang-chat-example
> go 
```

And, until the client has been created, in another terminal window

```
> telnet localhost 5000
> /user joe
> /chat hello
> /quit
```
