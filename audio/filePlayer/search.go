package filePlayer

import (
	"os"
	"strings"
	"unicode"

	"github.com/masatana/go-textdistance"
	"golang.org/x/text/unicode/norm"
)

const searchSimilarityThreshold = 0.6

// cache mapping original filenames to their normalized counterparts
// note: 50 is just the expected max number of files; minimizes resizing/rehashing overhead
var normalizedCache = make(map[string]string, 50)

func GetClosestMatch(term string) string {
	var highestName string
	highest := 0.0
	//normalizing input
	q := normalize(term)
	loadFileNames(audioPath)
	for originalFileName, cachedFileName := range normalizedCache {
		if strings.Contains(cachedFileName, q) {
			return originalFileName
		}
		current := similarity(cachedFileName, q)
		if highest < current {
			highest = current
			highestName = originalFileName
		}
	}
	if highest < searchSimilarityThreshold {
		return ""
	}
	return highestName
}

func similarity(a, b string) float64 {
	distance := textdistance.LevenshteinDistance(a, b)
	maxLen := max(len(a), len(b))
	return 1.0 - (float64(distance) / float64(maxLen))
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
	files, _ := os.ReadDir(path)

	for _, file := range files {
		currentFileName = file.Name()
		//cache result, if not cached already
		if _, hasCache := normalizedCache[currentFileName]; !hasCache {
			normalizedCache[currentFileName] = normalize(currentFileName)
		}
	}
}
