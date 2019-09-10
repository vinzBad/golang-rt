package rt

import (
	"fmt"
	"io/ioutil"

	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// Login to RequestTracker
func (rt *Tracker) Login() error {
	v := url.Values{}
	v.Add("user", rt.user)
	v.Add("pass", rt.password)
	resp, err := rt.client.PostForm(rt.apiURL, v)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	header, err := parseRtResponseHeader(body)
	if err != nil {
		return err
	}
	if header.status != http.StatusOK {
		return fmt.Errorf("Failed to login: %q", header.message)
	}

	rt.RTVersion = header.version
	rt.isLoggedIn = true

	return nil
}

// GetTicket fetches a ticket from RT
func (rt *Tracker) GetTicket(id int) (*Ticket, error) {
	header, body, err := rt.get("ticket/%v/show", id)
	if err != nil {
		return nil, err
	}
	// TODO: Check for 404 and other states
	if header.status != http.StatusOK {
		return nil, fmt.Errorf("Failed to get ticket: %q", header.message)
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

// CreateTicket creates a new ticket in RT and returns its id
func (rt *Tracker) CreateTicket(newTicket Ticket) (int, error) {
	content := make(map[string]string)

	content["ID"] = "ticket/new"
	content["Queue"] = newTicket.Queue
	content["Owner"] = newTicket.Owner
	content["Subject"] = newTicket.Subject
	content["Status"] = newTicket.Status
	content["Priority"] = newTicket.Priority
	content["Requestors"] = strings.Join(newTicket.Requestors, ", ")
	content["Cc"] = strings.Join(newTicket.Cc, ", ")
	content["AdminCc"] = strings.Join(newTicket.AdminCc, ", ")

	c := strings.Builder{}
	for k, v := range content {
		c.WriteString(fmt.Sprintf("%v: %v\n", k, v))
	}

	data := url.Values{}
	data.Add("content", c.String())

	header, body, err := rt.postForm(data, "ticket/new")

	if err != nil {
		return 0, err
	}

	if header.status != http.StatusOK {
		return 0, ErrCredentialsNeeded
	}

	comments, err := parseRTResponseComments(body)

	if err != nil {
		return 0, err
	}

	return comments[0].id, nil
}

// New RequestTracker client
func New(apiURL string, user string, password string) (*Tracker, error) {
	parsedURL, err := url.ParseRequestURI(apiURL)
	if err != nil {
		return nil, ErrInvalidAPIURL
	}
	if !strings.HasSuffix(parsedURL.Path, "/REST/1.0/") {
		return nil, ErrInvalidAPIURL
	}
	jar, _ := cookiejar.New(nil)
	rt := Tracker{
		apiURL:   parsedURL.String(),
		user:     user,
		password: password,
		client: &http.Client{
			Jar: jar,
		},
		isLoggedIn: false,
	}

	return &rt, nil
}
