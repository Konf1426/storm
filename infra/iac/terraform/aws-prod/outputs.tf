output "alb_dns_name" {
  value       = aws_lb.storm.dns_name
  description = "ALB public DNS"
}

output "vpc_id" {
  value       = aws_vpc.storm.id
  description = "VPC id"
}

output "rds_endpoint" {
  value       = aws_db_instance.storm.address
  description = "Postgres endpoint"
}

output "redis_endpoint" {
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
  description = "Redis endpoint"
}

output "ecr_gateway" {
  value       = aws_ecr_repository.gateway.repository_url
  description = "ECR repo for gateway"
}

output "ecr_messages" {
  value       = aws_ecr_repository.messages.repository_url
  description = "ECR repo for messages"
}
