package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	TimezoneVolumePath  = "/etc/timezone"
	LocaltimeVolumePath = "/etc/localtime"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Config
		cfg := config.New(ctx, "")
		ghostDatabasePasswordSecret := cfg.RequireSecret("ghost_database_password")
		directusKeySecret := cfg.RequireSecret("directus_key")
		directusSecretSecret := cfg.RequireSecret("directus_secret")
		directusAdminPasswordSecret := cfg.RequireSecret("directus_admin_password")

		// Set up provider
		provider, err := docker.NewProvider(ctx, "khorne", &docker.ProviderArgs{
			Host: pulumi.String("tcp://192.168.0.200:2375"),
		})
		if err != nil {
			return err
		}

		// Create networks
		ghostNetwork, err := docker.NewNetwork(ctx, "daniel_tlach_cz_ghost_network", &docker.NetworkArgs{
			Name: pulumi.String("daniel_tlach_cz_ghost_network"),
		}, pulumi.Provider(provider))

		// Get images
		ghostImage, err := docker.NewRemoteImage(ctx, "ghost_image", &docker.RemoteImageArgs{
			Name: pulumi.String("ghost:latest"),
		})
		if err != nil {
			return err
		}

		ghostDatabaseImage, err := docker.NewRemoteImage(ctx, "ghost_database_image", &docker.RemoteImageArgs{
			Name: pulumi.String("mysql:8.0"),
		})
		if err != nil {
			return err
		}

		directusImage, err := docker.NewRemoteImage(ctx, "directus_image", &docker.RemoteImageArgs{
			Name: pulumi.String("directus/directus"),
		})
		if err != nil {
			return err
		}

		// Create volumes
		timezoneVolume := &docker.ContainerVolumeArgs{
			ContainerPath: pulumi.String(TimezoneVolumePath),
			HostPath:      pulumi.String(TimezoneVolumePath),
			ReadOnly:      pulumi.Bool(true),
		}
		localtimeVolume := &docker.ContainerVolumeArgs{
			ContainerPath: pulumi.String(LocaltimeVolumePath),
			HostPath:      pulumi.String(LocaltimeVolumePath),
			ReadOnly:      pulumi.Bool(true),
		}

		ghostDatabaseVolume, err := docker.NewVolume(ctx, "daniel_tlach_cz_ghost_database_volume", &docker.VolumeArgs{
			Name:   pulumi.String("daniel_tlach_cz_ghost_database_volume"),
			Driver: pulumi.String("local"),
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		ghostDataVolume, err := docker.NewVolume(ctx, "daniel_tlach_cz_ghost_volume", &docker.VolumeArgs{
			Name:   pulumi.String("daniel_tlach_cz_ghost_data_volume"),
			Driver: pulumi.String("local"),
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		directusDatabaseVolume, err := docker.NewVolume(ctx, "daniel_tlach_cz_directus_database_volume", &docker.VolumeArgs{
			Name: pulumi.String("daniel_tlach_cz_directus_database_volume"),
		})
		if err != nil {
			return err
		}

		directusUploadsVolume, err := docker.NewVolume(ctx, "daniel_tlach_cz_directus_uploads_volume", &docker.VolumeArgs{
			Name: pulumi.String("daniel_tlach_cz_directus_uploads_volume"),
		})
		if err != nil {
			return err
		}

		// Create containers
		ghostDatabaseContainer, err := docker.NewContainer(ctx, "daniel_tlach_cz_ghost_database_container", &docker.ContainerArgs{
			Image:   ghostDatabaseImage.RepoDigest,
			Name:    pulumi.String("daniel_tlach_cz_ghost_database"),
			Restart: pulumi.String("always"),
			NetworksAdvanced: &docker.ContainerNetworksAdvancedArray{
				&docker.ContainerNetworksAdvancedArgs{
					Name: ghostNetwork.Name,
				},
			},
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/lib/mysql"),
					VolumeName:    ghostDatabaseVolume.Name,
				},
				timezoneVolume,
				localtimeVolume,
			},
			Envs: pulumi.StringArray{
				pulumi.String("MYSQL_ROOT_PASSWORD=ghost"),
				pulumi.String("MYSQL_DATABASE=ghost"),
				pulumi.String("MYSQL_USER=ghost"),
				ghostDatabasePasswordSecret.ApplyT(func(password string) string {
					return "MYSQL_PASSWORD=" + password
				}).(pulumi.StringOutput),
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		ghostContainer, err := docker.NewContainer(ctx, "daniel_tlach_cz_ghost_container", &docker.ContainerArgs{
			Image:   ghostImage.RepoDigest,
			Name:    pulumi.String("daniel_tlach_cz_ghost"),
			Restart: pulumi.String("always"),
			NetworksAdvanced: &docker.ContainerNetworksAdvancedArray{
				&docker.ContainerNetworksAdvancedArgs{
					Name: ghostNetwork.Name,
				},
			},
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/lib/ghost/content"),
					VolumeName:    ghostDataVolume.Name,
				},
				timezoneVolume,
				localtimeVolume,
			},
			Ports: &docker.ContainerPortArray{
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(2368),
					External: pulumi.Int(8103),
				},
			},
			Envs: pulumi.StringArray{
				pulumi.String("database__client=mysql"),
				pulumi.String("database__connection__host=daniel_tlach_cz_ghost_database"),
				pulumi.String("database__connection__user=ghost"),
				ghostDatabasePasswordSecret.ApplyT(func(password string) string {
					return "database__connection__password=" + password
				}).(pulumi.StringOutput),
				pulumi.String("database__connection__database=ghost"),
				pulumi.String("url=http://ghost.daniel_tlach_cz.home"),
				pulumi.String("NODE_ENV=development"),
			},
			Labels: &docker.ContainerLabelArray{
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.enable"),
					Value: pulumi.String("true"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.ghost_daniel_tlach_cz.entrypoints"),
					Value: pulumi.String("http"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.ghost_daniel_tlach_cz.rule"),
					Value: pulumi.String("Host(`ghost.daniel_tlach_cz.home`)"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.services.ghost_daniel_tlach_cz.loadbalancer.server.port"),
					Value: pulumi.String("2368"),
				},
			},
		}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{
			ghostDatabaseContainer,
		}))
		if err != nil {
			return err
		}

		directusContainer, err := docker.NewContainer(ctx, "daniel_tlach_cz_directus_container", &docker.ContainerArgs{
			Image:   directusImage.RepoDigest,
			Name:    pulumi.String("daniel_tlach_cz_directus"),
			Restart: pulumi.String("always"),
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/directus/database"),
					VolumeName:    directusDatabaseVolume.Name,
				},
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/directus/uploads"),
					VolumeName:    directusUploadsVolume.Name,
				},
				timezoneVolume,
				localtimeVolume,
			},
			Ports: &docker.ContainerPortArray{
				&docker.ContainerPortArgs{
					Internal: pulumi.Int(8055),
					External: pulumi.Int(8104),
				},
			},
			Envs: pulumi.StringArray{
				directusKeySecret.ApplyT(func(secret string) string {
					return "KEY=" + secret
				}).(pulumi.StringOutput),
				directusSecretSecret.ApplyT(func(secret string) string {
					return "SECRET=" + secret
				}).(pulumi.StringOutput),
				pulumi.String("ADMIN_EMAIL=daniel@tlach.cz"),
				directusAdminPasswordSecret.ApplyT(func(secret string) string {
					return "ADMIN_PASSWORD=" + secret
				}).(pulumi.StringOutput),
				pulumi.String("DB_CLIENT=sqlite3"),
				pulumi.String("DB_FILENAME=/directus/database/database.sqlite"),
				pulumi.String("WEBSOCKETS_ENABLED=true"),
			},
			Labels: &docker.ContainerLabelArray{
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.enable"),
					Value: pulumi.String("true"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.directus_daniel_tlach_cz.entrypoints"),
					Value: pulumi.String("http"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.directus_daniel_tlach_cz.rule"),
					Value: pulumi.String("Host(`directus.daniel_tlach_cz.home`)"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.services.directus_daniel_tlach_cz.loadbalancer.server.port"),
					Value: pulumi.String("8055"),
				},
			},
		}, pulumi.Provider(provider))

		ctx.Export("ghost_daniel_tlach_cz_database_container_id", ghostDatabaseContainer.ID())
		ctx.Export("ghost_daniel_tlach_cz_container_id", ghostContainer.ID())
		ctx.Export("directus_daniel_tlach_cz_container_id", directusContainer.ID())

		return nil
	})
}
