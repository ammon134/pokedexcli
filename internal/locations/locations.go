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
	"io"
	"log"
	"net/http"
)

type CurrentPage struct {
	Previous string
	Next     string
	Results  []location
}

type location struct {
	name string
	url  string
}

func get_locations(c CurrentPage, arg string) (cp CurrentPage) {
	// call GET here
	// if arg == "map" {
	// 	arg = c.next
	// } else if arg == "bmap" {
	// 	arg = c.previous
	// } else {
	// 	log.Fatal("Not valid arg.")
	// }

	switch arg {
	case "map":
		arg = c.Next
	case "bmap":
		arg = c.Previous
	default:
		log.Fatal("Not valid arg")

	}
	res, err := http.Get(arg)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Failed with status code %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatal(err)
	}

	cp = CurrentPage{}
	err = json.Unmarshal(body, &cp)
	if err != nil {
		log.Fatal(err)
	}
	return
}
