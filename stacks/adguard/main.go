package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
	"os"
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

		// Get the AdGuard image
		adGuardImage, err := docker.NewRemoteImage(ctx, "adguard_image", &docker.RemoteImageArgs{
			Name: pulumi.String("adguard/adguardhome:latest"),
		})
		if err != nil {
			return err
		}

		// Create the configuration file
		adGuardConfiguration, err := createConfig()
		if err != nil {
			return err
		}

		// Create the AdGuard container
		adGuardContainer, err := docker.NewContainer(ctx, "adguard_container", &docker.ContainerArgs{
			Image:   adGuardImage.RepoDigest,
			Name:    pulumi.String("adguard"),
			Restart: pulumi.String("always"),
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/opt/adguardhome/conf"),
					HostPath:      pulumi.String("/opt/adguardhome/conf"),
				},
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/opt/adguardhome/work"),
					HostPath:      pulumi.String("/opt/adguardhome/work"),
				},
			},
			Uploads: &docker.ContainerUploadArray{
				&docker.ContainerUploadArgs{
					Content: pulumi.String(adGuardConfiguration),
					File:    pulumi.String("/opt/adguardhome/conf/AdGuardHome.yaml"),
				},
			},
			Ports: &docker.ContainerPortArray{
				&docker.ContainerPortArgs{ // DNS
					Internal: pulumi.Int(53),
					External: pulumi.Int(53),
					Protocol: pulumi.String("tcp"),
				},
				&docker.ContainerPortArgs{ // DNS
					Internal: pulumi.Int(53),
					External: pulumi.Int(53),
					Protocol: pulumi.String("udp"),
				},
				&docker.ContainerPortArgs{ // Dashboard
					Internal: pulumi.Int(80),
					External: pulumi.Int(8100),
					Protocol: pulumi.String("tcp"),
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
					Label: pulumi.String("traefik.http.routers.adguard.entrypoints"),
					Value: pulumi.String("http"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.routers.adguard.rule"),
					Value: pulumi.String("Host(`adguard.home`)"),
				},
				&docker.ContainerLabelArgs{
					Label: pulumi.String("traefik.http.services.adguard.loadbalancer.server.port"),
					Value: pulumi.String("80"),
				},
			},
		}, pulumi.Provider(provider))

		if err != nil {
			return err
		}

		// Export the container's information
		ctx.Export("adguard_container_id", adGuardContainer.ID())

		return nil
	})
}

func mergeMaps(dst, src map[interface{}]interface{}) map[interface{}]interface{} {
	for k, v := range src {
		if _, exists := dst[k]; !exists {
			dst[k] = v
		} else {
			srcMap, srcOk := v.(map[interface{}]interface{})
			dstMap, dstOk := dst[k].(map[interface{}]interface{})
			if srcOk && dstOk {
				dst[k] = mergeMaps(dstMap, srcMap)
			} else {
				dst[k] = v
			}
		}
	}
	return dst
}

func createConfig() ([]byte, error) {
	distData, err := os.ReadFile("etc/AdGuardHome.dist.yaml")
	if err != nil {
		return nil, err
	}

	overrideData, err := os.ReadFile("etc/AdGuardHome.yaml")
	if err != nil {
		return nil, err
	}

	var dist, override map[interface{}]interface{}
	err = yaml.Unmarshal(distData, &dist)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(overrideData, &override)
	if err != nil {
		return nil, err
	}
	merged := mergeMaps(dist, override)
	mergedData, err := yaml.Marshal(merged)
	if err != nil {
		return nil, err
	}
	return mergedData, nil
}
