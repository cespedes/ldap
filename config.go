package main

import (
	"log"
	"os"

	ini "github.com/glacjay/goini"
)

var Config ini.Dict

func init() {
	filename := os.Getenv("HOME") + "/.ldaporg"
	var err error
	Config, err = ini.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
}
