package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/franela/goblin"
	"github.com/pauldambra/filesaver/api"
)

var (
	server    *httptest.Server
	reader    io.Reader //Ignore this for now
	thingsURL string
)

func init() {
	server = httptest.NewServer(api.Handlers()) //Creating new server with the user handlers

	thingsURL = fmt.Sprintf("%s/things", server.URL) //Grab the address for the API endpoint
}

func createThing(json string) (resp *http.Response, err error) {

	reader = strings.NewReader(json) //Convert string to reader

	request, err := http.NewRequest("POST", thingsURL, reader) //Create request with JSON body

	return http.DefaultClient.Do(request)
}

func TestCreateThing(t *testing.T) {
	g := Goblin(t)
	g.Describe("when creating a thing", func() {
		g.It("should return 201", func() {
			thingJSON := `{"name": "dennis", "area": 200}`
			res, err := createThing(thingJSON)

			if err != nil {
				g.Fail(err)
			}
			g.Assert(res.StatusCode).Equal(201)
		})
	})
}

func contains(name string, list []api.ThingParams) bool {
	for _, b := range list {
		if b.Name == name {
			return true
		}
	}
	return false
}

func TestGetThing(t *testing.T) {
	g := Goblin(t)
	g.Describe("when there is at least one thing", func() {
		g.It("should return at least that thing in the list", func() {
			_, err := createThing(`{"name": "definitely", "area": 200}`)

			request, err := http.NewRequest("GET", thingsURL, strings.NewReader(""))
			res, err := http.DefaultClient.Do(request)

			if err != nil {
				g.Fail(err)
			}
			g.Assert(res.StatusCode).Equal(200)
			defer res.Body.Close()
			contents, err := ioutil.ReadAll(res.Body)

			p := []api.ThingParams{}
			err = json.Unmarshal(contents, &p)
			if err != nil {
				g.Fail(err)
			}
			fmt.Printf("contents: %s", string(contents))
			g.Assert(contains("definitely", p)).Equal(true)
		})
	})
}
