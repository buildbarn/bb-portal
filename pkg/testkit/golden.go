package testkit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var updateGoldenFilesCommand = "bazel run //test/integrationtest:integrationtest_test -- --update-golden"

// CompareOptions struct
type CompareOptions struct {
	// Whether to ignore date/times when comparing or writing the golden
	// file to disk
	DateTimeAgnostic bool
}

// CheckAgainstGoldenFile sees if the JSON serialization of what we got matches a golden file.
//
// got is the object we're testing. It will be rendered as JSON & compared with the golden file
// goldenDir defines the directory that holds golden files.
// testName is the name of the test, which will be sluggified to determine the golden file to use
// update is a flag that determines if we should update the golden file.
func CheckAgainstGoldenFile(t *testing.T, got map[string]interface{}, goldenDir, testName string, update *bool, opts *CompareOptions) {
	testSlug := strings.ReplaceAll(testName, " ", "-")
	golden := filepath.Join(goldenDir, fmt.Sprintf("%s.golden.json", testSlug))

	if *update {
		dir := filepath.Dir(golden)
		require.NoError(t, os.MkdirAll(dir, os.ModePerm))
		require.NoError(t, os.WriteFile(golden, prettyJSON(got), 0o640))
	}

	// get golden file
	want, err := os.ReadFile(golden)
	if err != nil {
		require.FailNow(
			t,
			fmt.Sprintf(
				"Failed to read golden file %s. If the file is missing, create it and then update its contents by running '%s'. Underlying error: %s",
				golden,
				updateGoldenFilesCommand,
				err.Error(),
			),
		)
	}

	gotStr := string(prettyJSON(got))
	wantStr := string(want)
	if opts != nil && opts.DateTimeAgnostic {
		var e error
		gotStr, e = replaceTimes(gotStr)
		require.NoError(t, e)
		wantStr, e = replaceTimes(wantStr)
		require.NoError(t, e)
	}

	require.JSONEq(t, wantStr, gotStr)
}

// replaceTimes finds all RFC3339 times and RFC7232 (section 2.2) times in the
// given string and replaces them with "0001-01-01T00:00:00Z" (for RFC3339) or
// "Sat, 01 Jan 0001 00:00:00 GMT" (for RFC7232) respectively.
func replaceTimes(str string) (string, error) {
	const hour = "([01][0-9]|2[0-3])"
	const minute = "([0-5][0-9])"
	const second = "([0-5][0-9]|60)"
	year := "([0-9]+)"
	month := "(0[1-9]|1[012])"
	day := "(0[1-9]|[12][0-9]|3[01])"
	datePattern := year + "-" + month + "-" + day

	subSecond := "(\\.[0-9]+)?"
	timePattern := hour + ":" + minute + ":" + second + subSecond

	timeZoneOffset := "(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))"

	pattern := datePattern + "[Tt]" + timePattern + timeZoneOffset

	rfc3339Pattern, err := regexp.Compile(pattern)
	if err != nil {
		return "", errors.Wrapf(err, "failed to compile RFC3339 regex pattern: %s", pattern)
	}
	res := rfc3339Pattern.ReplaceAllString(str, `0001-01-01T00:00:00Z`)

	dayName := "(Mon|Tue|Wed|Thu|Fri|Sat|Sun)"
	day = "[0-9]{2}"
	month = "(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)"
	year = "[0-9]{4}"
	tz := "(GMT|CEST|UTC|IST|[A-Z]+)"
	pattern = dayName + ", " + day + " " + month + " " + year + " " + hour + ":" + minute + ":" + second + " " + tz

	lastModifiedPattern, err := regexp.Compile(pattern)
	if err != nil {
		return "", errors.Wrapf(err, "failed to compile RFC7232 last-modified regex pattern: %s", pattern)
	}

	return lastModifiedPattern.ReplaceAllString(res, `Mon, 01 Jan 0001 00:00:00 GMT`), nil
}

// prettyJSON
func prettyJSON(data interface{}) []byte {
	jsonText, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	return jsonText
}
