echo ACCOUNT_ID=$ACCOUNT_ID

# -----------------------------------------------
# 1. Create admin role and a group
# 2. Apply policy to the group, enabling it to assume
#    the role
# 3. Add users to the group
# -----------------------------------------------
POLICY=$(echo -n '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"arn:aws:iam::'; echo -n "$ACCOUNT_ID"; echo -n ':root"},"Action":"sts:AssumeRole","Condition":{}}]}')
echo POLICY=$POLICY
aws iam create-role \
  --role-name k8sAdmin \
  --description "Kubernetes administrator role (for AWS IAM Authenticator for Kubernetes)." \
  --assume-role-policy-document "$POLICY" \
  --output text \
  --query 'Role.Arn'

aws iam create-group --group-name k8sAdmin

ADMIN_GROUP_POLICY=$(echo -n '{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAssumeOrganizationAccountRole",
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Resource": "arn:aws:iam::'; echo -n "$ACCOUNT_ID"; echo -n ':role/k8sAdmin"
    }
  ]
}')
echo ADMIN_GROUP_POLICY=$ADMIN_GROUP_POLICY
aws iam put-group-policy \
--group-name k8sAdmin \
--policy-name k8sAdmin-policy \
--policy-document "$ADMIN_GROUP_POLICY"

aws iam add-user-to-group --group-name k8sAdmin --user-name mudit
aws iam add-user-to-group --group-name k8sAdmin --user-name shivam