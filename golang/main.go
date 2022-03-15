package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"strings"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

const (
	key = "swGh889KyxjWyz"
)

func main() {
	fileLog, err := os.OpenFile("./logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer fileLog.Close()
	logger := NewLogger("internet_speed", 5, fileLog)
	conf, err := initConfig()
	if err != nil {
		logger.Panic("error while initializing config: %s", err, map[string]interface{}{})
	}
	logger.Info("initialized config", map[string]interface{}{})

	changeWindowsConsoleLanguage(logger)
	logger.Info("changed terminal language", map[string]interface{}{})
	ctx := context.Background()
	user, _ := speedtest.FetchUserInfo()
	serverList, _ := speedtest.FetchServerListContext(ctx, user)
	targets, _ := serverList.FindServer([]int{})
	logger.Info("speed-test config recived", map[string]interface{}{})

	go func(ctx context.Context, targets speedtest.Servers, conf Config, logger Logger) {
		for {
			//раз в час проверяем скорость интернета
			is := checkInternetSpeed(ctx, targets, conf)
			err := sendSpeedToServer(is, conf.User)
			if err != nil {
				logger.Error("error sending internet speed to server: %s", err, map[string]interface{}{})
			}
			time.Sleep(time.Hour)
		}
	}(ctx, targets, conf, logger)

	//все что связанно с ping
	writer := &customWriter{
		config: &conf,
		logger: logger,
	}

	//первый пинг шлем через три минуты
	for {
		//раз в минуту пингуем
		cmd := exec.Command("ping", conf.Server, "-n", "1") //85.192.32.12
		cmd.Stdout = writer
		err := cmd.Run()
		if err != nil {
			writer.Error()
			logger.Error("error pinging: %s", err, map[string]interface{}{})
		}
		cmd.Process.Release()
		time.Sleep(time.Second * 60)
	}

}

type customWriter struct {
	config *Config
	logger Logger
}

func (c customWriter) Write(p []byte) (int, error) {
	str := string(p)
	if string(str[0]) == "b" {
		i := strings.Index(str, "time=") + 5 //time= для macos
		pingS := ""
		for string(str[i]) != "m" {
			pingS += string(str[i])
			i++
		}
		t := time.Now()
		err := sendPingToServer(fmt.Sprintf("%s_ping=%s", t.Format("15:04:05"), pingS), c.config.User)
		if err != nil {
			c.logger.Error("error sending ping to server: %s", err, map[string]interface{}{})
		}
	}
	return len(p), nil
}

func (c customWriter) Error() {
	t := time.Now()
	sendPingToServer(fmt.Sprintf("%s: ping error", t.Format("2006-01-02 15:04:05")), c.config.User)
}

func checkInternetSpeed(ctx context.Context, targets speedtest.Servers, conf Config) InternetSpeed {
	s := targets[0]
	s.PingTest()
	s.DownloadTest(false)
	s.UploadTest(false)

	speed := InternetSpeed{
		Latency:  s.Latency.String(),
		Upload:   s.ULSpeed,
		Download: s.DLSpeed,
	}
	return speed
}

type InternetSpeed struct {
	Latency  string
	Upload   float64
	Download float64
}

type Config struct {
	Loss        int
	InputSpeed  int
	OutputSpeed int
	Ping        int
	Server      string
	User        string
}

type FileSettting struct {
	User string `json:"user"`
}

func initConfig() (Config, error) {
	conf := Config{}
	resp, err := http.Get("http://mail.leadactiv.ru/settings.json")
	if err != nil {
		return conf, errors.New("error geting configuration from server")
	}

	bytesString, err := io.ReadAll(resp.Body)
	if err != nil {
		return conf, fmt.Errorf("error decoding configuration from server %s", err.Error())
	}
	err = json.Unmarshal(bytesString, &conf)
	if err != nil {
		return conf, fmt.Errorf("error unmarshaling configuration from server %s", err.Error())
	}

	file, _ := os.Open("settings.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	fileConf := FileSettting{}
	err = decoder.Decode(&fileConf)
	if err != nil {
		fmt.Println("error:", err)
	}
	conf.User = fileConf.User
	return conf, nil
}

func sendPingToServer(msg string, user string) error {
	_, err := http.Get(fmt.Sprintf("http://mail.leadactiv.ru/newApi.php?key=%s&type=ping&msg=%s&user=%s", key, msg, user))
	if err != nil {
		return errors.New("error sending ping to server")
	}
	return nil
}

func sendSpeedToServer(is InternetSpeed, user string) error {
	b, err := json.Marshal(is)
	if err != nil {
		return err
	}
	_, err = http.Get(fmt.Sprintf("http://mail.leadactiv.ru/newApi.php?key=%s&type=speed&msg=%s&user=%s", key, string(b), user))
	if err != nil {
		return errors.New("error sending ping to server")
	}
	return nil
}

func changeWindowsConsoleLanguage(logger Logger) {
	cmdlan := exec.Command("chcp", "437")
	err := cmdlan.Run()
	if err != nil {
		logger.Error("termainal language change error: %s", err, map[string]interface{}{})
	}
}
