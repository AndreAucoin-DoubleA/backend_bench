package wikiconsumer

import (
	"backend_bench/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func StartWikiConsumer(streamURL string, session *gocql.Session) {
	fmt.Println("Connecting to:", streamURL)

	req, err := http.NewRequest("GET", streamURL, nil)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("User-Agent", "backend-bench-dev")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Connected! Reading stream...")

	seen := make(map[int64]struct{})
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		jsonData := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		var change model.RecentChange
		if err := json.Unmarshal([]byte(jsonData), &change); err != nil {
			continue
		}

		if _, ok := seen[change.ID]; ok {
			continue
		}
		seen[change.ID] = struct{}{}

		if err := UpdateTotalStatsWithDB(change, session); err != nil {
			fmt.Println("Failed updating total stats:", err)
		}
		if err := UpdateUserStatsWithDB(change, session); err != nil {
			fmt.Println("Failed updating user stats:", err)
		}
		if err := UpdateUrlStatsWithDB(change, session); err != nil {
			fmt.Println("Failed updating URL stats:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
	}
}

func UpdateTotalStatsWithDB(change model.RecentChange, session *gocql.Session) error {
	today := time.Now().Format("2006-01-02")
	increment := int64(1)
	if change.Bot {
		return session.Query(`
        UPDATE wiki_total_stats
        SET total_changes = total_changes + ?, num_bots = num_bots + ?
        WHERE stat_date = ?`,
			increment, increment, today,
		).Exec()
	} else {
		return session.Query(`
        UPDATE wiki_total_stats
        SET total_changes = total_changes + ?, num_non_bots = num_non_bots + ?
        WHERE stat_date = ?`,
			increment, increment, today,
		).Exec()
	}
}

func UpdateUserStatsWithDB(change model.RecentChange, session *gocql.Session) error {
	today := time.Now().Format("2006-01-02")
	increment := int64(1)
	return session.Query(`
		UPDATE wiki_users_stats
		SET count = count + ?
		WHERE stat_date = ? AND username = ?`,
		increment, today, change.User,
	).Exec()
}

func UpdateUrlStatsWithDB(change model.RecentChange, session *gocql.Session) error {
	today := time.Now().Format("2006-01-02")
	increment := int64(1)
	return session.Query(`
		UPDATE wiki_url_stats
		SET count = count + ?
		WHERE stat_date = ? AND url = ?`,
		increment, today, change.ServerURL,
	).Exec()
}
