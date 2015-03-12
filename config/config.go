// Shareable configuration properties module
// Load properties from config.json in the root directory
package config

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "../util"
)

type Properties struct {
  Hostname string
  Port string
  HasEnteredTheRoomMessage string
  HasLeftTheRoomMessage string
  HasEnteredTheLobbyMessage string
  HasLeftTheLobbyMessage string
  ReceivedAMessage string
  LogFile string
}

func Load() Properties {
  pwd, _ := os.Getwd()

  payload, err := ioutil.ReadFile(pwd + "/config.json")
  util.CheckForError(err, "Unable to read config file")

  var dat map[string]interface{}
  err = json.Unmarshal(payload, &dat)
  util.CheckForError(err, "Invalid JSON in config file")

  // probably a better way to unmarshall directly in the Properties struct but I haven't found it
  return Properties {
    Hostname: dat["Hostname"].(string),
    Port: dat["Port"].(string),
    HasEnteredTheRoomMessage: dat["HasEnteredTheRoomMessage"].(string),
    HasLeftTheRoomMessage: dat["HasLeftTheRoomMessage"].(string),
    HasEnteredTheLobbyMessage: dat["HasEnteredTheLobbyMessage"].(string),
    HasLeftTheLobbyMessage: dat["HasLeftTheLobbyMessage"].(string),
    ReceivedAMessage: dat["ReceivedAMessage"].(string),
    LogFile: dat["LogFile"].(string),
  }
}
