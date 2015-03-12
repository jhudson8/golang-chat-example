package json

import (
  "net/http"
  "encoding/json"
  "../.././util"
)

const SEARCH_PATH = "/messages/search/"
const USER_PATH = "/messages/user/"
const ALL_PATH = "/messages/all"

func Start() {
  properties := util.LoadConfig();

  http.HandleFunc(SEARCH_PATH, searchMessages)
  http.HandleFunc(USER_PATH, userMessages)
  http.HandleFunc(ALL_PATH, allMessages)

  err := http.ListenAndServe(":" + properties.JSONEndpointPort, nil)
  util.CheckForError(err, "Can't create JSON endpoint")
}

func searchMessages(w http.ResponseWriter, r *http.Request) {
  var searchTerm = r.URL.Path[len(SEARCH_PATH):]

  returnQuery("message", searchTerm, "", w, r)
}

func userMessages(w http.ResponseWriter, r *http.Request) {
  var username = r.URL.Path[len(USER_PATH):]

  returnQuery("message", "", username, w, r)
}

func allMessages(w http.ResponseWriter, r *http.Request) {
  returnQuery("message", "", "", w, r)
}

func returnQuery(actionType string, search string, username string,
    w http.ResponseWriter, r *http.Request) {

  actions := util.QueryMessages(actionType, search, username);
  payload, err := json.Marshal(actions)
  util.CheckForError(err, "Can't create JSON response")

  w.Header().Set("Content-Type", "text/json")
  w.Write(payload);
}