package restapi

type Project struct {
	Id          int64
	Name        string
	Identifier  string
	Description string
	Created_on  string
	Updated_on  string
	Is_public   bool
}

type ProjectList struct {
	Projects []Project
}

type NameAndId struct {
	Id   int64
	Name string
}

type Issue struct {
	Id              int64
	Project         NameAndId
	Tracker         NameAndId
	Status          NameAndId
	Priority        NameAndId
	Author          NameAndId
	Category        NameAndId
	Subject         string
	Description     string
	Start_date      string
	Due_date        string
	Done_ratio      int
	Estimated_hours float32
	Is_private      bool
	Tags            []string
	Custom_fields   []struct {
		Id    int64
		Name  string
		Value string
	}
	Created_on string
	Updated_on string
	Closed_on  string
}

type IssueList struct {
	Issues []Issue
}
