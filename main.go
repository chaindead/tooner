package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
)

var (
	buildTag    = "dev"
	showVersion = pflag.BoolP("version", "v", false, "show version")
)

func main() {
	pflag.Parse()

	if *showVersion {
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

	waitMap := newWait()
	// Client → Server
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()

			if gjson.Get(line, "method").String() == "tools/call" {
				id := gjson.Get(line, "id").String()
				logger.Println("!!! call:", id)

				waitMap.Add(id)
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
			ok := waitMap.Take(id)

			if ok && id != "" {
				logger.Println("!!! detect:", id)
				line = convert(logger, line)
			}

			logger.Println("<<<", line)
			os.Stdout.WriteString(line + "\n")
		}
	}()

	// Log stderr too
	go io.Copy(logFile, serverErr)

	cmd.Wait()
}
