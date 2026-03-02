package model

import (
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

func UpdateTotalStatsWithDB(change RecentChange, session *gocql.Session) error {
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

func UpdateUserStatsWithDB(change RecentChange, session *gocql.Session) error {
	today := time.Now().Format("2006-01-02")
	increment := int64(1)
	return session.Query(`
		UPDATE wiki_users_stats
		SET count = count + ?
		WHERE stat_date = ? AND username = ?`,
		increment, today, change.User,
	).Exec()
}

func UpdateUrlStatsWithDB(change RecentChange, session *gocql.Session) error {
	today := time.Now().Format("2006-01-02")
	increment := int64(1)
	return session.Query(`
		UPDATE wiki_url_stats
		SET count = count + ?
		WHERE stat_date = ? AND url = ?`,
		increment, today, change.ServerURL,
	).Exec()
}
