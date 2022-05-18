#create cluster
eksctl create cluster -f ./yaml/cluster/clusterconfig.yaml
eksctl create cluster -f ./yaml/cluster/eks-admin-service-account.yaml

#------ give cluster access permissions to users #------ 
# 1. Mudit
eksctl create iamidentitymapping \
    --cluster anton-test \
    --region=us-east-2 \
    --arn arn:aws:iam::301129966109:user/mudit \
    --group "system:masters" \
    --no-duplicate-arns

# 2. Shivam
eksctl create iamidentitymapping \
    --cluster anton-test \
    --region=us-east-2 \
    --arn arn:aws:iam::301129966109:user/shivam \
    --group "system:masters" \
    --no-duplicate-arns