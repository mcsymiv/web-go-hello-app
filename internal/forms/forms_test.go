package forms

import (
	"net/url"
	"testing"
)

func TestForm_Required(t *testing.T) {

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	form := New(postData)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}
}

func TestForm_HasRequired(t *testing.T) {
	postData := url.Values{}
	postData.Add("a", "a")

	form := New(postData)

	form.HasRequired("a")

	if !form.Valid() {
		t.Error("Form shows valid when required field is missing")
	}
}
