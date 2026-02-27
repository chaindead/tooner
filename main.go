package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	buildTag = "dev"
)

func main() {
	if os.Getenv("TOONER_SHOW_VERSION") != "" {
		fmt.Println("tooner", buildTag)
		os.Exit(0)
	}

	logPath := os.Getenv("TOONER_LOG_PATH")
	if logPath == "" {
		logPath = "tooner.log"
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		log.Fatal("Cannot create log file", err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	if len(os.Args) <= 1 {
		log.Fatalln("Please provide at least one command")
	}

	args := os.Args[1:]

	cmd := exec.Command(args[0], args[1:]...)
	serverIn, _ := cmd.StdinPipe()
	serverOut, _ := cmd.StdoutPipe()
	serverErr, _ := cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		log.Fatalln("stating provided command", err)
	}

	waitCall := newWait()
	waitList := newWait()
	// Client → Server
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()

			if gjson.Get(line, "method").String() == "tools/call" {
				id := gjson.Get(line, "id").String()
				logger.Println("!!! call tool:", id)

				waitCall.Add(id)
			}

			if gjson.Get(line, "method").String() == "tools/list" {
				id := gjson.Get(line, "id").String()
				logger.Println("!!! call list:", id)

				waitList.Add(id)
			}

			logger.Println(">>>", line)
			io.WriteString(serverIn, line+"\n")
		}
	}()

	// Server → Client
	go func() {
		scanner := bufio.NewScanner(serverOut)
		for scanner.Scan() {
			line := scanner.Text()

			id := gjson.Get(line, "id").String()
			ok := waitCall.Take(id)

			if ok && id != "" {
				logger.Println("!!! detect:", id)
				line = convert(logger, line)
			}

			ok = waitList.Take(id)
			if ok && id != "" {
				logger.Println("!!! detect list:", id)

				r := gjson.Get(line, "result.tools").Array()
				for i := range r {
					path := fmt.Sprintf("result.tools.%d.outputSchema", i)

					line, err = sjson.Delete(line, path)
					if err != nil {
						logger.Println("!!! detect list index:", i, err)
					}
				}
			}

			logger.Println("<<<", line)
			os.Stdout.WriteString(line + "\n")
		}
	}()

	// Log stderr too
	go io.Copy(logFile, serverErr)

	cmd.Wait()
}
