# amictl

AMIs can be tricky to work with. AMI IDs are difficult to find and remember. SSM Aliases to AMI IDs are long paths that aren't very intuitive to remember or construct. On top of that, working with AMIs for K8s adds another dimension of complexity.

`amictl` aims to make looking up AMI information easy by providing some easy aliases and standardized query filters to lookup AMIs for EKS Optimized AL2, Bottlerocket, Ubuntu, and Windows.

## Usage:

```
> amictl --help
Usage:
  amictl [command]

Available Commands:
  get         finds information about an ami
  help        Help about any command

Flags:
  -f, --file string     YAML Config File
  -h, --help            help for amictl
  -o, --output string   Output mode: [short wide yaml] (default "short")
      --verbose         Verbose output
      --version         version

Use "amictl [command] --help" for more information about a command.
```

```
> amictl get --help
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
  -f, --file string     YAML Config File
  -o, --output string   Output mode: [short wide yaml] (default "short")
      --verbose         Verbose output
      --version         version
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
wget https://github.com/bwagner5/amictl/releases/download/v0.0.2/amictl_0.0.2_${OS}_${ARCH}.deb
dpkg --install amictl_0.0.2_linux_amd64.deb
amictl --help
```

RedHat:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
rpm -i https://github.com/bwagner5/amictl/releases/download/v0.0.2/amictl_0.0.2_${OS}_${ARCH}.rpm
```

Download Binary Directly:

```
[[ `uname -m` == "aarch64" ]] && ARCH="arm64" || ARCH="amd64"
OS=`uname | tr '[:upper:]' '[:lower:]'`
wget -qO- https://github.com/bwagner5/amictl/releases/download/v0.0.2/amictl_0.0.2_${OS}_${ARCH}.tar.gz | tar xvz
chmod +x amictl
```

## Examples: 

```
## Get all the eks-optimized K8s AMIs for the latest version of k8s
> amictl get eks
NAME                                                          	ALIAS           	VERSION  	AMI-ID               	ARCHITECTURE
amazon-eks-arm64-node-1.27-v20230607                          	eks-al2         	v20230607	ami-012fb2a3ce1880d5d	arm64
amazon-eks-gpu-node-1.27-v20230607                            	eks-al2         	v20230607	ami-0bb1342f3bfff2b36	x86_64 / amd64
amazon-eks-node-1.27-v20230607                                	eks-al2         	v20230607	ami-07c9c86f18d0ff01e	x86_64 / amd64
bottlerocket-aws-k8s-1.27-aarch64-v1.14.1-842c7134            	eks-bottlerocket	v1.14.1  	ami-0e0ecf30ef5b80554	arm64
bottlerocket-aws-k8s-1.27-nvidia-aarch64-v1.14.1-842c7134     	eks-bottlerocket	v1.14.1  	ami-0eba4734d387a432c	arm64
bottlerocket-aws-k8s-1.27-nvidia-x86_64-v1.14.1-842c7134      	eks-bottlerocket	v1.14.1  	ami-02b8bf3eed09b7e7a	x86_64 / amd64
bottlerocket-aws-k8s-1.27-x86_64-v1.14.1-842c7134             	eks-bottlerocket	v1.14.1  	ami-074055b5c6f0923d4	x86_64 / amd64
Windows_Server-2022-English-Core-EKS_Optimized-1.27-2023.06.06	eks-windows     	2022     	ami-0bb26c27a567a73a5	x86_64 / amd64
```

```
## Get a specific version of the EKS Optimized AL2 AMIs w/ a specific CPU arch and K8s version
> amictl get eks-al2 --ami-version v20230607 --cpu-arch arm64 --k8s-version 1.27
NAME                                	ALIAS  	VERSION  	AMI-ID               	ARCHITECTURE
amazon-eks-arm64-node-1.27-v20230607	eks-al2	v20230607	ami-012fb2a3ce1880d5d	arm64
amazon-eks-gpu-node-1.27-v20230607  	eks-al2	v20230607	ami-0bb1342f3bfff2b36	x86_64 / amd64
```

