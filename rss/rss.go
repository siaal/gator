package rss

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/siaal/gator/internal/database"
	"github.com/siaal/gator/internal/state"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("illegal request: %w", err)
	}
	req.Header.Add("User-Agent", "gator")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var feed RSSFeed
	if err = xml.Unmarshal(bytes, &feed); err != nil {
		return nil, fmt.Errorf("malformed body error: %w", err)
	}

	cleaned := cleanFeed(feed)

	return &cleaned, nil
}

func cleanFeed(feed RSSFeed) RSSFeed {
	feed.Channel.Link = html.UnescapeString(feed.Channel.Link)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Link = html.UnescapeString(feed.Channel.Item[i].Link)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].PubDate = html.UnescapeString(feed.Channel.Item[i].PubDate)
	}
	return feed
}

func parseDatetime(timestring string) (time.Time, error) {
	timestrings := []string{"Mon, 02 Jan 2006 15:04:05 -0700"}
	for _, format := range timestrings {
		date, err := time.Parse(format, timestring)
		if err == nil {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse timestring")
}

func ScrapeFeeds(s *state.State) error {
	now := time.Now()
	for {
		ctx := context.Background()
		feed, err := s.DB.NextToFetch(ctx)
		if err != nil {
			return fmt.Errorf("db err: %w", err)
		}
		if feed.LastFetchedAt.Valid && now.Sub(feed.LastFetchedAt.Time) < time.Hour {
			return nil
		}
		ctx = context.Background()
		slog.Debug("pulling next feed", "feed.Name", feed.Name)
		data, err := FetchFeed(ctx, feed.Url)
		if err != nil {
			return fmt.Errorf("fetch err: %w", err)
		}
		ctx = context.Background()
		err = s.DB.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{Now: sql.NullTime{Time: now, Valid: true}, ID: feed.ID})
		if err != nil {
			return fmt.Errorf("db err when marking fetched: %w", err)
		}
		for _, item := range data.Channel.Item {
			ctx := context.Background()
			pubDate, err := parseDatetime(item.PubDate)
			if err != nil {
				slog.Error("could not parse datetime string", "item.PubDate", item.PubDate)
				continue
			}
			arg := database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   now,
				UpdatedAt:   now,
				PublishedAt: pubDate,
				Title:       item.Title,
				Url:         item.Link,
				Description: item.Description,
				FeedID:      feed.ID,
			}
			post, err := s.DB.CreatePost(ctx, arg)
			switch {
			case err == nil:
				slog.Info("Posted", "post.Title", post.Title)
			case strings.Contains(err.Error(), "UNIQUE"):
				continue
			default:
				return fmt.Errorf("db error: %w", err)
			}
		}
	}

}
