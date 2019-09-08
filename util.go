package rt

import (
	"regexp"
)

var reStatusCode = regexp.MustCompile(`^RT\/([\d\.]+) (\d\d\d) (.+)`)
var reResponseKV = regexp.MustCompile(`(\w+):( (.*))?`)

func parseRtResponseHeader(message []byte) (string, string, string, error) {
	match := reStatusCode.FindSubmatch(message)
	if match == nil {
		return "", "", "", ErrParseRTMessageError
	}
	return string(match[1]), string(match[2]), string(match[3]), nil
}

func parseRTResponseKVs(message []byte) (result map[string]string, err error) {
	matches := reResponseKV.FindAllSubmatch(message, -1)

	if matches == nil {
		return nil, ErrParseRTMessageError
	}
	result = make(map[string]string)
	for _, match := range matches {
		result[string(match[1])] = string(match[3])
	}

	return result, nil
}
