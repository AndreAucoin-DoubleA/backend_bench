package model

import (
	"sync"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type RecentChange struct {
	ID        int64  `json:"id"`
	User      string `json:"user"`
	Bot       bool   `json:"bot"`
	ServerURL string `json:"server_url"`
}

type WikiStats struct {
	sync.Mutex                        // embed mutex
	TotalChanges  int                 `json:"total_changes"`
	DistinctUsers map[string]struct{} `json:"-"` // internal set for counting
	NumBots       int                 `json:"num_bots"`
	NumNonBots    int                 `json:"num_non_bots"`
	DistinctUrl   map[string]int      `json:"distinct_url"`
}

type WikiRepository struct {
	Session *gocql.Session
}

type WikiStatsNew struct {
	TotalChanges  int            `json:"total_changes"`
	DistinctUsers int            `json:"-"` // internal set for counting
	NumBots       int            `json:"num_bots"`
	NumNonBots    int            `json:"num_non_bots"`
	DistinctUrl   map[string]int `json:"distinct_url"`
}

func (r *WikiRepository) GetRecentChanges() (*WikiStatsNew, error) {
	var changes WikiStatsNew
	today := time.Now().Format("2006-01-02")
	urlMap := make(map[string]int)
	var count int64

	err := r.Session.Query(`
        SELECT total_changes, num_bots, num_non_bots
        FROM wiki_total_stats
        WHERE stat_date = ?`,
		today,
	).Scan(&changes.TotalChanges, &changes.NumBots, &changes.NumNonBots)

	if err != nil {
		return nil, err
	}

	err = r.Session.Query(`
        SELECT COUNT(*) 
        FROM wiki_users_stats 
        WHERE stat_date = ?
    `, today).Scan(&count)
	changes.DistinctUsers = int(count)

	if err != nil {
		return nil, err
	}

	iter := r.Session.Query(`
        SELECT url, count
        FROM wiki_url_stats
        WHERE stat_date = ?
    `, today).Iter()

	var url string

	for iter.Scan(&url, &count) {
		urlMap[url] = int(count)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	changes.DistinctUrl = urlMap
	return &changes, nil
}
