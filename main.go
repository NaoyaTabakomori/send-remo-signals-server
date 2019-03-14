package main

import (
	"io/ioutil"
	"net/http"
	"time"

	remocloud "github.com/NaoyaTabakomori/go-nature-remo/cloud"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Token   string   `yaml:"token"`
	Signals []signal `yaml:"signals"`
}

type signal struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

var remoclient *remocloud.Client
var commandhash map[string]string

func testHandler(w http.ResponseWriter, r *http.Request) {
	remoclient.SendSignal(commandhash["command1"])
	sleepInterval()
	remoclient.SendSignal(commandhash["command2"])
}

func sleepInterval() {
	time.Sleep(3 * time.Second)
}

func main() {
	buf, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	var config config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	token := config.Token
	remoclient = remocloud.NewClient(token)

	commandhash = map[string]string{}
	for _, v := range config.Signals {
		commandhash[v.Name] = v.ID
	}

	http.HandleFunc("/test", testHandler)
	http.ListenAndServe(":8000", nil)
}
