
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
  "time"
  "os"
  "./config"
  "./util"
)

const GLOBAL_ROOM = "global"
const TIME_LAYOUT = "Jan 2 2006 15.04.05 -0700 MST"

// Container for client username and connection details
type Client struct {
  // the client's connection
  Connection net.Conn
  // the client's username
  Username string
  // the current room or "global"
  Room string
  // list of usernames we are ignoring
  Ignore []string
  //
  Properties config.Properties
}

// Close the client connection and clenup
func (client *Client) Close(doSendMessage bool) {
  if (doSendMessage) {
    // if we send the close command, the connection will terminate causing another close
    // which will send the message
    sendMessage("disconnect", "", client, false, client.Properties)
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
  properties := config.Load()
  psock, err := net.Listen("tcp", ":" + properties.Port)
  util.CheckForError(err, "Can't create server")

  fmt.Printf("Chat server started on port %v...\n", properties.Port)
 
  for {
    // accept connections
    conn, err := psock.Accept()
    util.CheckForError(err, "Can't accept connections")

    // keep track of the client details
    client := Client{Connection: conn, Room: GLOBAL_ROOM, Properties: properties}
    client.Register();

    // allow non-blocking client request handling
    channel := make(chan string)
    go waitForInput(channel, &client)
    go handleInput(channel, &client, properties)

    sendMessage("ready", properties.Port, &client, true, properties)
  }
}

// wait for client input (buffered by newlines) and signal the channel
func waitForInput(out chan string, client *Client) {
  defer close(out)
 
  reader := bufio.NewReader(client.Connection)
  for {
    line, err := reader.ReadBytes('\n')
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
func handleInput(in <-chan string, client *Client, props config.Properties) {

  for {
    message := <- in
    if (message != "") {
      message = strings.TrimSpace(message)
      action, body := getAction(message)

      if (action != "") {
        switch action {

          // the user has submitted a message
          case "message":
            sendMessage("message", body, client, false, props)

          // the user has provided their username (initialization handshake)
          case "user":
            client.Username = body
            sendMessage("connect", "", client, false, props)

          // the user is disconnecting
          case "disconnect":
            client.Close(false);

          // the user is entering a room
          case "enter":
            if (body != "") {
              client.Room = body
              sendMessage("enter", body, client, false, props)
            }

          // the user is leaving the current room
          case "leave":
            if (client.Room != GLOBAL_ROOM) {
              sendMessage("leave", client.Room, client, false, props)
              client.Room = GLOBAL_ROOM
            }

          default:
            sendMessage("unrecognized", action, client, true, props)
        }
      }
    }
  }
}

// sent a message to all clients (except the sender)
func sendMessage(messageType string, message string, client *Client, thisClientOnly bool, props config.Properties) {

  if (thisClientOnly) {
    // this message is only for the provided client
    message = fmt.Sprintf("/%v", messageType);
    fmt.Fprintln(client.Connection, message)

  } else if (client.Username != "") {
    // this message is for all but the provided client
    logAction(messageType, message, client, props);

    // construct the payload to be sent to clients
    payload := fmt.Sprintf("/%v [%v] %v", messageType, client.Username, message);

    for _, _client := range clients {
      // write the message to the client
      if ((thisClientOnly && _client.Username == client.Username) ||
          (!thisClientOnly && _client.Username != "")) {

        // you should only see a message if you are in the same room
        if (messageType == "message" && client.Room != _client.Room) {
          continue;
        }

        // you won't hear any activity if you are anonymous unless thisClientOnly
        // when current client will *only* be messaged
        fmt.Fprintln(_client.Connection, payload)
      }
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

// log an action to the log file
// action: the action
//   - "enter": enter a room
//   - "leave": leave a room
//   - "connect": connect to the lobby
//   - "disconnect": disconnect from the lobby
//   - "message": post a message
//   - "ignore": ignore a user
// message: message/context appropriate for the action
// client: the initiating client
func logAction(action string, message string, client *Client, props config.Properties) {
  if (props.LogFile != "") {
    if (message == "") {
      message = "N/A"
    }
    fmt.Printf("logging values %s, %s, %s\n", action, message, client.Username);

    logMessage := fmt.Sprintf("username: %s, action: %s, value: %s, timestamp: %s, ip: %s\n",
      util.Encode(client.Username), util.Encode(action), util.Encode(message),
        util.Encode(time.Now().Format(TIME_LAYOUT)), util.Encode(client.Connection.RemoteAddr().String()))

    f, createErr := os.OpenFile(props.LogFile, os.O_RDWR|os.O_APPEND, 0666)
    util.CheckForError(createErr, "Can't open or create log file")
    defer f.Close()

    _, writeErr := f.Write([]byte(logMessage))
    util.CheckForError(writeErr, "Can't write to log file")
  }
}
