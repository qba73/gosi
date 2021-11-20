package gosi

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const (
	userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/604.1"
)

type SportEvent struct {
	ID        string
	Date      string
	EventType string
	Title     string
	Status    string
}

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

func GetEvents() []SportEvent {
	c := colly.NewCollector(
		colly.AllowedDomains("www.sientries.co.uk", "sientries.co.uk"),
		colly.CacheDir("."),
		colly.UserAgent(userAgent),
	)
	return CollectEvents(c, "https://www.sientries.co.uk/index.php?page=L")
}

func CollectEvents(c *colly.Collector, URL string) []SportEvent {
	var events []SportEvent

	c.OnHTML("div.eti_wrap", func(h *colly.HTMLElement) {
		id := eventID(h.ChildAttr("div.eti_title > a", "href"))
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

func eventID(s string) string {
	out := strings.Split(s, "event_id=")
	if len(out) < 2 {
		return ""
	}
	return out[1]
}

func RunCLI() {
	events := GetEvents()
	for _, e := range events {
		fmt.Printf("%+v\n", e)
	}
}
