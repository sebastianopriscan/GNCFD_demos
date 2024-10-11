module "vpc" {
  source = "./vpc-deploy"
}

module "security-groups" {
    source = "./security-groups"
    west_vpc_id = module.vpc.vpc_us_west_id
    east_vpc_id = module.vpc.vpc_us_east_id
}

module "ec2" {
  source = "./instance-deploy"
  vpc_us_east_id = module.vpc.vpc_us_east_id
  vpc_us_east_subn_id = module.vpc.subnet_us_east_id
  vpc_us_west_id = module.vpc.vpc_us_west_id
  vpc_us_west_subn_id = module.vpc.subnet_us_west_id
  west_security_id = module.security-groups.west_security_id
  east_security_id = module.security-groups.east_security_id
}