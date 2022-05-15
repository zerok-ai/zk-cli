# Permissions


### Assumption

`anton-test` is the name of the cluster

## On Admin machine. 
### 1. using kubectl

##### Open AWS config map

```
kubectl edit -n kube-system configmap/aws-auth
```
##### Add user details
Add user under `data -> mapUsers` (mapUsers is a sibling of mapRoles)

```
data:
  mapRoles: |
    ......
  mapUsers: |
    - groups:
      - system:masters
      userarn: arn:aws:iam::301129966109:user/mudit
      username: mudit
```
### OR
### 2. Using eksctl
```
eksctl create iamidentitymapping \
    --cluster anton-test \
    --region=us-east-2 \
    --arn arn:aws:iam::301129966109:user/mudit \
    --group "system:masters" \
    --no-duplicate-arns
```

---
## On user machine
```
aws eks update-kubeconfig --region us-east-2 --name anton-test
```
