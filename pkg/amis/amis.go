package amis

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/samber/lo"
)

const (
	EKSAL2Alias          = "eks-al2"
	EKSUbuntuAlias       = "eks-ubuntu"
	EKSBottlerocketAlias = "eks-bottlerocket"
	EKSWindowsAlias      = "eks-windows"
)
const (
	AMD64Architecture = "amd64"
	ARM64Architecture = "arm64"
)

type Client struct {
	ec2Client *ec2.Client
	eksClient *eks.Client
	ssmClient *ssm.Client
}

type Query struct {
	Alias                string
	ID                   string
	AMIVersion           string
	Architecture         string
	K8sMajorMinorVersion string
	GPUCompatible        bool
}

func New(ec2Client *ec2.Client, ssmClient *ssm.Client) *Client {
	return &Client{
		ec2Client: ec2Client,
		ssmClient: ssmClient,
	}
}

func (c Client) Get(ctx context.Context, query Query) ([]types.Image, error) {
	if query.ID != "" {
		return c.GetByID(ctx, query.ID)
	}
	return c.GetByAlias(ctx, query)
}

func (c Client) GetByAlias(ctx context.Context, query Query) ([]types.Image, error) {
	var resolvedAliases []string
	switch query.Alias {
	case EKSAL2Alias:
		resolvedAliases = append(resolvedAliases, c.eksAL2AMISSMPath(query)...)
	case EKSBottlerocketAlias:
		resolvedAliases = append(resolvedAliases, c.eksBottlerocketAMISSMPath(query)...)
	case EKSUbuntuAlias:
		resolvedAliases = append(resolvedAliases, c.eksUbuntuAMISSMPath(query)...)
	case EKSWindowsAlias:
	default:
	}
	if len(resolvedAliases) == 0 {
		return nil, fmt.Errorf("no AMIs found")
	}
	paramsOut, err := c.ssmClient.GetParameters(ctx, &ssm.GetParametersInput{
		Names: resolvedAliases,
	})
	if err != nil {
		return nil, err
	}
	if len(paramsOut.Parameters) == 0 {
		return nil, fmt.Errorf("no AMIs found")
	}
	return c.GetByID(ctx, lo.Map(paramsOut.Parameters, func(param ssmtypes.Parameter, _ int) string {
		return *param.Value
	})...)
}

func (c Client) GetByID(ctx context.Context, id ...string) ([]types.Image, error) {
	imageOut, err := c.ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds:          id,
		IncludeDeprecated: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	return imageOut.Images, nil
}

func (c Client) eksAL2AMISSMPath(query Query) []string {
	var aliases []string
	version := "recommended"
	if query.AMIVersion != "" {
		version = query.AMIVersion
	}
	if query.GPUCompatible {
		gpuVersion := fmt.Sprintf("amazon-eks-gpu-node-%s-%s", query.K8sMajorMinorVersion, version)
		return []string{fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2-gpu/%s/image_id", query.K8sMajorMinorVersion, gpuVersion)}
	}
	if query.Architecture == "arm64" || query.Architecture == "" {
		arm64Version := fmt.Sprintf("amazon-eks-arm64-node-%s-%s", query.K8sMajorMinorVersion, version)
		aliases = append(aliases, fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2-arm64/%s/image_id", query.K8sMajorMinorVersion, arm64Version))
	}
	if query.Architecture == "amd64" || query.Architecture == "" {
		amd64Version := fmt.Sprintf("amazon-eks-node-%s-%s", query.K8sMajorMinorVersion, version)
		aliases = append(aliases, fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/%s/image_id", query.K8sMajorMinorVersion, amd64Version))
	}
	return aliases
}

func (c Client) eksBottlerocketAMISSMPath(query Query) []string {
	var aliases []string
	version := "latest"
	if query.AMIVersion != "" {
		version = strings.ReplaceAll(query.AMIVersion, "v", "")
	}
	gpuSuffix := ""
	if query.GPUCompatible {
		gpuSuffix = "-nvidia"
	}
	if query.Architecture == "arm64" || query.Architecture == "" {
		aliases = append(aliases, fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s%s/arm64/%s/image_id", query.K8sMajorMinorVersion, gpuSuffix, version))
	}
	if query.Architecture == "amd64" || query.Architecture == "" {
		aliases = append(aliases, fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s%s/x86_64/%s/image_id", query.K8sMajorMinorVersion, gpuSuffix, version))
	}
	return aliases
}

func (c Client) eksUbuntuAMISSMPath(query Query) []string {
	var aliases []string
	version := "current"
	if query.AMIVersion != "" {
		version = strings.ReplaceAll(query.AMIVersion, "v", "")
	}
	if query.GPUCompatible {
		return nil
	}
	if query.Architecture == "arm64" || query.Architecture == "" {
		aliases = append(aliases, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/20.04/%s/stable/%s/arm64/hvm/ebs-gp2/ami-id", query.K8sMajorMinorVersion, version))
	}
	if query.Architecture == "amd64" || query.Architecture == "" {
		aliases = append(aliases, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/20.04/%s/stable/%s/amd64/hvm/ebs-gp2/ami-id", query.K8sMajorMinorVersion, version))
	}
	return aliases
}

func (c Client) EKSSupportedVersions(ctx context.Context) ([]string, error) {
	addonsOut, err := c.eksClient.DescribeAddonVersions(ctx, &eks.DescribeAddonVersionsInput{
		AddonName:  aws.String("vpc-cni"),
		MaxResults: aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(addonsOut.Addons) != 1 {
		return nil, fmt.Errorf("unable to find eks supported versions by inspecting add-on versions")
	}
	clusterVersions := lo.Map(addonsOut.Addons[0].AddonVersions[0].Compatibilities, func(compat ekstypes.Compatibility, _ int) string {
		return *compat.ClusterVersion
	})
	return clusterVersions, nil
}
