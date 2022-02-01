package restapi

type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Description string `json:"descriprion"`
}

type ProjectList struct {
	Projects []Project `json:"projects"`
}

type NameAndID struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

//FIXME redundant struct?
//type ID struct {
//	ID int64 `json:"id"`
//}

type Issue struct {
	ID          int64     `json:"id"`
	Project     NameAndID `json:"project"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
}

type IssueList struct {
	Issues     []Issue `json:"issues"`
	TotalCount int     `json:"total_count"`
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
	ProjectID  int64   `json:"project_id"`
}

type TimeEntryInner struct {
	IssueID  int64   `json:"issue_id"`
	SpentOn  string  `json:"spent_on"`
	Hours    float32 `json:"hours"`
	Comments string  `json:"comments"`
	UserID   int64   `json:"user_id"`
}

type TimeEntryRequest struct {
	TimeEntry TimeEntryInner `json:"time_entry"`
}

type TimeEntryResponse struct {
	ID        int64     `json:"id"`
	Project   NameAndID `json:"project"`
	Issue     int64     `json:"issue"`
	User      NameAndID `json:"user"`
	Activity  NameAndID `json:"activity"`
	Hours     float32   `json:"hours"`
	Comments  string    `json:"comments"`
	SpentOn   string    `json:"spent_on"`
	CreatedOn string    `json:"created_on"`
	UpdatedOn string    `json:"updated_on"`
}

type TimeEntryListResponse struct {
	TimeEntries []TimeEntryResponse `json:"time_entries"`
	TotalCount  int                 `json:"total_count"`
	Offset      int                 `json:"offset"`
	Limit       int                 `json:"limit"`
}

type UserInner struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	CreatedOn   string `json:"created_on"`
	LastLoginOn string `json:"last_login_on"`
	APIKey      string `json:"api_key"`
}

type User struct {
	User UserInner `json:"user"`
}
