# setup necessary bits for running k8s:
sudo swapoff -a
sudo modprobe br_netfilter

echo 1 | sudo tee -a /proc/sys/net/ipv4/ip_forward > /dev/null


sudo systemctl daemon-reload && sudo systemctl restart kubelet && sudo systemctl restart containerd

# initialize cluster
sudo -E kubeadm init --config=./kubeadm.yaml

# setup kubeconfig
rm -rf "${HOME}/.kube"
mkdir -p "${HOME}/.kube"
sudo cp -i /etc/kubernetes/admin.conf "${HOME}/.kube/config"
sudo chown "$(id -u):$(id -g)" "${HOME}/.kube/config"

# taint master node:
kubectl taint nodes --all node-role.kubernetes.io/master-

# setup canal
kubectl apply -f https://docs.projectcalico.org/v3.10/manifests/canal.yaml
