package rt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
)

var (
	reStatusCode      = regexp.MustCompile(`^RT\/([\d\.]+) (\d\d\d) (.+)`)
	reResponseKV      = regexp.MustCompile(`(\w+):( (.*))?`)
	reResponseComment = regexp.MustCompile(`# (\w+) (\d+) (.*)\.`)
)

type rtResponseHeader struct {
	version string
	status  int
	message string
}

type rtResponseComment struct {
	object string
	id     int
	verb   string
}

func (rt *Tracker) get(path string, a ...interface{}) (*rtResponseHeader, []byte, error) {
	resp, err := rt.client.Get(fmt.Sprintf(rt.apiURL+path, a...))
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	header, err := parseRtResponseHeader(body)
	if err != nil {
		return nil, nil, err
	}
	return header, body, nil
}

func (rt *Tracker) post(data *bytes.Buffer, path string, a ...interface{}) (*rtResponseHeader, []byte, error) {
	resp, err := rt.client.Post(fmt.Sprintf(rt.apiURL+path, a...), "text/plain", data)
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	header, err := parseRtResponseHeader(body)
	if err != nil {
		return nil, nil, err
	}
	return header, body, nil
}

func (rt *Tracker) postForm(data url.Values, path string, a ...interface{}) (*rtResponseHeader, []byte, error) {
	resp, err := rt.client.PostForm(fmt.Sprintf(rt.apiURL+path, a...), data)
	if err != nil {
		return nil, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	header, err := parseRtResponseHeader(body)
	if err != nil {
		return nil, nil, err
	}
	return header, body, nil
}

func parseRtResponseHeader(message []byte) (*rtResponseHeader, error) {
	match := reStatusCode.FindSubmatch(message)
	if match == nil {
		return nil, ErrParseRTMessageError
	}
	status, err := strconv.Atoi(string(match[2]))
	if err != nil {
		return nil, ErrParseRTMessageError
	}
	return &rtResponseHeader{
		version: string(match[1]),
		status:  status,
		message: string(match[3]),
	}, nil
}

func parseRTResponseKVs(message []byte) (result map[string]string, err error) {
	matches := reResponseKV.FindAllSubmatch(message, -1)

	if matches == nil || len(matches) == 0 {
		return nil, ErrParseRTMessageError
	}
	result = make(map[string]string)
	for _, match := range matches {
		result[string(match[1])] = string(match[3])
	}

	return result, nil
}

func parseRTResponseComments(message []byte) ([]rtResponseComment, error) {
	matches := reResponseComment.FindAllSubmatch(message, -1)

	if matches == nil || len(matches) == 0 {
		return nil, ErrParseRTMessageError
	}

	result := make([]rtResponseComment, len(matches))

	for i, match := range matches {
		id, err := strconv.Atoi(string(match[2]))
		if err != nil {
			return nil, ErrParseRTMessageError
		}
		result[i] = rtResponseComment{
			object: string(match[1]),
			id:     id,
			verb:   string(match[3]),
		}
	}

	return result, nil
}
