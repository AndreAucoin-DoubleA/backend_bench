package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	Client *kgo.Client
	Topic  string
	Group  string
}

func (consumer *Consumer) Close() {
	consumer.Client.Close()
}

func (c *Consumer) Consume(ctx context.Context, session *gocql.Session) {
	const maxSeenIDs = 100000
	seen := make(map[int64]struct{})
	for {
		if ctx.Err() != nil {
			return
		}

		pollCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		fetches := c.Client.PollFetches(pollCtx)
		cancel()

		if pollCtx.Err() == context.DeadlineExceeded {
			continue
		}

		if ctx.Err() != nil {
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Println("Fetch error:", err)
			}
			if ctx.Err() != nil {
				return
			}
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			var change RecentChange
			if err := json.Unmarshal(record.Value, &change); err != nil {
				return
			}

			if len(seen) >= maxSeenIDs {
				seen = make(map[int64]struct{}, maxSeenIDs)
			}

			if _, ok := seen[change.ID]; ok {
				return
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
		})
	}
}
