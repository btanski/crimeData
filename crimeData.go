package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"

	//"encoding/json"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-martini/martini"
)

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

type CrimeData struct {
	CrimeDataBook []*CrimeDataEntry
}

func NewCrimeDataBook() *CrimeData {
	return &CrimeData{
		make([]*CrimeDataEntry, 0),
	}
}

type WebService interface {
	GetPath() string
	WebDelete(params martini.Params) (int, string)
	WebGet(params martini.Params) (int, string)
	WebPost(params martini.Params, req *http.Request) (int, string)
}

func (c CrimeData) GetPath() string {
	return "/crimebook"
}

func (c CrimeData) WebDelete(params martini.Params) (int, string) {
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

func (c CrimeData) WebGet(params martini.Params) (int, string) {
	/*fmt.Println(len(params))
	fmt.Println(params)
	return http.StatusOK, "OK Get"*/
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

func (c CrimeData) WebPost(params martini.Params, req *http.Request) (int, string) {
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

	fmt.Println(requestBody)
	bodyString := string(requestBody)
	fmt.Println(bodyString)
	fmt.Println()
	fmt.Println(crimeDataEntry)

	/*newEntry := &CrimeDataEntry{
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
	}*/

	fmt.Println(len(c.CrimeDataBook))
	c.CrimeDataBook = append(c.CrimeDataBook, &crimeDataEntry)
	fmt.Println(len(c.CrimeDataBook))

	return http.StatusOK, "OK Post"
}
func (c CrimeData) GetAllEntries() []*CrimeDataEntry {
	entries := make([]*CrimeDataEntry, 0)
	for _, entry := range c.CrimeDataBook {
		if entry != nil {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (c CrimeData) GetEntry(id int) (*CrimeDataEntry, error) {
	if id < 0 || id >= len(c.CrimeDataBook) {
		return nil, fmt.Errorf("invalid id")
	}
	return c.CrimeDataBook[id], nil
}

func (c CrimeData) RemoveAllEntries() {
	fmt.Println(len(c.CrimeDataBook))
	c.CrimeDataBook = []*CrimeDataEntry{}
	fmt.Println(len(c.CrimeDataBook))
}

func (c CrimeData) RemoveEntry(id int) error {
	if id < 0 || id >= len(c.CrimeDataBook) {
		return fmt.Errorf("invalid id")
	}

	c.CrimeDataBook[id] = nil

	return nil
}

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
	fmt.Println(len(crimeBook.CrimeDataBook))
	fmt.Println(crimeBook.CrimeDataBook[len(crimeBook.CrimeDataBook)-1])

	martiniClassic := martini.Classic()
	RegisterWebService(crimeBook, martiniClassic)
	martiniClassic.Run()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

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
