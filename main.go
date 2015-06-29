package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"handlers"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type flags struct {
	isProxy    bool
	listenIP   string
	listenPort string
	proxyIP    string
	proxyPort  string
	name       string
}

type heartbeat struct {
	fromIP   string
	fromName string
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func printMux() {
	for key, value := range mux {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func getCmdLineArgs() flags {
	isProxy := flag.Bool("x", true, "run as proxy")
	listenIP := flag.String("lip", "0.0.0.0", "listen IP address")
	listenPort := flag.String("lp", "8005", "listen port")
	proxyIP := flag.String("pip", "0.0.0.0", "proxy IP address")
	proxyPort := flag.String("pp", "8006", "proxy port")
	name := flag.String("n", "NoName", "Name of this agent")
	flag.Parse()

	return flags{*isProxy, *listenIP, *listenPort, *proxyIP, *proxyPort, *name}
}

func setupRoutes() {
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = handlers.HandleRequest
	mux["/kakki"] = handlers.HandleRequest1
	mux["/foo"] = handlers.HandleRequest2
}

func sendHeartbeat(f flags) {
	url := "http://" + f.proxyIP + ":" + f.proxyPort
	fmt.Println("URL:>", url)

	hb := map[string]string{"fromIP": f.listenIP, "fromName": f.name}
	json, _ := json.Marshal(hb)
	var jsonStr = []byte(string(json))
	//var hb heartbeat = []byte(`{"fromIP":"123.345.234.111"}, {"fromName":"antti"}`)
	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func startHeartbeat(sec time.Duration, f flags) {
	ticker := time.NewTicker(sec * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				sendHeartbeat(f)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {

	flags := getCmdLineArgs()
	fmt.Println("Command Line Flags:")
	if flags.isProxy {
		fmt.Println("[isProxy: true]")
	} else {
		fmt.Println("[isProxy: false]")
	}
	fmt.Println("[listenUrl: " + flags.listenIP + "]")
	fmt.Println("[listenPort: " + flags.listenPort + "]")
	fmt.Println("[proxyUrl: " + flags.proxyIP + "]")
	fmt.Println("[proxyPort: " + flags.proxyPort + "]\n")

	server := http.Server{
		Addr:    flags.listenIP + ":" + flags.listenPort,
		Handler: &myHandler{},
	}

	setupRoutes()

	fmt.Println("Listening to " + server.Addr)
	if !flags.isProxy {
		go startHeartbeat(5, flags)
	}
	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request: " + r.URL.String())
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}
