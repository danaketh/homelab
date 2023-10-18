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

		// Get the Traefik image
		traefikImage, err := docker.NewRemoteImage(ctx, "traefik_image", &docker.RemoteImageArgs{
			Name: pulumi.String("traefik:v3.0"),
		})

		if err != nil {
			return err
		}

		// Create the Traefik container
		traefikContainer, err := docker.NewContainer(ctx, "traefik_container", &docker.ContainerArgs{
			Image:       traefikImage.ImageId,
			Name:        pulumi.String("traefik"),
			Restart:     pulumi.String("always"),
			NetworkMode: pulumi.String("host"),
			Volumes: docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/run/docker.sock"),
					HostPath:      pulumi.String("/var/run/docker.sock"),
					ReadOnly:      pulumi.Bool(true),
				},
			},
			Ports: docker.ContainerPortArray{
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(80),
					External: pulumi.Int(80),
					Ip:       pulumi.String("0.0.0.0"),
					Protocol: pulumi.String("tcp"),
				},
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(443),
					External: pulumi.Int(443),
					Ip:       pulumi.String("0.0.0.0"),
					Protocol: pulumi.String("tcp"),
				},
			},
			Envs: pulumi.StringArray{
				pulumi.String("TRAEFIK_API_INSECURE=true"),
				pulumi.String("TRAEFIK_API_DASHBOARD=true"),
				pulumi.String("TRAEFIK_GLOBAL_CHECKNEWVERSION=true"),
				pulumi.String("TRAEFIK_GLOBAL_SENDANONYMOUSUSAGE=true"),
				pulumi.String("TRAEFIK_ENTRYPOINTS_HTTP_ADDRESS=:80"),
				pulumi.String("TRAEFIK_ENTRYPOINTS_HTTPS_ADDRESS=:443"),
				pulumi.String("TRAEFIK_PROVIDERS_DOCKER=true"),
				pulumi.String("TRAEFIK_PROVIDERS_DOCKER_EXPOSEDBYDEFAULT=false"),
			},
		}, pulumi.Provider(provider))

		if err != nil {
			return err
		}

		// Export the container's information
		ctx.Export("traefik_container_id", traefikContainer.ID())

		return nil
	})
}
