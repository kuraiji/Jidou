resource "aws_ecs_cluster" "cluster" {
  name = var.cluster_name
}

resource "aws_iam_role" "role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
  })

  tags = {
    tag-key = "jidou"
  }
}

resource "aws_ecs_task_definition" "task" {
  family = "jidou-backend"
  requires_compatibilities = ["EC2"]
  network_mode = "awsvpc"
  task_role_arn = aws_iam_role.role.arn
  container_definitions = jsonencode([
    {
      name = "main"
      image = var.image_uri
      essential = true
      memory = 512
      cpu = 10
      portMappings = [
        {
          containerPort = 8080
          hostPort = 8080
        }
      ]
    }
  ])
}