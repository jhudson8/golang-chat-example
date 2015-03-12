package endpoint

import (
  "http"
  "encoding/json"
  "../.././config"
  "../.././util"
)

func Start() {
  properties := config.Load();

  http.HandleFunc("/messages", recentMessages)
  http.ListenAndServe(":" + properties.JSONEndpointPort, nil)
}

func recentMessages(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
