package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

type DiscoveryService struct {
	Uri string
}

type Bridge struct {
	ID                string `json:"id"`
	Internalipaddress string `json:"internalipaddress"`
}

type LightState struct {
	Alert      string    `json:"alert"`
	Brightness int       `json:"bri"`
	Colormode  string    `json:"colormode"`
	Ct         int       `json:"ct"`
	Effect     string    `json:"effect"`
	Hue        int       `json:"hue"`
	On         bool      `json:"on"`
	Reachable  bool      `json:"reachable"`
	Saturation int       `json:"sat"`
	XY         []float64 `json:"xy"`
}

type Light struct {
	State            LightState `json:"state"`
	Type             string     `json:"type"`
	Name             string     `json:"name"`
	Modelid          string     `json:"modelid"`
	Manufacturername string     `json:"manufacturername"`
	Uniqueid         string     `json:"uniqueid"`
	Swversion        string     `json:"swversion"`
}

type LightContainer struct {
	Pool map[string]Light
}

func (service DiscoveryService) GetBridges() []Bridge {
	resp, err := http.Get(service.Uri)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var bridgeContainer []Bridge

	json.Unmarshal(body, &bridgeContainer)

	return bridgeContainer
}

func (bridge Bridge) GetLights() LightContainer {
	developerKey := viper.GetString("developer_key")
	resp, err := http.Get("http://" + string(bridge.Internalipaddress) + "/api/" + developerKey + "/lights")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var lightContainer LightContainer

	err = json.Unmarshal(body, &lightContainer.Pool)
	if err != nil {
		log.Fatal(err)
	}
	return lightContainer
}

func main() {
	viper.SetConfigType("yml")
	viper.SetConfigName("local")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	service := DiscoveryService{Uri: "https://www.meethue.com/api/nupnp"}

	bridges := service.GetBridges()

	for i := range bridges {
		lightContainer := bridges[i].GetLights()
		for j := range lightContainer.Pool {
			fmt.Println("%v", lightContainer.Pool[j].Name)
		}

	}
}
