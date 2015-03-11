
// Simple chat server which uses connection properties from "connection.json" in same directory
// Simple chat client is intended to be used but a standard telnet connection can be used
// > telnet {host} {port}
// > /user {username}
// > /message {message}
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
func (client *Client) Close(doSendMessage bool) {
  if (doSendMessage) {
    // if we send the close command, the connection will terminate causing another close
    // which will send the message
    sendMessage("leave", "", client, false)
  }
  client.Connection.Close();
  clients = removeEntry(client, clients);
}

// Register the connection and cache it
func (client *Client) Register() {
  clients = append(clients, client);
}


// static client list
var clients []*Client


// program main
func main() {
  // start the server
  psock, err := net.Listen("tcp", ":5000")
  if err != nil {
    fmt.Printf("Can't create server %v\n", err)
    return
  }
  println("Chat server started...")
 
  for {
    // accept connections
    conn, err := psock.Accept()
    if err != nil {
      fmt.Printf("Can't accept connections %v\n", err)
      return
    }

    // keep track of the client details
    client := Client{Connection: conn}
    client.Register();

    // allow non-blocking client request handling
    channel := make(chan string)
    go waitForInput(channel, &client)
    go handleInput(channel, &client)

    println("User connection")
  }
}

// wait for client input (buffered by newlines) and signal the channel
func waitForInput(out chan string, client *Client) {
  defer close(out)
 
  for {
    line, err := bufio.NewReader(client.Connection).ReadBytes('\n')
    if err != nil {
      // connection has been closed, remove the client
      client.Close(true);
      return
    }
    out <- string(line)
  }
}

// listen for channel updates for a client and handle the message
// messages must be in the format of /{action} {content} where content is optional depending on the action
// supported actions are "user", "chat", and "quit".  the "user" must be set before any chat messages are allowed
func handleInput(in <-chan string, client *Client) {
  fmt.Printf("input received \"%v\"\n", in);

  for {
    message := <- in
    message = strings.TrimSpace(message)
    action, body := getAction(message)

    if (action != "") {
      switch action {
        case "message":
          sendMessage("message", body, client, false)
        case "user":
          client.Username = body
          sendMessage("enter", "", client, false)
        case "leave":
          client.Close(false);
        default:
          sendMessage("unrecognized", action, client, true)
      }
    }
  }
}

// sent a message to all clients (except the sender)
func sendMessage(messageType string, message string, client *Client, thisClientOnly bool) {
  message = fmt.Sprintf("/%v [%v] %v", messageType, client.Username, message);
  if (thisClientOnly) {
      fmt.Printf("sending message to only %v \"%v\"\n", client.Username, message)
    } else {
      fmt.Printf("sending message to all but %v \"%v\"\n", client.Username, message)
    }
  

  for _, _client := range clients {
    // write the message to the client
    if ((thisClientOnly && _client.Username == client.Username) ||
        (!thisClientOnly && _client != client && _client.Username != "")) {
      // you won't hear any activity if you are anonymous unless thisClientOnly
      // when current client will *only* be messaged
      fmt.Fprintln(_client.Connection, message)
    }
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

// remove client entry from stored clients
func removeEntry(client *Client, arr []*Client) []*Client {
  rtn := arr
  index := -1
  for i, value := range arr {
    if (value == client) {
      index = i;
      break;
    }
  }

  if (index >= 0) {
    // we have a match, create a new array without the match
    rtn = make([]*Client, len(arr)-1)
    copy(rtn, arr[:index])
    copy(rtn[index:], arr[index+1:])
  }

  return rtn;
}
