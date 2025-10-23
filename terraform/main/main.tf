provider "aws" {
  region = var.aws_region
}


resource "aws_iam_role" "ecs_execution_role" {
  name = "stockvision-ecs-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "sts:AssumeRole"
        ],
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        },
        Effect = "Allow",
        Sid    = ""
      }
    ]
  })
}


resource "aws_iam_role_policy_attachment" "ecs_execution_role_policy" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# -----------------------------
# Security Groups
# -----------------------------
# SG para ALB 
resource "aws_security_group" "alb_sg" {
  name        = "stockvision-alb-sg"
  description = "ALB security group"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP from anywhere"
    from_port   = 80
    to_port     = 80
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

# SG to ECS tasks (API + APP)
resource "aws_security_group" "ecs_sg" {
  name        = "stockvision-ecs-sg"
  description = "ECS tasks security group"
  vpc_id      = aws_vpc.main.id

  # allow ALB to ECS
  ingress {
    description      = "ALB to API"
    from_port        = 8080
    to_port          = 8080
    protocol         = "tcp"
    security_groups  = [aws_security_group.alb_sg.id]
  }

  ingress {
    description      = "ALB to APP"
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    security_groups  = [aws_security_group.alb_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}



# -----------------------------
# VPC and subnets 
# -----------------------------
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "public_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "us-east-1b"
  map_public_ip_on_launch = true
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "stockvision-igw"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "stockvision-public-rt"
  }
}

resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_b" {
  subnet_id      = aws_subnet.public_b.id
  route_table_id = aws_route_table.public.id
}

# -----------------------------
# Redis
# -----------------------------
resource "aws_security_group" "redis_sg" {
  name        = "redis-security-group"
  description = "Security group for Redis cluster"
    vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    security_groups = [aws_security_group.ecs_sg.id] 
  }
}

# Subnet group for ElastiCache
resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "stockvision-redis-subnet-group-2"
  subnet_ids = [aws_subnet.public_a.id, aws_subnet.public_b.id]
}


resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "stockvision-redis"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  engine_version       = "7.0"
  apply_immediately    = true
  port                 = 6379

  subnet_group_name  = aws_elasticache_subnet_group.redis_subnet_group.name
  security_group_ids = [aws_security_group.redis_sg.id]

  tags = {
    Name = "stockvision-redis"
  }
}


# -----------------------------
# ECS Cluster
# -----------------------------
resource "aws_ecs_cluster" "main" {
  name = "stockvision-cluster"
}

# -----------------------------
# ALB
# -----------------------------
resource "aws_lb" "main" {
  name               = "stockvision-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.public_a.id, aws_subnet.public_b.id]
}

# Target Group API
resource "aws_lb_target_group" "api_tg" {
  name        = "stockvision-api-tg"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id

  health_check {
    path                = "/"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-399"
  }
}

# Target Group APP
resource "aws_lb_target_group" "app_tg" {
  name        = "stockvision-app-tg"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id

  health_check {
    path                = "/"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-399"
  }
}

# Listener HTTP 80
resource "aws_lb_listener" "http_listener" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app_tg.arn
  }
}

# Listener para API path /api
resource "aws_lb_listener_rule" "api_rule" {
  listener_arn = aws_lb_listener.http_listener.arn
  priority     = 10

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_tg.arn
  }

  condition {
    path_pattern {
      values = ["/api/*"]
    }
  }
}


resource "aws_cloudwatch_log_group" "client" {
  name              = "/ecs/stockvision-client"
  retention_in_days = 14
}

resource "aws_cloudwatch_log_group" "api" {
  name              = "/ecs/stockvision-api"
  retention_in_days = 14
}


# -----------------------------
# ECS Task Definitions
# -----------------------------
# API
resource "aws_ecs_task_definition" "api_task" {
  family                   = "stockvision-api-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "stockvision-api"
      image     = var.ecr_api_image
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      environment = [
        { name = "ENV"     , value = "production" },
        { name = "DB_HOST"     , value = var.db_host },
        { name = "DB_NAME"     , value = var.db_name },
        { name = "DB_USER"     , value = var.db_user },
        { name = "DB_PASSWORD" , value = var.db_password },
        { name = "DB_PORT"     , value = var.db_port },
        { name = "DB_SSL"     , value = var.db_ssl },
        { name = "DB_SCHEMA"      , value = "stocksvision" },
        { name = "REDIS_HOST",  value = aws_elasticache_cluster.redis.cache_nodes[0].address },
        { name = "REDIS_PORT",  value = tostring(aws_elasticache_cluster.redis.cache_nodes[0].port) }, 
        { name = "CLIENT_HOST",  value = "http://${aws_lb.main.dns_name}" },
        { name = "LOG_LEVEL",  value = "INFO" },
        { name = "STOCK_API_URL",  value = var.stock_api_url },
        { name = "STOCK_API_TOKEN",  value = var.stock_api_token },
        { name = "FINANCIAL_BASE_URL",  value = var.financial_base_url },
        { name = "FINANCIAL_TOKEN",  value = var.financial_token },
        { name = "FINHUB_BASE_URL",  value = var.finhub_base_url },
        { name = "FINHUB_TOKEN",  value = var.finhub_token },
        { name = "GEMINI_API_KEY",  value = var.gemini_api_key }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/stockvision-api"
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

# Cliente
resource "aws_ecs_task_definition" "client_task" {
  family                   = "stockvision-client-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "stockvision-client"
      image     = var.ecr_app_image
      essential = true
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
        }
      ]
      environment = [
        { name = "VITE_API_URL", value = "http://${aws_lb.main.dns_name}" }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/stockvision-client"
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])
}

# -----------------------------
# ECS Services
# -----------------------------
# API Service
resource "aws_ecs_service" "api_service" {
  name            = "stockvision-api-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.public_a.id, aws_subnet.public_b.id]
    security_groups = [aws_security_group.ecs_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api_tg.arn
    container_name   = "stockvision-api"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener_rule.api_rule]
}

# Cliente Service
resource "aws_ecs_service" "client_service" {
  name            = "stockvision-client-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.client_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.public_a.id, aws_subnet.public_b.id]
    security_groups = [aws_security_group.ecs_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.app_tg.arn
    container_name   = "stockvision-client"
    container_port   = 80
  }

  depends_on = [aws_lb_listener.http_listener]
}
