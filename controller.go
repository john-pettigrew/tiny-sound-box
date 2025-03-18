package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/john-pettigrew/workerpool/workerpool"
)

const DEFAULT_LOOP_TIMES int = 1
const DEFAULT_DELAY_SECONDS int = 0

func parseIntParam(param string, defaultValue int) (int, error) {
	if param == "" {
		return defaultValue, nil
	}

	paramParsed, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return paramParsed, nil
}

type Controller struct {
	SoundManager *workerpool.WorkerPool
	SoundsDir    string
}

func (c *Controller) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) StopAllHandler(w http.ResponseWriter, r *http.Request) {
	c.SoundManager.StopAllRunningTasks()
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) PlayHandler(w http.ResponseWriter, r *http.Request) {
	soundNamePattern := regexp.MustCompile(`^[\w]*$`)
	params := r.URL.Query()

	sound := params.Get("sound")
	if sound == "" || !soundNamePattern.MatchString(sound) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid sound name")
		return
	}

	loopTimes, err := parseIntParam(params.Get("loop"), DEFAULT_LOOP_TIMES)
	if err != nil || loopTimes < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid loop argument value")
		return
	}

	delaySeconds, err := parseIntParam(params.Get("delay"), DEFAULT_DELAY_SECONDS)
	if err != nil || delaySeconds < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid delay argument value")
		return
	}

	soundFilePath := filepath.Join(c.SoundsDir, sound+".wav")

	_, err = os.Stat(soundFilePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "missing sound file")
		return
	}

	c.SoundManager.StartNewTask(
		loopTimes,
		time.Duration(delaySeconds)*time.Second,
		SoundInstruction{
			Filepath: soundFilePath,
		},
	)

	w.WriteHeader(http.StatusOK)
}
