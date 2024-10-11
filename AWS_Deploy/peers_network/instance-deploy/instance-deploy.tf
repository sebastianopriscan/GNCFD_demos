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

variable "vpc_us_east_id" {
  type = string
}

variable "vpc_us_east_subn_id" {
  type = string
}

variable "vpc_us_west_id" {
  type = string
}

variable "vpc_us_west_subn_id" {
  type = string
}

variable "demo_west_ami" {
    type = string
    default = "ami-0b3cb98570977cf06"
}

variable "demo_east_ami" {
    type = string
    default = "ami-02cfbd070b4ad04cb"
}

variable "west_ssh_key" {
    type = string
    default = "oregon_access_key"
}

variable "east_ssh_key" {
    type = string
    default = "first_ec2_pair"
}

variable "demo_server_private_ip" {
  type = string
  default = "10.0.0.10"
}

variable "demo_client_serv_port" {
  type = string
  default = "6405"
}

variable "demo_discoverer_port" {
  type = string
  default = "6401"
}

variable "demo_client_west_replicas" {
    type = number
    default = 4
}

variable "demo_client_east_replicas" {
    type = number
    default = 4
}

variable "west_security_id" {
    type = string
}

variable "east_security_id" {
  type = string
}

resource "aws_instance" "demo_server" {
    provider = aws.us-east-1
    ami = var.demo_east_ami
    instance_type = "t2.micro"
    subnet_id = var.vpc_us_east_subn_id
    private_ip = var.demo_server_private_ip
    security_groups = [ var.east_security_id ]
    associate_public_ip_address = true
    key_name = var.east_ssh_key

    tags = {
      "Name" = "gncfd_demo_network_discoverer",
    }

    user_data = <<-EOF
      #!/bin/bash
      echo "export DISCOVERER_PORT=${var.demo_discoverer_port}" >> /home/ec2-user/.bashrc
      sudo systemctl start discoverer_network.service
    EOF
}

resource "aws_instance" "demo_client_east" {
    provider = aws.us-east-1
    ami = var.demo_east_ami
    count = var.demo_client_east_replicas
    instance_type = "t2.micro"
    subnet_id = var.vpc_us_east_subn_id
    security_groups = [ var.east_security_id ]
    associate_public_ip_address = true
    key_name = var.east_ssh_key

    tags = {
      "Name" = "gncfd_demo_network_peer",
    }

    user_data = <<-EOF
      #!/bin/bash
      echo "export CLIENT_SERV_PORT=${var.demo_client_serv_port}" >> /home/ec2-user/.bashrc
      echo "export DISCOVERER_ADDR=${var.demo_server_private_ip}" >> /home/ec2-user/.bashrc
      echo "export DISCOVERER_PORT=${var.demo_discoverer_port}" >> /home/ec2-user/.bashrc
      sudo systemctl start client_network.service
      EOF
}

resource "aws_instance" "demo_client_west" {
    provider = aws.us-west-2
    count = var.demo_client_west_replicas
    ami = var.demo_west_ami
    instance_type = "t2.micro"
    subnet_id = var.vpc_us_west_subn_id
    security_groups = [ var.west_security_id ]
    associate_public_ip_address = true
    key_name = var.west_ssh_key

    tags = {
      "Name" = "gncfd_demo_network_peer",
    }

    user_data = <<-EOF
      #!/bin/bash
      echo "export CLIENT_SERV_PORT=${var.demo_client_serv_port}" >> /home/ec2-user/.bashrc
      echo "export DISCOVERER_ADDR=${var.demo_server_private_ip}" >> /home/ec2-user/.bashrc
      echo "export DISCOVERER_PORT=${var.demo_discoverer_port}" >> /home/ec2-user/.bashrc
      sudo systemctl start client_network.service
      EOF
}
