package rt

import (
	"net/http"
	"net/url"
)

//RT is a RequestTracker client
type RT struct {
	url        *url.URL
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
