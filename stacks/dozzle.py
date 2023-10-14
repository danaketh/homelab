import pulumi
import pulumi_docker as docker


def definition(provider: docker.Provider, network: docker.Network):
    stack_prefix = "dozzle"
    stack_domain = "dozzle.home"
    stack_public_port = 8102
    stack_service_dozzle = "dozzle"

    try:
        # Pull a remote Docker image
        image = docker.RemoteImage(stack_prefix+"_image",
                                   name="amir20/dozzle:latest")

        # Start a container
        container = docker.Container(stack_prefix+"_container",
                                     name=stack_prefix+"_container",
                                     image=image.name,
                                     opts=pulumi.ResourceOptions(provider=provider),
                                     networks_advanced=[docker.ContainerNetworksAdvancedArgs(name=network.name)],
                                     ports=[{
                                         "internal": 8080,
                                         "external": stack_public_port
                                     }],
                                     volumes=[{
                                         "host_path": "/var/run/docker.sock",
                                         "container_path": "/var/run/docker.sock",
                                         "read_only": True,
                                     }],
                                     envs=[
                                         f"PUID=1000",
                                         f"PGID=1000",
                                     ],
                                     labels=[{
                                         "label": "traefik.enable",
                                         "value": "true",
                                     }, {
                                         "label": "traefik.http.routers."+stack_service_dozzle+".entrypoints",
                                         "value": "http",
                                     }, {
                                         "label": "traefik.http.routers."+stack_service_dozzle+".rule",
                                         "value": "Host(`"+stack_domain+"`)",
                                     }, {
                                         "label": "traefik.http.services."+stack_service_dozzle+".loadbalancer.server.port",
                                         "value": "8080",
                                     }])

        # Export container information
        pulumi.export(name=stack_prefix+"_container_id", value=container.id)

    except Exception as ex:
        pulumi.log.error("Error creating infrastructure: " + str(ex))
