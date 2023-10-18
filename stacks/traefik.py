import pulumi
import pulumi_docker as docker


def definition(provider: docker.Provider, network: docker.Network):
    try:
        # Pull a remote Docker image
        image = docker.RemoteImage("traefik_image",
                                   name="traefik:v3.0")

        # Start a container
        container = docker.Container("traefik_container",
                                     name="traefik_container",
                                     image=image.name,
                                     opts=pulumi.ResourceOptions(provider=provider),
                                     network_mode="host",
                                     volumes=[{
                                         "host_path": "/var/run/docker.sock",
                                         "container_path": "/var/run/docker.sock",
                                         "read_only": True,
                                     }],
                                     ports=[{
                                         "internal": 80,
                                         "external": 80,
                                         "ip": "0.0.0.0",
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 443,
                                         "external": 443,
                                         "ip": "0.0.0.0",
                                         "protocol": "tcp",
                                     }],
                                     envs=[
                                         "TRAEFIK_API_INSECURE=true",
                                         "TRAEFIK_API_DASHBOARD=true",
                                         "TRAEFIK_GLOBAL_CHECKNEWVERSION=true",
                                         "TRAEFIK_GLOBAL_SENDANONYMOUSUSAGE=true",
                                         "TRAEFIK_ENTRYPOINTS_HTTP_ADDRESS=:80",
                                         "TRAEFIK_ENTRYPOINTS_HTTPS_ADDRESS=:443",
                                         "TRAEFIK_PROVIDERS_DOCKER=true",
                                         "TRAEFIK_PROVIDERS_DOCKER_EXPOSEDBYDEFAULT=false",
                                     ])

        # Export container information
        pulumi.export(name="traefik_container_id", value=container.id)

    except Exception as ex:
        pulumi.log.error("Error creating infrastructure: " + str(ex))
