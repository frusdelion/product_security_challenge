package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

func main() {
	logRusInstance := logrus.New()
	srv := NewServer(logRusInstance)
	srv.Run()
}
