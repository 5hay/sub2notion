package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kjk/notionapi"
)

const VERSION = "1.0.0"

type App struct {
	client *notionapi.Client
}

func main() {
	app := &App{
		client: &notionapi.Client{},
	}

	pageID := os.Getenv("NOTION_PAGEID")
	if pageID == "" {
		log.Fatalln("Notion page id was not specified, specify env var NOTION_PAGEID.")
	}
	if os.Getenv("NOTION_CMD") == "" {
		log.Fatalln("No command set to run on changes, please set env var NOTION_CMD")
	}

	// First start, get the initial value
	pageTitle, lastChange := app.getPageLastChange(pageID)
	log.Printf("Starting sub2notion (v%s) | Title: %s, last change: %s", VERSION, pageTitle, lastChange.Format(time.RFC850))

	doneC := make(chan bool)
	signalC := make(chan os.Signal)
	signal.Notify(signalC, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalC
		log.Printf("Received %s signal, exiting sub2notion.", sig)
		doneC <- true
	}()

	seconds, err := strconv.Atoi(os.Getenv("NOTION_CMD_INTERVAL_SECONDS"))
	if err != nil {
		seconds = 60
		log.Printf("Checking for changes interval could not be retrieved: %s.\nYou can customize this by setting env var NOTION_CMD_INTERVAL_SECONDS\nSetting default value of %d seconds.\n", err, seconds)
	}

	go func() {
		for {
			time.Sleep(time.Duration(seconds) * time.Second)

			title, t := app.getPageLastChange(pageID)
			log.Printf("Current check | Title: %s, last change: %s\n", title, t.Format(time.RFC850))

			if t.After(lastChange) {
				pageTitle = title
				lastChange = t
				runCommand()
			}
		}
	}()

	<-doneC
}

func (app *App) getPageLastChange(pageID string) (string, time.Time) {
	// We need to download the page every time because otherwise it would just be cached
	page, err := app.client.DownloadPage(pageID)
	if err != nil {
		log.Fatalf("DownloadPage() method failed with %s\n", err)
	}
	return page.Root().Title, page.Root().LastEditedOn()
}

func runCommand() {
	cmdToRun := os.Getenv("NOTION_CMD")
	cmdArgs := os.Getenv("NOTION_CMD_ARGS")

	fullCmdPath, err := exec.LookPath(cmdToRun)
	if err != nil {
		log.Fatalf("Could not get full path for command %s, error: %s\n", cmdToRun, err)
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("%s %s", fullCmdPath, cmdArgs))
	cmd.Env = os.Environ()

	log.Printf("Command to run: %s", cmd.Args[0:])

	bytes, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error calling command: %s\n", err)
	}

	log.Printf("Command returned:\n%s", string(bytes))
}
