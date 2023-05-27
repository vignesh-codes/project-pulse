# pulse-service

docker build -t test-client .
docker build -t agg-v1 .
docker build -t pulse .

docker compose up redis postgres pulse agg fluent -d
docker compose up test 