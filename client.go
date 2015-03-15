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
  "./util"
)

// input message regular expression (look for a command /whatever)
var standardInputMessageRegex, _ = regexp.Compile(`^\/([^\s]*)\s*(.*)$`)
// chat server command /command [username] body contents
var chatServerResponseRegex, _ = regexp.Compile(`^\/([^\s]*)\s?(?:\[([^\]]*)\])?\s*(.*)$`)

// container for chat server Command details
type Command struct {
  // "leave", "message", "enter"
  Command, Username, Body string
}

// program main
func main() {
  username, properties := getConfig();

  conn, err := net.Dial("tcp", properties.Hostname + ":" + properties.Port)
  util.CheckForError(err, "Connection refused")
  defer conn.Close()

  // we're listening to chat server commands *and* user terminal commands
  go watchForConnectionInput(username, properties, conn)
  for true {
    watchForConsoleInput(conn)
  }
}

// parse out the arguments to be used when connecting to the chat server
func getConfig() (string, util.Properties) {
  if (len(os.Args) >= 2) {
    username := os.Args[1]
    properties := util.LoadConfig()
    return username, properties
  } else {
    println("You must provide the username as the first parameter ")
    os.Exit(1)
    return "", util.Properties{}
  }
}

// keep watching for console input
// send the "message" command to the chat server when we have some
func watchForConsoleInput(conn net.Conn) {
  reader := bufio.NewReader(os.Stdin)

  for true {
    message, err := reader.ReadString('\n')
    util.CheckForError(err, "Lost console connection")

    message = strings.TrimSpace(message)
    if (message != "") {
      command := parseInput(message)

      if (command.Command == "") {
        // there is no command so treat this as a simple message to be sent out
        sendCommand("message", message, conn);
      } else {
        switch command.Command {

          // enter a room
          case "enter":
            sendCommand("enter", command.Body, conn)

          // ignore someone
          case "ignore":
            sendCommand("ignore", command.Body, conn)

          // leave a room
          case "leave":
            // leave the current room (we aren't allowing multiple rooms)
            sendCommand("leave", "", conn)

          // disconnect from the chat server
          case "disconnect":
            sendCommand("disconnect", "", conn)

          default:
            fmt.Printf("Unknown command \"%s\"\n", command.Command)
        }
      }
    }
  }
}

// listen for any commands that come from the chat server
// like someone entered the room, said something, or left the room
func watchForConnectionInput(username string, properties util.Properties, conn net.Conn) {
  reader := bufio.NewReader(conn)

  for true {
    message, err := reader.ReadString('\n')
    util.CheckForError(err, "Lost server connection");
    message = strings.TrimSpace(message)
    if (message != "") {
      Command := parseCommand(message)
      switch Command.Command {

        // the handshake - send out our username
        case "ready":
          sendCommand("user", username, conn)

        // the user has connected to the chat server
        case "connect":
          fmt.Printf(properties.HasEnteredTheLobbyMessage + "\n", Command.Username)

        // the user has disconnected
        case "disconnect":
          fmt.Printf(properties.HasLeftTheLobbyMessage + "\n", Command.Username)

        // the user has entered a room
        case "enter":
          fmt.Printf(properties.HasEnteredTheRoomMessage + "\n", Command.Username, Command.Body)

        // the user has left a room
        case "leave":
          fmt.Printf(properties.HasLeftTheRoomMessage + "\n", Command.Username, Command.Body)

        // the user has sent a message
        case "message":
          if (Command.Username != username) {
            fmt.Printf(properties.ReceivedAMessage + "\n", Command.Username, Command.Body)
          }

        // the user has connected to the chat server
        case "ignoring":
          fmt.Printf(properties.IgnoringMessage + "\n", Command.Body)
      }
    }
  }
}

// send a command to the chat server
// commands are in the form of /command {command specific body content}\n
func sendCommand(command string, body string, conn net.Conn) {
  message := fmt.Sprintf("/%v %v\n", util.Encode(command), util.Encode(body));
  conn.Write([]byte(message))
}

// parse the input message and return an Command
// if there is a command the "Command" will != "", otherwise just Body will exist
func parseInput(message string) Command {
  res := standardInputMessageRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    // there is a command
    return Command {
      Command: res[0][1],
      Body: res[0][2],
    }
  } else {
    return Command {
      Body: util.Decode(message),
    }
  }
}

// look for "/Command [name] body contents" where [name] is optional
func parseCommand(message string) Command {
  res := chatServerResponseRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    // we've got a match
    return Command {
      Command: util.Decode(res[0][1]),
      Username: util.Decode(res[0][2]),
      Body: util.Decode(res[0][3]),
    }
  } else {
    // it's irritating that I can't return a nil value here - must be something I'm missing
    return Command{}
  }
}
