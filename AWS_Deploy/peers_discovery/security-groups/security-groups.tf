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

variable "east_vpc_id" {
    type = string
}

variable "west_vpc_id" {
    type = string
}

# West Security

resource "aws_security_group" "west_security" {
    provider = aws.us-west-2
    vpc_id = var.west_vpc_id

    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        from_port = 6400
        to_port = 6500
        protocol = "tcp"
        cidr_blocks = [ "10.1.0.0/24", "10.0.0.0/24" ]
    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = [ "0.0.0.0/0" ]
    }

    tags = {
        "Name" = "gncfd_west_security",
    }
}

# East security

resource "aws_security_group" "east_security" {
    provider = aws.us-east-1
    vpc_id = var.east_vpc_id

    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        from_port = 6400
        to_port = 6500
        protocol = "tcp"
        cidr_blocks = [ "10.1.0.0/24", "10.0.0.0/24" ]
    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = [ "0.0.0.0/0" ]
    }

    tags = {
        "Name" = "gncfd_east_security",
    }
}

# Outputs

output "west_security_id" {
  value = aws_security_group.west_security.id
}

output "east_security_id" {
  value = aws_security_group.east_security.id
}