package util

import (
  "os"
  "strings"
  "encoding/json"
  "io/ioutil"
  "net"
  "time"
  "fmt"
)

// time format for log files and JSON response
const TIME_LAYOUT = "Jan 2 2006 15.04.05 -0700 MST"
// thins we are encoding when sending stuff over the wire to clients
var ENCODING_UNENCODED_TOKENS = []string{"%", ":", "[", "]", ",", "\""}
var ENCODING_ENCODED_TOKENS = []string{"%25", "%3A", "%5B", "%5D", "%2C", "%22"}
var DECODING_UNENCODED_TOKENS = []string{":", "[", "]", ",", "\"", "%"}
var DECODING_ENCODED_TOKENS = []string{"%3A", "%5B", "%5D", "%2C", "%22", "%25"}

// Container for client username and connection details
type Client struct {
  // the client's connection
  Connection net.Conn
  // the client's username
  Username string
  // the current room or "global"
  Room string
  // list of usernames we are ignoring
  ignoring []string
  // the config properties
  Properties Properties
}
// Close the client connection and clenup
func (client *Client) Close(doSendMessage bool) {
  if (doSendMessage) {
    // if we send the close command, the connection will terminate causing another close
    // which will send the message
    SendClientMessage("disconnect", "", client, false, client.Properties)
  }
  client.Connection.Close();
  clients = removeEntry(client, clients);
}

// Register the connection and cache it
func (client *Client) Register() {
  clients = append(clients, client);
}

func (client *Client) Ignore(username string) {
  client.ignoring = append(client.ignoring, username)
}

func (client *Client) IsIgnoring(username string) bool {
  for _, value := range client.ignoring {
    if (value == username) {
      return true;
    }
  }
  return false;
}

// log content container
type Action struct {
  // "message", "leave", "enter", "connect", "disconnect"
  Command string      `json:"command"`
  // action specific content - either the chat message or room that was entered/left
  Content string      `json:"content"`
  // the username that performed the action
  Username string     `json:"username"`
  // ip address of the uwer
  IP string           `json:"ip"`
  // timestamp of the activity
  Timestamp string    `json:"timestamp"`
}

// general configuration properties
type Properties struct {
  // chat server hostname (for client connection)
  Hostname string
  // chat server port (for server execution and client connection)
  Port string
  // port used for JSON server
  JSONEndpointPort string
  // message format for when someone enters a private room
  HasEnteredTheRoomMessage string
  // message format for when someone leaves a private room
  HasLeftTheRoomMessage string
  // message format for when someone connects
  HasEnteredTheLobbyMessage string
  // message format for when someone disconnects
  HasLeftTheLobbyMessage string
  // message format for when someone sends a chat
  ReceivedAMessage string
  // message received when the user is ignoring someone else
  IgnoringMessage string
  // the absolute log file location
  LogFile string
}

// all actions (chats, enter/leave private room, connect/disconnect)
// that have occured while the server has been running
var actions = []Action{}
// cached config properties
var config = Properties{}
// static client list
var clients []*Client

// load the configuration properties from the "config.json" file
func LoadConfig() Properties {
  if (config.Port != "") {
    return config;
  }
  pwd, _ := os.Getwd()

  payload, err := ioutil.ReadFile(pwd + "/config.json")
  CheckForError(err, "Unable to read config file")

  var dat map[string]interface{}
  err = json.Unmarshal(payload, &dat)
  CheckForError(err, "Invalid JSON in config file")

  // probably a better way to unmarshall directly in the Properties struct but I haven't found it
  var rtn = Properties {
    Hostname: dat["Hostname"].(string),
    Port: dat["Port"].(string),
    JSONEndpointPort: dat["JSONEndpointPort"].(string),
    HasEnteredTheRoomMessage: dat["HasEnteredTheRoomMessage"].(string),
    HasLeftTheRoomMessage: dat["HasLeftTheRoomMessage"].(string),
    HasEnteredTheLobbyMessage: dat["HasEnteredTheLobbyMessage"].(string),
    HasLeftTheLobbyMessage: dat["HasLeftTheLobbyMessage"].(string),
    ReceivedAMessage: dat["ReceivedAMessage"].(string),
    IgnoringMessage: dat["IgnoringMessage"].(string),
    LogFile: dat["LogFile"].(string),
  }
  config = rtn;
  return rtn;
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

// sent a message to all clients (except the sender)
func SendClientMessage(messageType string, message string, client *Client, thisClientOnly bool, props Properties) {

  if (thisClientOnly) {
    // this message is only for the provided client
    message = fmt.Sprintf("/%v", messageType);
    fmt.Fprintln(client.Connection, message)

  } else if (client.Username != "") {
    // this message is for all but the provided client
    LogAction(messageType, message, client, props);

    // construct the payload to be sent to clients
    payload := fmt.Sprintf("/%v [%v] %v", messageType, client.Username, message);

    for _, _client := range clients {
      // write the message to the client
      if ((thisClientOnly && _client.Username == client.Username) ||
          (!thisClientOnly && _client.Username != "")) {

        // you should only see a message if you are in the same room
        if (messageType == "message" && client.Room != _client.Room || _client.IsIgnoring(client.Username)) {
          continue;
        }

        // you won't hear any activity if you are anonymous unless thisClientOnly
        // when current client will *only* be messaged
        fmt.Fprintln(_client.Connection, payload)
      }
    }
  }
}

// fail if an error is provided and print out the message
func CheckForError(err error, message string) {
  if err != nil {
      println(message + ": ", err.Error())
      os.Exit(1)
  }
}

// double quote the single quotes
func EncodeCSV(value string) (string) {
  return strings.Replace(value, "\"", "\"\"", -1)
}

// simple http-ish encoding to handle special characters
func Encode(value string) (string) {
  return replace(ENCODING_UNENCODED_TOKENS, ENCODING_ENCODED_TOKENS, value)
}

// simple http-ish decoding to handle special characters
func Decode(value string) (string) {
  return replace(DECODING_ENCODED_TOKENS, DECODING_UNENCODED_TOKENS, value)
}

// replace the from tokens to the to tokens (both arrays must be the same length)
func replace(fromTokens []string, toTokens []string, value string) (string) {
  for i:=0; i<len(fromTokens); i++ {
      value = strings.Replace(value, fromTokens[i], toTokens[i], -1)
  }
  return value;
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
func LogAction(action string, message string, client *Client, props Properties) {
  ip := client.Connection.RemoteAddr().String()
  timestamp := time.Now().Format(TIME_LAYOUT)

  // keep track of the actions to query against for the JSON endpoing
  actions = append(actions, Action {
    Command: action,
    Content: message,
    Username: client.Username,
    IP: ip,
    Timestamp: timestamp,
  })

  if (props.LogFile != "") {
    if (message == "") {
      message = "N/A"
    }
    fmt.Printf("logging values %s, %s, %s\n", action, message, client.Username);

    logMessage := fmt.Sprintf("\"%s\", \"%s\", \"%s\", \"%s\", \"%s\"\n",
      EncodeCSV(client.Username), EncodeCSV(action), EncodeCSV(message),
        EncodeCSV(timestamp), EncodeCSV(ip))

    f, err := os.OpenFile(props.LogFile, os.O_APPEND|os.O_WRONLY, 0600)
    if (err != nil) {
      // try to create it
      err = ioutil.WriteFile(props.LogFile, []byte{}, 0600)
      f, err = os.OpenFile(props.LogFile, os.O_APPEND|os.O_WRONLY, 0600)
      CheckForError(err, "Cant create log file")
    }

    defer f.Close()
    _, err = f.WriteString(logMessage)
    CheckForError(err, "Can't write to log file")
  }
}

func QueryMessages(actionType string, search string, username string) ([]Action) {

  isMatch := func(action Action) (bool) {
    if (actionType != "" && action.Command != actionType) {
      return false;
    }
    if (search != "" && !strings.Contains(action.Content, search)) {
      return false;
    }
    if (username != "" && action.Username != username) {
      return false;
    }
    return true;
  }

  rtn := make([]Action, 0, len(actions))

  // find out which items match the search criteria and add them to what we will be returning
  for _, value := range actions {
    if (isMatch(value)) {
      rtn = append(rtn, value)
    }
  }

  return rtn;
}
