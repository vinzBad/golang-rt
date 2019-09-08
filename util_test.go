package rt

import (
	"testing"
)

func TestParseRtResponseHeader(t *testing.T) {
	_, err := parseRtResponseHeader([]byte(`<!DOCTYPE html>`))
	if err != ErrParseRTMessageError {
		t.Errorf("parseRtResponseHeader didn't detect faulty message")
	}

	header, _ := parseRtResponseHeader([]byte(`RT/4.4.4 200 Ok`))
	expectedHeader := rtResponseHeader{
		version: "4.4.4",
		status:  200,
		message: "Ok",
	}

	if *header != expectedHeader {
		t.Error("Parsed header didn' match expected header")
	}
}

func TestParseRTResponseKVs(t *testing.T) {
	_, err := parseRTResponseKVs([]byte(`<!DOCTYPE html>`))
	if err != ErrParseRTMessageError {
		t.Errorf("parseRTResponseKVs didn't detect faulty message")
	}

	result, _ := parseRTResponseKVs([]byte(`
	foo: bar
	yes: yaa
	`))

	expectedResult := map[string]string{
		"foo": "bar",
		"yes": "yaa",
	}

	for k, v := range expectedResult {
		if result[k] != v {
			t.Errorf("Expected %q: %q got %q", k, v, result[k])
		}
	}

}
