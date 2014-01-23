package views_test

import "io/ioutil"
import "net/http/httptest"
import "testing"
import "github.com/zaiuz/testutil"
import . "github.com/zaiuz/views"
import a "github.com/stretchr/testify/assert"

var _ View = &HtmlView{}

func TestNewHtmlView(t *testing.T) {
	test := func() { NewHtmlView() }
	a.Panics(t, test, "does not throw error even when no filename given.")

	result := NewHtmlView(singleFile)
	a.NotNil(t, result, "constructor returns nil errorneously.")
}

func TestSubview(t *testing.T) {
	parent := NewHtmlView(parentFile)
	test := func() { parent.Subview() }
	a.Panics(t, test, "does not throw error even when no filename given.")

	result := parent.Subview(childFile)
	a.NotNil(t, result, "subview is nil errorneously.")
}

func TestRenderContentType(t *testing.T) {
	context := testutil.NewTestContext()
	recorder := context.ResponseWriter.(*httptest.ResponseRecorder)

	singleView := NewHtmlView(singleFile)
	singleView.Render(context, nil)

	contentType := recorder.HeaderMap["Content-Type"]
	a.NotEmpty(t, contentType, "Content-Type header was nil or empty.")
	a.Contains(t, contentType[0], "text/html", "Content-Type not text/html.")
}

func TestRenderSingle(t *testing.T) {
	singleView := NewHtmlView(singleFile)
	output, e := ioutil.ReadFile(singleFileOutput)
	a.NoError(t, e)

	renderEqual(t, singleView, output)
}

func TestRenderParent(t *testing.T) {
	renderFail(t, NewHtmlView(parentFile))
}

func TestRenderChild(t *testing.T) {
	renderFail(t, NewHtmlView(childFile))
}

func TestRenderCombined(t *testing.T) {
	output, e := ioutil.ReadFile(combinedFileOutput)
	a.NoError(t, e)

	parent := NewHtmlView(parentFile)
	child := parent.Subview(childFile)

	result, e := renderToString(child, nil)
	a.NoError(t, e)
	a.Equal(t, string(result), string(output), "combined result wrong.")
}
