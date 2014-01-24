package views_test

import "testing"
import . "github.com/zaiuz/views"

var _ View = &HtmlView{}

var htmlViewInfo = &viewInfo{
	ctor:               NewHtmlView,
	contentType:        "text/html",
	singleFile:         "./testviews/single.html",
	singleFileOutput:   "./testviews/single.output.html",
	parentFile:         "./testviews/parent.html",
	childFile:          "./testviews/child.html",
	combinedFileOutput: "./testviews/combined.output.html",
}

func TestNewHtmlView(t *testing.T)                { CheckView_Ctor(t, htmlViewInfo) }
func TestHtmlView_Subview(t *testing.T)           { CheckView_Subview(t, htmlViewInfo) }
func TestHtmlView_Render_Single(t *testing.T)     { CheckView_Render_Single(t, htmlViewInfo) }
func TestHtmlView_Render_WithData(t *testing.T)   { CheckView_Render_WithData(t, htmlViewInfo) }
func TestHtmlView_Render_Uncombined(t *testing.T) { CheckView_Render_Uncombined(t, htmlViewInfo) }
func TestHtmlView_Render_Combined(t *testing.T)   { CheckView_Render_Combined(t, htmlViewInfo) }
