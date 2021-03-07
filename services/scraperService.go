package services

import (
	"fmt"
	"strings"

	"github.com/Jeff-All/Dohyo/scrapers"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ScraperService - Service for handling scraping data
type ScraperService struct {
	log          *logrus.Logger
	db           *gorm.DB
	matchScraper scrapers.MatchScraper
}

// NewScraperService - Instantiates a new ScraperService
func NewScraperService(
	log *logrus.Logger,
	db *gorm.DB,
	matchScraper scrapers.MatchScraper,
) ScraperService {
	return ScraperService{
		log:          log,
		db:           db,
		matchScraper: matchScraper,
	}
}

// ScrapeMatches - Scrapes matches and populates the database
func (s *ScraperService) ScrapeMatches(month string, day uint) error {
	s.log.Infof("scraping matches for day %d", day)
	var matches *scrapers.ScrapedDayMatches
	var err error
	if matches, err = s.matchScraper.Scrape(1, day); err != nil {
		s.log.Errorf("error while scraping matches: %s", err)
		return err
	}
	scrapedMonth := strings.ToLower(matches.DayHead)
	if !strings.Contains(scrapedMonth, month) {
		err := fmt.Errorf("scraped month(%s) does not match expected month '%s'", scrapedMonth, month)
		s.log.Error(err)
		return err
	}
	s.log.Infof("scraped %d matches for day %d of month %s", len(matches.Matches), matches.Day, month)
	for _, match := range matches.Matches {
		s.log.Infof("scraped east{id=%d name=%s, rank=%s, won=%d, lost=%d} west{id=%d name=%s, rank=%s, won=%d, lost=%d}",
			match.East.RikishiID,
			match.East.Rikishi,
			match.East.Rank,
			match.East.Won,
			match.East.Lost,
			match.West.RikishiID,
			match.West.Rikishi,
			match.West.Rank,
			match.West.Won,
			match.West.Lost,
		)
	}
	return nil
}
