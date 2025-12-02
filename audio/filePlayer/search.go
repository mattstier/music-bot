// This source file is a part of the Discord bot 'Audiophile' made by Máté Stier
//
// Copyright (C) 2025 Máté Stier
//
// This software is licensed under the "GPLv3" License as described in the "LICENSE" file,
// which should be included with this package. The terms are also available at
// http://www.gnu.org/licenses/gpl-3.0.html

package filePlayer

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"github.com/masatana/go-textdistance"
	_ "github.com/masatana/go-textdistance"
	"golang.org/x/text/unicode/norm"
)

const searchSimilarityThreshold = 0.15
const weightArtist = 1
const weightAlbum = 1
const weightTitle = 2
const normalizedCacheStore = ".cache.gob"
const cacheBuffer = 50

type SongData struct {
	Album    string `json:"album"`
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Filename string
}

// cache mapping original filenames to their normalized SongData
var normalizedCache map[string]SongData

func GetClosestMatch(term string) string {
	var highestName string
	highest := 0.0
	//normalizing input
	normalizedTerm := normalize(term)
	fmt.Println(mediaDir)
	loadFileNames(mediaDir)

	for fileName, cachedFileMetadata := range normalizedCache {
		cachedFileMetadata.Filename = fileName
		current := similarity(cachedFileMetadata, normalizedTerm)
		fmt.Printf("\nSimilarity of %v to %v is %v", cachedFileMetadata, term, current)
		if highest < current {
			highest = current
			highestName = fileName
		}
	}
	if highest < searchSimilarityThreshold {
		return ""
	}
	return highestName
}

// using a weighed average of the similarities of the title, artist and album
// which use the levenshtein distance
func similarity(metadata SongData, term string) float64 {
	// fallback to filename if title is missing
	title := metadata.Title
	if title == "" {
		title = metadata.Filename
	}

	// map of field value → weight
	fields := map[string]float64{
		title:           weightTitle,
		metadata.Album:  weightAlbum,
		metadata.Artist: weightArtist,
	}

	var sum, total float64
	for val, weight := range fields {
		if val == "" {
			continue
		}
		maxLen := max(len(val), len(term))
		if maxLen == 0 {
			continue
		}
		sum += (1 - float64(textdistance.LevenshteinDistance(val, term))/float64(maxLen)) * weight
		total += weight
	}

	if total == 0 {
		return 0
	}
	return sum / total
}

func normalize(str string) string {

	str = strings.ReplaceAll(str, "_", " ") //underscore
	str = strings.ReplaceAll(str, ".", " ")
	str = strings.ReplaceAll(str, "-", " ") // hypen
	str = strings.ReplaceAll(str, "–", " ") // en dash
	str = strings.ToLower(str)
	return removeAccents(str)
}

func removeAccents(str string) string {
	//decompose characters
	decomposed := norm.NFD.String(str)

	//filtering out combining marks
	filtered := make([]rune, 0, len(decomposed))
	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue //skip accent mark
		}
		filtered = append(filtered, r)
	}
	return string(filtered)
}

// loads filenames from disk to the cache (along with their normalized versions)
func loadFileNames(path string) {
	var currentFileName string
	//if the cache is not loaded in and exists saved
	if normalizedCache == nil && hasCacheOnDisk() {
		cache, err := loadCache()
		if err == nil {
			normalizedCache = cache
			return
		}
	}
	files, _ := os.ReadDir(path)
	//cache buffer is just the expected max number of files; minimizes resizing/rehashing overhead
	normalizedCache = make(map[string]SongData, cacheBuffer)

	for _, file := range files {
		currentFileName = file.Name()
		//cache result, if not cached already
		if _, cacheHit := normalizedCache[currentFileName]; !cacheHit {
			normalizedCache[currentFileName] = getNormalizedSongData(mediaDir + "/" + currentFileName)
		}
	}
	dumpCache()
}

func dumpCache() error {
	os.MkdirAll(cacheDir, 0755)
	f, _ := os.Create(cacheDir + "/" + normalizedCacheStore)
	defer f.Close()
	fmt.Println("Dumping Cache")
	enc := gob.NewEncoder(f)
	return enc.Encode(normalizedCache)
}

func loadCache() (map[string]SongData, error) {
	file, err := os.Open(cacheDir + "/" + normalizedCacheStore)
	fmt.Println("Loading cache")
	defer file.Close()

	cache := make(map[string]SongData, cacheBuffer)
	dec := gob.NewDecoder(file)
	err = dec.Decode(&cache)
	return cache, err
}

func hasCacheOnDisk() bool {
	_, err := os.Stat(cacheDir + "/" + normalizedCacheStore)
	return err == nil
}

// returns a string of the artist and the album
func getSongData(filePath string) SongData {
	var out struct {
		Format struct {
			Tags SongData `json:"tags"`
		} `json:"format"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	//ffprobe checks metadata of a file
	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_entries", "format_tags=artist,title,album",
		filePath)
	output, err := cmd.Output()
	if err != nil {
		return SongData{}
	}
	if err = json.Unmarshal(output, &out); err != nil {
		return SongData{}
	}
	return out.Format.Tags
}

func getNormalizedSongData(filePath string) SongData {
	metadata := getSongData(filePath)
	metadata.Album = normalize(metadata.Album)
	metadata.Artist = normalize(metadata.Artist)
	metadata.Title = normalize(metadata.Title)
	return metadata
}
