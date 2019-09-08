package rt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
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

// Login to RequestTracker
func (rt *RT) Login() error {
	v := url.Values{}
	v.Add("user", rt.user)
	v.Add("pass", rt.password)
	resp, err := rt.client.PostForm(rt.url.String(), v)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	version, status, message, err := parseRtResponseHeader(body)
	if err != nil {
		return err
	}
	if status != "200" {
		return fmt.Errorf("Failed to authorize: %q", message)
	}

	rt.Version = version
	rt.isLoggedIn = true

	return nil
}

// GetTicket fetches a ticket from RT
func (rt *RT) GetTicket(id int) (*Ticket, error) {
	resp, err := rt.client.Get(rt.url.String() + "ticket/" + strconv.Itoa(id) + "/show")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_, status, message, err := parseRtResponseHeader(body)
	if err != nil {
		return nil, err
	}
	// TODO: Check for 404 and other states
	if status != "200" {
		return nil, fmt.Errorf("Failed to authorize: %q", message)
	}

	result, err := parseRTResponseKVs(body)
	if err != nil {
		return nil, err
	}

	return &Ticket{
		ID:         id,
		Queue:      result["Queue"],
		Creator:    result["Creator"],
		Subject:    result["Subject"],
		Status:     result["Status"],
		Priority:   result["Priority"],
		Requestors: strings.Split(result["Requestors"], ", "),
		Cc:         strings.Split(result["Cc"], ", "),
		AdminCc:    strings.Split(result["AdminCc"], ", "),
	}, nil
}

// New RequestTracker client
func New(apiURL string, user string, password string) (*RT, error) {
	parsedURL, err := url.ParseRequestURI(apiURL)
	if err != nil {
		return nil, ErrInvalidAPIURL
	}
	if !strings.HasSuffix(parsedURL.Path, "/REST/1.0/") {
		return nil, ErrInvalidAPIURL
	}
	jar, _ := cookiejar.New(nil)
	rt := RT{
		url:      parsedURL,
		user:     user,
		password: password,
		client: &http.Client{
			Jar: jar,
		},
		isLoggedIn: false,
	}

	return &rt, nil
}
