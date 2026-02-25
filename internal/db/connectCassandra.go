package db

import (
	"log"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"golang.org/x/crypto/bcrypt"
)

func ConnectToCassandra(host string, port int, keyspace string, email, password string) *gocql.Session {
	var session *gocql.Session
	var err error

	// 1️⃣ Connect WITHOUT keyspace to create it
	for {
		cluster := gocql.NewCluster(host)
		cluster.Port = port
		cluster.Consistency = gocql.Quorum

		session, err = cluster.CreateSession()
		if err != nil {
			log.Println("Connecting to:", host)
			log.Println("Waiting for Cassandra container to be ready...")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	log.Println("Connected to Cassandra (no keyspace)")

	// Create keyspace
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
	session.Close() // close temporary session

	// 2️⃣ Connect WITH keyspace to create tables
	for {
		cluster := gocql.NewCluster(host)
		cluster.Port = port
		cluster.Keyspace = keyspace
		cluster.Consistency = gocql.Quorum

		session, err = cluster.CreateSession()
		if err != nil {
			log.Println("Waiting for keyspace to be ready...")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	log.Println("Connected to Cassandra with keyspace:", keyspace)

	err = session.Query(`
		CREATE TABLE IF NOT EXISTS user_by_email (
			email TEXT,
			id UUID,
			password_hash TEXT,
			created_at TEXT,
			PRIMARY KEY(email)
		);
	`).Exec()
	if err != nil {
		log.Fatal("Failed creating table:", err)
	}

	// Create tables
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS wiki_url_stats (
			stat_date TEXT,
			url TEXT,
			count counter,
			PRIMARY KEY (stat_date, url)
		);
	`).Exec()
	if err != nil {
		log.Fatal("Failed creating table:", err)
	}

	err = session.Query(`
		CREATE TABLE IF NOT EXISTS wiki_users_stats (
			stat_date TEXT,
    		username TEXT,
			count counter,
    		PRIMARY KEY (stat_date, username)
		);
	`).Exec()
	if err != nil {
		log.Fatal("Failed creating table:", err)
	}

	err = session.Query(`
		CREATE TABLE IF NOT EXISTS wiki_total_stats (
			stat_date TEXT,
			total_changes counter,
			num_bots counter,
			num_non_bots counter,
			PRIMARY KEY (stat_date)
		);
	`).Exec()
	if err != nil {
		log.Fatal("Failed creating table:", err)
	}
	log.Println("Cassandra tables initialized")

	seedInitialUser(session, email, password)
	seedInitialTotalStats(session)
	return session
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
