variable "region" {
  type        = string
  description = "AWS region"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR"
  default     = "10.20.0.0/16"
}

variable "public_subnet_count" {
  type        = number
  description = "Number of public subnets"
  default     = 2
}

variable "private_subnet_count" {
  type        = number
  description = "Number of private subnets"
  default     = 2
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t3a.medium"
}

variable "key_name" {
  type        = string
  description = "EC2 key pair name"
}

variable "enable_ssh" {
  type        = bool
  description = "Allow SSH ingress"
  default     = false
}

variable "ssh_cidr" {
  type        = string
  description = "CIDR allowed to SSH"
  default     = "0.0.0.0/0"
}

variable "repo_url" {
  type        = string
  description = "Git repo to clone on instances"
  default     = "https://github.com/Konf1426/storm.git"
}

variable "git_ref" {
  type        = string
  description = "Git ref to deploy"
  default     = "feature/front"
}

variable "compose_file" {
  type        = string
  description = "Compose file path"
  default     = "infra/docker/docker-compose.cloud.yml"
}

variable "asg_min" {
  type        = number
  description = "ASG min size"
  default     = 1
}

variable "asg_max" {
  type        = number
  description = "ASG max size"
  default     = 2
}

variable "asg_desired" {
  type        = number
  description = "ASG desired capacity"
  default     = 1
}

variable "db_username" {
  type        = string
  description = "Postgres username"
  default     = "storm"
}

variable "db_password" {
  type        = string
  description = "Postgres password (leave empty for auto-generated)"
  default     = ""
  sensitive   = true
}

variable "db_name" {
  type        = string
  description = "Postgres database name"
  default     = "storm"
}

variable "db_instance_class" {
  type        = string
  description = "RDS instance class"
  default     = "db.t4g.micro"
}

variable "db_engine_version" {
  type        = string
  description = "Postgres engine version"
  default     = "16.3"
}

variable "db_storage_gb" {
  type        = number
  description = "RDS storage in GB"
  default     = 20
}

variable "redis_node_type" {
  type        = string
  description = "ElastiCache Redis node type"
  default     = "cache.t4g.micro"
}

variable "image_tag" {
  type        = string
  description = "Docker image tag to deploy"
  default     = "latest"
}

variable "jwt_secret" {
  type        = string
  description = "JWT secret (leave empty for auto-generated)"
  default     = ""
  sensitive   = true
}

variable "jwt_refresh_secret" {
  type        = string
  description = "JWT refresh secret (leave empty for auto-generated)"
  default     = ""
  sensitive   = true
}

variable "cors_origin" {
  type        = string
  description = "CORS origin for gateway"
  default     = "https://example.com"
}

variable "domain_name" {
  type        = string
  description = "Public domain for HTTPS (optional)"
  default     = ""
}

variable "hosted_zone_id" {
  type        = string
  description = "Route53 hosted zone id for domain validation"
  default     = ""
}
