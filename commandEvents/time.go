// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package commandEvents

import (
	"errors"
	"fmt"
	"math"
	"music-bot/audio/types"
	"strconv"
	"strings"
	"time"
)

func formatTimestamp(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func generateLoadingBar(timestamp time.Duration, songLength time.Duration, size int) string {
	if songLength == 0 {
		return ""
	}
	loadingBar := ""
	conversionRatio := float64(timestamp.Milliseconds()) / float64(songLength.Milliseconds())
	loaded := int(math.Ceil(float64(size) * conversionRatio))

	for i := 0; i < size; i++ {
		if i <= loaded {
			loadingBar += "▓"
		} else {
			loadingBar += "░"
		}
	}
	return loadingBar
}

func parseTimeStamp(userArg string, player types.Player) (time.Duration, error) {
	// cleaning up timestamp for spaces
	userArg = strings.TrimSpace(userArg)
	currentSong := player.CurrentSong()

	timestamp, err := time.ParseDuration(userArg)
	//handle going out bounds with the songs duration or a parsing error
	if err != nil {
		// in case it cannot be parsed but is a valid integer, we interpret them as seconds
		if userInt, atoiErr := strconv.Atoi(userArg); atoiErr == nil {
			timestamp = time.Duration(userInt) * time.Second
		} else if t, timeParseErr := time.Parse("04:05", userArg); timeParseErr == nil {
			timestamp = time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
		} else {
			return time.Duration(0), errors.New("Invalid time format")
		}
	}

	// handles negative user argument, only allowing it if it is within the songs bounds
	if (player.Timestamp()+timestamp) < 0 || timestamp > currentSong.GetDuration() {
		return time.Duration(0), errors.New("Invalid timestamp")
	}
	return timestamp, nil
}
