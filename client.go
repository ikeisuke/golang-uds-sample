package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	if len(os.Args) != 3 {
		return
	}
	file := os.Args[1]
	message := os.Args[2]
	conn, err := net.Dial("unix", file)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	log.Printf("send: %s\n", message)
	err = conn.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	data := make([]byte, 0)
	for {
		// サンプルコードなのでバッファサイズを小さめに
		buf := make([]byte, 10)
		nr, err := conn.Read(buf)
		if err != nil {
			break
		}
		buf = buf[:nr]
		// 受信データのログ出力
		log.Printf("receive: %s\n", string(buf))
		data = append(data, buf...)
	}
	// 受信データの出力
	fmt.Printf("%s\n", string(data))
}
