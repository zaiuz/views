package views_test

import "testing"
import . "github.com/zaiuz/views"

var textViewInfo = &viewInfo{
	ctor:               NewTextView,
	contentType:        "text/plain",
	singleFile:         "./testviews/single.txt",
	singleFileOutput:   "./testviews/single.output.txt",
	parentFile:         "./testviews/parent.txt",
	childFile:          "./testviews/child.txt",
	combinedFileOutput: "./testviews/combined.output.txt",
}

func TestNewTextView(t *testing.T)                { CheckView_Ctor(t, textViewInfo) }
func TestTextView_Subview(t *testing.T)           { CheckView_Subview(t, textViewInfo) }
func TestTextView_Render_Single(t *testing.T)     { CheckView_Render_Single(t, textViewInfo) }
func TestTextView_Render_WithData(t *testing.T)   { CheckView_Render_WithData(t, textViewInfo) }
func TestTextView_Render_Uncombined(t *testing.T) { CheckView_Render_Uncombined(t, textViewInfo) }
func TestTextView_Render_Combined(t *testing.T)   { CheckView_Render_Combined(t, textViewInfo) }
