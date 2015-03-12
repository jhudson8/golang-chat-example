// A simple chat client to talk to the simple chat server (./server.go)
// To run the client, use "go run client.go {username}" where username is your username
// This will listen for chat room events and display them
// (someone entered the room, left the room, or chatted something)
// To chat a message, simply type something after you run the program and press the enter key
//
// reference https://gist.github.com/iwanbk/2295233
//           http://golang.org/pkg/net/
package main

import (
  "fmt"
  "os"
  "net"
  "bufio"
  "regexp"
  "strings"
  "./config"
)

// container for chat server action details
type Action struct {
  // "leave", "message", "enter"
  ActionType string
  Username string
  Body string
}

// program main
func main() {
  username, properties := getConfig();

  conn, err := net.Dial("tcp", properties.Hostname + ":" + properties.Port)
  checkForError(err, "Connection refused")
  defer conn.Close()

  // we're listening to chat server commands *and* user terminal commands
  go watchForConnectionInput(username, properties, conn)
  for true {
    watchForConsoleInput(conn)
  }
}

// parse out the arguments to be used when connecting to the chat server
func getConfig() (string, config.Properties) {
  if (len(os.Args) >= 2) {
    username := os.Args[1]
    properties := config.Load()
    return username, properties
  } else {
    println("You must provide the username as the first parameter ")
    os.Exit(1)
    return "", config.Properties{}
  }
}

// fail if an error is provided and print out the message
func checkForError(err error, message string) {
  if err != nil {
      println(message + ": ", err.Error())
      os.Exit(1)
  }
}

// keep watching for console input
// send the "message" command to the chat server when we have some
func watchForConsoleInput(conn net.Conn) {
  reader := bufio.NewReader(os.Stdin)

  for true {
    message, err := reader.ReadString('\n')
    checkForError(err, "Lost console connection")

    message = strings.TrimSpace(message)
    if (message != "") {
      sendCommand("message", message, conn);
    }
  }
}

// listen for any commands that come from the chat server
// like someone entered the room, said something, or left the room
func watchForConnectionInput(username string, properties config.Properties, conn net.Conn) {
  reader := bufio.NewReader(conn)

  for true {
    message, err := reader.ReadString('\n')
    checkForError(err, "Lost server connection");
    message = strings.TrimSpace(message)
    if (message != "") {
      action := parseAction(message)
      switch action.ActionType {
        case "ready":
          // the handshake - send out our username
          sendCommand("user", username, conn)
          fmt.Printf(properties.HasEnteredTheRoomMessage + "\n", username)
        case "enter":
          fmt.Printf(properties.HasEnteredTheRoomMessage + "\n", action.Username)
        case "leave":
          fmt.Printf(properties.HasLeftTheRoomMessage + "\n", action.Username)
        case "message":
          fmt.Printf(properties.ReceivedAMessage + "\n" + "\n", action.Username, action.Body)
      }
    }
  }
}

// send a command to the chat server
// commands are in the form of /action {command specific body content}\n
func sendCommand(action string, body string, conn net.Conn) {
  message := fmt.Sprintf("/%v %v\n", action, body);
  conn.Write([]byte(message))
}

// look for "/action [name] body contents" where [name] is optional
func parseAction(message string) Action {
  actionRegex, _ := regexp.Compile(`^\/([^\s]*)\s?(?:\[([^\]]*)\])?\s*(.*)$`)
  res := actionRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    // we've got a match
    return Action{
      ActionType: res[0][1],
      Username: res[0][2],
      Body: res[0][3],
    }
  } else {
    // it's irritating that I can't return a nil value here - must be something I'm missing
    return Action{}
  }
}
