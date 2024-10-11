# Configure AWS us-east-1 provider
provider "aws" {
  alias = "us-east-1"
  region = "us-east-1" # Example region for the first VPC
}

# Configure AWS second region provider
provider "aws" {
  alias = "us-west-2"
  region = "us-west-2"
}

# Create the first VPC
resource "aws_vpc" "vpc-us-east-1" {
  provider = aws.us-east-1
  cidr_block = "10.0.0.0/16"
  tags = {
    "Name" = "gncfd_east_vpc",
  }
}

# Create a subnet in the first VPC
resource "aws_subnet" "subnet1" {
  provider = aws.us-east-1
  vpc_id = aws_vpc.vpc-us-east-1.id
  cidr_block = "10.0.0.0/24"
  availability_zone = "us-east-1a"
  tags = {
    "Name" = "gncfd_east_subnet",
  }
}

# Create the second VPC
resource "aws_vpc" "vpc-us-west-2" {
  provider = aws.us-west-2
  cidr_block = "10.1.0.0/16"
  tags = {
    "Name" = "gncfd_west_vpc",
  }
}

# Create a subnet in the second VPC
resource "aws_subnet" "subnet2" {
  provider = aws.us-west-2
  vpc_id = aws_vpc.vpc-us-west-2.id
  cidr_block = "10.1.0.0/24"
  availability_zone = "us-west-2a"
  tags = {
    "Name" = "gncfd_west_subnet",
  }
}

# Create a VPC peering connection
resource "aws_vpc_peering_connection" "peering" {
  vpc_id = aws_vpc.vpc-us-east-1.id
  peer_vpc_id = aws_vpc.vpc-us-west-2.id
  peer_region = "us-west-2"
  tags = {
    "Name" = "gncfd_peering",
  }
}

# Accept the VPC peering connection request
resource "aws_vpc_peering_connection_accepter" "accepter" {
  provider = aws.us-west-2
  vpc_peering_connection_id = aws_vpc_peering_connection.peering.id
  auto_accept = true
  tags = {
    "Name" = "gncfd_peering_accepter",
  }
}

# Connection of the clouds to the internet

resource "aws_internet_gateway" "vpc_us_east_gw" {
  provider = aws.us-east-1
  vpc_id = aws_vpc.vpc-us-east-1.id
  tags = {
    "Name" = "gncfd_east_vpc_igw",
  }
}

resource "aws_internet_gateway" "vpc_us_west_gw" {  
  provider = aws.us-west-2
  vpc_id = aws_vpc.vpc-us-west-2.id
  tags = {
    "Name" = "gncfd_west_vpc_igw",
  }
}

# Interconnection of the clouds

resource "aws_route_table" "us-east-table" {
  provider = aws.us-east-1
  vpc_id = aws_vpc.vpc-us-east-1.id

  route {
    cidr_block = "10.1.0.0/24"
    gateway_id = aws_vpc_peering_connection.peering.id
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vpc_us_east_gw.id
  }

  tags = {
    "Name" = "gncfd_east_vpc_table",
  }
}

resource "aws_route_table_association" "subnet1-assosiation" {
  provider = aws.us-east-1
  route_table_id = aws_route_table.us-east-table.id
  subnet_id = aws_subnet.subnet1.id
}

resource "aws_route_table" "us-west-table" {
  provider = aws.us-west-2
  vpc_id = aws_vpc.vpc-us-west-2.id

  route {
    cidr_block = "10.0.0.0/24"
    gateway_id = aws_vpc_peering_connection.peering.id
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.vpc_us_west_gw.id
  }

  tags = {
    "Name" = "gncfd_west_vpc_table",
  }
}

resource "aws_route_table_association" "subnet2-assosiation" {
  provider = aws.us-west-2
  route_table_id = aws_route_table.us-west-table.id
  subnet_id = aws_subnet.subnet2.id
}

# Outputs

output "vpc_us_east_id" {
  value = aws_vpc.vpc-us-east-1.id
}

output "vpc_us_west_id" {
  value = aws_vpc.vpc-us-west-2.id
}

output "subnet_us_east_id" {
  value = aws_subnet.subnet1.id
}

output "subnet_us_west_id" {
  value = aws_subnet.subnet2.id
}