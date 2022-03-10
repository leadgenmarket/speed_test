package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"strings"
	"sync"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
	"golang.org/x/text/encoding/charmap"
)

const (
	key = "swGh889KyxjWyz"
)

func main() {
	conf, err := initConfig()
	if err != nil {
		panic(err)
	}

	changeWindowsConsoleLanguage()
	ctx := context.Background()
	user, _ := speedtest.FetchUserInfo()
	serverList, _ := speedtest.FetchServerListContext(ctx, user)
	targets, _ := serverList.FindServer([]int{})

	go func(ctx context.Context, targets speedtest.Servers, conf Config) {
		for {
			//раз в час проверяем скорость интернета
			is := checkInternetSpeed(ctx, targets, conf)
			sendSpeedToServer(is, conf.User)
			time.Sleep(time.Hour)
		}
	}(ctx, targets, conf)

	//все что связанно с ping
	writer := &customWriter{
		chann: &messagesChan{
			chann: make(chan string),
		},
	}

	go func(writer *customWriter) {
		ch := writer.chann.Subscribe()
		for {
			select {
			case msg := <-ch:
				//вычитываем пинг из канала и шлем на сервер
				sendPingToServer(msg, conf.User)
			}
		}
	}(writer)

	for {
		//раз в минуту пингуем
		cmd := exec.Command("ping", conf.Server, "-n", "1") //85.192.32.12
		cmd.Stdout = writer
		err := cmd.Run()
		if err != nil {
			writer.Error()
		}
		time.Sleep(time.Second * 60)
	}

}

type messagesChan struct {
	chann chan string
	mu    sync.RWMutex
}

func (m *messagesChan) AddMessage(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chann <- msg
}

func (m *messagesChan) Subscribe() <-chan string {
	return m.chann
}

type customWriter struct {
	chann *messagesChan
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
		c.chann.AddMessage(fmt.Sprintf("%s_ping=%s", t.Format("15:04:05"), pingS))
	}
	return len(p), nil
}

func (c customWriter) Error() {
	t := time.Now()
	c.chann.AddMessage(fmt.Sprintf("%s: ping error", t.Format("2006-01-02 15:04:05")))
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

func decodeBytes(input []byte) string {
	decoder := charmap.Windows1251.NewDecoder()
	reader := decoder.Reader(strings.NewReader(string(input)))
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func changeWindowsConsoleLanguage() {
	cmdlan := exec.Command("chcp", "437")
	cmdlan.Run()
}
