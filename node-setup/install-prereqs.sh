#!/bin/bash

# install pre-reqs
sudo apt-get update && sudo apt-get -y upgrade && sudo apt install -y curl 

# Install Kubernetes:

sudo -E apt install -y curl
sudo bash -c "cat <<EOF > /etc/apt/sources.list.d/kubernetes.list
deb http://apt.kubernetes.io/ kubernetes-xenial-unstable main
EOF"
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -

sudo -E apt update
sudo -E apt install -y kubelet kubeadm kubectl

# install containerd
VERSION="1.3.0"
echo "Install Containerd ${VERSION}"
wget -q https://storage.googleapis.com/cri-containerd-release/cri-containerd-${VERSION}.linux-amd64.tar.gz
sudo tar -C / -xzf cri-containerd-${VERSION}.linux-amd64.tar.gz
sudo systemctl start containerd

# setup necessary bits for running k8s:
sudo swapoff -a
sudo modprobe br_netfilter
echo 1 | sudo tee -a /proc/sys/net/ipv4/ip_forward > /dev/null
