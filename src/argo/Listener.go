package argo

import (
	"crypto/tls"
	"fmt"
	"github.com/r3labs/sse"
	"net/http"
)

func Listen() {
	fmt.Println("Start listening argo events")

	client := sse.NewClient("https://34.71.103.174/api/v1/stream/applications?name=task")
	client.Headers = map[string]string{
		"Cookie": "argocd.token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTQ5OTEzNDcsImlzcyI6ImFyZ29jZCIsIm5iZiI6MTU5NDk5MTM0Nywic3ViIjoiYWRtaW4ifQ.2XAKP67ovnkCqlX4nEdqMTJH894LPt7QPJo-fQzgyAY",
	}
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	err := client.Subscribe("message", func(msg *sse.Event) {
		// Got some data!
		fmt.Println(string(msg.Data))
	})

	client.OnDisconnect(func(c *sse.Client) {
		fmt.Println("Disconnect")
	})

	fmt.Println(err)
}
