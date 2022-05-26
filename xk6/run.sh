ulimit -n 65536
K6_PROMETHEUS_REMOTE_URL=http://localhost:9090/api/v1/write \
./k6 run script.js -o output-prometheus-remote
