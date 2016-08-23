package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":7777")
	checkError(err)
	fmt.Println("Start Listening.....")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("Accepted %#p\n", conn)
		go handleClient(conn)
	}
}

type SingleData struct {
	Head    uint8
	Type    uint8
	Len     uint16
	CmdId   uint8
	NTime   uint32
	Float64 float64
	FLoat32 float32
	Byte10  [10]uint8
	//stringField []byte
}

func processData(b []byte) {
	var s SingleData
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &s)
	checkError(err)
	fmt.Printf("%s\n%#v\n", err, s)
	konwLen := uint16(1 + 4 + 8 + 4 + 10)
	leftLen := s.Len - konwLen
	if leftLen > 0 {
		fmt.Println("lef data len:", leftLen)
		sLeft := make([]byte, leftLen)
		binary.Read(buf, binary.LittleEndian, &sLeft)
		fmt.Println(string(sLeft))
	}
}

func processDataOrgin(b []byte) {
	var s SingleData
	buf := bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &s.Head)
	binary.Read(buf, binary.LittleEndian, &s.Type)
	binary.Read(buf, binary.LittleEndian, &s.Len)
	binary.Read(buf, binary.LittleEndian, &s.CmdId)
	binary.Read(buf, binary.LittleEndian, &s.NTime)
	binary.Read(buf, binary.LittleEndian, &s.Float64)
	binary.Read(buf, binary.LittleEndian, &s.FLoat32)
	binary.Read(buf, binary.LittleEndian, &s.Byte10)

	fmt.Printf("%#v\n", s)
	konwLen := uint16(1 + 4 + 8 + 4 + 10)
	leftLen := s.Len - konwLen
	if leftLen > 0 {
		fmt.Println("lef data len:", leftLen)
		sLeft := make([]byte, leftLen)
		binary.Read(buf, binary.LittleEndian, &sLeft)
		fmt.Println(string(sLeft))
	}

}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		head := make([]byte, 1)
		bType := make([]byte, 1)
		bLen := make([]byte, 2)
		n, err := conn.Read(head)
		fmt.Printf("head: %d, %#x, %s\n", n, head, err)
		n, err = conn.Read(bType)
		fmt.Printf("type: %d, %#x, %s\n", n, bType, err)
		n, err = conn.Read(bLen)
		l := binary.LittleEndian.Uint16(bLen)
		fmt.Printf("data len: %d, %d, %d,  %s\n", n, l, binary.BigEndian.Uint16(bLen), err)
		buf := make([]byte, l)
		n, err = conn.Read(buf[0:])
		if err == io.EOF {
			fmt.Println("Trying to read")
			return
		} else {
			fmt.Println("read from client:", buf[0:n], "\nread n=", n, "must be:", l)
		}
		if err != nil {
			fmt.Println(err)
		}
		all := append(head, bType...)
		all = append(all, bLen...)
		all = append(all, buf...)
		processDataOrgin(all)
		processData(all)
		_, err = conn.Write(all)
		if err != nil {
			fmt.Println("write err", err)
			return
		}
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
