package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// 원격 서버 주소 및 포트
const remoteAddr = "aodd.xyz:7445"

// 로컬 포트
const localAddr = ":445"

func handleConnection(localConn net.Conn) {
	defer localConn.Close()

	// 원격 서버에 연결
	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("원격 서버 연결 실패: %v", err)
		return
	}
	defer remoteConn.Close()

	log.Printf("연결됨: %s -> %s", localConn.RemoteAddr(), remoteAddr)

	// 데이터를 양방향으로 복사 (로컬 <-> 원격)
	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}

func main() {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("포트 %s에서 리스닝 실패: %v", localAddr, err)
		pauseConsole()
		return
	}
	defer listener.Close()

	log.Printf("포트 포워딩 시작: %s -> %s", localAddr, remoteAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("연결 수락 실패: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func pauseConsole() {
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}