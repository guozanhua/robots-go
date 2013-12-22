package main

import (
	"crypto/aes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

var host string = "10.0.253.60"
var port string = "9118"
var aesKey string = "P%2BViyZLtO^gRT2Huxqx#5Vygbfl$8m"

func main() {
	fmt.Println("Robos-go starting ...")

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Connect to game server failed : ", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connects to game server success.")

	bytes := []byte{1, 0, 0, 0}
	_, err2 := conn.Write(bytes)
	if err2 != nil {
		log.Fatal("Send login packet failed : ", err2)
		return
	}

	resp := make([]byte, 10240)
	nread, err3 := conn.Read(resp)
	if err3 != nil {
		log.Fatal("Read login response failed : ", err3)
		return
	}
	fmt.Println("Client received :")
	fmt.Println(hex.Dump(resp[:nread]))

	dataLen := binary.LittleEndian.Uint16(resp[:2])
	encrypted := (resp[2] != 0)
	fmt.Printf("Packet Len : %v, Data Len : %v, is encrypted : %v\n", nread, dataLen, encrypted)

	c, err4 := aes.NewCipher([]byte(aesKey))
	if err4 != nil {
		panic(err4)
	}

	rawBytes := make([]byte, dataLen-4)
	c.Decrypt(rawBytes, resp[4:])
	fmt.Println("Decrypted data :")
	fmt.Println(hex.Dump(rawBytes))

	msgType := binary.LittleEndian.Uint32(rawBytes[:4])
	state := rawBytes[4]
	key := binary.LittleEndian.Uint64(rawBytes[5:])
	fmt.Printf("msgType = %v, state = %v, key = %v\n", msgType, state, key)

	fmt.Println("Robos-go exiting ...")
}
