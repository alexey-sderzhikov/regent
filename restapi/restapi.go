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

type respStruct struct {
	ByteListBody []byte
	Status       string
}

type TimeEntryParam struct {
	Limit      int
	User_id    int64
	Project_id int64
	Spent_on   string
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

// create request with request type, url, body etc. before send to server
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

// send before created request to server and return respons like bytes slice
func (r RmClient) doRequest(req *http.Request) (respStruct, error) {
	respHttp, err := r.HttpClient.Do(req)
	if err != nil {
		return respStruct{}, err
	}

	defer respHttp.Body.Close()
	resp := respStruct{}
	resp.ByteListBody, err = ioutil.ReadAll(respHttp.Body)
	if err != nil {
		return respStruct{}, err
	}
	resp.Status = respHttp.Status

	return resp, nil
}

// TODO add handling error status codes
func (r RmClient) GetProjects() (ProjectList, error) {
	req, err := r.makeRequest("GET", "/projects.json", nil, nil)
	if err != nil {
		return ProjectList{}, err
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return ProjectList{}, err
	}

	projects := ProjectList{}
	err = json.Unmarshal(resp.ByteListBody, &projects)
	if err != nil {
		return ProjectList{}, err
	}

	return projects, nil
}

// TODO add handling error status codes
func (r RmClient) GetIssues(projectId int64) (IssueList, error) {
	var projectIdParam string
	if projectId != 0 {
		projectIdParam = fmt.Sprintf("&project_id=%v", projectId)
	}

	req, err := r.makeRequest("GET", "/issues.json", []string{projectIdParam}, nil)
	if err != nil {
		return IssueList{}, err
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return IssueList{}, err
	}

	issues := IssueList{}
	err = json.Unmarshal(resp.ByteListBody, &issues)
	if err != nil {
		return IssueList{}, err
	}

	return issues, nil

}

func (r RmClient) CreateTimeEntry(issueId int64, date string, comment string, hours float32) (string, error) {
	timeEntry := TimeEntryRequest{
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
		return "", err
	}

	reqBody := bytes.NewBuffer(byteList)
	req, err := r.makeRequest("POST", "/time_entries.json", nil, reqBody)
	if err != nil {
		return "", err
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Status, nil

}

func (r RmClient) GetTimeEntryList(teparam TimeEntryParam) (TimeEntryListResponse, error) {
	params := make([]string, 0)
	if teparam.Limit != 0 {
		params = append(params, fmt.Sprintf("&limit=%v", teparam.Limit))
	}
	if teparam.Project_id != 0 {
		params = append(params, fmt.Sprintf("&project_id=%v", teparam.Project_id))
	}
	if teparam.User_id != 0 {
		params = append(params, fmt.Sprintf("&user_id=%v", teparam.User_id))
	}

	req, err := r.makeRequest("GET", "/time_entries.json", params, nil)
	if err != nil {
		return TimeEntryListResponse{}, err
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return TimeEntryListResponse{}, err
	}

	timeEntries := TimeEntryListResponse{}
	err = json.Unmarshal(resp.ByteListBody, &timeEntries)
	if err != nil {
		return TimeEntryListResponse{}, err
	}

	return timeEntries, nil
}

// get user data from api key
func (r RmClient) getCurrentUser() (UserInner, error) {
	req, err := r.makeRequest("GET", "/users/current.json", nil, nil)
	if err != nil {
		return UserInner{}, err
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return UserInner{}, err
	}

	userResp := User{}
	err = json.Unmarshal(resp.ByteListBody, &userResp)
	if err != nil {
		return UserInner{}, err
	}

	if userResp.User.Id == 0 {
		return UserInner{}, fmt.Errorf("user can not have user id - %v", userResp.User.Id)
	}

	return userResp.User, nil
}
