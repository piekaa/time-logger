package main


import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:4503/flights/gb/en/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", "localhost:9982")
	req.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "2468a737-312e-4f31-bb83-9fbb89aa335e")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}