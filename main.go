package main

import (
	"qok.com/identity/config"
	"qok.com/identity/router/httprouter"
)

func main() {
	httprouter.Run(config.Load().Http_port)
}
