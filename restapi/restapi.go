package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type RmClient struct {
	SourceURL  string
	APIKey     string
	User       UserInner
	HTTPClient *http.Client
}

type respStruct struct {
	ByteListBody []byte
	Status       string
}

type Params map[string]interface{}

type TimeEntryParam struct {
	Limit     int
	UserID    int64
	ProjectID int64
	SpentOn   string
}

func NewRm(source string, apiKey string) (*RmClient, error) {
	r := &RmClient{}

	r.SourceURL = source

	r.APIKey = apiKey

	r.HTTPClient = &http.Client{}

	var err error
	r.User, err = r.getCurrentUser()
	if err != nil {
		return &RmClient{}, err
	}

	return r, nil
}

func (p Params) makeRequestParameters() string {
	var params string
	for key, value := range p {
		params += "&" + key + "=" + fmt.Sprintf("%v", value)
	}

	return params
}

// create request with request type, url, body etc. before send to server
func (r RmClient) makeRequest(reqType string, endPoint string, params string, body io.Reader) (*http.Request, error) {
	url := r.SourceURL + endPoint + "?key=" + r.APIKey + params

	req, err := http.NewRequest(reqType, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// send before created request to server and return respons like bytes slice
func (r RmClient) doRequest(req *http.Request) (respStruct, error) {
	respHTTP, err := r.HTTPClient.Do(req)
	if err != nil {
		return respStruct{}, err
	}

	defer respHTTP.Body.Close()

	resp := respStruct{}
	resp.ByteListBody, err = ioutil.ReadAll(respHTTP.Body)
	if err != nil {
		return respStruct{}, err
	}
	if respHTTP.StatusCode < 200 || respHTTP.StatusCode > 299 {
		return respStruct{}, fmt.Errorf("status code not in 2xx range, url-%+v", req.URL)
	}
	resp.Status = respHTTP.Status

	return resp, nil
}

// TODO add handling error status codes
func (r RmClient) GetProjects() (ProjectList, error) {
	req, err := r.makeRequest("GET", "/projects.json", "", nil)
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
func (r RmClient) GetIssues(params Params) (IssueList, error) {
	// TODO extract struct with parameters for request, like GetTimeEntryList
	p := params.makeRequestParameters()
	req, err := r.makeRequest("GET", "/issues.json", p, nil)
	if err != nil {
		return IssueList{}, fmt.Errorf("error occured during creating request - %q", err)
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return IssueList{}, fmt.Errorf("error occured during do request\n %q", err)
	}

	issues := IssueList{}
	err = json.Unmarshal(resp.ByteListBody, &issues)
	if err != nil {
		return IssueList{}, fmt.Errorf("error occured during unmurshaling response from redmine server - %q\nResponse structure:\n%+v", err, resp)
	}

	var ok bool
	issues.ProjectID, ok = params["project_id"].(int64)
	if !ok {
		return IssueList{}, fmt.Errorf("error occured during convert %v (project id) to int64", params["project_id"])
	}

	return issues, nil
}

// TODO refactor params like GetTimeEntryList
func (r RmClient) CreateTimeEntry(issueID int64, date string, comment string, hours float32) (string, error) {
	timeEntry := TimeEntryRequest{
		TimeEntry: TimeEntryInner{
			IssueID:  issueID,
			SpentOn:  date,
			Hours:    hours,
			Comments: comment,
			UserID:   r.User.ID,
		},
	}

	byteList, err := json.Marshal(timeEntry)
	if err != nil {
		return "", err
	}

	reqBody := bytes.NewBuffer(byteList)
	req, err := r.makeRequest("POST", "/time_entries.json", "", reqBody)
	if err != nil {
		return "", err
	}

	resp, err := r.doRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Status, nil

}

func (r RmClient) GetTimeEntryList(params Params) (TimeEntryListResponse, error) {
	p := params.makeRequestParameters()

	req, err := r.makeRequest("GET", "/time_entries.json", p, nil)
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
	req, err := r.makeRequest("GET", "/users/current.json", "", nil)
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

	if userResp.User.ID == 0 {
		return UserInner{}, fmt.Errorf("user can not have user id - %v", userResp.User.ID)
	}

	return userResp.User, nil
}
