# amictl

AMIs can be tricky to work with. AMI IDs are difficult to find and remember. SSM Aliases to AMI IDs are long paths that aren't very intuitive to remember or construct. On top of that, working with AMIs for K8s adds another dimension of complexity.

`amictl` aims to make looking up AMI information easy by providing some easy aliases and standardized query filters to lookup AMIs for EKS Optimized AL2, Bottlerocket, Ubuntu, and Windows.

## Usage:

```
Usage:
  amictl [command]

Available Commands:
  get         finds information about an ami
  help        Help about any command

Flags:
  -f, --file string   YAML Config File
  -h, --help          help for amictl
      --verbose       Verbose output
      --version       version

Use "amictl [command] --help" for more information about a command.
```

```
Finds information about an AMI. Valid AMI aliases are: [eks-al2 eks-bottlerocket eks-ubuntu eks-windows]

Usage:
  amictl get [ami or alias] [flags]

Flags:
  -a, --ami-version string   AMI Version; if empty use latest (i.e. v20230607 for eks-al2 or 1.6 for Bottlerocket)
  -c, --cpu-arch string      CPU Architecture [amd64 or arm64]
  -g, --gpu-compatible       GPU Compatible
  -h, --help                 help for get
  -k, --k8s-version string   K8s Major Minor version (i.e. 1.27)

Global Flags:
  -f, --file string   YAML Config File
      --verbose       Verbose output
      --version       version
```

## Installation:

```
brew install bwagner5/wagner/amictl
```

Packages, binaries, and archives are published for all major platforms (Mac amd64/arm64 & Linux amd64/arm64):

Debian / Ubuntu:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
wget https://github.com/bwagner5/amictl/releases/download/v0.0.1/amictl_0.0.1_${OS}_${ARCH}.deb
dpkg --install amictl_0.0.2_linux_amd64.deb
amictl --help
```

RedHat:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
rpm -i https://github.com/bwagner5/amictl/releases/download/v0.0.1/amictl_0.0.1_${OS}_${ARCH}.rpm
```

Download Binary Directly:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
wget -qO- https://github.com/bwagner5/amictl/releases/download/v0.0.1/amictl_0.0.1_${OS}_${ARCH}.tar.gz | tar xvz
chmod +x amictl
```

## Examples: 


```
## Get the latest released EKS Optimized AL2 AMIs (arm64 and amd64) w/ the latest version of K8s
> amictl get eks-al2
[
    {
        "Architecture": "arm64",
        "BlockDeviceMappings": [
            {
                "DeviceName": "/dev/xvda",
                "Ebs": {
                    "DeleteOnTermination": true,
                    "Encrypted": false,
                    "Iops": null,
                    "KmsKeyId": null,
                    "OutpostArn": null,
                    "SnapshotId": "snap-08574db993db823a7",
                    "Throughput": null,
                    "VolumeSize": 20,
                    "VolumeType": "gp2"
                },
                "NoDevice": null,
                "VirtualName": null
            }
        ],
        "BootMode": "uefi",
```

```
## Get a specific version of the EKS Optimized AL2 AMIs w/ a specific CPU arch and K8s version
> amictl get eks-al2 --ami-version v20230607 --cpu-arch arm64 --k8s-version 1.27
[
    {
        "Architecture": "arm64",
        "BlockDeviceMappings": [
            {
                "DeviceName": "/dev/xvda",
                "Ebs": {
                    "DeleteOnTermination": true,
                    "Encrypted": false,
                    "Iops": null,
                    "KmsKeyId": null,
                    "OutpostArn": null,
                    "SnapshotId": "snap-08574db993db823a7",
                    "Throughput": null,
                    "VolumeSize": 20,
                    "VolumeType": "gp2"
                },
                "NoDevice": null,
                "VirtualName": null
            }
        ],
        "BootMode": "uefi",
```

