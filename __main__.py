import pulumi
import pulumi_docker as docker
from stacks import traefik

# Set up provider
provider = docker.Provider(
    'khorne',
    host="tcp://192.168.0.200:2375"
)

# Create network for containers
network = docker.Network(
    'global_network',
    opts=pulumi.ResourceOptions(provider=provider)
)

# Create containers
traefik.definition(provider, network)
