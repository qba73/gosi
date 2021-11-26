package gosi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gocolly/colly"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/qba73/gosi"
)

func TestListEvents(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/response-all-events.html")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(rw, f)
	}))
	defer ts.Close()

	c := colly.NewCollector()

	got := gosi.ListEvents(c, ts.URL)

	want := 468

	if !cmp.Equal(len(got), want, cmpopts.IgnoreFields(gosi.SportEvent{}, "Title")) {
		t.Errorf("gosi.CollectEvents(c, %s) \n%s\n", ts.URL, cmp.Diff(len(got), want))
	}
}

func TestGetEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/response-single-event-open.html")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(rw, f)
	}))
	defer ts.Close()

	c := colly.NewCollector()

	got := gosi.GetEvent(c, ts.URL)

	want := gosi.SportEvent{
		Title:        "Bollihope Carrs Fell Race",
		EntriesOpen:  "Tuesday 19th October 2021",
		EntriesClose: "Thursday 9th December 2021 at 23:00",
		EntriesSoFar: "45 Participants",
		Organizer: gosi.Organizer{
			Name:    "Andy Blackett",
			Email:   "andyblackett@googlemail.com",
			Website: "https://www.durhamfellrunners.org/bollihope-carrs/",
		},
		SocialMedia: gosi.SocialMedia{
			Facebook: "https://www.facebook.com/45809454858",
			Twitter:  "https://twitter.com/durhamfellrun",
			Website:  "https://www.durhamfellrunners.org/bollihope-carrs/",
		},
	}

	if !cmp.Equal(got, want) {
		t.Errorf("gosi.GetEvent(c, \"8957\") \n%s\n", cmp.Diff(got, want))
	}
}
