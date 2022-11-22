package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func getIP() string {
	url := "https://ifconfig.me/ip"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(res)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func main() {
	CLOUDFLARE_API_KEY := os.Getenv("CLOUDFLARE_API_KEY")
	CLOUDFLARE_API_EMAIL := os.Getenv("CLOUDFLARE_API_EMAIL")
	CLOUDFLARE_DOMAIN := os.Getenv("CLOUDFLARE_DOMAIN")
	CLOUDFLARE_RECORD_ID := os.Getenv("CLOUDFLARE_RECORD_ID")

	api, err := cloudflare.New(CLOUDFLARE_API_KEY, CLOUDFLARE_API_EMAIL)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	id, err := api.ZoneIDByName(CLOUDFLARE_DOMAIN)
	if err != nil {
		log.Fatal(err)
	}
	records, err := api.DNSRecord(ctx, id, CLOUDFLARE_RECORD_ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	IP := getIP()
	if len(IP) == 0 {
		panic("IP is none.")
	}
	if len(records.Content) == 0 {
		panic("records.Content is none.")
	}
	log.Println(records.Content)
	log.Println(IP)
	if records.Content != IP {
		records.Content = IP

		err := api.UpdateDNSRecord(ctx, id, CLOUDFLARE_RECORD_ID, records)
		if err != nil {
			panic(err)
		}
	}
}
