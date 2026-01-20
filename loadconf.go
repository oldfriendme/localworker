package main

import (
    "encoding/json"
	"os"
)

type SSLConfig struct {
    Crt string `json:"crt"`
    Key string `json:"key"`
}

type Config struct {
    Listen    string    `json:"listen"`
    EnableSSL bool      `json:"enable_ssl"`
    SSLConfig SSLConfig `json:"ssl-config"`
    WorkDir   string    `json:"workdir"`
    AppFile   string    `json:"appfile"`
    User      string    `json:"user"`
    Pass      string    `json:"pass"`
}

func checkconf(f string) (Config,error) {
var c Config
d, err := os.ReadFile(f)
	if err != nil {
return c,err
}
	if err := json.Unmarshal(d, &c); err != nil {
return c,err
}
return c,nil
}