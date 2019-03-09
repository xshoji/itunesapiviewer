package main

import (
	"encoding/json"
	"fmt"
	"github.com/zserge/lorca"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	ui, _ := lorca.New("", "", 1200, 640)
	defer ui.Close()

	//// - [web applications - How do I serve CSS and JS in Go Lang - Stack Overflow](https://stackoverflow.com/questions/43601359/how-do-i-serve-css-and-js-in-go-lang)
	ln, err := net.Listen("tcp", "127.0.0.1:0") //監視するポートを設定します。
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))
	ui.Bind("getResponseSearch", func(query string) interface{} {
		response, err := http.Get("https://itunes.apple.com/search?term=" + query + "&media=music")
		if err != nil {
			fmt.Println("itunes.apple.com http error")
			os.Exit(1)
		}
		// body to string
		var bodyString string
		body, err := ioutil.ReadAll(response.Body)
		bodyString = string(body)

		// string to object
		var bodyObject interface{}
		defer response.Body.Close()
		json.Unmarshal([]byte(bodyString), &bodyObject)

		// extract values
		//fmt.Println(bodyObject.(map[string]interface{})["resultCount"])
		//fmt.Println(bodyObject.(map[string]interface{})["results"].([]interface{})[1].(map[string]interface{})["trackName"])
		//fmt.Println(bodyObject.(map[string]interface{})["results"].([]interface{})[1].(map[string]interface{})["artistName"])
		return bodyObject
	})

	<-ui.Done()
}
