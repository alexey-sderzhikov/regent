package restapi

type Project struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"descriprion"`
}

type ProjectList struct {
	Projects []Project `json:"projects"`
}

type NameAndId struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	Id          int64     `json:"id"`
	Project     NameAndId `json:"project"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
}

type IssueList struct {
	Issues []Issue `json:"issues"`
}

type TimeEntryInner struct {
	Issue_id int64  `json:"issue_id"`
	Spent_on string `json:"spent_on"`
	Hours    int    `json:"hours"`
	Comments string `json:"comments"`
	User_id  int64  `json:"user_id"`
}
type TimeEntry struct {
	Time_entry TimeEntryInner `json:"time_entry"`
}
type UserInner struct {
	Id            int64  `json:"id"`
	Login         string `json:"login"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Created_on    string `json:"created_on"`
	Last_login_on string `json:"last_login_on"`
	Api_key       string `json:"api_key"`
}

type User struct {
	User UserInner `json:"user"`
}
