package tiss

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

///
// This package searches for tiss rooms by their name. 2 different functions are available: 
//  - GetTissRoom: quick lookup from provided csv file. Much, much faster, but could get out of date
//  - SearchTissRoom: makes web-requests to get the current information. Very slow and will probably break later...
///

type room struct {
	name              string
	capacity          int
	approvalInstitute string
	approvalMode      string
	availability      string
	building          string
	address           string
	roomNumber        string
	tissLink          string
}

const csvComma = ';'

//go:embed rooms.csv
var roomsContent string
var rooms []room = parseRooms(roomsContent)

type RoomNotFoundError struct {
	Inner error
}

func (e *RoomNotFoundError) Error() string {
	return fmt.Sprintf("tiss room not found: %v", e.Inner)
}

func (e *RoomNotFoundError) Unwrap() error {
	return e.Inner
}

func GetTissRoom(icsName string) (string, error) {
	for _, room := range rooms {
		if room.name == icsName {
			return room.tissLink, nil
		}
	}

	return "", &RoomNotFoundError{fmt.Errorf("couldn't find room in csv-cache")}
}

func SearchTISSRoom(icsName string) (string, error) {
	link, err := searchTISSRoom(icsName)
	if err != nil {
		return "", &RoomNotFoundError{err}
	}
	return link, nil
}

// Do not question this function, please
func searchTISSRoom(icsName string) (string, error) {
	// Constants
	dswid := 6122
	dsrid := 236

	// Create HTTP client with cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{
		Jar: jar,
	}

	// Set initial cookies
	cookieURL, err := url.Parse("https://tiss.tuwien.ac.at")
	if err != nil || cookieURL == nil {
		return "", err
	}
	cookie := &http.Cookie{
		Name:   fmt.Sprintf("dsrwid-%d", dsrid),
		Value:  fmt.Sprintf("%d", dswid),
		Domain: cookieURL.Host,
		Path:   "/",
	}
	client.Jar.SetCookies(cookieURL, []*http.Cookie{cookie})

	// First GET request
	firstURL := fmt.Sprintf("https://tiss.tuwien.ac.at/events/selectRoom.xhtml?dswid=%d&dsrid=%d", dswid, dsrid)
	resp, err := client.Get(firstURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Extract ViewState
	viewStateRegex := regexp.MustCompile(`ViewState:2" value="(.*?)"`)
	body := new(strings.Builder)
	_, err = io.Copy(body, resp.Body)
	if err != nil {
		return "", err
	}
	viewStateMatch := viewStateRegex.FindStringSubmatch(body.String())
	if len(viewStateMatch) != 2 {
		return "", err
	}
	viewState := viewStateMatch[1]

	// Encode search string

	// POST request data
	form := url.Values{}
	form.Add("filterForm:quickSearchTerm", icsName)
	form.Add("filterForm:roomFilter:minimumCapacity", "0")
	form.Add("filterForm:roomFilter:selectBuildingLb", "")
	form.Add("filterForm:roomFilter:selectRoomLb", "")
	form.Add("filterForm_SUBMIT", "1")
	form.Add("jakarta.faces.ViewState", viewState)
	form.Add("jakarta.faces.ClientWindow", "6122")
	form.Add("jakarta.faces.behavior.event", "action")
	form.Add("jakarta.faces.partial.event", "click")
	form.Add("jakarta.faces.source", "filterForm:quickSearchButton")
	form.Add("jakarta.faces.partial.ajax", "true")
	form.Add("filterForm", "filterForm")
	form.Add("jakarta.faces.partial.execute", "filterForm:quickSearchButton filterForm:quickSearchTerm")
	form.Add("jakarta.faces.partial.render", "tableForm")
	form.Add("filterForm:quickSearchButton", "Suchen")

	// Prepare POST request
	postURL := fmt.Sprintf("https://tiss.tuwien.ac.at/events/selectRoom.xhtml?dswid=%d&dsrid=%d", dswid, dsrid)
	req, err := http.NewRequest("POST", postURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.PostForm = form

	// Execute POST request
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Extract relative link
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	relLinkRegex := regexp.MustCompile(fmt.Sprintf(`<a href="(.*?)">(.*?)(%s)(.*?)</a>`, regexp.QuoteMeta(icsName)))
	relLinkMatch := relLinkRegex.FindStringSubmatch(string(bodyBytes))
	if len(relLinkMatch) != 5 {
		return "", err
	}
	relLink := relLinkMatch[1]

	return fmt.Sprintf("https://tiss.tuwien.ac.at%s", relLink), nil
}

// This function reads from the embedded csv-string. Since this is provided at built-time it is assumed that it is a valid input.
// Therefore, this function panics when an error occurs
func parseRooms(roomsContent string) []room {
	strReader := strings.NewReader(roomsContent)
	csvReader := csv.NewReader(strReader)
	csvReader.Comma = csvComma
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("Failed to parse roomsContent csv: %w", err))
	}

	rooms := make([]room, 0, len(records))
	for i, record := range records {
		capacity, err := strconv.Atoi(record[1])
		if err != nil {
			panic(fmt.Errorf("Failed to parse room capacity at line %d: %w", i+1, err))
		}

		room := room{
			name:              record[0],
			capacity:          capacity,
			approvalInstitute: record[2],
			approvalMode:      record[3],
			availability:      record[4],
			building:          record[5],
			address:           record[6],
			roomNumber:        record[7],
			tissLink:          record[8],
		}

		rooms = append(rooms, room)
	}

	return rooms
}
