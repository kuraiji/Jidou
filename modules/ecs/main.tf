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

resource "aws_iam_role_policy_attachment" "ecs_attach" {
  policy_arn = aws_iam_policy.ecs_policy.arn
  role       = aws_iam_role.ecs_role.name
}

resource "aws_ecs_task_definition" "task" {
  family = "${var.cluster_name}-backend"
  requires_compatibilities = ["EC2"]
  network_mode = "awsvpc"
  task_role_arn = aws_iam_role.ecs_role.arn
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
  skip_destroy = true
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
//ami-04b6b70c8a2ea3023
resource "aws_iam_instance_profile" "ec2_instance_profile" {
  name = "${var.cluster_name}_profile"
  role = aws_iam_role.ec2_role.name
}