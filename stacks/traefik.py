import pulumi
import pulumi_docker as docker


def definition(provider: docker.Provider, network: docker.Network):
    try:
        # Pull a remote Docker image
        image = docker.RemoteImage("traefik_image",
                                   name="traefik:latest")

        # Start a container
        container = docker.Container("traefik_container",
                                     name="traefik_container",
                                     image=image.name,
                                     opts=pulumi.ResourceOptions(provider=provider),
                                     networks_advanced=[docker.ContainerNetworksAdvancedArgs(name=network.name)],
                                     volumes=[{
                                         "host_path": "/var/run/docker.sock",
                                         "container_path": "/var/run/docker.sock",
                                         "read_only": True,
                                     }],
                                     ports=[{
                                         "internal": 80,
                                         "external": 80,
                                     }, {
                                         "internal": 443,
                                         "external": 443,
                                     }, {
                                         "internal": 8080,
                                         "external": 8100,
                                     }],
                                     envs=[
                                         f"TRAEFIK_API_INSECURE=true",
                                         f"TRAEFIK_API_DASHBOARD=true",
                                         f"TRAEFIK_GLOBAL_CHECKNEWVERSION=true",
                                         f"TRAEFIK_GLOBAL_SENDANONYMOUSUSAGE=true",
                                         f"TRAEFIK_ENTRYPOINTS_HTTP_ADDRESS=:80",
                                         f"TRAEFIK_ENTRYPOINTS_HTTPS_ADDRESS=:443",
                                         f"TRAEFIK_PROVIDERS_DOCKER=true",
                                         f"TRAEFIK_PROVIDERS_DOCKER_EXPOSEDBYDEFAULT=false",
                                     ])

        # Export container information
        pulumi.export(name="traefik_container_id", value=container.id)

    except Exception as ex:
        pulumi.log.error("Error creating infrastructure: " + str(ex))
