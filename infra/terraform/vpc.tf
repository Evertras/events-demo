resource "aws_vpc" "events-demo" {
  cidr_block = "10.0.0.0/16"

  tags = "${
    map(
      "Name", "events-demo-eks-node",
      "kubernetes.io/cluster/${var.cluster-name}", "shared",
    )
  }"
}

resource "aws_subnet" "events-demo" {
  vpc_id            = "${aws_vpc.events-demo.id}"
  availability_zone = "${data.aws_availability_zones.available.names[count.index]}"

  count      = 2
  cidr_block = "10.0.${count.index}.0/24"

  tags = "${
    map(
      "Name", "events-demo-eks-node",
      "kubernetes.io/cluster/${var.cluster-name}", "shared",
    )
  }"
}

resource "aws_internet_gateway" "events-demo" {
  vpc_id = "${aws_vpc.events-demo.id}"

  tags = {
    Name = "events-demo"
  }
}

resource "aws_route_table" "events-demo" {
  vpc_id = "${aws_vpc.events-demo.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.events-demo.id}"
  }
}

resource "aws_route_table_association" "events-demo" {
  count = 2

  subnet_id      = "${aws_subnet.events-demo.*.id[count.index]}"
  route_table_id = "${aws_route_table.events-demo.id}"
}
