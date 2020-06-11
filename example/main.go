package main

import (
	"fmt"
	"github.com/JimmyTsai16/tcpclient"
	"github.com/JimmyTsai16/tcpclient/errorcode"
	"github.com/JimmyTsai16/tcpclient/errorhandler"
	logging "github.com/z9905080/gloger"
	"log"
	"time"
)

func main() {
	logging.SetLogMode(logging.Stdout)
	logging.SetCurrentLevel(logging.DEBUG)

	client := tcpclient.New()
	err := client.Connect("127.0.0.1:1234", true)
	errorhandler.HandleError(err)
	if err != nil {
		if err != nil {
			log.Println("client close error:", err)
		}
	}

	client.ConnStateHandler = func(state tcpclient.ConnState) {
		log.Println("state changed:", state)
	}

	client.ErrorHandler = func(code errorcode.ErrorCode, err error) {
		log.Println("error occur:", err)
	}

	client.ReadHandler = func(bytes []byte, i int) {
		fmt.Println(string(bytes[:i]))
	}

	t := time.NewTicker(time.Second*3)
	for {
		<-t.C
		log.Println(client.IsClosed())
		log.Println(client.WriteByte([]byte("qqqqqqqqqwe")))
	}
}
