package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/john-pettigrew/workerpool/workerpool"
)

type SoundInstruction struct {
	Filepath string
}

func playSound(taskData interface{}) error {
	soundInstruction, ok := taskData.(SoundInstruction)
	if !ok || soundInstruction.Filepath == "" {
		return errors.New("invalid data")
	}

	cmd := exec.Command("aplay", soundInstruction.Filepath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func parseParams() (string, string, int, int, error) {
	addr := flag.String("addr", ":8500", "server address")
	soundsDir := flag.String("sounds-dir", "", "directory containing sound files")
	numWorkers := flag.Int("num-workers", 3, "number of workers")
	timeoutSeconds := flag.Int("timeout-seconds", 0, "max sound timeout seconds")

	flag.Parse()

	// addr
	if *addr == "" {
		return "", "", 0, 0, errors.New("invalid 'addr' value")
	}

	// soundsDir
	if *soundsDir == "" {
		return "", "", 0, 0, errors.New("'sounds-dir' required")
	}
	fi, err := os.Stat(*soundsDir)
	if err != nil {
		return "", "", 0, 0, err
	}
	if !fi.IsDir() {
		return "", "", 0, 0, errors.New("invalid 'sounds-dir' value")
	}

	// numWorkers
	if *numWorkers < 1 {
		return "", "", 0, 0, errors.New("invalid 'num-workers' value")
	}

	// timeoutSeconds
	if *timeoutSeconds <= 0 {
		return "", "", 0, 0, errors.New("invalid 'timeout-seconds' value")
	}

	return *addr, *soundsDir, *numWorkers, *timeoutSeconds, nil
}

func main() {
	addr, soundsDir, numWorkers, timeoutSeconds, err := parseParams()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	controller := Controller{
		SoundsDir: soundsDir,
		SoundManager: workerpool.NewWorkerPool(
			numWorkers,
			time.Duration(timeoutSeconds)*time.Second,
			playSound,
		),
	}

	http.HandleFunc("GET /health", controller.HealthHandler)
	http.HandleFunc("GET /play", controller.PlayHandler)
	http.HandleFunc("GET /stop-all", controller.StopAllHandler)

	fmt.Printf("server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
