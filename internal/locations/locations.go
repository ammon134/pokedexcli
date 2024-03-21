package locations

// Create 2 commands in the commands list, map & mapb
//
//
// map displays the next 20 locations
// - send GET request to the api
// - save the next and previous result into a struct
// - iterate through the result list and print
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type LocationResult struct {
	Previous string
	Next     string
	Results  []location
}

type location struct {
	Name string
	Url  string
}

func GetLocations(url *string) (*LocationResult, error) {
	// call GET here
	// if arg == "map" {
	// 	arg = c.next
	// } else if arg == "bmap" {
	// 	arg = c.previous
	// } else {
	// 	log.Fatal("Not valid arg.")
	// }
	locationRes := LocationResult{}
	getURL := baseURL + "/location-area"
	if url != nil {
		getURL = *url
	}

	res, err := http.Get(getURL)
	if err != nil {
		return &locationRes, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return &locationRes, fmt.Errorf("failed with status code %d and\nbody: %s", res.StatusCode, body)
	}

	if err != nil {
		return &locationRes, err
	}

	err = json.Unmarshal(body, &locationRes)
	if err != nil {
		return &locationRes, err
	}
	return &locationRes, nil
}
