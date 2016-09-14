# Nested KVM

In KVM - you can "pass along" CPU capabilities to the guest. This poses interesting scenarios where you can have a VM within a VM running a VM (sort of a "dream within a dream" type of thing)

This allows to create a pretty cool oVirt virtual lab :-)

__Prep Work__
First; check if your machine has Hardware Virt

	
	root@host# egrep '(vmx|svm)' --color=always /proc/cpuinfo


If nothing is displayed, then your processor doesn't support hardware virtualization.

Now, check to see if nested KVM support is enabled on your machine by running:

	
	root@host# cat /sys/module/kvm_intel/parameters/nested


If the answer is “N,” you can enable it by running:

	
	root@host# echo "options kvm-intel nested=1" > /etc/modprobe.d/kvm-intel.conf



After adding that // kvm-intel.conf // file, reboot your machine, after which // cat /sys/module/kvm_intel/parameters/nested // should return // Y //

You can add virtulazation on the UI or on the command line like so

	
	[root@dlp ~]# virsh edit www
	# edit a virtual machine "www"
	# add following lines
	
	`<cpu mode='custom' match='exact'>`
	 
	# CPU model
	
	  `<model fallback='allow'>`SandyBridge`</model>`
	 
	# CPU vendor
	
	  `<vendor>`Intel`</vendor>`
	  `<feature policy='require' name='vmx'/>`
	`</cpu>`


Now I am going to use 4 nodes for the following


*  oVirt - The oVirt engine server (i.e. RHEV-M)

*  Hypervisors - 2 nodes (running VDSM) that will "house" the VMs

*  IPA/NFS - And IPA server to act as a DNS server, Authentication Server, and Storage System

All nodes will have 2 vNICs - one for the oVirt/Storage Network and one for "serving" (note in "real" life...you'll need at least 3 networks)

Here is a drawing of what it will look like.

{{:ovirt_-_lab.png|}}

For this set up...I will be following this guide [http://wiki.centos.org/HowTos/oVirt](http://wiki.centos.org/HowTos/oVirt)

__IPA / Storage node Setup__

oVirt works best if DNS is set up properly. ALSO, you can "hook up" oVirt to an authentication system - so I figured IPA would be the best choice since it offers both.

Here is an overview of my domain.

Domain: example.org - 172.16.1.0/24 \\
IPA: ipa.example.org - 172.16.1.149 \\
NFS: nfs.example.org (CNAME to ipa.example.org) \\
Ovirt: ovirt.example.org - 172.16.1.150 \\
VR1: vr1.example.org - 172.16.1.151 \\
VR2: vr2.example.org - 172.16.1.152 \\

