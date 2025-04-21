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
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const localPort = ":53"

var configFilePath string

type Config struct {
	DiscordWebhookURL string `json:"discordWebhookURL"`
	RemoteAddr        string `json:"remoteAddr"`
	LocalAddr         string `json:"localAddr"`
}

var defaultConfig = Config{
	DiscordWebhookURL: "your_webhook_url",
	RemoteAddr:        "aodd.xyz:53",
	LocalAddr:         "ns1.aodd.xyz",
}

func init() {
	configFilePath = getConfigPath()
	ensureConfigDir()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-config" {
		loadConfig()
		editConfigFile()
		return
	}

	ensureRoot()
	config := loadConfig()

	// TCP 포워딩 시작
	go startTCPForwarding(config)

	// UDP 포워딩 시작
	startUDPForwarding(config)
}

func startTCPForwarding(config Config) {
	listener, err := net.Listen("tcp", localPort)
	if err != nil {
		log.Printf("TCP 포트 %s에서 리스닝 실패: %v", localPort, err)
		pauseConsole()
		return
	}
	defer listener.Close()

	log.Printf("TCP 포워딩 시작: %s -> %s", config.LocalAddr+localPort, config.RemoteAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP 연결 수락 실패: %v", err)
			continue
		}
		go handleTCPConnection(conn, config)
	}
}

func handleTCPConnection(localConn net.Conn, config Config) {
	defer localConn.Close()

	remoteConn, err := net.Dial("tcp", config.RemoteAddr)
	if err != nil {
		log.Printf("원격 TCP 서버 연결 실패: %v", err)
		return
	}
	defer remoteConn.Close()

	log.Printf("TCP 연결됨: %s -> %s -> %s | %d", localConn.RemoteAddr(), config.LocalAddr+localPort, config.RemoteAddr, time.Now().Unix())
	go sendToDiscord(fmt.Sprintf("TCP 연결됨: %s -> %s -> %s | <t:%d:F>", localConn.RemoteAddr(), config.LocalAddr+localPort, config.RemoteAddr, time.Now().Unix()))

	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
}

func startUDPForwarding(config Config) {
	localAddr, err := net.ResolveUDPAddr("udp", localPort)
	if err != nil {
		log.Fatalf("로컬 UDP 주소 오류: %v", err)
	}
	remoteAddr, err := net.ResolveUDPAddr("udp", config.RemoteAddr)
	if err != nil {
		log.Fatalf("원격 UDP 주소 오류: %v", err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalf("UDP 포트 %s 리스닝 실패: %v", localPort, err)
	}
	defer conn.Close()

	log.Printf("UDP 포워딩 시작: %s -> %s", config.LocalAddr+localPort, config.RemoteAddr)

	buf := make([]byte, 65535)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("UDP 읽기 오류: %v", err)
			continue
		}

		go func(data []byte, addr *net.UDPAddr) {
			remoteConn, err := net.DialUDP("udp", nil, remoteAddr)
			if err != nil {
				log.Printf("UDP 원격 연결 실패: %v", err)
				return
			}
			defer remoteConn.Close()

			log.Printf("UDP 연결됨: %s -> %s -> %s | %d", addr, config.LocalAddr+localPort, config.RemoteAddr, time.Now().Unix())
    		go sendToDiscord(fmt.Sprintf("UDP 연결됨: %s -> %s -> %s | <t:%d:F>", addr, config.LocalAddr+localPort, config.RemoteAddr, time.Now().Unix()))

			_, err = remoteConn.Write(data)
			if err != nil {
				log.Printf("UDP 쓰기 실패: %v", err)
				return
			}

			remoteConn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, _, err := remoteConn.ReadFrom(buf)
			if err != nil {
				log.Printf("UDP 응답 없음 또는 읽기 실패: %v", err)
				return
			}

			conn.WriteToUDP(buf[:n], addr)
		}(buf[:n], clientAddr)
	}
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

func ensureRoot() {
	if runtime.GOOS == "windows" {

	} else {
		if os.Geteuid() != 0 {
			log.Fatalln("이 프로그램은 sudo 권한으로 실행되어야 합니다.")
		}
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

func getConfigPath() string {
	if runtime.GOOS == "windows" {
		return `C:\godns\configGodns.json`
	}
	return "/etc/godns/configGodns.json"
}

func ensureConfigDir() {
	dir := filepath.Dir(configFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("설정 디렉토리 생성 실패: %v", err)
		}
	}
}

func editConfigFile() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("notepad", configFilePath)
	} else {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" //기본값
		}
		cmd = exec.Command("sudo", editor, configFilePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		log.Fatalf("설정 파일 편집기 실행 실패: %v", err)
	}
}

func pauseConsole() {
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
