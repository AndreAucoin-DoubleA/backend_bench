package db

import (
	"log"
	"os"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"golang.org/x/crypto/bcrypt"
)

func ConnectToCassandra(host string, port int, keyspace string, email, password string) *gocql.Session {
	var session *gocql.Session
	var err error

	timeout := 60 * time.Second
	start := time.Now()

	for {
		cluster := gocql.NewCluster(host)
		cluster.Port = port
		cluster.Consistency = gocql.Quorum

		session, err = cluster.CreateSession()
		if err != nil {
			if os.Getenv("RUN_CI") == "true" && time.Since(start) > timeout {
				log.Fatalf("CI Timeout waiting for Cassandra: %v", err)
			}
			log.Println("Connecting to:", host)
			log.Println("Waiting for Cassandra container to be ready...")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	log.Println("Connected to Cassandra (no keyspace)")

	err = session.Query(`
		CREATE KEYSPACE IF NOT EXISTS ` + keyspace + `
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		};
	`).Exec()
	if err != nil {
		log.Fatal("Failed creating keyspace:", err)
	}
	log.Println("Keyspace created or already exists")
	session.Close()

	start = time.Now()
	for {
		cluster := gocql.NewCluster(host)
		cluster.Port = port
		cluster.Keyspace = keyspace
		cluster.Consistency = gocql.Quorum

		session, err = cluster.CreateSession()
		if err != nil {
			if os.Getenv("RUN_CI") == "true" && time.Since(start) > timeout {
				log.Fatalf("CI Timeout waiting for keyspace: %v", err)
			}
			log.Println("Waiting for keyspace to be ready...")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	log.Println("Connected to Cassandra with keyspace:", keyspace)

	createTables(session)
	seedInitialUser(session, email, password)
	seedInitialTotalStats(session)

	return session
}

func createTables(session *gocql.Session) {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS user_by_email (
			email TEXT,
			id UUID,
			password_hash TEXT,
			created_at TEXT,
			PRIMARY KEY(email)
		);`,
		`CREATE TABLE IF NOT EXISTS wiki_url_stats (
			stat_date TEXT,
			url TEXT,
			count counter,
			PRIMARY KEY (stat_date, url)
		);`,
		`CREATE TABLE IF NOT EXISTS wiki_users_stats (
			stat_date TEXT,
			username TEXT,
			count counter,
			PRIMARY KEY (stat_date, username)
		);`,
		`CREATE TABLE IF NOT EXISTS wiki_total_stats (
			stat_date TEXT,
			total_changes counter,
			num_bots counter,
			num_non_bots counter,
			PRIMARY KEY (stat_date)
		);`,
	}

	for _, q := range tables {
		if err := session.Query(q).Exec(); err != nil {
			log.Fatalf("Failed creating table: %v", err)
		}
	}
	log.Println("Cassandra tables initialized")
}

func seedInitialUser(session *gocql.Session, email, password string) {
	id := gocql.TimeUUID()
	createdAt := time.Now().Format(time.RFC3339)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed hashing password:", err)
	}

	err = session.Query(`
		INSERT INTO user_by_email (email, id, password_hash, created_at)
		VALUES (?, ?, ?, ?)
	`, email, id, hashedPassword, createdAt).Exec()
	if err != nil {
		log.Fatal("Failed seeding initial user:", err)
	}
}

func seedInitialTotalStats(session *gocql.Session) {
	today := time.Now().Format("2006-01-02")
	err := session.Query(`
		UPDATE wiki_total_stats
		SET total_changes = total_changes + 0, num_bots = num_bots + 0, num_non_bots = num_non_bots + 0
		WHERE stat_date = ?
	`, today).Exec()
	if err != nil {
		log.Fatal("Failed seeding initial stats:", err)
	}
}
