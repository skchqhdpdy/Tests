package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var configFilePath string // OS에 따라 동적으로 설정

type Config struct {
	LocalAddr         string `json:"localAddr"`
	LocalPort         string `json:"localPort"`
	RemoteAddr        string `json:"remoteAddr"`
	RemotePort        string `json:"remotePort"`
}

var defaultConfig = Config{
	LocalAddr:         "127.0.0.1",
	LocalPort:         "6969",
	RemoteAddr:        "spt-eft.aodd.xyz",
	RemotePort:        "6969",
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

	listener, err := net.Listen("tcp", config.LocalAddr+":"+config.LocalPort)
	if err != nil {
		log.Printf("포트 %s에서 리스닝 실패: %v", config.LocalPort, err)
		pauseConsole()
		return
	}
	defer listener.Close()

	log.Printf("포트 포워딩 시작: %s -> %s", config.LocalAddr+":"+config.LocalPort, config.RemoteAddr+":"+config.RemotePort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("연결 수락 실패: %v", err)
			continue
		}
		go handleConnection(conn, config)
	}
}

func handleConnection(localConn net.Conn, config Config) {
	defer localConn.Close()

	remoteConn, err := net.Dial("tcp", config.RemoteAddr+":"+config.RemotePort)
	if err != nil {
		log.Printf("원격 서버 연결 실패: %v", err)
		return
	}
	defer remoteConn.Close()

	log.Printf("연결됨: %s -> %s -> %s | %d", localConn.RemoteAddr(), config.LocalAddr+":"+config.LocalPort, config.RemoteAddr+":"+config.RemotePort, time.Now().Unix())

	go io.Copy(remoteConn, localConn)
	io.Copy(localConn, remoteConn)
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
		return `C:\gospt\configSPT.json`
	}
	return "/etc/gospt/configSPT.json"
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
