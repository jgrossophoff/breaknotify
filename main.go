// breaknotify requires Linux or Darwin with
// either notify-send or osascript installed.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var (
	notifyOnStart  = flag.Bool("onstart", false, "Show notification on start")
	notifyInterval = time.Minute * 90
)

const (
	checkInterval       = time.Minute * 5
	defaultNotification = `Time to take a break!`
	os                  = runtime.GOOS
)

type proverbs []string

func main() {
	if os != "linux" && os != "darwin" {
		log.Fatalf("OS %s not supported. Only linux and darwin for now.\n", os)
	}

	if !notificationBinExists() {
		log.Fatalf("make sure that either notify-send or osascript are installed\n")
	}

	flag.Parse()

	var lastChk time.Time
	if !*notifyOnStart {
		lastChk = time.Now()
	}

	for {
		if now := time.Now(); now.Sub(lastChk) >= notifyInterval {
			notify()
			lastChk = now
		}

		time.Sleep(checkInterval)
	}
}

func notificationBinExists() bool {
	var err error
	if os == "linux" {
		_, err = exec.LookPath("notify-send")
	} else if os == "darwin" {
		_, err = exec.LookPath("osascript")
	}
	return err == nil
}

// notify will call either notify-send or osascript
// to trigger a OS notification.
func notify() {
	msg, err := fetchRandomProverb()

	if os == "linux" {
		if err == nil {
			err := exec.Command("notify-send", msg+"\nSo take a break!").Run()
			if err != nil {
				log.Println(err)
			}
		} else {
			err := exec.Command("notify-send", defaultNotification).Run()
			if err != nil {
				log.Println(err)
			}
		}
	} else if os == "darwin" {
		arg := fmt.Sprintf("display notification \"%s\" with title \"%s\"", msg, defaultNotification)
		err := exec.Command("osascript", "-e", arg).Run()
		if err != nil {
			log.Println(err)
		}
	}
}

// fetchRandomProverb will try to fetch a random proverb
// from http://proverbs-app.antjan.us.
func fetchRandomProverb() (string, error) {
	var cl http.Client
	cl.Timeout = time.Second * 7
	resp, err := cl.Get("http://proverbs-app.antjan.us")
	if err != nil {
		log.Println(err)
		return "", nil
	}

	var prvbs proverbs

	err = json.NewDecoder(resp.Body).Decode(&prvbs)
	if err != nil {
		log.Println(err)
		return "", err
	}

	l := len(prvbs)
	if l == 0 {
		return "", err
	}

	return prvbs[rand.New(rand.NewSource(time.Now().Unix())).Intn(l)-1], nil
}
