package rt

import (
	"regexp"
	"strconv"
)

var (
	reStatusCode = regexp.MustCompile(`^RT\/([\d\.]+) (\d\d\d) (.+)`)
	reResponseKV = regexp.MustCompile(`(\w+):( (.*))?`)
)

type rtResponseHeader struct {
	version string
	status  int
	message string
}



func parseRtResponseHeader(message []byte) (*rtResponseHeader, error) {
	match := reStatusCode.FindSubmatch(message)
	if match == nil {
		return nil, ErrParseRTMessageError
	}
	status, err :=  strconv.Atoi(string(match[2]))
	if err != nil {
		return nil, ErrParseRTMessageError
	}
	return &rtResponseHeader{
		version: string(match[1]),
		status: status,
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
