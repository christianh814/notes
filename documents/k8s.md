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
