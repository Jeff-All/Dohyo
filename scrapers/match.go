package scrapers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// ScrapedMatchRikishi - The scraped data representing a rikishi in a match
type ScrapedMatchRikishi struct {
	RikishiID string `json:"rikishi_id"`
	Rank      string `json:"banzuke_name_eng"`
	Rikishi   string `json:"shikona_eng"`
	Won       uint   `json:"won_number"`
	Lost      uint   `json:"lost_number"`
}

// ScrapedDayMatches - The scraped data for a the matches of a day of a tournament
type ScrapedDayMatches struct {
	DayHead string `json:"dayHead"`
	Day     string
	Matches []struct {
		East ScrapedMatchRikishi `json:"east"`
		West ScrapedMatchRikishi `json:"west"`
	} `json:"TorikumiData"`
}

// MatchScraper - Scraper for pulling matches
type MatchScraper struct {
	log  *logrus.Logger
	url  string
	host string
}

// NewMatchScraper - Instantiates and returns a new MatchScraper
func NewMatchScraper(
	log *logrus.Logger,
	url string,
	host string,
) MatchScraper {
	return MatchScraper{
		log:  log,
		url:  url,
		host: host,
	}
}

// Scrape - Scrapes the GrandSumo page for data
func (s MatchScraper) Scrape(division uint, day uint) (*ScrapedDayMatches, error) {
	var err error
	var resp *http.Response
	var req *http.Request
	s.log.Info(fmt.Sprintf("%s/%d/%d/", s.url, division, day))
	if req, err = http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%d/%d/", s.url, division, day),
		strings.NewReader((url.Values{}).Encode()),
	); err != nil {
		s.log.Error("error while building request to auth0 server: %s", err)
		return nil, err
	}
	req.Header.Add("Host", s.host)
	if resp, err = http.DefaultClient.Do(req); err != nil {
		s.log.Error("error while scraping matches: %s", err)
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		s.log.Infof("unexpected status code while scraping matches: %s", resp.StatusCode)
		return nil, fmt.Errorf("unexpected error code %d while scraping matches", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	dayMatches := ScrapedDayMatches{}
	if err = decoder.Decode(&dayMatches); err != nil {
		s.log.Errorf("error while decoding scaped match response: %s", err)
		return nil, err
	}
	return &dayMatches, nil
}
