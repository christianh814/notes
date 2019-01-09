# Kubernetes Installation

This guide will show you how to "easily" (as easy as possible) install a Multimaster system. You can also modify this to have a similer one contoller/one node setup (steps are almost exactly the same).

Also, this is here mainly for my notes; so I may be brief in some sections. I TRY to keep this up to date but you're better off looking at the [official documentation](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

This is broken off into sections

* [Prerequisites And Assumptions](#prerequisites-and-assumptions)
* [Infrastructure](#infrastructure)

# Prerequisites And Assumptions

Out of all sections, this is probably the most important one. Here I will list prereqs and assumptions. This list isn't complete and I assume you know your way around Linux and are familiar with containers and Kubernetes already.

Prereqs:

* You have control of DNS
  * Also that DNS is in place for the hosts (forward and reverse)
  * You can use [nip.io](http://nip.io) if you wish
* You are using Tmux and/or Ansible
  * A lot of these tasks are done on all machines at the same time
  * I used mostly Ansible with some Tmux here and there (where it made sense)
* All hosts are on one network segment
  * All my VMs were on the same vlan
  * I didn't have to create firewall rules on a router/firewall


Assumptions:

* You are well vested in the art of Linux
* You're running this on virtual machines
* You either have preassinged DHCP or static ips
* DNS is in place (forward and reverse)

# Infrastructure

I have spun up 7 VMs

* 1 LB VM
* 3 VMs for the controllers (will also run etcd)
* 3 VMs for the workers


The contollers/workers are set up as follows

* Latest CentOS 7
* 50GB HD for the OS
* 4 GB of RAM
* 4 vCPUs
* Swap has been disabled
* SELinux has been disabled (part of the howto)
* Firewall has been disabled (part of the howto)

## HAProxy Setup

Although not 100% within the scope of this document (really any LB will do); this is how I set up my LB on my LB server. These are "high level" notes.

Install HAProxy package

```
yum -y install haproxy
```

Then I enable `6443` (kube api) and `9000` (haproxy stats) on the firewall

```
firewall-cmd --add-port=6443/tcp --add-port=6443/udp --add-port=9000/tcp --add-port=9000/udp --permanent
firewall-cmd --reload
```

If using selinux, make sure "connect any" is on

```
setsebool -P haproxy_connect_any=on
```

Edit the `/etc/haproxy/haproxy.cfg` file and add the following

```
frontend  kube-api
    bind *:6443
    default_backend k8s-api
    mode tcp
    option tcplog

backend k8s-api
    balance source
    mode tcp
    server      controller-0 192.168.1.98:6443 check
    server      controller-1 192.168.1.99:6443 check
    server      controller-2 192.168.1.5:6443 check
```

If you want a copy of the whole file, you can find an example [here](k8s-resources/k8s-haproxy.cfg)

Now you can start/enable the haproxy service

```
systemctl enable --now haproxy
```

I like to keep `<ip of lb>:9000` up so I can monitor when my control plane servers has come up
