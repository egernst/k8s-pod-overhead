# k8s-pod-overhead

Helper files for setting up cluster with PodOverhead enabled. This will create a single-node cluster using K8S 1.16.x + PodOverhead feature gate set with containerd 1.3.0, Kata 1.9.1 with RuntimeClasses registered for QEMU-virtiofs VMM and Firecracker VMM (runc is the default). 

## setup

This assumes you are on an Ubuntu VM which supports VMX. I validated the below flow on an Azure standard d8s_v3 instance. (d*s_v3 supports virtualization)

### install pre-requisites, setup node:


Install curl, Kubernetes, containerd 1.3.0, disable swap, enable br_netfilter and ip_forwarding for ipv4:
```bash 
./node-setup/install-prereqs.sh
```

With containerd on the system, update the configuration to make use of devmapper snapshotter. This script will create data and meta-data disks for the containerd devmapper snapshotter to consume, and back them with a thinpool of size 20GB. ```/etc/containerd/config.toml``` will be updated to use the devmapper snapshotter and the newly created `contd-thin-pool`:

```bash
./node-setup/containerd_devmapper_setup.sh
```

If you ever need to clean-up the devmapper (or if things go sideways), take a look @ ./node-setup/cleanup-devmapper.sh

### Start a cluster

We use a kubeadm configuration yaml for assisting in the init process.  To get started:
```bash
./create_stack.sh
```

*note, this is a simplified version of [this ClearLinux example stack](https://github.com/clearlinux/cloud-native-setup/tree/master/clr-k8s-examples)

### Bootstrap all the things:

Setup Kata:
```
kubectl apply -f kata-setup/kata-rbac.yaml
kubectl apply -f kata-setup/kata-deploy.yaml
kubectl apply -f kata-setup/kata-runtimeclasses.yaml
```

After kata-deploy is up and running, you'll see kata artifacts are installed at /opt/kata, and that 'etc/containerd/config.toml' was updated to register the kata handlers.

### Example pods:

Launch a few pods, utilizing runc, kata-fc and kata-qemu:

```
kubectl apply -f test-workloads/kata-fc-pod.yaml
kubectl apply -f test-workloads/kata-qemu-pod.yaml
kubectl apply -f test-workloads/runc-pod.yaml
```



## Hacking:

helpful watcher for containerd devmapper:
```bash
watch -d ls -al /dev/mapper && sudo du -khs /var/lib/containerd/devmapper/*
```
