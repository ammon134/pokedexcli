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

type LocationPage struct {
	Previous string
	Next     string
	Results  []location
}

type location struct {
	name string
	url  string
}

func GetLocations(p *LocationPage, arg string) (*LocationPage, error) {
	// call GET here
	// if arg == "map" {
	// 	arg = c.next
	// } else if arg == "bmap" {
	// 	arg = c.previous
	// } else {
	// 	log.Fatal("Not valid arg.")
	// }
	zeroLocationPage := LocationPage{}

	switch arg {
	case "map":
		arg = p.Next
	case "bmap":
		arg = p.Previous
	default:
		return p, nil

	}
	res, err := http.Get(arg)
	if err != nil {
		return &zeroLocationPage, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return &zeroLocationPage, fmt.Errorf("failed with status code %d and\nbody: %s", res.StatusCode, body)
	}

	if err != nil {
		return &zeroLocationPage, err
	}

	err = json.Unmarshal(body, p)
	if err != nil {
		return &zeroLocationPage, err
	}
	return p, nil
}
