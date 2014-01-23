package views_test

import "net/http/httptest"
import "testing"
import "github.com/zaiuz/testutil"
import "regexp"
import . "github.com/zaiuz/views"
import z "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

const (
	singleFile         = "./testviews/single.html"
	singleFileOutput   = "./testviews/single.output.html"
	parentFile         = "./testviews/parent.html"
	childFile          = "./testviews/child.html"
	combinedFileOutput = "./testviews/combined.output.html"
)

type renderFunc func(*z.Context, interface{}) error

func renderEqual(t *testing.T, view View, expected []byte) {
	result, e := renderToString(view, nil)
	a.NoError(t, e)
	a.Equal(t, string(result), string(expected), "render result mismatch.")
}

func renderMatch(t *testing.T, view View, data interface{}, pattern string) {
	re := regexp.MustCompile(pattern)

	result, e := renderToString(view, data)
	a.NoError(t, e)
	a.NotNil(t, re.FindString(result), "render output does not match pattern.")
}

func renderFail(t *testing.T, view View) {
	_, e := renderToString(view, nil)
	a.Error(t, e, "fail rendering should return an error.")
}

func renderToString(view View, data interface{}) (string, error) {
	context := testutil.NewTestContext()
	e := view.Render(context, data)
	if e != nil {
		return "", e
	}

	resp := context.ResponseWriter.(*httptest.ResponseRecorder)
	return string(resp.Body.Bytes()), nil
}
