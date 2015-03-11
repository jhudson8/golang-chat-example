package config

type Properties struct {
  Hostname string
  Port string
}

func Load() Properties {
  return Properties {
    Hostname: "localhost",
    Port: "5555",
  }
}
