# IaC (Skeleton)

This directory contains a minimal Terraform skeleton for a single-VM deployment.
It is not production-ready. Fill variables and harden security before use.

## Layout
- `terraform/aws-ec2/` minimal EC2 deployment
- `terraform/aws-prod/` VPC + ALB + ASG deployment (realistic baseline)

## Usage (example)
```
cd infra/iac/terraform/aws-ec2
terraform init
terraform apply -var "region=eu-west-1" -var "instance_type=t3.small" -var "key_name=YOUR_KEY"
```

## Usage (aws-prod example)
```
cd infra/iac/terraform/aws-prod
terraform init
terraform apply \
  -var "region=eu-west-1" \
  -var "key_name=YOUR_KEY" \
  -var "git_ref=feature/front"
```

Notes:
- The user_data clones the repo and runs `docker compose` on the instance.
- The aws-prod stack pulls pre-built images from ECR (see `scripts/ecr-push.sh`).
- Optional HTTPS: set `domain_name` and `hosted_zone_id` to enable ACM + ALB 443.
- This stack exposes the gateway on port 80 (and 443 if configured).
- If you change `image_tag`, push matching tags to ECR before applying.
