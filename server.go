
// Simple chat server which uses connection properties from "connection.json" in same directory
// Simple chat client is intended to be used but a standard telnet connection can be used
// > telnet {host} {port}
// > /user {username}
// > /chat {message}
//
// reference: https://parroty00.wordpress.com/2013/07/18/golang-tcp-server-example/
package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
  "regexp"
)

// Container for client username and connection details
type Client struct {
  Connection net.Conn
  Username string
}

// Close the client connection and clenup
func (c Client) Close() {
  c.Connection.Close();
}

// Register the client in the global cache
func (c Client) Register() {
  fmt.Printf("Adding client %v\n", c)
  numClients := len(clients)
  availableClients[numClients] = c;
  clients = availableClients[0:numClients+1]
}

// static client list
var availableClients [256]Client
var clients []Client

// program main
func main() {
  // start the server
  psock, err := net.Listen("tcp", ":5000")
  if err != nil {
    fmt.Printf("Can't create server %v\n", err)
    return
  }
 
  for {
    // accept connections
    conn, err := psock.Accept()
    if err != nil { return }

    // keep track of the client details
    client := Client{Connection: conn}
    client.Register();

    // allow non-blocking client request handling
    channel := make(chan string)
    go waitForInput(channel, client)
    go handleInput(channel, client)
  }
}

// wait for client input (buffered by newlines) and signal the channel
func waitForInput(out chan string, client Client) {
  defer close(out)
 
  for {
    line, err := bufio.NewReader(client.Connection).ReadBytes('\n')
    if err != nil {
      // connection has been closed, remove the client
      client.Close();
      return
    }
    out <- string(line)
  }
}

// listen for channel updates for a client and handle the message
// messages must be in the format of /{action} {content} where content is optional depending on the action
// supported actions are "user", "chat", and "quit".  the "user" must be set before any chat messages are allowed
func handleInput(in <-chan string, client Client) {
  for {
    message := <- in
    message = strings.TrimSpace(message)
    action, rest := getAction(message)

    if (action != "") {
      switch action {
        case "chat":
          sendMessage(rest, client)
        case "user":
          client.Username = rest
        case "quit":
          client.Close();
      }
    }
  }
}

// sent a message to all clients (except the sender)
func sendMessage(message string, client Client) {
  message = fmt.Sprintf("[%v] %v\n", client.Username, message);

  for _, client := range clients {
    // write the message to the client
    fmt.Fprintf(client.Connection, message)
  }
}

// parse out message contents (/{action} {message}) and return individual values
func getAction(message string) (string, string) {
  actionRegex, _ := regexp.Compile(`^\/([^\s]*)\s*(.*)$`)
  res := actionRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    return res[0][1], res[0][2]
  }
  return "", ""
}
