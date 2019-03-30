package server

import (
    "os"

    "encoding/json"
)

// Config is configuration
type Config struct {
    SMTP struct { 
        Server      string  `json:"smtp_server"`
        Port        int     `json:"smtp_port"`
        Email       string  `json:"smtp_email"`
        Password    string  `json:"smtp_password"`
    } `json:"smtp"`

    Server struct {
        BindAddr        string  `json:"bind_addr"`
    } `json:"server"`
}

// GetConfig loads config stuff
func GetConfig() (config *Config, err error) {
    // Open a config file, or return error if it doesn't exist.
    file, err := os.Open("config.json")
    if err != nil {
        return nil, err
    }

    config = &Config{}
    err = json.NewDecoder(file).Decode(config)
    if err != nil {
        return nil, err
    }

    return config, nil
}