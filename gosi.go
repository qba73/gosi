package gosi

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

const (
	userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/604.1"
)

// SportEvent is a short description of the event.
type SportEvent struct {
	ID        string
	Date      string
	EventType string
	Title     string
	Status    string
}

// SportEventInfo describes details of the event.
type SportEventInfo struct {
	Title       string
	SocialMedia struct {
		Facebook  string
		Twitter   string
		Instagram string
		Website   string
	}
	Date          string
	EntriesOpen   string
	EntriesClosed string
	EntriesSoFar  string
	Organizer     struct {
		Name    string
		Email   string
		Website string
	}
}

// GetEvents knows how to retrive sport events from the SiEntries.
//
// GetEnevts uses pre-configured collector and default user-agent.
func GetEvents() []SportEvent {
	c := colly.NewCollector(
		colly.AllowedDomains("www.sientries.co.uk", "sientries.co.uk"),
		colly.UserAgent(userAgent),
	)
	return CollectEvents(c, "https://www.sientries.co.uk/index.php?page=L")
}

// CollectEvents knows how to retrieve sort events data from
// the provided website URL.
func CollectEvents(c *colly.Collector, URL string) []SportEvent {
	var events []SportEvent
	c.OnHTML("div.eti_wrap", func(h *colly.HTMLElement) {
		id, _ := eventID(h.ChildAttr("div.eti_title > a", "href"))
		date := fmt.Sprintf("%s %s %s",
			h.ChildText("div.eti_date > .eti_day"),
			h.ChildText("div.eti_date > .eti_num"),
			h.ChildText("div.eti_date > .eti_month"),
		)
		eventType := h.ChildAttr("div.eti_type > img", "title")
		title := h.ChildText("div.eti_title > a")
		status := h.ChildText("div.eti_status > a > div.eti_button > span")

		event := SportEvent{
			ID:        id,
			Date:      date,
			EventType: eventType,
			Title:     title,
			Status:    status,
		}
		events = append(events, event)
	})
	c.Visit(URL)
	return events
}

func eventID(s string) (string, error) {
	path := strings.ReplaceAll(s, "&amp;", "&")
	u, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("parsing eventID string %s", path)
	}
	vals, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", fmt.Errorf("parsing query %s", u.RawQuery)
	}
	if vals.Get("event_id") == "" {
		return "", errors.New("missing event_id")
	}
	return vals.Get("event_id"), nil
}

// RunCLI runs the main machinery when the gosi
// is used as a command line utility.
func RunCLI() {
	events := GetEvents()
	for _, e := range events {
		fmt.Fprintf(os.Stdout, "%s, %s, %s, %s\n", e.ID, e.Date, e.EventType, e.Title)
	}
}
