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

		// Get images
		image, err := docker.NewRemoteImage(ctx, "act_image", &docker.RemoteImageArgs{
			Name: pulumi.String("gitea/act_runner:latest"),
		})
		if err != nil {
			return err
		}

		// Create volume
		volume, err := docker.NewVolume(ctx, "gitea_runner_data_volume", &docker.VolumeArgs{
			Name:   pulumi.String("gitea_runner_data"),
			Driver: pulumi.String("local"),
		}, pulumi.Provider(provider))

		// Create the configuration file
		runnerConfiguration, err := createConfig()
		if err != nil {
			return err
		}

		// Create container
		container, err := docker.NewContainer(ctx, "gitea_runner_container", &docker.ContainerArgs{
			Image:   image.RepoDigest,
			Name:    pulumi.String("gitea_runner"),
			Restart: pulumi.String("always"),
			Dns: pulumi.StringArray{
				pulumi.String("192.168.0.200"),
			},
			Volumes: &docker.ContainerVolumeArray{
				&docker.ContainerVolumeArgs{
					VolumeName:    volume.Name,
					ContainerPath: pulumi.String("/data"),
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
				&docker.ContainerVolumeArgs{
					ContainerPath: pulumi.String("/var/run/docker.sock"),
					HostPath:      pulumi.String("/var/run/docker.sock"),
					ReadOnly:      pulumi.Bool(true),
				},
			},
			Uploads: &docker.ContainerUploadArray{
				&docker.ContainerUploadArgs{
					Content: pulumi.String(runnerConfiguration),
					File:    pulumi.String("/config.yaml"),
				},
			},
			Envs: pulumi.StringArray{
				pulumi.String("CONFIG_FILE=/config.yaml"),
				pulumi.String("GITEA_URL=http://gitea.home"),
				pulumi.String("GITEA_INSTANCE_URL=http://192.168.0.200:8102"),
				pulumi.String("GITEA_RUNNER_REGISTRATION_TOKEN=MeaEi8pB2MUSIi9GDDTVbEoOwDCq68XGSTv4gNQh"),
				pulumi.String("GITEA_RUNNER_NAME=act-runner"),
			},
		}, pulumi.Provider(provider))

		if err != nil {
			return err
		}

		// Export the container's information
		ctx.Export("gitea_runner_container_id", container.ID())

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
	distData, err := os.ReadFile("etc/config.dist.yaml")
	if err != nil {
		return nil, err
	}

	overrideData, err := os.ReadFile("etc/config.yaml")
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
