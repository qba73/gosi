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

func TestCollectEvents(t *testing.T) {

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

	got := gosi.CollectEvents(c, ts.URL)

	want := 468

	if !cmp.Equal(len(got), want, cmpopts.IgnoreFields(gosi.SportEvent{}, "Title")) {
		t.Errorf("gosi.CollectEvents(c, %s) \n%s\n", ts.URL, cmp.Diff(len(got), want))
	}
}
