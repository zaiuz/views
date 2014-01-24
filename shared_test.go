package views_test

import "net/http/httptest"
import "testing"
import "github.com/zaiuz/testutil"
import "io/ioutil"
import "regexp"
import "code.google.com/p/go-uuid/uuid"
import . "github.com/zaiuz/views"
import z "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

type testViewData map[string]string

func NewTestViewData() testViewData {
	return testViewData(map[string]string{"Data": uuid.New()})
}

func (t testViewData) SearchPattern() string {
	return t["Data"]
}

type viewInfo struct {
	ctor               func(filenames ...string) View
	contentType        string
	singleFile         string
	singleFileOutput   string
	parentFile         string
	childFile          string
	combinedFileOutput string
}

func CheckView_Ctor(t *testing.T, info *viewInfo) {
	test := func() { info.ctor() }
	a.Panics(t, test, "view ctor does not panic when given no filename.")

	result := info.ctor(info.singleFile)
	a.NotNil(t, result, "view ctor returns nil even when filename given.")
}

func CheckView_Subview(t *testing.T, info *viewInfo) {
	parent := info.ctor(info.parentFile)
	child := parent.Subview(info.childFile)
	a.NotNil(t, child, "subview should never returns nil.")
}

func CheckView_Render_Single(t *testing.T, info *viewInfo) {
	view := info.ctor(info.singleFile)
	Render(t, view, nil).ExpectOutput(info.contentType, info.singleFileOutput)
}

func CheckView_Render_WithData(t *testing.T, info *viewInfo) {
	view := info.ctor(info.singleFile)
	data := NewTestViewData()

	Render(t, view, data).ExpectOutputMatch(info.contentType, data.SearchPattern())
}

func CheckView_Render_Uncombined(t *testing.T, info *viewInfo) {
	view := info.ctor(info.parentFile)
	Render(t, view, nil).ShouldFail()

	view = info.ctor(info.childFile)
	Render(t, view, nil).ShouldFail()
}

func CheckView_Render_Combined(t *testing.T, info *viewInfo) {
	parent := info.ctor(info.parentFile)
	child := parent.Subview(info.childFile)

	Render(t, child, nil).ExpectOutput(info.contentType, info.combinedFileOutput)
}

type RenderExpectable struct {
	t    *testing.T
	view View

	context *z.Context
	result  []byte
	err     error
}

func Render(t *testing.T, view View, data interface{}) *RenderExpectable {
	context := testutil.NewTestContext()
	e := view.Render(context, data)

	recorder := context.ResponseWriter.(*httptest.ResponseRecorder)
	result := recorder.Body.Bytes()
	return &RenderExpectable{t, view, context, result, e}
}

func (r *RenderExpectable) ShouldFail() *RenderExpectable {
	a.NotNil(r.t, r.err, "error was expected, but no error occurs.")
	return r
}

func (r *RenderExpectable) ExpectOutput(contentType, filename string) *RenderExpectable {
	a.NoError(r.t, r.err)

	recorder := r.context.ResponseWriter.(*httptest.ResponseRecorder)
	contentTypes := recorder.HeaderMap["Content-Type"]
	a.NotEmpty(r.t, contentTypes, "Content-Type header was nil or empty.")
	a.Contains(r.t, contentTypes[0], contentType, "Content-Type not `%s`.", contentType)

	output, e := ioutil.ReadFile(filename)
	a.NoError(r.t, e)
	a.Equal(r.t, string(output), string(r.result), "rendered result incorrect.")
	return r
}

func (r *RenderExpectable) ExpectOutputMatch(contentType, pattern string) *RenderExpectable {
	a.NoError(r.t, r.err)

	recorder := r.context.ResponseWriter.(*httptest.ResponseRecorder)
	contentTypes := recorder.HeaderMap["Content-Type"]
	a.NotEmpty(r.t, contentTypes, "Content-Type header was nil or empty.")
	a.Contains(r.t, contentTypes[0], contentType, "Content-Type not `%s`.", contentType)

	rx := regexp.MustCompile(pattern)
	a.True(r.t, rx.Match(r.result), "rendered result does not match expected regex.")
	return r
}
