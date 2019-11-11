#!/bin/bash


function aggressive_cleanup() {

  sudo rm -rf /var/lib/containerd/devmapper/data-disk.img
  sudo rm -rf /var/lib/containerd/devmapper/meta-disk.img

  sudo systemctl stop containerd-devmapper
  sudo systemctl disable containerd-devmapper

  sudo losetup -d /dev/loop20
  sudo losetup -d /dev/loop21

  sudo dmsetup remove_all

  echo "you'll want to update /etc/containerd/config.toml to remove devmapper usage"
}

echo "this will remove all thin-pools on this node, (dmsetup remove_all)."
read -r -p "Are you sure? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY])
        aggressive_cleanup
        ;;
    *)
        return 0
        ;;
esac
