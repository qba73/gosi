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

// SportEventItem is a short description of the event.
type SportEventItem struct {
	ID        string
	Date      string
	EventType string
	Title     string
	Status    string
}

// SportEvent describes details of the event.
type SportEvent struct {
	ID           string
	Type         string
	Title        string
	Status       string
	Date         string
	Summaries    []string
	EntriesOpen  string
	EntriesClose string
	EntriesSoFar string
	SocialMedia  SocialMedia
	Organizer    Organizer
}

type SocialMedia struct {
	Facebook  string
	Twitter   string
	Instagram string
	Website   string
}

type Organizer struct {
	Name    string
	Email   string
	Website string
}

type SportEvents struct {
	Event []SportEvent
}

func (s SportEvent) GetByType(eventType string) []SportEvent {
	return []SportEvent{}
}

// GetEvents knows how to retrive sport events from the SiEntries.
//
// GetEnevts uses pre-configured collector and default user-agent.
func GetEvents() []SportEvent {
	c := colly.NewCollector(
		colly.AllowedDomains("www.sientries.co.uk", "sientries.co.uk"),
		colly.UserAgent(userAgent),
	)
	return ListEvents(c, "https://www.sientries.co.uk/index.php?page=L")
}

// ListEvents knows how to get sport events data from given website url.
func ListEvents(c *colly.Collector, url string) []SportEvent {
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
			ID:     id,
			Date:   date,
			Type:   eventType,
			Title:  title,
			Status: status,
		}
		events = append(events, event)
	})
	c.Visit(url)
	return events
}

// GetEvent knows how to retrieve individual sport event data.
func GetEvent(c *colly.Collector, url string) SportEvent {
	var (
		event  SportEvent
		social SocialMedia

		// Event title
		title string

		// Sport event social media platforms
		facebook, twitter, web, instagram string

		// Sport event dates and participants
		entriesOpen, entriesClose, entriesSoFar string

		// Sport event organizer contact
		contact, email, website string
	)

	c.OnHTML("div.event-header-right", func(h *colly.HTMLElement) {
		h.ForEach("div.event-icon > a", func(_ int, e *colly.HTMLElement) {
			switch e.ChildAttr("img", "title") {
			case "Event Facebook Page":
				facebook = e.Attr("href")
			case "Event Twitter Feed":
				twitter = e.Attr("href")
			case "Event Instagram Feed":
				instagram = e.Attr("href")
			case "Event Website":
				web = e.Attr("href")
			}
		})

		social = SocialMedia{
			Facebook:  facebook,
			Twitter:   twitter,
			Instagram: instagram,
			Website:   web,
		}
	})

	c.OnHTML("div.event-summary", func(h *colly.HTMLElement) {
		title = h.ChildAttr("div.event-summary-logo > div.event-image > img", "alt")

		h.ForEach("div.event-summary-detail > div.event-summary-row", func(_ int, e *colly.HTMLElement) {
			switch e.ChildText("div.event-summary-label") {
			case "Entries Open":
				entriesOpen = e.ChildText("div.event-summary-text")
			case "Entries Close":
				entriesClose = e.ChildText("div.event-summary-text")
			case "Entries so Far":
				entriesSoFar = e.ChildText("div.event-summary-text")
			case "Contact":
				contact = e.ChildText("div.event-summary-text")
			case "Email":
				email = e.ChildText("div.event-summary-text > a")
			case "Website":
				website = e.ChildText("div.event-summary-text > a")
			}
		})

		event = SportEvent{
			Title:        title,
			EntriesOpen:  entriesOpen,
			EntriesClose: entriesClose,
			EntriesSoFar: entriesSoFar,
			Organizer: Organizer{
				Name:    contact,
				Email:   email,
				Website: website,
			},
		}
	})

	c.Visit(url)
	event.SocialMedia = social
	return event
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
		fmt.Fprintf(os.Stdout, "%s, %s, %s, %s\n", e.ID, e.Date, e.Type, e.Title)
	}
}
