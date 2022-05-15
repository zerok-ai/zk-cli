# Install application

### 1. Apply namespace, deployment and services
`kubectl apply -f node-deployment.yaml
`
### 2. Get the Ingress Controller's public IP address

`
kubectl get services --namespace ingress | grep 'app-ingress-ingress' | grep -v 'admission' | awk '{print $4}'
`

### 3. Configure app-ingress.yaml

Open app-ingress.yaml and put the IP printed in the previous step against the `- host:` entry

`
open app-ingress.yaml
`

### 4. Apply ingress object
`
kubectl apply -f app-ingress.yaml
`

### 4. Apply ingress object
`
kubectl apply -f app-ingress.yaml
`

## Labeling the pods

### Check the labels of the pods
`
kubectl get pods --namespace app --show-labels
`

### Label all the pods of the app
`
kubectl label pods --namespace app --selector=app=load-test load=high
`

### Set label to a pod
kubectl label pods --namespace app <name of the pod> load=low --overwrite=true

