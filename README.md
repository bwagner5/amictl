# amictl

DESCRIPTION HERE

## Usage:


```
Put Usage here
Usage:
  amictl [command]
...
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

EXAMPLES HERE