package main

import (
	"net/http"
	"log"
	"orbs/orbserver"
	"strconv"
	"io/ioutil"
	"flag"
	"encoding/json"
	"gopkg.in/natefinch/lumberjack.v2"

)

func writeLog(ip string, payload string, errorcode int) {
	log.Printf("%v \"%v\" %v\n", ip, payload, errorcode)
}

func main() {
	config_file := flag.String("config", "config.yml", "Path to the configuration file")
	flag.Parse()

	config := orbserver.ParseConfig(*config_file)

	res_index_data, err := ioutil.ReadFile(config.IndexPath)
	if err != nil {
		log.Fatal(err)
	}

	var res_index interface{}

	err = json.Unmarshal(res_index_data, &res_index)
	if err != nil {
		log.Fatal(err)
	}

	//list of valid game character sprite resource keys
	var spriteNames []string
	for k := range res_index.(map[string]interface{})["cache"].(map[string]interface{})["charset"].(map[string]interface{}) {
		if k != "_dirname" {
			spriteNames = append(spriteNames, k)
		}
	}

	//list of valid system resource keys
	var systemNames []string
	for k := range res_index.(map[string]interface{})["cache"].(map[string]interface{})["system"].(map[string]interface{}) {
		if k != "_dirname" {
			systemNames = append(systemNames, k)
		}
	}

	var roomNames []string

	for i:=0; i < config.NumRooms; i++ {
		roomNames = append(roomNames, strconv.Itoa(i))
	}

	orbserver.CreateAllHubs(roomNames, spriteNames, systemNames)

	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Logging.File,
		MaxSize:    config.Logging.MaxSize,
		MaxBackups: config.Logging.MaxBackups,
		MaxAge:     config.Logging.MaxAge,
	})
	log.SetFlags(log.Ldate | log.Ltime)

	http.Handle("/", http.FileServer(http.Dir("public/")))
	log.Fatalf("%v %v \"%v\" %v", config.IP, "server", http.ListenAndServe(":" + strconv.Itoa(config.Port), nil), 500)
}
