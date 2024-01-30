package exporter

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var timeRegex = regexp.MustCompile(`\d\d:\d\d:\d\d`)

func parseTimeSeconds(timeStr string) *float64 {
	match := timeRegex.FindStringSubmatch(timeStr)
	if len(match) < 1 {
		return nil
	}
	parts := strings.Split(match[0], ":")
	duration := parts[0] + "h" + parts[1] + "m" + parts[2] + "s"
	t, _ := time.ParseDuration(duration)
	seconds := t.Seconds()
	return &seconds
}

func parseBool(b bool) float64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func retrieveData(frmAddress string, details any) error {
	resp, err := http.Get(frmAddress)

	if err != nil {
		log.Printf("error fetching statistics from FRM: %s\n", err)
		return err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&details)

	// Try to found if it is an empty object
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		details = []any{}
		return nil
	}

	return err
}
