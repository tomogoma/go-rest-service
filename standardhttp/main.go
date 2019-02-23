package main

import (
	"encoding/json"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/tomogoma/go-rest-service/model"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
)

var persons = make(map[string]model.Person)

func main() {

	http.HandleFunc("/persons", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleNewPerson(w, r)
		case http.MethodGet:
			handleGetPersons(w, r)
		default:
			http.Error(w, "unknown method for request", http.StatusNotFound)
		}
	})

	addr := ":9000"
	log.Printf("Listening on: %s", addr)
	err := http.ListenAndServe(addr, nil)
	log.Fatalf("Quiting with error: %v", err)
}

// handleGetPersons handles HTTP requests of the form:
//     GET /persons?pageNumber=1&pageSize=300
// 1. extractsPagination in request
// 2. .sorts persons
// 3. responds with subset of persons in page.
func handleGetPersons(w http.ResponseWriter, r *http.Request) {

	page := extractPagination(r)

	pas := personsMapToSlice(persons)
	sort.Slice(pas, func(i, j int) bool {
		return pas[i].Name < pas[j].Name
	})

	startIdx := page.StartIndex()
	if len(pas) <= startIdx {
		log.Printf("No person record found in selected page: %+v", page)
		http.Error(w, fmt.Sprintf("No person record found in selected pageNumber (%d)", page.Number), http.StatusNotFound)
		return
	}

	endPos := int(math.Min(float64(len(pas)), float64(page.EndPosition())))

	respPas := pas[startIdx:endPos]

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(respPas); err != nil {
		log.Printf("Unable to encode response: %v", err)
	}
}

// handleGetPerson handles HTTP requests of the form:
//     GET /persons/{personid}
func handleGetPerson(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This feature is coming soon", http.StatusNotImplemented)
}

// handleGetPerson handles HTTP requests of the form:
//     POST /persons
//         {"name": "John Doe", "Age": 22}
// 1. reads request body into person
// 2. assigns an ID
// 3. validates
// 4. adds person to list of persons
// 5. responds with inserted person.
func handleNewPerson(w http.ResponseWriter, r *http.Request) {
	newPa := model.Person{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read request: %v", err)
		http.Error(w, "Something wicked happened, please try again", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &newPa); err != nil {
		log.Printf("Invalid json in request: %v", err)
		http.Error(w, "Invalid json in request: %v", http.StatusBadRequest)
		return
	}
	newPa.ID = uuid.New()

	if err := newPa.Valid(); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	persons[newPa.ID] = newPa

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newPa); err != nil {
		log.Printf("Unable to encode response: %v", err)
	}
}

// handleUpdatePerson handles HTTP requests of the form:
//     POST /persons/{personID}
//         {"name": "John Doe", "Age": 22}
func handleUpdatePerson(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This feature is coming soon", http.StatusNotImplemented)
}

func extractPagination(r *http.Request) model.Page {

	page := model.Page{}
	var err error

	page.Number, err = requestParamAsInt(r, "pageNumber")
	if err != nil {
		page.Number = 1
	}

	page.Size, err = requestParamAsInt(r, "pageSize")
	if err != nil {
		page.Size = 100
	}

	return page
}

func requestParamAsInt(r *http.Request, key string) (int, error) {
	valStr := r.FormValue(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func personsMapToSlice(mp map[string]model.Person) []model.Person {
	var ps []model.Person
	for _, p := range mp {
		ps = append(ps, p)
	}
	return ps
}
