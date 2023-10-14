import pulumi
import pulumi_docker as docker


def definition(provider: docker.Provider, network: docker.Network):
    stack_prefix = "gitea"
    stack_domain = "gitea.home"
    stack_public_port = 8103
    stack_service_gitea = "gitea"

    try:
        # Network to connect containers
        local_network = docker.Network(stack_prefix + "_network",
                                       opts=pulumi.ResourceOptions(provider=provider))

        # Pull a remote Docker image
        image = docker.RemoteImage(stack_prefix + "_image",
                                   name="gitea/gitea:1.20.4-rootless")
        database_image = docker.RemoteImage(stack_prefix + "_database_image",
                                            name="postgres:14")

        # Create volumes
        database_volume = docker.Volume(stack_prefix + "_database_volume",
                                        name="gitea_database",
                                        driver="local",
                                        opts=pulumi.ResourceOptions(provider=provider))
        config_volume = docker.Volume(stack_prefix + "_config_volume",
                                      name="gitea_config",
                                      driver="local",
                                      opts=pulumi.ResourceOptions(provider=provider))
        data_volume = docker.Volume(stack_prefix + "_data_volume",
                                    name="gitea_data",
                                    driver="local",
                                    opts=pulumi.ResourceOptions(provider=provider))

        # Start a container
        database_container = docker.Container(stack_prefix + "_database_container",
                                              name=stack_prefix + "_database_container",
                                              image=database_image.name,
                                              opts=pulumi.ResourceOptions(provider=provider),
                                              networks_advanced=[
                                                  docker.ContainerNetworksAdvancedArgs(name=local_network.name),
                                              ],
                                              volumes=[{
                                                  "volume_name": database_volume.name,
                                                  "container_path": "/var/lib/postgresql/data",
                                              }],
                                              envs=[
                                                  f"POSTGRES_USER=gitea",
                                                  f"POSTGRES_PASSWORD=gitea",
                                                  f"POSTGRES_DB=gitea",
                                              ])
        container = docker.Container(stack_prefix + "_container",
                                     name=stack_prefix + "_container",
                                     image=image.name,
                                     opts=pulumi.ResourceOptions(provider=provider, depends_on=[database_container]),
                                     networks_advanced=[
                                         docker.ContainerNetworksAdvancedArgs(name=network.name),
                                         docker.ContainerNetworksAdvancedArgs(name=local_network.name)
                                     ],
                                     volumes=[{
                                         "volume_name": config_volume.name,
                                         "container_path": "/etc/gitea",
                                     }, {
                                         "volume_name": data_volume.name,
                                         "container_path": "/var/lib/gitea",
                                     }, {
                                         "host_path": "/etc/timezone",
                                         "container_path": "/etc/timezone",
                                         "read_only": True,
                                     }, {
                                         "host_path": "/etc/localtime",
                                         "container_path": "/etc/localtime",
                                         "read_only": True,
                                     }],
                                     ports=[{
                                         "internal": 3000,
                                         "external": stack_public_port,
                                     }],
                                     envs=[
                                         f"GITEA__database__DB_TYPE=postgres",
                                         f"GITEA__database__HOST=gitea_database_container:5432",
                                         f"GITEA__database__NAME=gitea",
                                         f"GITEA__database__USER=gitea",
                                         f"GITEA__database__PASSWD=gitea",
                                     ],
                                     labels=[{
                                         "label": "traefik.enable",
                                         "value": "true",
                                     }, {
                                         "label": "traefik.http.routers." + stack_service_gitea + ".entrypoints",
                                         "value": "http",
                                     }, {
                                         "label": "traefik.http.routers." + stack_service_gitea + ".rule",
                                         "value": "Host(`" + stack_domain + "`)",
                                     }, {
                                         "label": "traefik.http.services." + stack_service_gitea + ".loadbalancer.server.port",
                                         "value": "3000",
                                     }])

        # Export container information
        pulumi.export(name=stack_prefix + "_container_id", value=container.id)
        pulumi.export(name=stack_prefix + "_database_container_id", value=database_container.id)

    except Exception as ex:
        pulumi.log.error("Error creating infrastructure: " + str(ex))
