package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var thingIDCounter uint32

var thingStore = []Thing{}

// Thing is some dto?
type Thing struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	Area uint32 `json:"area"`
}

// ThingParams is what you post?
type ThingParams struct {
	Name string `json:"name"`
	Area uint32 `json:"area"`
}

// Handlers is the router?
func Handlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/things", createThingHandler).Methods("POST")

	r.HandleFunc("/things", listThingsHandler).Methods("GET")

	return r
}

func createThingHandler(w http.ResponseWriter, r *http.Request) {
	p := ThingParams{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = validateUniqueness(p.Name)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := Thing{
		ID:   thingIDCounter,
		Name: p.Name,
		Area: p.Area,
	}
	toSave, _ := json.Marshal(s)

	err = ioutil.WriteFile(fmt.Sprintf("./things/%s.json", s.Name), []byte(toSave), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thingIDCounter++

	w.WriteHeader(http.StatusCreated)
}

func allThingNames() []string {
	files, _ := ioutil.ReadDir("./things")
	fileNames := []string{}

	for _, f := range files {
		fileNames = append(fileNames, strings.Replace(f.Name(), ".json", "", -1))
	}
	return fileNames
}

func allThingContents() []string {
	contents := []string{}
	for _, f := range allThingNames() {
		b, err := ioutil.ReadFile(fmt.Sprintf("./things/%s.json", f))
		if err != nil {
			fmt.Printf("error: %s", err)
		} else {
			contents = append(contents, string(b))
		}
	}
	return contents
}

func validateUniqueness(name string) error {
	files := allThingNames()
	for _, f := range files {
		if f == name {
			return fmt.Errorf("Name %s is already used", f)
		}
	}

	return nil
}

func listThingsHandler(w http.ResponseWriter, r *http.Request) {
	things, err := json.Marshal(allThingContents())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(things)
}
