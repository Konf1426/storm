# IaC (Skeleton)

This directory contains a minimal Terraform skeleton for a single-VM deployment.
It is not production-ready. Fill variables and harden security before use.

## Layout
- `terraform/aws-ec2/` minimal EC2 deployment

## Usage (example)
```
cd infra/iac/terraform/aws-ec2
terraform init
terraform apply -var "region=eu-west-1" -var "instance_type=t3.small" -var "key_name=YOUR_KEY"
```
