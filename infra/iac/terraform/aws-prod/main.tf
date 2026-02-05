terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.region
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "storm" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "storm-vpc"
  }
}

resource "aws_internet_gateway" "storm" {
  vpc_id = aws_vpc.storm.id
  tags = {
    Name = "storm-igw"
  }
}

resource "aws_subnet" "public" {
  count                   = var.public_subnet_count
  vpc_id                  = aws_vpc.storm.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = {
    Name = "storm-public-${count.index}"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.storm.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.storm.id
  }
  tags = {
    Name = "storm-public-rt"
  }
}

resource "aws_route_table_association" "public" {
  count          = var.public_subnet_count
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_subnet" "private" {
  count             = var.private_subnet_count
  vpc_id            = aws_vpc.storm.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 10)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = "storm-private-${count.index}"
  }
}

resource "aws_security_group" "alb" {
  name        = "storm-alb-sg"
  description = "ALB ingress"
  vpc_id      = aws_vpc.storm.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "app" {
  name        = "storm-app-sg"
  description = "Gateway ingress from ALB and optional SSH"
  vpc_id      = aws_vpc.storm.id

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  dynamic "ingress" {
    for_each = var.enable_ssh ? [1] : []
    content {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = [var.ssh_cidr]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "db" {
  name        = "storm-db-sg"
  description = "Postgres ingress from app"
  vpc_id      = aws_vpc.storm.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "redis" {
  name        = "storm-redis-sg"
  description = "Redis ingress from app"
  vpc_id      = aws_vpc.storm.id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_ami" "al2" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

resource "aws_iam_role" "storm_ec2" {
  name = "storm-ec2-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "storm_ec2_ecr" {
  name = "storm-ec2-ecr"
  role = aws_iam_role.storm_ec2.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchGetImage",
          "ecr:GetDownloadUrlForLayer"
        ]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_instance_profile" "storm_ec2" {
  name = "storm-ec2-profile"
  role = aws_iam_role.storm_ec2.name
}

resource "aws_ecr_repository" "gateway" {
  name = "storm-gateway"
}

resource "aws_ecr_repository" "messages" {
  name = "storm-messages"
}

resource "random_password" "jwt" {
  length  = 32
  special = false
}

resource "random_password" "jwt_refresh" {
  length  = 32
  special = false
}

resource "random_password" "db" {
  length  = 24
  special = false
}

locals {
  db_password    = var.db_password != "" ? var.db_password : random_password.db.result
  jwt_secret     = var.jwt_secret != "" ? var.jwt_secret : random_password.jwt.result
  jwt_refresh    = var.jwt_refresh_secret != "" ? var.jwt_refresh_secret : random_password.jwt_refresh.result
  rds_endpoint   = aws_db_instance.storm.address
  redis_endpoint = aws_elasticache_cluster.redis.cache_nodes[0].address
  postgres_dsn   = "postgres://${var.db_username}:${local.db_password}@${local.rds_endpoint}:5432/${var.db_name}?sslmode=disable"
}

resource "aws_db_subnet_group" "storm" {
  name       = "storm-db-subnet"
  subnet_ids = aws_subnet.private[*].id
}

resource "aws_db_instance" "storm" {
  identifier              = "storm-postgres"
  engine                  = "postgres"
  engine_version          = var.db_engine_version
  instance_class          = var.db_instance_class
  allocated_storage       = var.db_storage_gb
  username                = var.db_username
  password                = local.db_password
  db_name                 = var.db_name
  publicly_accessible     = false
  skip_final_snapshot     = true
  vpc_security_group_ids  = [aws_security_group.db.id]
  db_subnet_group_name    = aws_db_subnet_group.storm.name
  backup_retention_period = 3
}

resource "aws_elasticache_subnet_group" "storm" {
  name       = "storm-redis-subnet"
  subnet_ids = aws_subnet.private[*].id
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "storm-redis"
  engine               = "redis"
  node_type            = var.redis_node_type
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.storm.name
  security_group_ids   = [aws_security_group.redis.id]
}

resource "aws_launch_template" "storm" {
  name_prefix   = "storm-"
  image_id      = data.aws_ami.al2.id
  instance_type = var.instance_type
  key_name      = var.key_name

  vpc_security_group_ids = [aws_security_group.app.id]

  user_data = base64encode(templatefile("${path.module}/user_data.sh.tftpl", {
    repo_url     = var.repo_url
    git_ref      = var.git_ref
    compose_file = var.compose_file
    region       = var.region
    ecr_gateway  = aws_ecr_repository.gateway.repository_url
    ecr_messages = aws_ecr_repository.messages.repository_url
    image_tag    = var.image_tag
    jwt_secret   = local.jwt_secret
    jwt_refresh  = local.jwt_refresh
    cors_origin  = var.cors_origin
    postgres_dsn = local.postgres_dsn
    redis_addr   = "${local.redis_endpoint}:6379"
  }))

  iam_instance_profile {
    name = aws_iam_instance_profile.storm_ec2.name
  }
}

resource "aws_lb" "storm" {
  name               = "storm-alb"
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.public[*].id
}

resource "aws_lb_target_group" "gateway" {
  name     = "storm-gateway-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.storm.id

  health_check {
    path                = "/healthz"
    interval            = 15
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    matcher             = "200"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.storm.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.gateway.arn
  }
}

resource "aws_acm_certificate" "storm" {
  count             = var.domain_name != "" ? 1 : 0
  domain_name       = var.domain_name
  validation_method = "DNS"
}

locals {
  cert_validation = var.domain_name != "" ? tolist(aws_acm_certificate.storm[0].domain_validation_options)[0] : null
}

resource "aws_route53_record" "storm_validation" {
  count   = var.domain_name != "" ? 1 : 0
  zone_id = var.hosted_zone_id
  name    = local.cert_validation.resource_record_name
  type    = local.cert_validation.resource_record_type
  records = [local.cert_validation.resource_record_value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "storm" {
  count                   = var.domain_name != "" ? 1 : 0
  certificate_arn         = aws_acm_certificate.storm[0].arn
  validation_record_fqdns = [aws_route53_record.storm_validation[0].fqdn]
}

resource "aws_lb_listener" "https" {
  count             = var.domain_name != "" ? 1 : 0
  load_balancer_arn = aws_lb.storm.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate_validation.storm[0].certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.gateway.arn
  }
}

resource "aws_autoscaling_group" "storm" {
  name                      = "storm-asg"
  max_size                  = var.asg_max
  min_size                  = var.asg_min
  desired_capacity          = var.asg_desired
  vpc_zone_identifier       = aws_subnet.public[*].id
  health_check_type         = "ELB"
  health_check_grace_period = 120

  launch_template {
    id      = aws_launch_template.storm.id
    version = "$Latest"
  }

  target_group_arns = [aws_lb_target_group.gateway.arn]
  tag {
    key                 = "Name"
    value               = "storm-node"
    propagate_at_launch = true
  }
}
