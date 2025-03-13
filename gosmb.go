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

const remoteAddr = "aodd.xyz:7445"        //원격 서버 주소 및 포트
const localAddr = ":445"                  //로컬 포트
const configFilePath = "configGosmb.json" //설정 파일 경로

// 설정 구조체
type Config struct {
	DiscordWebhookURL string `json:"discordWebhookURL"`
}

// 기본값
var defaultConfig = Config{
	DiscordWebhookURL: "https://discord.com/api/webhooks/your-webhook-url",
}

func sendToDiscord(message string) {
	config := loadConfig() //Discord 웹훅 URL을 설정 파일에서 읽어옴
	// Discord 웹훅으로 메시지 전송
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
	// 설정 파일이 있으면 읽고, 없으면 기본값으로 생성
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// 설정 파일이 없으면 기본값으로 설정 파일 생성
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

	file, err := os.Open(configFilePath) //설정 파일 읽기
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
	sendToDiscord(fmt.Sprintf("연결됨: %s -> %s | <t:%d:F>", localConn.RemoteAddr(), remoteAddr, time.Now().Unix()))

	// 데이터를 양방향으로 복사 (로컬 <-> 원격)
	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}

func main() {
	loadConfig()

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
	// fmt.Scanln()을 사용하여 Enter를 눌러 종료
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}