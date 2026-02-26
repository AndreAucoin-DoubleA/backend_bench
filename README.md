# BenchBackEnd

docker-compose up --build   : To run application 

CASSANDRA_HOST=localhost TEST_KEYSPACE=wiki_test go test ./test/integration/... -v : To run integration test locally

http://localhost:7000/status: To see all the stats that have been consumed

http://localhost:7000/login: To login as status endpoint is authenticated

act pull_request -j build-test-publish --secret-file .env: To run CI/CD pipeline


# 1. GoLang

- [x] Create a basic Go application that listens on port 7000 and has a status endpoint
- [x] Create a process to consume the wikipedia recent changes stream https://stream.wikimedia.org/v2/stream/recentchange and log these to stdout.
- [x] Replace the logs, with an in-memory /stats endpoint that a user can hit to get the latest stats on what we’ve processed
- [x] Create the following stats
    Number of messages consumed
    Number of distinct users
    Number of bots & Number of non-bots
    Count by distinct server URLs
- [x] Create tests for your application (if you didn’t already)
- [x] Run tests with the race detector on (-race)


# 2. Docker

- [x] Create a  Dockerfile for your application
- [x] Build & Run your dockerized application
- [x] Build a scratch container image of your application
- [x] Use a file to set all the configurable items like ports, URLs and anything else that can be       dynamic, load these configs via the file


# 3. Databases

- [x] Bring up a Scylla or Cassandra DB Docker
- [x] Build out a user login API endpoint
- [x] Build in a simple auth middleware, with bearer auth or JWT
- [x] Design a data model that supports efficient reading and writing of statistics to the DB, A. Start persisting your stats to the database, B. Allow your application to use an in-memory database
- [x] Bring up your application + DB with docker-compose
- [x] Create integration tests for your DB interactions

# 4. CI/CD

- [x] Test your application when there are PRs to main, builds should run a form of Linux (ubuntu). A.Run all unit tests within the project B.Run all integration tests for your application against the DB, this is most likely using docker-compose from above 
- [x] Run go vet tools to confirm good standards are being adhered to
- [x] Run https://github.com/golangci/golangci-lint
- [x] Create a new image for your application as part of a docker file
- [x] Publish this image to ghcr.io or dockerhub when all the checks & stages pass