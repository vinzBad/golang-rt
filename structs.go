package rt

import (
	"net/http"
)

//Tracker is a RequestTracker client
type Tracker struct {
	apiURL     string
	user       string
	password   string
	client     *http.Client
	isLoggedIn bool
	RTVersion  string
}

//Ticket is a RequestTracker Ticket
type Ticket struct {
	ID         int
	Queue      string
	Owner      string
	Creator    string
	Subject    string
	Status     string
	Priority   string
	Requestors []string
	Cc         []string
	AdminCc    []string
}
