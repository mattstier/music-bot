## Searching files

### Introduction
This module takes care of the calculation of the similarity, retrieval and caching of filenames 
<br> and metadata within a given [Discord guild](https://discord.com/developers/docs/resources/guild).

### Getting file metadata
```go
type SongData struct {
	Album    string `json:"album"`
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Filename string
}

```
What it is:
- These fields are usually present in media files (especially audio formats)
- Therefore, we use a **weighed average** of the similarity of these fields to determine closest matching file
- With filename as a fallback for a missing title field (this is added in a [different function](#))
- We get these using in a json format from the `ffprobe` built-in command of FFMPEG  
```go
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
```
What it does:
- Runs `ffprobe` (with a 0.25 second timeout as a safety measure)
- Which outputs to following tags in JSON
- Which is unmarshalled/read into our SongData struct
- Note: the `var out struct` above is just metadata boilerplate

### Calculating Similarity
- To calculate similarity we use the following formula: <br>

```go
sum += (1 - float64(textdistance.LevenshteinDistance(val, term))/float64(maxLen)) * weight
```
What it does:
- It uses the **complement of the [Levenshtein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)** between the user's search term and the current field
- Divided by the **length of the longer string** (which is either the song data field or the search term) 
  - This normalizes the distance (keeps it always between 0-1, therefore the complement will be >0)
- Multiplied by the weight of the field
- All of which is **added to a sum** to get the **weighed average** at the end
```go
func similarity(metadata SongData, term string) float64 {
    ...
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

```

### Caching and Field Normalization
- In order for us to make searching convenient we have make similarity case and accent insensitive search 
- I.e. we have to **normalize the fields** of the SongData, 
  - however this (as well as running ffprobe on each file to get the fields) is resource and I/O heavy
  - especially if there are many files for many guilds
  - therefore we **cache them** by **mapping _filenames_ to _normalized song data_**
```go
// cache mapping original filenames to their normalized SongData
var normalizedCache map[string]SongData
```
-
  - which is also **saved as disk cache** , for a more consistent search speed
  - It is **encoded in Go's .gob** binary format to the following directory (_cacheDir_)
    - `audio/filePlayer/files/:discord_guild_id/cache/.cache.gob`
    - this is **loaded when no in-memory cache** is found, but there is on the disk
```go
func dumpCache() error {
	os.MkdirAll(cacheDir, 0755)
	f, _ := os.Create(cacheDir + "/" + normalizedCacheStore)
	defer f.Close()
	fmt.Println("Dumping Cache")
	enc := gob.NewEncoder(f)
	return enc.Encode(normalizedCache)
}
```
---
## Getting results
```go
func GetClosestMatch(term string) string {
	var highestName string
	highest := 0.0
	//normalizing input
	normalizedTerm := normalize(term)
	fmt.Println(mediaDir)
	loadFileNames(mediaDir)

```
What this does:
- **Normalizes the user's search term**, to easily approximate match with the SongData
- Then **loads the disk cache** (if there is any) to the in-memory cache `normalizedCache`
- **Returns the filename** of the most similar file

```go
...
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
```
What it does:
- Initializes in-memory cache
- For each filename in the directory it checks if it's a key in that cache
  - if it is not found it adds it, and updates the disk cache