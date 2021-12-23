package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BERGEN_URL = "https://support.bergen.tech"
const USER_API_KEY = "c370a381d4bc709c419094f8a63f78b64f7a1b56"

func GetProjects() (ProjectList, error) {
	req, err := http.NewRequest("GET", BERGEN_URL+"/projects.json"+"?key="+USER_API_KEY, nil)
	if err != nil {
		return ProjectList{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return ProjectList{}, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ProjectList{}, err
	}

	projects := ProjectList{}
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return ProjectList{}, err
	}

	return projects, nil
}

func GetIssues(projectId int64) (IssueList, error) {
	var projectIdParam string
	if projectId != 0 {
		projectIdParam = fmt.Sprintf("&project_id=%v", projectId)
	}
	req, err := http.NewRequest("GET", BERGEN_URL+"/issues.json"+"?key="+USER_API_KEY+projectIdParam, nil)
	if err != nil {
		return IssueList{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return IssueList{}, err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IssueList{}, err
	}

	issues := IssueList{}
	err = json.Unmarshal(bytes, &issues)
	if err != nil {
		return IssueList{}, err
	}

	return issues, nil

}
