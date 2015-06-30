// handlers.go
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func HandleSpa(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Body: " + string(body))
	var data string
	json.Unmarshal(body, &data)
	fmt.Println("JSON: " + data)
	io.WriteString(w, "Hello spa!")
}

type Heartbeat struct {
	repoTags string `json:"RepoTags"`
}

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Body: " + string(body))
	var hb Heartbeat
	json.Unmarshal(body, &hb)

	fmt.Println(hb)
	io.WriteString(w, "Hello hb!")
}

type Images []struct {
	Created     uint64   `json:"Created"`
	Id          string   `json:"Id"`
	Labels      string   `json:"Labels"`
	ParentId    string   `json:"ParentId"`
	RepoDigests []string `json:"RepoDigests"`
	RepoTags    []string `json:"RepoTags"`
	Size        uint64   `json:"Size"`
	VirtualSize uint64   `json:"VirtualSize"`
}

func HandleImages(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Body: " + string(body))

	var images Images
	json.Unmarshal(body, &images)

	for _, image := range images {
		fmt.Println("%v", image)
	}
	io.WriteString(w, "Hello images!")
}

//	hb := data.(map[string]interface{})

//for _, item := range strings.Split(string(body), ",") {
//		if strings.Contains(string(item), "Created") {
//			fmt.Println("%s", item)
//		}
//		if strings.Contains(string(item), "RepoTags") {
//			fmt.Println("%s", item)
//		}
//	}

//for k, v := range images {
//		switch vv := v.(type) {
//		case string:
//			fmt.Println(k, "is string", vv)
//		case int:
//			fmt.Println(k, "is int", vv)
//		case []interface{}:
//			fmt.Println(k, "is an array:")
//			for i, u := range vv {
//				fmt.Println(i, u)
//			}
//		default:
//			fmt.Println(k, "is of a type I don't know how to handle")
//		}
//	}
