package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const localPort = ":445"                  //로컬 포트
const configFilePath = "configGosmb.json" // 설정 파일 경로

// 설정 구조체
type Config struct {
	DiscordWebhookURL string `json:"discordWebhookURL"`
	RemoteAddr        string `json:"remoteAddr"`
	LocalAddr         string `json:"localAddr"`
}

// 기본값
var defaultConfig = Config{
	DiscordWebhookURL: "your_webhook_url",
	RemoteAddr:        "aodd.xyz:7445",
	LocalAddr:         "ns2.aodd.xyz",
}

func sendToDiscord(message string) {
	config := loadConfig()
	payload := map[string]string{"content": message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Discord 메시지 전송 실패: %v", err)
		return
	}

	_, err = http.Post(config.DiscordWebhookURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		log.Printf("Discord 웹훅 요청 실패: %v", err)
	}
}

func loadConfig() Config {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Printf("설정 파일 없음. 기본 설정 파일을 생성합니다: %s", configFilePath)
		file, err := os.Create(configFilePath)
		if err != nil {
			log.Fatalf("설정 파일 생성 실패: %v", err)
		}
		defer file.Close()

		configData, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			log.Fatalf("설정 파일 데이터 생성 실패: %v", err)
		}
		file.Write(configData)
		return defaultConfig
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("설정 파일 열기 실패: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("설정 파일 읽기 실패: %v", err)
	}
	return config
}

func handleConnection(localConn net.Conn, config Config) {
	defer localConn.Close()

	remoteConn, err := net.Dial("tcp", config.RemoteAddr)
	if err != nil {
		log.Printf("원격 서버 연결 실패: %v", err)
		return
	}
	defer remoteConn.Close()

	log.Printf("연결됨: %s -> %s -> %s", localConn.RemoteAddr(), config.LocalAddr+localPort, config.RemoteAddr)
	sendToDiscord(fmt.Sprintf("연결됨: %s -> %s -> %s | <t:%d:F>", localConn.RemoteAddr(), config.LocalAddr+localPort, config.RemoteAddr, time.Now().Unix()))

	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}

func main() {
	config := loadConfig()

	listener, err := net.Listen("tcp", localPort)
	if err != nil {
		log.Fatalf("포트 %s에서 리스닝 실패: %v", localPort, err)
		pauseConsole()
		return
	}
	defer listener.Close()

	log.Printf("포트 포워딩 시작: %s -> %s", config.LocalAddr+localPort, config.RemoteAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("연결 수락 실패: %v", err)
			continue
		}
		go handleConnection(conn, config)
	}
}

func pauseConsole() {
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
