// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package filePlayer

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Song struct {
	name     string
	duration *time.Duration
	path     string
}

func (s Song) GetName() string {
	return s.name
}

func (s Song) GetDuration() time.Duration {
	if s.duration != nil {
		return *s.duration
	}
	//if the song has no duration associated, we try and fetch it from the file
	var duration time.Duration
	out, _ := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		mediaDir+"/"+s.GetName()).Output()

	f, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		duration = time.Duration(0)
	}
	// Convert float to time.Duration
	duration = time.Duration(f * float64(time.Second))
	s.duration = &duration
	return duration
}

func (s Song) GetFilePath() string {
	return s.path
}
