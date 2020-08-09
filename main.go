package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/hawwwdi/TCP-server/model"
)

func main() {
	server, err := net.Listen("tcp", "localhost:4323")
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; true; i++ {
		fmt.Println("wating for client ", i)
		client, err := server.Accept()
		fmt.Println("client ", i, " accepted")
		if err != nil {
			log.Println(err)
		}
		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()
	isGET := parseReq(client)
	if isGET {
		sendRes(client)
	}

}

func parseReq(reader io.Reader) bool {
	var flag bool
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), "GET") {
			flag = true
		} else if scanner.Text() == "" {
			break
		}
	}
	return flag
}

func sendRes(writer io.Writer) {
	body := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title></title></head><body><strong>hello world :)</strong></body></html>`
	res := model.NewBuilder().Protocol("HTTP/1.1").Headers(
		"Content-Length", fmt.Sprint(len(body)),
		"Content-Type", " text/html").Body(body).Status("200 OK").Build()
	//fmt.Println(res.String())
	fmt.Fprintln(writer, res.String())
}
