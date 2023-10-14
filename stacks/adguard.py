import pulumi
import pulumi_docker as docker


def definition(provider: docker.Provider, network: docker.Network):
    try:
        # Pull a remote Docker image
        image = docker.RemoteImage("adguard_image",
                                   name="adguard/adguardhome:latest")

        # Start a container
        container = docker.Container("adguard_container",
                                     name="adguard_container",
                                     image=image.name,
                                     opts=pulumi.ResourceOptions(provider=provider),
                                     networks_advanced=[docker.ContainerNetworksAdvancedArgs(name=network.name)],
                                     dns=[
                                         "127.0.0.1",
                                         "94.140.14.14",  # AdGuard DNS
                                         "94.140.15.15",  # AdGuard DNS
                                         "1.1.1.1",  # Cloudflare
                                         "1.0.0.1",  # Cloudflare
                                         "8.8.8.8",  # Google
                                         "8.8.4.4",  # Google
                                         "76.76.2.0",  # Control D
                                         "76.76.10.0",  # Control D
                                         "9.9.9.9",  # Quad9
                                         "149.112.112.112",  # Quad9
                                         "208.67.222.222",  # OpenDNS
                                         "208.67.220.220",  # OpenDNS
                                         "185.228.168.9",  # CleanBrowsing
                                         "185.228.169.9",  # CleanBrowsing
                                         "76.76.19.19",  # Alternate DNS
                                         "76.223.122.150",  # Alternate DNS
                                     ],
                                     ports=[{
                                         "internal": 53,
                                         "external": 53,
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 53,
                                         "external": 53,
                                         "protocol": "udp",
                                     }, {
                                         "internal": 67,
                                         "external": 67,
                                         "protocol": "udp",
                                     }, {
                                         "internal": 68,
                                         "external": 68,
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 68,
                                         "external": 68,
                                         "protocol": "udp",
                                     }, {
                                         "internal": 784,
                                         "external": 784,
                                         "protocol": "udp",
                                     }, {
                                         "internal": 853,
                                         "external": 853,
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 853,
                                         "external": 853,
                                         "protocol": "udp",
                                     }, {  # Web UI installation
                                         "internal": 3000,
                                         "external": 3000,
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 5443,
                                         "external": 5443,
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 5443,
                                         "external": 5443,
                                         "protocol": "udp",
                                     }, {
                                         "internal": 80,
                                         "external": 8101,
                                         "protocol": "tcp",
                                     }, {
                                         "internal": 8853,
                                         "external": 8853,
                                         "protocol": "udp",
                                     }],
                                     envs=[
                                         f"PUID=1000",
                                         f"PGID=1000",
                                     ],
                                     labels=[{
                                         "label": "traefik.enable",
                                         "value": "true",
                                     }, {
                                         "label": "traefik.http.routers.adguard.entrypoints",
                                         "value": "http",
                                     }, {
                                         "label": "traefik.http.routers.adguard.rule",
                                         "value": "Host(`adguard.home`)",
                                     }, {
                                         "label": "traefik.http.services.adguard.loadbalancer.server.port",
                                         "value": "80",
                                     }])

        # Export container information
        pulumi.export(name="adguard_container_id", value=container.id)

    except Exception as ex:
        pulumi.log.error("Error creating infrastructure: " + str(ex))
