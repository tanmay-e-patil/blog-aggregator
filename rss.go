package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchRSSFeeds(url string) (*RSSFeed, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rss RSSFeed
	if err := xml.Unmarshal(data, &rss); err != nil {
		return nil, err
	}
	return &rss, nil
}

func (cfg *apiConfig) worker(concurrency int) error {
	feedsToFetch, err := cfg.DB.GetNextFeedsToFetch(context.Background(), int32(concurrency))
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, feed := range feedsToFetch {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rss, err := fetchRSSFeeds(feed.Url)
			if err != nil {
				log.Printf("Error fetching RSS feed for %s: %s", feed.Url, err)
			}
			log.Printf("Fetched RSS feed for %s", feed.Url)

			for _, item := range rss.Channel.Item {
				log.Printf("Fetching %s", item.Title)
			}

		}()
	}
	wg.Wait()
	return nil
}
