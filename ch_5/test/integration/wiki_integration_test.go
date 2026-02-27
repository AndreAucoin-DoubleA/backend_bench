package integration

import "testing"

func TestWikiCounterIncrement(t *testing.T) {
	statDate := "2026-02-24"
	url := "https://example.com"

	// Reset table for isolation
	if err := testSession.Query("TRUNCATE wiki_url_stats").Exec(); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	// Increment counter 3 times
	for i := 0; i < 3; i++ {
		if err := testSession.Query(`
			UPDATE wiki_url_stats
			SET count = count + 1
			WHERE stat_date = ? AND url = ?`,
			statDate, url,
		).Exec(); err != nil {
			t.Fatalf("Failed to increment counter: %v", err)
		}
	}

	// Read back count
	var count int64
	if err := testSession.Query(`
		SELECT count FROM wiki_url_stats
		WHERE stat_date = ? AND url = ?`,
		statDate, url,
	).Scan(&count); err != nil {
		t.Fatalf("Failed to fetch counter: %v", err)
	}

	if count != 3 {
		t.Fatalf("Expected count 3, got %d", count)
	}
}

func TestWikiDistinctUrlsAggregation(t *testing.T) {
	statDate := "2026-02-24"
	urls := []string{"a.com", "b.com", "a.com"}

	// Reset table

	if err := testSession.Query("TRUNCATE wiki_url_stats").Exec(); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	// Increment counters for multiple URLs
	for _, url := range urls {
		if err := testSession.Query(`
			UPDATE wiki_url_stats
			SET count = count + 1
			WHERE stat_date = ? AND url = ?`,
			statDate, url,
		).Exec(); err != nil {
			t.Fatalf("Failed to increment counter: %v", err)
		}
	}

	// Fetch all counts
	iter := testSession.Query(`
		SELECT url, count FROM wiki_url_stats
		WHERE stat_date = ?`,
		statDate,
	).Iter()

	distinctUrl := make(map[string]int)
	var url string
	var count int64
	for iter.Scan(&url, &count) {
		distinctUrl[url] = int(count)
	}
	iter.Close()

	if len(distinctUrl) != 2 {
		t.Fatalf("Expected 2 distinct URLs, got %d", len(distinctUrl))
	}
	if distinctUrl["a.com"] != 2 || distinctUrl["b.com"] != 1 {
		t.Fatalf("Unexpected counts: %+v", distinctUrl)
	}
}

func TestWikiTotalChangesAggregation(t *testing.T) {
	statDate := "2026-02-24"
	urls := []string{"a.com", "b.com", "c.com"}

	if err := testSession.Query("TRUNCATE wiki_url_stats").Exec(); err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}

	// Increment each URL
	for _, url := range urls {
		for i := 0; i < len(url); i++ { // variable increments
			if err := testSession.Query(`
				UPDATE wiki_url_stats
				SET count = count + 1
				WHERE stat_date = ? AND url = ?`,
				statDate, url,
			).Exec(); err != nil {
				t.Fatalf("Failed to increment counter: %v", err)
			}
		}
	}

	// Aggregate total
	iter := testSession.Query(`
		SELECT count FROM wiki_url_stats
		WHERE stat_date = ?`,
		statDate,
	).Iter()

	var total int64
	var count int64
	for iter.Scan(&count) {
		total += count
	}
	iter.Close()

	if total == 0 {
		t.Fatal("Expected total > 0")
	}
}
