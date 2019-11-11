# k8s-pod-overhead

Helper files for setting up cluster with PodOverhead enabled


## Hacking:

helpful watcher for containerd devmapper:
```bash
watch -d ls -al /dev/mapper && sudo du -khs /var/lib/containerd/devmapper/*
```
