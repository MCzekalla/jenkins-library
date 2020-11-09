// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/spf13/cobra"
)

type kanikoExecuteOptions struct {
	BuildOptions                []string `json:"buildOptions,omitempty"`
	ContainerBuildOptions       string   `json:"containerBuildOptions,omitempty"`
	ContainerImage              string   `json:"containerImage,omitempty"`
	ContainerImageName          string   `json:"containerImageName,omitempty"`
	ContainerImageTag           string   `json:"containerImageTag,omitempty"`
	ContainerPreparationCommand string   `json:"containerPreparationCommand,omitempty"`
	ContainerRegistryURL        string   `json:"containerRegistryUrl,omitempty"`
	CustomTLSCertificateLinks   []string `json:"customTlsCertificateLinks,omitempty"`
	DockerConfigJSON            string   `json:"dockerConfigJSON,omitempty"`
	DockerfilePath              string   `json:"dockerfilePath,omitempty"`
}

type kanikoExecuteCommonPipelineEnvironment struct {
	container struct {
		registryURL  string
		imageNameTag string
	}
}

func (p *kanikoExecuteCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    interface{}
	}{
		{category: "container", name: "registryUrl", value: p.container.registryURL},
		{category: "container", name: "imageNameTag", value: p.container.imageNameTag},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		log.Entry().Fatal("failed to persist Piper environment")
	}
}

// KanikoExecuteCommand Executes a [Kaniko](https://github.com/GoogleContainerTools/kaniko) build for creating a Docker container.
func KanikoExecuteCommand() *cobra.Command {
	const STEP_NAME = "kanikoExecute"

	metadata := kanikoExecuteMetadata()
	var stepConfig kanikoExecuteOptions
	var startTime time.Time
	var commonPipelineEnvironment kanikoExecuteCommonPipelineEnvironment

	var createKanikoExecuteCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Executes a [Kaniko](https://github.com/GoogleContainerTools/kaniko) build for creating a Docker container.",
		Long:  `Executes a [Kaniko](https://github.com/GoogleContainerTools/kaniko) build for creating a Docker container.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.DockerConfigJSON)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetryData.ErrorCategory = log.GetErrorCategory().String()
				telemetry.Send(&telemetryData)
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			kanikoExecute(stepConfig, &telemetryData, &commonPipelineEnvironment)
			telemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addKanikoExecuteFlags(createKanikoExecuteCmd, &stepConfig)
	return createKanikoExecuteCmd
}

func addKanikoExecuteFlags(cmd *cobra.Command, stepConfig *kanikoExecuteOptions) {
	cmd.Flags().StringSliceVar(&stepConfig.BuildOptions, "buildOptions", []string{`--skip-tls-verify-pull`}, "Defines a list of build options for the [kaniko](https://github.com/GoogleContainerTools/kaniko) build.")
	cmd.Flags().StringVar(&stepConfig.ContainerBuildOptions, "containerBuildOptions", os.Getenv("PIPER_containerBuildOptions"), "Deprected, please use buildOptions. Defines the build options for the [kaniko](https://github.com/GoogleContainerTools/kaniko) build.")
	cmd.Flags().StringVar(&stepConfig.ContainerImage, "containerImage", os.Getenv("PIPER_containerImage"), "Defines the full name of the Docker image to be created including registry, image name and tag like `my.docker.registry/path/myImageName:myTag`. If left empty, image will not be pushed.")
	cmd.Flags().StringVar(&stepConfig.ContainerImageName, "containerImageName", os.Getenv("PIPER_containerImageName"), "Name of the container which will be built - will be used instead of parameter `containerImage`")
	cmd.Flags().StringVar(&stepConfig.ContainerImageTag, "containerImageTag", os.Getenv("PIPER_containerImageTag"), "Tag of the container which will be built - will be used instead of parameter `containerImage`")
	cmd.Flags().StringVar(&stepConfig.ContainerPreparationCommand, "containerPreparationCommand", `rm -f /kaniko/.docker/config.json`, "Defines the command to prepare the Kaniko container. By default the contained credentials are removed in order to allow anonymous access to container registries.")
	cmd.Flags().StringVar(&stepConfig.ContainerRegistryURL, "containerRegistryUrl", os.Getenv("PIPER_containerRegistryUrl"), "http(s) url of the Container registry where the image should be pushed to - will be used instead of parameter `containerImage`")
	cmd.Flags().StringSliceVar(&stepConfig.CustomTLSCertificateLinks, "customTlsCertificateLinks", []string{}, "List containing download links of custom TLS certificates. This is required to ensure trusted connections to registries with custom certificates.")
	cmd.Flags().StringVar(&stepConfig.DockerConfigJSON, "dockerConfigJSON", os.Getenv("PIPER_dockerConfigJSON"), "Path to the file `.docker/config.json` - this is typically provided by your CI/CD system. You can find more details about the Docker credentials in the [Docker documentation](https://docs.docker.com/engine/reference/commandline/login/).")
	cmd.Flags().StringVar(&stepConfig.DockerfilePath, "dockerfilePath", `Dockerfile`, "Defines the location of the Dockerfile relative to the Jenkins workspace.")

}

// retrieve step metadata
func kanikoExecuteMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "kanikoExecute",
			Aliases:     []config.Alias{},
			Description: "Executes a [Kaniko](https://github.com/GoogleContainerTools/kaniko) build for creating a Docker container.",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:        "buildOptions",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "containerBuildOptions",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name:        "containerImage",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "containerImageNameAndTag"}},
					},
					{
						Name:        "containerImageName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "dockerImageName"}},
					},
					{
						Name: "containerImageTag",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "artifactVersion",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "artifactVersion"}},
					},
					{
						Name:        "containerPreparationCommand",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name: "containerRegistryUrl",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "container/registryUrl",
							},
						},
						Scope:     []string{"GENERAL", "PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "dockerRegistryUrl"}},
					},
					{
						Name:        "customTlsCertificateLinks",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
					},
					{
						Name: "dockerConfigJSON",
						ResourceRef: []config.ResourceReference{
							{
								Name: "dockerConfigJsonCredentialsId",
								Type: "secret",
							},

							{
								Name:  "",
								Paths: []string{"$(vaultPath)/docker-config", "$(vaultBasePath)/$(vaultPipelineName)/docker-config", "$(vaultBasePath)/GROUP-SECRETS/docker-config"},
								Type:  "vaultSecretFile",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{},
					},
					{
						Name:        "dockerfilePath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "dockerfile"}},
					},
				},
			},
			Containers: []config.Container{
				{Image: "gcr.io/kaniko-project/executor:debug", Options: []config.Option{{Name: "-u", Value: "0"}, {Name: "--entrypoint", Value: "''"}}},
			},
		},
	}
	return theMetaData
}
