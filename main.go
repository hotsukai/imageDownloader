package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	rawJSON, err := ioutil.ReadFile("data.json")
	var emojiData map[string]string
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(rawJSON, &emojiData)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir("images", 0777)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Chdir("images")
	c := make(chan bool, 5)
	wg := &sync.WaitGroup{}
	for imageName, url := range emojiData {
		wg.Add(1)
		go func(imageName string, url string) {
			c <- true
			defer func() {
				<-c
			}()
			defer wg.Done()
			fmt.Println(imageName)
			if strings.Contains(url, "alias:") {
				return
			}
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			fp, err := os.Create(imageName + url[len(url)-4:])
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(fp, resp.Body)
			if err != nil {
				fmt.Println("error", err)
			}
		}(imageName, url)
	}
	fmt.Println("waiting")
	wg.Wait()
	fmt.Println("end")
}
