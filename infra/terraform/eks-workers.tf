resource "aws_iam_role" "events-demo-node" {
  name = "events-demo-node"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "events-demo-node-AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = "${aws_iam_role.events-demo-node.name}"
}

resource "aws_iam_role_policy_attachment" "events-demo-node-AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = "${aws_iam_role.events-demo-node.name}"
}

resource "aws_iam_role_policy_attachment" "events-demo-node-AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = "${aws_iam_role.events-demo-node.name}"
}

resource "aws_iam_instance_profile" "events-demo-node" {
  name = "terraform-eks-demo"
  role = "${aws_iam_role.events-demo-node.name}"
}

resource "aws_security_group" "events-demo-node" {
  name        = "events-demo-node"
  description = "Security group for worker nodes in the EKS cluster"
  vpc_id      = "${aws_vpc.events-demo.id}"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = "${
    map(
      "Name", "events-demo-node",
      "kubernetes.io/cluster/${var.cluster-name}", "owned",
    )
  }"
}

resource "aws_security_group_rule" "events-demo-node-ingress-self" {
  description = "Allows worker nodes to talk to each other"
  from_port   = 0
  to_port     = 65535
  type        = "ingress"
  protocol    = "-1"

  security_group_id        = "${aws_security_group.events-demo-node.id}"
  source_security_group_id = "${aws_security_group.events-demo-node.id}"
}

resource "aws_security_group_rule" "events-demo-node-ingress-cluster" {
  description = "Allow worker nodes to receive communication from control plane"
  from_port   = 1025
  to_port     = 65535
  type        = "ingress"
  protocol    = "tcp"

  security_group_id        = "${aws_security_group.events-demo-node.id}"
  source_security_group_id = "${aws_security_group.events-demo-cluster.id}"
}

resource "aws_security_group_rule" "events-demo-node-ingress-https" {
  description = "Allow pods to communicate with the cluster API server"
  from_port   = 443
  to_port     = 443
  type        = "ingress"
  protocol    = "tcp"

  security_group_id        = "${aws_security_group.events-demo-cluster.id}"
  source_security_group_id = "${aws_security_group.events-demo-node.id}"
}

# EKS currently documents this required userdata for EKS worker nodes to
# properly configure Kubernetes applications on the EC2 instance.
# We implement a Terraform local here to simplify Base64 encoding this
# information into the AutoScaling Launch Configuration.
# More information: https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
locals {
  demo-node-userdata = <<USERDATA
#!/bin/bash
set -o xtrace
/etc/eks/bootstrap.sh --apiserver-endpoint '${aws_eks_cluster.events-demo.endpoint}' --b64-cluster-ca '${aws_eks_cluster.events-demo.certificate_authority.0.data}' '${var.cluster-name}'
USERDATA
}

resource "aws_launch_configuration" "events-demo" {
  associate_public_ip_address = true

  iam_instance_profile = "${aws_iam_instance_profile.events-demo-node.id}"
  image_id             = "${data.aws_ami.eks-worker.id}"
  instance_type        = "m4.large"
  name_prefix          = "events-demo-eks"
  security_groups      = ["${aws_security_group.events-demo-node.id}"]
  user_data_base64     = "${base64encode(local.demo-node-userdata)}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "events-demo-workers" {
  name                 = "events-demo-workers"
  desired_capacity     = 2
  max_size             = 2
  min_size             = 1
  vpc_zone_identifier  = "${aws_subnet.events-demo[*].id}"
  launch_configuration = "${aws_launch_configuration.events-demo.id}"

  tag {
    key   = "Name"
    value = "events-demo-workers"

    propagate_at_launch = true
  }

  tag {
    key   = "kubernetes.io/cluster/${var.cluster-name}"
    value = "owned"

    propagate_at_launch = true
  }
}
