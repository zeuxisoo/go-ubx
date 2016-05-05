package main

import (
    "errors"
)

type Config struct {
    EventId string
}

func NewConfig() (c *Config) {
    return &Config{}
}

func (c *Config) Check() (string, error) {
    if c.EventId == "" {
        return "", errors.New("Please enter event id")
    }

    return "", nil
}
