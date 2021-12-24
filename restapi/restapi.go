package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const BERGEN_URL = "https://support.bergen.tech"
const USER_API_KEY = "c370a381d4bc709c419094f8a63f78b64f7a1b56"

type RmClient struct {
	SourceUrl  string
	ApiKey     string
	User       UserInner
	HttpClient *http.Client
}

func NewRm(source string, apiKey string) (*RmClient, error) {
	r := &RmClient{}

	if source == "" {
		r.SourceUrl = BERGEN_URL
	} else {
		r.SourceUrl = source
	}

	if apiKey == "" {
		r.ApiKey = USER_API_KEY
	} else {
		r.ApiKey = apiKey
	}

	r.HttpClient = &http.Client{}

	var err error
	r.User, err = r.getCurrentUser()
	if err != nil {
		return &RmClient{}, err
	}

	return r, nil
}

func (r RmClient) makeRequest(reqType string, endPoint string, params []string, body io.Reader) (*http.Request, error) {
	url := r.SourceUrl + endPoint + "?key=" + r.ApiKey

	for _, p := range params {
		url += p
	}

	req, err := http.NewRequest(reqType, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (r RmClient) doRequest(req *http.Request) ([]byte, error) {
	resp, err := r.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	byteList, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return byteList, nil
}

func (r RmClient) GetProjects() (ProjectList, error) {
	req, err := r.makeRequest("GET", "/projects.json", nil, nil)
	if err != nil {
		return ProjectList{}, err
	}

	byteList, err := r.doRequest(req)
	if err != nil {
		return ProjectList{}, err
	}

	projects := ProjectList{}
	err = json.Unmarshal(byteList, &projects)
	if err != nil {
		return ProjectList{}, err
	}

	return projects, nil
}

func (r RmClient) GetIssues(projectId int64) (IssueList, error) {
	var projectIdParam string
	if projectId != 0 {
		projectIdParam = fmt.Sprintf("&project_id=%v", projectId)
	}

	req, err := r.makeRequest("GET", "/issues.json", []string{projectIdParam}, nil)
	if err != nil {
		return IssueList{}, err
	}

	byteList, err := r.doRequest(req)
	if err != nil {
		return IssueList{}, err
	}

	issues := IssueList{}
	err = json.Unmarshal(byteList, &issues)
	if err != nil {
		return IssueList{}, err
	}

	return issues, nil

}

func (r RmClient) CreateTimeEntry(issueId int64, date string, comment string, hours int) error {
	timeEntry := TimeEntry{
		Time_entry: TimeEntryInner{
			Issue_id: issueId,
			Spent_on: date,
			Hours:    hours,
			Comments: comment,
			User_id:  r.User.Id,
		},
	}

	byteList, err := json.Marshal(timeEntry)
	if err != nil {
		return err
	}

	reqBody := bytes.NewBuffer(byteList)
	req, err := r.makeRequest("POST", "/time_entries.json", nil, reqBody)
	fmt.Print(req)
	if err != nil {
		return err
	}

	byteList, err = r.doRequest(req)
	if err != nil {
		return err
	}

	fmt.Print(string(byteList))
	return nil

}

func (r RmClient) getCurrentUser() (UserInner, error) {
	req, err := r.makeRequest("GET", "/users/current.json", nil, nil)
	if err != nil {
		return UserInner{}, err
	}

	byteList, err := r.doRequest(req)
	if err != nil {
		return UserInner{}, err
	}

	userResp := User{}
	err = json.Unmarshal(byteList, &userResp)
	if err != nil {
		return UserInner{}, nil
	}

	if userResp.User.Id == 0 {
		return UserInner{}, fmt.Errorf("user can not have user id - %v", userResp.User.Id)
	}

	return userResp.User, nil
}
