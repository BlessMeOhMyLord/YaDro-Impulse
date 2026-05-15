package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const errMessage = "Use: go run ./cmd/client --help, for information of command list"

var restServer = ""

type DNSRequest struct {
	DNS string `json:"dns"`
}

type DNSResponse struct {
	DNS string `json:"dns"`
}

type DNSListResponse struct {
	Items []string `json:"items"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	restServer = os.Getenv("CLIENT_ADDR")

	if len(os.Args) < 2 {
		fmt.Println(errMessage)
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		if len(os.Args) != 3 {
			fmt.Println("expected argument dns server")
			os.Exit(1)
		}
		addDNS(os.Args[2])
	case "list":
		listDNS()
	case "del":
		if len(os.Args) != 3 {
			fmt.Println("expected argument dns server")
			os.Exit(1)
		}
		deleteDNS(os.Args[2])
	case "--help":
		help()
	default:
		fmt.Println("Unknown command")
		fmt.Println(errMessage)
		os.Exit(1)
	}
}

func addDNS(newDns string) {
	c := &http.Client{Timeout: 2 * time.Second}
	body, _ := json.Marshal(DNSRequest{newDns})
	request, _ := http.NewRequest(http.MethodPost, restServer, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response, err := c.Do(request)
	if err != nil {
		log.Fatal("Error add dns:", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		printError(response)
		os.Exit(1)
	}

	var result DNSResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	log.Println("DNS added: ", result.DNS)
}

func listDNS() {
	c := &http.Client{Timeout: 2 * time.Second}
	request, _ := http.NewRequest(http.MethodGet, restServer, nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.Do(request)
	if err != nil {
		log.Fatal("Error list dns:", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		printError(response)
		os.Exit(1)
	}

	var result DNSListResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Error list dns:", err)
	}

	if len(result.Items) == 0 {
		fmt.Println("DNS list is empty")
	} else {
		for _, dns := range result.Items {
			fmt.Println(dns)
		}
	}
}

func deleteDNS(dnsToDelete string) {
	c := &http.Client{Timeout: 2 * time.Second}
	body, _ := json.Marshal(DNSRequest{dnsToDelete})
	request, _ := http.NewRequest(http.MethodDelete, restServer, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response, err := c.Do(request)
	if err != nil {
		log.Fatal("Error delete dns:", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		printError(response)
		os.Exit(1)
	}

	var result DNSResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Fatal("Error delete dns:", err)
	}

	fmt.Println("dns deleted: ", result.DNS)
}

func help() {
	fmt.Println("add <dns server>   |   add dns server to resolv.conf")
	fmt.Println("del <dns server>   |   del dns server from resolv.conf if it exists")
	fmt.Println("list               |   list dns servers from resolv.conf")
}

func printError(response *http.Response) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("server error: status=%d\n", response.StatusCode)
		return
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		fmt.Printf("server error: status=%d body=%s\n", response.StatusCode, string(body))
		return
	}

	fmt.Printf("server error: status=%d error=%s\n", response.StatusCode, errResp.Error)
}
