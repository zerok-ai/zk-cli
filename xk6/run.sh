
ulimit -n 65536

HOSTNAME=$(aws elb describe-load-balancers | jq -r '.LoadBalancerDescriptions[0].CanonicalHostedZoneName')
echo $HOSTNAME

kubectl port-forward service/prom-kube-prometheus-stack-prometheus -n monitoring 9090:9090 &

K6_PROMETHEUS_REMOTE_URL=http://localhost:9090/api/v1/write \
./k6 run script.js \
    -o output-prometheus-remote \
    -e MY_HOSTNAME=${HOSTNAME}
