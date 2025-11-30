package soundCloudPlayer

import (
	"fmt"
	"io"
	"net/http"
)

const tracksEndpoint = "https://api.soundcloud.com/tracks"

func getJSONResponse(term string) string {
	var query = fmt.Sprintf("%v?client_id=%v&q=%v&limit=5", tracksEndpoint, soundCloudClientID, term)
	response, err := http.Get(query)
	defer response.Body.Close()
	if err != nil {
		return ""
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("Raw JSON string:\n", bodyString)
	return bodyString
}
