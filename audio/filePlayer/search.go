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
	"golang.org/x/text/unicode/norm"
)

const searchSimilarityThreshold = 0.15
const weightArtist = 1
const weightAlbum = 1
const weightTitle = 2

const normalizedCacheStore = "cache.gob"
const cacheBuffer = 50

type SongData struct {
	Album  string `json:"album"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

// cache mapping original filenames to their normalized counterparts
// note: 50 is just the expected max number of files; minimizes resizing/rehashing overhead
var normalizedCache = make(map[string]SongData, cacheBuffer)

func GetClosestMatch(term string) string {
	var highestName string
	highest := 0.0
	//normalizing input
	q := normalize(term)
	loadFileNames(audioPath)

	for fileName, cachedFileMetadata := range normalizedCache {
		current := similarity(cachedFileMetadata, q)
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
	titleDistance := textdistance.LevenshteinDistance(metadata.Title, term)
	albumDistance := textdistance.LevenshteinDistance(metadata.Album, term)
	artistDistance := textdistance.LevenshteinDistance(metadata.Artist, term)

	titleSimilarity := 1.0 - (float64(titleDistance) / float64(max(len(metadata.Title), len(term))))
	albumSimilarity := 1.0 - (float64(albumDistance) / float64(max(len(metadata.Album), len(term))))
	artistSimilarity := 1.0 - (float64(artistDistance) / float64(max(len(metadata.Artist), len(term))))
	sumSimilarity := titleSimilarity*weightTitle + albumSimilarity*weightAlbum + artistSimilarity*weightArtist

	return sumSimilarity / 3
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
		fmt.Println("Loading cache from disk")
		cache, err := loadCache()
		if err == nil {
			normalizedCache = cache
			return
		}
	}
	files, _ := os.ReadDir(path)

	for _, file := range files {
		currentFileName = file.Name()
		//cache result, if not cached already
		if _, cacheHit := normalizedCache[currentFileName]; !cacheHit {
			fmt.Println("making cache")
			normalizedCache[currentFileName] = getNormalizedSongData(audioPath + "/" + currentFileName)
			fmt.Println(normalizedCache)
		}
	}
}

func DumpCache() error {
	f, _ := os.Create(audioPath + normalizedCacheStore)
	defer f.Close()
	fmt.Println("Dumping Cache")
	enc := gob.NewEncoder(f)
	return enc.Encode(normalizedCache)
}

func loadCache() (map[string]SongData, error) {
	file, err := os.Open(audioPath + normalizedCacheStore)
	fmt.Println("Loading cache")
	defer file.Close()

	cache := make(map[string]SongData, cacheBuffer)
	dec := gob.NewDecoder(file)
	err = dec.Decode(&cache)
	return cache, err
}

func hasCacheOnDisk() bool {
	_, err := os.Stat(audioPath + normalizedCacheStore)
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
