
ulimit -n 65536

HOSTNAME=$(aws elb describe-load-balancers | jq -r '.LoadBalancerDescriptions[0].CanonicalHostedZoneName')

K6_PROMETHEUS_REMOTE_URL=http://localhost:9090/api/v1/write \
./k6 run script.js \
    -o output-prometheus-remote \
    -e MY_HOSTNAME=${HOSTNAME}
