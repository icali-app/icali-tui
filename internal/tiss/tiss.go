package tiss

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

// Do not question this function, please
func SearchTISSRoom(icsName string) (string, error) {
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
        Name: fmt.Sprintf("dsrwid-%d", dsrid),
        Value: fmt.Sprintf("%d", dswid),
        Domain: cookieURL.Host,
        Path: "/",
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

