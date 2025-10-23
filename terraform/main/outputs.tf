# DNS público del ALB (frontend y API)
output "alb_dns" {
  value       = aws_lb.main.dns_name
  description = "load balancer DNS"
}


output "redis_endpoint" {
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
  description = "Redis cluster endpoint"
}

output "redis_port" {
  value       = tostring(aws_elasticache_cluster.redis.cache_nodes[0].port)
  description = "Redis cluster port"
}
