package main

import (
	"log"
	"video-stream/internal/server"
)

func main(){
	if err := server.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}