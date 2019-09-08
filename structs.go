package rt

import (
	"net/http"
)

//RT is a RequestTracker client
type RT struct {
	apiURL     string
	user       string
	password   string
	client     *http.Client
	isLoggedIn bool
	Version    string
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
