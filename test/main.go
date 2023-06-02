package main

import (
	"bufio"
	"fmt"
	"net"
)

func handlerRead(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	//var buffer [100]byte

	for {
		readString, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read: %s", readString)
		back := []byte(readString)
		conn.Write(back)
		//_, err := conn.Read(buffer[:])
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Printf("Read: %s", string(buffer[:]))
		//
		//conn.Write(buffer[:])
	}

}

func main() {
	listen, err := net.Listen("tcp", ":8888")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		go handlerRead(conn)
	}

}
