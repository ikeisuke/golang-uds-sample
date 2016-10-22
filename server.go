package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	file := "/tmp/uds-sample-server.sock"
	// 削除はdeferで予約しておく
	// シグナルハンドラ定義しないと終了時に呼ばれない？
	defer os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go func() {
			// Closeはdeferで予約しておく
			defer conn.Close()
			data := make([]byte, 0)
			for {
				// サンプルコードなのでバッファサイズを小さめに
				buf := make([]byte, 10)
				// 全データ受信後EOFがこない限りここで止まる
				nr, err := conn.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Printf("error: %v", err)
					}
					break
				}
				// 実際に読み込んだバイト以降のデータを除去したデータに変換
				buf = buf[:nr]
				log.Printf("receive: %s\n", string(buf))
				// slice同士の結合は二つ目のsliceの後ろに...をつける
				data = append(data, buf...)
				// レスポンスデータ送信
			}
			_, err = conn.Write(data)
    	if err != nil {
    		log.Printf("error: %v\n", err)
    		return
    	}
			_, err = conn.Write(data)
    	if err != nil {
    		log.Printf("error: %v\n", err)
    		return
    	}
			log.Printf("send: %s\n", string(append(data, data...)))
		}()
	}
}
