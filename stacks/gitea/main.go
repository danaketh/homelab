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

		// Get images
		giteaImage, err := docker.NewRemoteImage(ctx, "gitea_image", &docker.RemoteImageArgs{
			Name: pulumi.String("gitea/gitea:1.20.4-rootless"),
		})
		if err != nil {
			return err
		}

		databaseImage, err := docker.NewRemoteImage(ctx, "gitea_database_image", &docker.RemoteImageArgs{
			Name: pulumi.String("postgres:14"),
		})
		if err != nil {
			return err
		}

		// Create networks
		network, err := docker.NewNetwork(ctx, "gitea_network", &docker.NetworkArgs{
			Name: pulumi.String("gitea_network"),
		}, pulumi.Provider(provider))

		// Create database container
		databaseContainer, err := docker.NewContainer(ctx, "gitea_database_container", &docker.ContainerArgs{
			Image:   databaseImage.RepoDigest,
			Name:    pulumi.String("gitea_database"),
			Restart: pulumi.String("always"),
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/lib/postgresql/data"),
					HostPath:      pulumi.String("/opt/gitea/database"),
				},
			},
			NetworksAdvanced: &docker.ContainerNetworksAdvancedArray{
				&docker.ContainerNetworksAdvancedArgs{
					Name: network.Name,
				},
			},
			Envs: pulumi.StringArray{
				pulumi.String("POSTGRES_USER=gitea"),
				pulumi.String("POSTGRES_PASSWORD=gitea"),
				pulumi.String("POSTGRES_DB=gitea"),
			},
		}, pulumi.Provider(provider))
		// Create Gitea container
		giteaContainer, err := docker.NewContainer(ctx, "gitea_container", &docker.ContainerArgs{
			Image:   giteaImage.RepoDigest,
			Name:    pulumi.String("gitea"),
			Restart: pulumi.String("always"),
			NetworksAdvanced: &docker.ContainerNetworksAdvancedArray{
				&docker.ContainerNetworksAdvancedArgs{
					Name: network.Name,
				},
			},
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/etc/gitea"),
					HostPath:      pulumi.String("/opt/gitea/etc"),
				},
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/lib/gitea"),
					HostPath:      pulumi.String("/opt/gitea/data"),
				},
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/etc/timezone"),
					HostPath:      pulumi.String("/etc/timezone"),
					ReadOnly:      pulumi.Bool(true),
				},
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/etc/localtime"),
					HostPath:      pulumi.String("/etc/localtime"),
					ReadOnly:      pulumi.Bool(true),
				},
			},
			Ports: &docker.ContainerPortArray{
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(3000),
					External: pulumi.Int(8102),
				},
			},
			Envs: pulumi.StringArray{
				pulumi.String("GITEA__database__DB_TYPE=postgres"),
				pulumi.String("GITEA__database__HOST=gitea_database:5432"),
				pulumi.String("GITEA__database__NAME=gitea"),
				pulumi.String("GITEA__database__USER=gitea"),
				pulumi.String("GITEA__database__PASSWD=gitea"),
				pulumi.String("GITEA__actions__ENABLED=true"),
			},
			Labels: &docker.ContainerLabelArray{
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.enable"),
					Value: pulumi.String("true"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.gitea.entrypoints"),
					Value: pulumi.String("http"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.gitea.rule"),
					Value: pulumi.String("Host(`gitea.home`)"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.services.gitea.loadbalancer.server.port"),
					Value: pulumi.String("3000"),
				},
			},
		}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{
			databaseContainer,
		}))

		if err != nil {
			return err
		}

		// Export the container's information
		ctx.Export("gitea_database_container_id", databaseContainer.ID())
		ctx.Export("gitea_container_id", giteaContainer.ID())

		return nil
	})
}
