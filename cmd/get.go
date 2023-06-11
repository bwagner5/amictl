/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/bwagner5/amictl/pkg/amis"
)

type GetOptions struct {
	K8sMajorMinorVersion string
	Architecture         string
	AMIVersion           string
	GPUCompataible       bool
}

type GetTableOutput struct {
	Name          string `table:"name"`
	Alias         string `table:"alias"`
	Version       string `table:"version"`
	AMIID         string `table:"ami-id"`
	Arch          string `table:"architecture"`
	GPUCompatible string `table:"gpu compatible,wide"`
	K8sVersion    string `table:"k8s version,wide"`
	OS            string `table:"os,wide"`
	Region        string `table:"region,wide"`
}

var (
	idRE = regexp.MustCompile(`ami-[0-9]+`)
)

var (
	getOpts = GetOptions{}
	cmdGet  = &cobra.Command{
		Use:   "get [ami or alias]",
		Short: "finds information about an ami",
		Long: fmt.Sprintf(`Finds information about an AMI. Valid AMI aliases are: %v `,
			[]string{amis.EKSAL2Alias, amis.EKSBottlerocketAlias, amis.EKSUbuntuAlias, amis.EKSWindowsAlias}),
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getAWSConfgOrDie(cmd.Context())
			ec2Client := ec2.NewFromConfig(cfg)
			ssmClient := ssm.NewFromConfig(cfg)
			eksClient := eks.NewFromConfig(cfg)
			query := amis.Query{
				K8sMajorMinorVersion: getOpts.K8sMajorMinorVersion,
				Architecture:         getOpts.Architecture,
				GPUCompatible:        lo.Ternary(cmd.Flag("gpu-compatible").Changed, lo.ToPtr(getOpts.GPUCompataible), nil),
				AMIVersion:           getOpts.AMIVersion,
			}
			if idRE.MatchString(args[0]) {
				query.ID = args[0]
			} else {
				query.Alias = args[0]
			}
			images, err := amis.New(ec2Client, ssmClient, eksClient).Get(cmd.Context(), query)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			switch globalOpts.Output {
			case OutputYAML:
				fmt.Println(PrettyEncode(images))
			case OutputTableShort, OutputTableWide:
				rows := lo.Map(images, func(image amis.ImageOutput, _ int) GetTableOutput {
					return GetTableOutput{
						Name:          *image.Name,
						Version:       image.Version,
						K8sVersion:    image.K8sVersion,
						Alias:         image.Alias,
						AMIID:         *image.ImageId,
						Arch:          lo.Ternary(string(image.Architecture) == "x86_64", "x86_64 / amd64", string(image.Architecture)),
						OS:            image.OS,
						GPUCompatible: lo.Ternary(image.GPUCompatible, "yes", "no"),
						Region:        cfg.Region,
					}
				})
				sort.SliceStable(rows, func(i, j int) bool {
					return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
				})
				fmt.Println(PrettyTable(rows, globalOpts.Output == OutputTableWide))
			default:
				fmt.Printf("unknown output options %s\n", globalOpts.Output)
				os.Exit(1)
			}
		},
	}
)

func init() {
	cmdGet.Flags().StringVarP(&getOpts.K8sMajorMinorVersion, "k8s-version", "k", "", "K8s Major Minor version (i.e. 1.27)")
	cmdGet.Flags().StringVarP(&getOpts.Architecture, "cpu-arch", "c", "", "CPU Architecture [amd64 or arm64]")
	cmdGet.Flags().BoolVarP(&getOpts.GPUCompataible, "gpu-compatible", "g", false, "GPU Compatible")
	cmdGet.Flags().StringVarP(&getOpts.AMIVersion, "ami-version", "a", "", "AMI Version; if empty use latest (i.e. v20230607 for eks-al2 or 1.6 for Bottlerocket)")
	rootCmd.AddCommand(cmdGet)
}
