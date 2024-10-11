# GNCFD Demos

There are two demos to run for the GNCFD library:

- The peers discovery demo, in which all peers find about each other's existance through a discovery service, and start gossiping. This demo showcases the convergence of the distributed _Network Coordinate System_

- The peers network demo is similar to the previous one, but the discovery service let's the peers only know about a subset of the others, in order to form a toy overlay network. This demo showcases the gossip algorithm's spreading capabilities

Each of these demos can be run by `cd`-ing in the respective directory, and they are very similar in their setup

## Docker Compose Setup

For the Docker Compose version, building uses the local checked out library, so it is necessary to be on branch master

For the first option, it is necessary to stay on branch master, and simply run:

In repo's root directory:

```bash
docker buildx build -t gncfd_embed:latest -f GNCFD.dockerfile .
```

In the demo of interest's directory:

```bash
docker compose up -d
```

To check the algorithm in actions, just select a peer and get its logs, e.g. :

```bash
docker container logs peers_discovery-peer-1
```

## AWS setup

For AWS, there are, in the AWS_deploy directory, some terraform scripts, but to run them some preliminary actions are needed:

1. Prepare the AMI:
    1. Create an EC2 instance in us-east-1
    2. Install uuidgen
    3. clone this repository, the _using\_release_ branch
    4. Inside the clone's root, run `make install-demos-systemd-services`
    5. This will install the systemd services
    6. Shudown the instance and generate an AMI
    7. Copy the AMI to the us-west-2 region
    7. Generate key pairs for the us-east-1 and us-west-2 regions
    8. Change the key names and ami IDs in the `instance-deploy.tf` file in both the demos terraform directories

2. Select the demo by `cd`-ing in its directory
3. Run `terraform init` and `terraform apply`
4. Check an instance's status by running `journalctl -u client_network.service` or `journalctl -u client_discovery.service` based on the chosen demo to run