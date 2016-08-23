package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func initData(frameindex uint8) []byte {
	buf := new(bytes.Buffer)
	t := uint32(time.Now().Unix())

	binary.Write(buf, binary.LittleEndian, uint8(0x88)) //head
	binary.Write(buf, binary.LittleEndian, uint8(0x07)) //type
	s := "I am a string data"
	var l uint16 //left data len
	l = uint16(1 + 4 + 8 + 4 + 10 + len([]byte(s)))
	binary.Write(buf, binary.LittleEndian, l) //data len
	var data = []interface{}{
		uint8(0x01),                           // 1 BYTE
		uint32(t),                             // 4 BYTE
		float64(1.11111111),                   // 8 BTTE
		float32(3.33333333),                   // 4 BYTE
		[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // 10 BYTE
		[]byte(s), // len([]byte(s))
	}
	for _, v := range data {
		binary.Write(buf, binary.LittleEndian, v) //data filed
	}
	fmt.Println("datalen:", l, "all:", len(buf.Bytes()), buf.Bytes())
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host:port")
		os.Exit(1)
	}
	service := os.Args[1]
	conn, err := net.Dial("tcp", service)
	checkError(err)
	for n := 0; n < 1; n++ {
		w := initData(uint8(n + 1))
		fmt.Println("write:", len(w), w)
		conn.Write(w)
		//conn.Write([]byte(fmt.Sprintf("hello %d", n+48)))
		var buf [1024]byte
		n, err := conn.Read(buf[0:])
		checkError(err)
		fmt.Println("receive:", len(buf), buf[0:n])
	}
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
