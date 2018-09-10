package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-martini/martini"
)

//CrimeDataEntry is the structure for each row of Crime Data
type CrimeDataEntry struct {
	ID                 int
	IncidentNumber     string
	OffenseCode        string
	OffenseCodeGroup   string
	OffenseDescription string
	District           string
	ReportingArea      string
	Shooting           string
	OccurredOnDate     string
	Year               string
	Month              string
	DayOfWeek          string
	Hour               string
	UcrPart            string
	Street             string
	Lat                string
	Long               string
	Location           string
}

//CrimeData is the slice array used to hold all of the Crime Data
type CrimeData struct {
	CrimeDataBook []*CrimeDataEntry
}

//NewCrimeDataBook creates the empty CrimeDataBook Slice
func NewCrimeDataBook() *CrimeData {
	return &CrimeData{
		make([]*CrimeDataEntry, 0),
	}
}

//WebService interfaces
type WebService interface {
	GetPath() string
	WebDelete(params martini.Params) (int, string)
	WebGet(params martini.Params, req *http.Request) (int, string)
	WebPost(params martini.Params, req *http.Request) (int, string)
}

//GetPath is the implmentation of the GetPath Interface
func (c *CrimeData) GetPath() string {
	return "/crimebook"
}

//WebDelete is the implmentation of the WebDelete Interface
func (c *CrimeData) WebDelete(params martini.Params) (int, string) {
	if len(params) == 0 {
		c.RemoveAllEntries()
		return http.StatusOK, "collection deleted"
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return http.StatusBadRequest, "invaild entry id"
	}

	err = c.RemoveEntry(id)
	if err != nil {
		return http.StatusNotFound, "entry not found"
	}

	return http.StatusOK, "entry deleted"
}

//WebGet is the implementation of the WebGet interface
func (c *CrimeData) WebGet(params martini.Params, req *http.Request) (int, string) {
	qString := req.URL.Query()
	if len(qString) > 1 {
		return http.StatusInternalServerError, "only one query parameter allowed"
	} else if len(qString) == 1 {
		jsonResults, err := json.Marshal(c.FilterAllEntries(req))
		if err != nil {
			return http.StatusInternalServerError, "internal error"
		}
		return http.StatusOK, string(jsonResults)
	}

	if len(params) == 0 {
		jsonResults, err := json.Marshal(c.GetAllEntries())
		if err != nil {
			return http.StatusInternalServerError, "internal error"
		}
		return http.StatusOK, string(jsonResults)
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return http.StatusNotFound, "entry not found"
	}

	entry, err := c.GetEntry(id)
	if err != nil {
		return http.StatusInternalServerError, "internal error"
	}

	jsonEntry, err := json.Marshal(entry)
	return http.StatusOK, string(jsonEntry)
}

//WebPost is the implementation of the WebPost Interface
func (c *CrimeData) WebPost(params martini.Params, req *http.Request) (int, string) {
	defer req.Body.Close()

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return http.StatusInternalServerError, "internal error"
	}

	if len(params) != 0 {
		return http.StatusMethodNotAllowed, "method not allowed"
	}

	var crimeDataEntry CrimeDataEntry
	err = json.Unmarshal(requestBody, &crimeDataEntry)
	if err != nil {
		checkError(err)
		return http.StatusBadRequest, "invalid JSON data"
	}

	c.CrimeDataBook = append(c.CrimeDataBook, &crimeDataEntry)

	return http.StatusOK, "new entry created"
}

//GetAllEntries is used by WebGet to get all the entries
func (c *CrimeData) GetAllEntries() []*CrimeDataEntry {
	entries := make([]*CrimeDataEntry, 0)
	for _, entry := range c.CrimeDataBook {
		if entry != nil {
			entries = append(entries, entry)
		}
	}
	return entries
}

//FilterAllEntries is used by WebGet to get all the entries
func (c *CrimeData) FilterAllEntries(req *http.Request) []*CrimeDataEntry {
	IncidentNumber, OffenseCode := req.FormValue("IncidentNumber"), req.FormValue("OffenseCode")
	District, OffenseCodeGroup := req.FormValue("District"), req.FormValue("OffenseCodeGroup")

	entries := make([]*CrimeDataEntry, 0)
	for _, entry := range c.CrimeDataBook {
		if entry != nil {
			if IncidentNumber == entry.IncidentNumber || OffenseCode == entry.OffenseCode || District == entry.District || OffenseCodeGroup == entry.OffenseCodeGroup {
				entries = append(entries, entry)
			}
		}
	}

	return entries
}

//GetEntry is used be WebGet to get a single Entry
func (c *CrimeData) GetEntry(id int) (*CrimeDataEntry, error) {
	if id < 0 || id >= len(c.CrimeDataBook) {
		return nil, fmt.Errorf("invalid id")
	}
	return c.CrimeDataBook[id], nil
}

//RemoveAllEntries is used by WebDelete to delete all the entries
func (c *CrimeData) RemoveAllEntries() {
	c.CrimeDataBook = []*CrimeDataEntry{}
}

//RemoveEntry is used by WebDelete to remove a single entry
func (c *CrimeData) RemoveEntry(id int) error {
	if id < 0 || id >= len(c.CrimeDataBook) {
		return fmt.Errorf("invalid id")
	}

	c.CrimeDataBook[id] = nil

	return nil
}

//RegisterWebService is used to register the web services
func RegisterWebService(webService WebService, classicMartini *martini.ClassicMartini) {
	path := webService.GetPath()

	classicMartini.Get(path, webService.WebGet)
	classicMartini.Get(path+"/:id", webService.WebGet)

	classicMartini.Post(path, webService.WebPost)
	classicMartini.Post(path+"/:id", webService.WebPost)

	classicMartini.Delete(path, webService.WebDelete)
	classicMartini.Delete(path+"/:id", webService.WebDelete)

}

func main() {
	csvFile, err := os.Open("crime10.csv")
	checkError(err)

	crimeBook := NewCrimeDataBook()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	//Throw away header
	_, err = reader.Read()
	checkError(err)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			checkError(err)
		}
		crimeBook.AddEntry(line)

	}

	martiniClassic := martini.Classic()
	RegisterWebService(crimeBook, martiniClassic)
	martiniClassic.Run()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//AddEntry is used to add an entry to the Crime Data Book
func (c *CrimeData) AddEntry(line []string) {
	ID := len(c.CrimeDataBook)

	newEntry := &CrimeDataEntry{
		ID,
		line[0],
		line[1],
		line[2],
		line[3],
		line[4],
		line[5],
		line[6],
		line[7],
		line[8],
		line[9],
		line[10],
		line[11],
		line[12],
		line[13],
		line[14],
		line[15],
		line[16],
	}

	c.CrimeDataBook = append(c.CrimeDataBook, newEntry)
}
