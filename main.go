package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	var oneLine, verboseMode bool
	var webhookURL, lines string

	flag.StringVar(&webhookURL, "u", "", "Discord Webhook URL")
	flag.BoolVar(&oneLine, "1", false, "Send message line by line")
	flag.BoolVar(&verboseMode, "v", false, "Verbose mode")
	flag.Parse()

	webhookENV := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookENV != "" {
		webhookURL = webhookENV
	}else{
		if webhookURL == "" {
		    if verboseMode {
		        fmt.Println("Discord Webhook URL not set!")
		    }
		}
	}

	if !isStdin() {
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
			lines += line
			lines += "\n"
	}
	wg.Add(1)
	go disc(webhookURL, lines)
	wg.Wait()
}

func isStdin() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if info.Mode()&os.ModeNamedPipe == 0 { //to make sure the input comes from the pipe
		return false
	}
	return true
}

type data struct {
	Content string `json:"content"`
}

func disc(url string, line string) {
	data, _ := json.Marshal(data{ Content: line})
	http.Post(url, "application/json", strings.NewReader(string(data)))
	wg.Done()
}