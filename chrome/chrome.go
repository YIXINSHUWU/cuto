package chrome

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/diiyw/chr/protocol/dom"
	"github.com/diiyw/chr/protocol/page"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type Chrome struct {
	bin        string
	remoteAddr string
	dataDir    string
	process    *os.Process
	timeout    time.Duration
}

// Create chrome client
func Create(opts ...Option) (*Chrome, error) {

	c := new(Chrome)
	for _, opt := range opts {
		opt(c)
	}

	if c.bin == "" {
		for _, filename := range defaultChrome {
			_, err := os.Stat(filename)
			if err == nil {
				c.bin = filename
				break
			}
		}
	}

	if c.remoteAddr == "" {
		c.remoteAddr = "127.0.0.1:9222"
	}

	if c.dataDir == "" {
		c.dataDir = defaultUserDataTmpDir
	}

	if c.timeout == 0 {
		c.timeout = 5 * time.Second
	}

	_, err := os.Stat(c.bin)
	if err != nil {
		return nil, errors.New("Chrome not found ")
	}

	cmd := &exec.Cmd{}
	cmd = exec.Command(c.bin, []string{
		"--remote-debugging-port=" + strings.Split(c.remoteAddr, ":")[1],
		"--user-data-dir=" + c.dataDir,
	}...)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Start chrome with error: %s ", err)
	}

	go func() {
		c.process = cmd.Process
		_ = cmd.Wait()
	}()

	http.DefaultClient.Timeout = c.timeout
	resp, err := http.Get("http://127.0.0.1:9222")
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	var s = make(chan os.Signal)
	go func() {
		signal.Notify(s, os.Interrupt, os.Kill)
		<-s
		if err := c.Close(); err != nil {
			log.Println(err)
		}
	}()
	return c, nil
}

// Open new tab
func (chrome *Chrome) Open(url string) (*Tab, error) {
	r, err := http.Get("http://" + chrome.remoteAddr + "/json/new?" + url)
	if err != nil {
		return nil, errors.New("Http request error:" + err.Error())
	}
	defer r.Body.Close()
	tab := &Tab{}
	if err := json.NewDecoder(r.Body).Decode(&tab); err != nil {
		return nil, errors.New("Open new window error:" + err.Error())
	}
	return newTab(tab)
}

// Close chrome
func (chrome *Chrome) Close() error {
	_ = os.RemoveAll(chrome.dataDir)
	return chrome.process.Kill()
}

func newTab(tab *Tab) (*Tab, error) {
	conn, _, err := websocket.DefaultDialer.Dial(tab.WebSocketDebuggerUrl, nil)
	if err != nil {
		return nil, err
	}
	tab.Ipc.Conn = conn
	tab.Ipc.events = make(chan []byte, 1024)
	tab.Ipc.returns = make(chan []byte, 1024)
	tab.Ipc.errors = make(chan []byte, 1024)
	_ = tab.Send(dom.Enable, dom.EnableParams{})
	_ = tab.Send(page.Enable, page.EnableParams{})
	go func() {
		if err := tab.handle(); err != nil {
			log.Println(err)
			return
		}
	}()
	return tab, nil
}

// Catch tab
func (chrome *Chrome) Find(kw string) (*Tab, error) {
	r, err := http.Get("http://" + chrome.remoteAddr + "/json")
	if err != nil {
		return nil, errors.New("Http request error:" + err.Error())
	}
	defer r.Body.Close()
	var tabs = make([]*Tab, 0, 16)
	if err = json.NewDecoder(r.Body).Decode(&tabs); err != nil {
		return nil, err
	}
	for _, tab := range tabs {
		if strings.Contains(tab.Url, kw) || strings.Contains(tab.Id, kw) || strings.Contains(tab.Title, kw) {
			tab.Ipc.events = make(chan []byte, 1024)
			tab.Ipc.returns = make(chan []byte, 1024)
			tab.Ipc.errors = make(chan []byte, 1024)
			return tab, nil
		}
	}
	return nil, errors.New("tab not found")
}
