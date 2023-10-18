package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Set up provider
		provider, err := docker.NewProvider(ctx, "khorne", &docker.ProviderArgs{
			Host: pulumi.String("tcp://192.168.0.200:2375"),
		})

		if err != nil {
			return err
		}

		// Get the Dozzle image
		dozzleImage, err := docker.NewRemoteImage(ctx, "dozzle_image", &docker.RemoteImageArgs{
			Name: pulumi.String("amir20/dozzle:latest"),
		})

		if err != nil {
			return err
		}

		// Create the Dozzle container
		dozzleContainer, err := docker.NewContainer(ctx, "dozzle_container", &docker.ContainerArgs{
			Image: dozzleImage.ImageId,
			Volumes: docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/run/docker.sock"),
					HostPath:      pulumi.String("/var/run/docker.sock"),
					ReadOnly:      pulumi.Bool(true),
				},
			},
			Ports: docker.ContainerPortArray{
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(8080),
					External: pulumi.Int(8101),
				},
			},
			Envs: pulumi.StringArray{
				pulumi.String("PUID=1000"),
				pulumi.String("PGID=1000"),
			},
			Labels: &docker.ContainerLabelArray{
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.enable"),
					Value: pulumi.String("true"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.dozzle.entrypoints"),
					Value: pulumi.String("http"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.dozzle.rule"),
					Value: pulumi.String("Host(`dozzle.home`)"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.services.dozzle.loadbalancer.server.port"),
					Value: pulumi.String("8080"),
				},
			},
		}, pulumi.Provider(provider))

		if err != nil {
			return err
		}

		// Export the container's information
		ctx.Export("dozzle_container_id", dozzleContainer.ID())

		return nil
	})
}
