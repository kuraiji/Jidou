resource "aws_ecs_cluster" "cluster" {
  name = var.cluster_name
}

data "aws_iam_policy_document" "assume_ecs_role" {
  statement {
    effect = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["ecs-tasks.amazonaws.com"]
      type = "Service"
    }
  }
}

resource "aws_iam_role" "ecs_role" {
  assume_role_policy = data.aws_iam_policy_document.assume_ecs_role.json
  name = "${var.cluster_name}_ecs_role"
  tags = {
    tag-key = var.cluster_name
  }
}

data "aws_iam_policy_document" "assume_ecs_execution_role" {
  statement {
    effect = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["ecs-tasks.amazonaws.com"]
      type = "Service"
    }
  }
}

resource "aws_iam_role" "ecs_execution_role" {
  assume_role_policy = data.aws_iam_policy_document.assume_ecs_execution_role.json
  name = "${var.cluster_name}_ecs_execution_role"
  tags = {
    tag-key = var.cluster_name
  }
}

data "aws_iam_policy_document" "ecs_policy_document" {
  statement {
    sid    = "DsqlBasicPermissions"
    effect = "Allow"
    actions = [
      "dsql:GetCluster",
      "dsql:UpdateCluster",
      "dsql:ListClusters",
      "dsql:DbConnectAdmin",
      "dsql:DbConnect",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "ecs_policy" {
  name = "${var.cluster_name}_ecs_dsql_policy"
  description = "The role for ecs role that will give the container access to the DSQL cluster."
  policy = data.aws_iam_policy_document.ecs_policy_document.json
}

resource "aws_iam_role_policy_attachment" "ecs_attach_dsql" {
  //policy_arn = aws_iam_policy.ecs_policy.arn
  policy_arn = "arn:aws:iam::aws:policy/AmazonAuroraDSQLFullAccess"
  role       = aws_iam_role.ecs_role.name
}

resource "aws_iam_role_policy_attachment" "ecs_attach_ssm" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMFullAccess"
  role       = aws_iam_role.ecs_role.name
}

resource "aws_iam_role_policy_attachment" "ecs_attach_execution_ecr" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role = aws_iam_role.ecs_execution_role.name
}

resource "aws_ecs_task_definition" "task" {
  family = "${var.cluster_name}-backend"
  requires_compatibilities = ["EC2"]
  network_mode = "awsvpc"
  task_role_arn = aws_iam_role.ecs_role.arn
  execution_role_arn = aws_iam_role.ecs_execution_role.arn
  container_definitions = jsonencode([
    {
      name = "main"
      image = var.image_uri
      essential = true
      memory = 512
      cpu = 512
      portMappings = [
        {
          containerPort = var.exposed_port
          hostPort = var.exposed_port
        }
      ]
    }
  ])
}

data "aws_iam_policy_document" "assume_ec2_role" {
  statement {
    effect = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["ec2.amazonaws.com"]
      type = "Service"
    }
  }
}

resource "aws_iam_role" "ec2_role" {
  assume_role_policy = data.aws_iam_policy_document.assume_ec2_role.json
  name = "${var.cluster_name}_ec2_instance_role"
  tags = {
    tag-key = var.cluster_name
  }
}

resource "aws_iam_role_policy_attachment" "ec2_attach" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
  role       = aws_iam_role.ec2_role.name
}

resource "aws_iam_instance_profile" "ec2_instance_profile" {
  name = "${var.cluster_name}_profile"
  role = aws_iam_role.ec2_role.name
}

resource "aws_security_group" "ec2_security" {
  name = "${var.cluster_name}_security_group"
  description = "Allow SSH Ingress, App Ports Ingress and Allow Anywhere Egress"
  tags = {
    Name = "${var.cluster_name}_security_group"
  }
}

resource "aws_vpc_security_group_ingress_rule" "ssh_ingress" {
  description = "Allows SSH"
  security_group_id = aws_security_group.ec2_security.id
  cidr_ipv4 = "0.0.0.0/0"
  from_port = 22
  to_port = 22
  ip_protocol = "tcp"
  tags = {
    Name = "SSH"
  }
}

resource "aws_vpc_security_group_ingress_rule" "http_ingress_ec2" {
  description = "Opens Specific App Ports Inbound"
  security_group_id = aws_security_group.ec2_security.id
  cidr_ipv4 = "0.0.0.0/0"
  from_port = var.exposed_port
  to_port = var.exposed_port
  ip_protocol = "tcp"
  tags = {
    Name = "App Ports"
  }
}

resource "aws_vpc_security_group_egress_rule" "ec2_egress" {
  description = "Allows Anywhere Egress"
  security_group_id = aws_security_group.ec2_security.id
  cidr_ipv4 = "0.0.0.0/0"
  ip_protocol = "-1"
  tags = {
    Name = "Egress Anywhere"
  }
}

data "aws_ssm_parameter" "ecs_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2023/recommended"
}

locals {
  ami_info = jsondecode(data.aws_ssm_parameter.ecs_ami.value)
  ami_id = local.ami_info.image_id
}

resource "aws_instance" "ec2" {
  tags = {
    Name = "${var.cluster_name}_ec2"
  }
  ami = local.ami_id
  instance_type = "t2.micro"
  key_name = var.ssh_key_name
  associate_public_ip_address = true
  security_groups = [aws_security_group.ec2_security.name]
  iam_instance_profile = aws_iam_instance_profile.ec2_instance_profile.name
  instance_market_options {
    market_type = "spot"
  }
  user_data = <<EOF
#!/bin/bash
echo ECS_CLUSTER="${var.cluster_name}" >> /etc/ecs/ecs.config
  EOF
}

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_security_group" "app_security" {
  name = "${var.cluster_name}_app_security_group"
  description = "Allow App Required Port Ingress and Allow Anywhere Egress"
  tags = {
    Name = "${var.cluster_name}_security_group"
  }
}

resource "aws_vpc_security_group_ingress_rule" "http_ingress" {
  description = "Opens Specific App Ports Inbound"
  security_group_id = aws_security_group.app_security.id
  cidr_ipv4 = "0.0.0.0/0"
  from_port = var.exposed_port
  to_port = var.exposed_port
  ip_protocol = "tcp"
  tags = {
    Name = "App Ports"
  }
}

resource "aws_vpc_security_group_egress_rule" "app_egress" {
  description = "Allows Anywhere Egress"
  security_group_id = aws_security_group.app_security.id
  cidr_ipv4 = "0.0.0.0/0"
  ip_protocol = "-1"
  tags = {
    Name = "Egress Anywhere"
  }
}

resource "aws_ecs_service" "app" {
  name = "${var.cluster_name}_app"
  cluster = aws_ecs_cluster.cluster.id
  task_definition = aws_ecs_task_definition.task.arn
  scheduling_strategy = "DAEMON"
  network_configuration {
    subnets = data.aws_subnets.default.ids
    security_groups = [aws_security_group.app_security.id]
  }
  depends_on = [aws_instance.ec2]
}