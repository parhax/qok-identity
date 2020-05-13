package logwrapper

import (
	"log"
	"os"
)

//Load retrun a signleton of the logger
func Load() *log.Logger {
	f, err := os.OpenFile("identity.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return log.New(f, "Identity :: ", log.LstdFlags)

}
