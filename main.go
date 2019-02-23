package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/tomogoma/go-rest-service/model"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var persons = make(map[string]model.Person)

func main() {

	r := mux.NewRouter()

	r.PathPrefix("/persons/{id}").Methods(http.MethodGet).HandlerFunc(HandleGetPerson)
	r.PathPrefix("/persons/{id}").Methods(http.MethodPut).HandlerFunc(HandleUpdatePerson)
	r.PathPrefix("/persons").Methods(http.MethodGet).HandlerFunc(HandleGetPersons)
	r.PathPrefix("/persons").Methods(http.MethodPost).HandlerFunc(HandleCreatePerson)

	addr := ":9000"
	log.Printf("Server is listening at %s", addr)
	err := http.ListenAndServe(addr, r)
	log.Printf("Server stopped with error: %v", err)

}

func HandleGetPersons(w http.ResponseWriter, r *http.Request) {

	pNmbrStr := r.FormValue("pageNumber")
	pSizeStr := r.FormValue("pageSize")

	pNmbr, err := strconv.Atoi(pNmbrStr)
	if err != nil {
		pNmbr = 1
	}

	pSize, err := strconv.Atoi(pSizeStr)
	if err != nil {
		pNmbr = 100
	}

	page := model.Page{Number: pNmbr, Size: pSize}

	var personsSlc []model.Person
	for _, person := range persons {
		personsSlc = append(personsSlc, person)
	}

	sort.SliceStable(personsSlc, func(i, j int) bool {
		return personsSlc[i].Age > personsSlc[j].Age
	})

	endPosn := len(personsSlc)
	if len(personsSlc) > page.EndPosition() {
		endPosn = page.EndPosition()
	}

	startIdx := len(personsSlc)
	if len(personsSlc) > page.StartIndex() {
		startIdx = page.StartIndex()
	}

	personsRslt := personsSlc[startIdx:endPosn]

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	err = json.NewEncoder(w).Encode(personsRslt)
	if err != nil {
		log.Printf("Error encoding results: %v", err)
	}
}

func HandleGetPerson(w http.ResponseWriter, r *http.Request) {

	personId := mux.Vars(r)["id"]

	person, exist := persons[personId]
	if !exist {
		log.Printf("Person with id %s does not exist", personId)
		http.Error(w, "Person with id %s does not exist", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(person)
	if err != nil {
		log.Printf("Error encoding results: %v", err)
	}
}

func HandleCreatePerson(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(newPa); err != nil {
		log.Printf("Unable to encode response: %v", err)
	}
}

func HandleUpdatePerson(w http.ResponseWriter, r *http.Request) {

	personId := mux.Vars(r)["id"]

	oldPerson, exist := persons[personId]
	if !exist {
		log.Printf("Person with id %s does not exist", personId)
		http.Error(w, "Person with id %s does not exist", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read request: %v", err)
		http.Error(w, "Something wicked happened, please try again", http.StatusInternalServerError)
		return
	}

	newPa := model.Person{}
	if err := json.Unmarshal(body, &newPa); err != nil {
		log.Printf("Invalid json in request: %v", err)
		http.Error(w, "Invalid json in request: %v", http.StatusBadRequest)
		return
	}
	newPa.ID = oldPerson.ID

	if err := newPa.Valid(); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	persons[newPa.ID] = newPa

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(newPa); err != nil {
		log.Printf("Unable to encode response: %v", err)
	}

}
