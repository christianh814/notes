# Installation

The installation of OpenShift Enterprise (OSE); will be done via scripts. More information can be found using the Red Hat [documentation site](https///access.redhat.com/beta/documentation/en/openshift-enterprise-30-administrator-guide/chapter-1-installation). 

## Infrastructure

For this installation we have the following

*  Wildcard DNS entry - *.cloudapps.example.com 172.16.1.247

*  Master
    * ose3-master.example.com
    * 172.16.1.247
    * Also acting as a node

*  Node1
    * ose3-node1.example.com
    * 172.16.1.246

*  Node2
    * ose3-node2.example.com
    * 172.16.1.245

Servers installed with RHEL 7.1 (Greater than 7.1 is required for openvswitch) at a "minimum" installation.

Forward/Reverse DNS is a MUST for master/nodes

Map of how OSEv3 works:

{{::osev3.jpg|}} 

## Host preparation

Each host must be registered using RHSM and have an active OSE subscription attached to access the required packages.

On each host, register with RHSM:

	
	subscription-manager register --username=${user_name} --password=${password}


List the available subscriptions:

	
	subscription-manager list --available


In the output for the previous command, find the pool ID for an OpenShift Enterprise subscription and attach it:

	
	subscription-manager attach --pool=${pool_id}


**__NOTE:__** You can have RHSM do this for you in one shot

	
	subscription-manager register --username=${user_name} --password=${password} --auto-attach


Disable all repositories and enable only the required ones:

	
	subscription-manager repos --disable="*"
	subscription-manager repos \
	    --enable="rhel-7-server-rpms" \
	    --enable="rhel-7-server-extras-rpms" \
	    --enable="rhel-7-server-ose-3.2-rpms"


Make sure the pre-req pkgs are installed/removed and make sure the system is updated

	
	yum -y install wget git net-tools bind-utils iptables-services bridge-utils bash-completion vim yum-versionlock
	yum -y update
	yum -y install atomic-openshift-utils docker-1.9


## Docker Configuration

Configure docker by editing the // /etc/sysconfig/docker // file and add // --insecure-registry 0.0.0.0/0 // to the OPTIONS parameter. For example: 

	
	OPTIONS=--selinux-enabled --insecure-registry 0.0.0.0/0


You can do this with a sed statement

	
	root@master# sed -i.bak "s/^OPTIONS='--selinux-enabled'/OPTIONS='--selinux-enabled --insecure-registry 0.0.0.0\/0'/g" /etc/sysconfig/docker


Next configure docker storage.

Docker’s default loopback storage mechanism is not supported for production use and is only appropriate for proof of concept environments. For production environments, you must create a thin-pool logical volume and re-configure docker to use that volume.

You can use the docker-storage-setup script to create a thin-pool device and configure docker’s storage driver after installing docker but before you start using it. The script reads configuration options from the */etc/sysconfig/docker-storage-setup* file.

Configure docker-storage-setup for your environment. There are three options available based on your storage configuration:

a) Create a thin-pool volume from the remaining free space in the volume group where your root filesystem resides; this requires no configuration:
`# docker-storage-setup`

b) Use an existing volume group, in this example docker-vg, to create a thin-pool:

	
	# echo `<<EOF >` /etc/sysconfig/docker-storage-setup
	VG=docker-vg
	SETUP_LVM_THIN_POOL=yes
	EOF
	# docker-storage-setup


c) Use an unpartitioned block device to create a new volume group and thinpool. In this example, the /dev/vdc device is used to create the docker-vg volume group:

	
	# cat `<<EOF >` /etc/sysconfig/docker-storage-setup
	DEVS=/dev/vdc
	VG=docker-vg
	EOF
	# docker-storage-setup


Verify your configuration. You should have dm.thinpooldev value in the /etc/sysconfig/docker-storage file and a docker-pool device:

	
	# lvs
	LV                  VG        Attr       LSize  Pool Origin Data%  Meta% Move Log Cpy%Sync Convert
	docker-pool         docker-vg twi-a-tz-- 48.95g             0.00   0.44
	# cat /etc/sysconfig/docker-storage
	DOCKER_STORAGE_OPTIONS=--storage-opt dm.fs=xfs --storage-opt dm.thinpooldev=/dev/mapper/docker--vg-docker--pool

Re-initialize docker.

**Warning** This will destroy any docker containers or images currently on the host.

	
	    # systemctl stop docker
	    # rm -rf /var/lib/docker/*
	    # systemctl restart docker



## Ansible Installer

On the master host, generate ssh keys to use for ansible press enter to accept the defaults

	
	root@master# ssh-keygen


Distribue these keys to all hosts (including the master)

	
	root@master# for host in ose3-master.example.com \
	    ose3-node1.example.com \
	    ose3-node2.example.com; \
	    do ssh-copy-id -i ~/.ssh/id_rsa.pub $host; \
	    done


Test passwordless ssh

	
	root@master# for host in ose3-master.example.com \
	    ose3-node1.example.com \
	    ose3-node2.example.com; \
	    do ssh $host hostname; \
	    done


Make a backup of the // /etc/ansible/hosts // file

	
	cp /etc/ansible/hosts{,.bak}


Next You must create an */etc/ansible/hosts* file for the playbook to use during the installation. The following is an example of a Bring Your Own (BYO) host inventory, based on the example host configuration. You can see these example hosts present in both the //[masters]// and //[nodes]// sections:

	
	# Create an OSEv3 group that contains the masters and nodes groups
	[OSEv3:children]
	masters
	nodes
	
	# Set variables common for all OSEv3 hosts
	[OSEv3:vars]
	# SSH user, this user should allow ssh based auth without requiring a password
	ansible_ssh_user=root
	
	# If ansible_ssh_user is not root, ansible_sudo must be set to true
	#ansible_sudo=true
	
	#product_type=openshift
	
	deployment_type=openshift-enterprise
	osm_default_subdomain=cloudapps.example.com
	
	osm_cluster_network_cidr=10.1.0.0/16
	osm_host_subnet_length=8
	openshift_master_portal_net=172.30.0.0/16
	osm_default_node_selector="region=primary"
	openshift_router_selector='region=infra'
	openshift_registry_selector='region=infra'
	
	# Configure dnsmasq for cluster dns, switch the host's local resolver to use dnsmasq
	# and configure node's dnsIP to point at the node's local dnsmasq instance. Defaults
	# to True for Origin 1.2 and OSE 3.2. False for 1.1 / 3.1 installs, this cannot
	# be used with 1.0 and 3.0.
	# openshift_use_dnsmasq=False
	
	
	# uncomment the following to enable htpasswd authentication; defaults to DenyAllPasswordIdentityProvider
	openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/openshift/openshift-passwd'}]
	
	# host group for masters
	[masters]
	ose3-master.example.com openshift_public_hostname=ose3-master.example.com openshift_ip=172.16.1.247 openshift_public_ip=172.16.1.247 openshift_hostname=ose3-master.example.com
	
	# host group for nodes, includes region info
	[nodes]
	ose3-master.example.com openshift_public_hostname=ose3-master.example.com openshift_ip=172.16.1.247 openshift_public_ip=172.16.1.247 openshift_hostname=ose3-master.example.com openshift_node_labels="{'region': 'infra', 'zone': 'default'}"
	ose3-node1.example.com openshift_public_hostname=ose3-node1.example.com openshift_ip=172.16.1.246 openshift_public_ip=172.16.1.246 openshift_hostname=ose3-node1.example.com openshift_node_labels="{'region': 'primary', 'zone': 'east'}"
	ose3-node2.example.com openshift_public_hostname=ose3-node2.example.com openshift_ip=172.16.1.245 openshift_public_ip=172.16.1.245 openshift_hostname=ose3-node2.example.com openshift_node_labels="{'region': 'primary', 'zone': 'west'}"
	##
	##


Make sure that you have *deployment_type=openshift-enterprise* after *ansible_ssh_user=root* if it's not already there.

Also change the subdomain to yours

	
	osm_default_subdomain=cloudapps.example.com



You can run the playbook (specifying a *-i* if you wrote the hosts file somewhere else)

	
	root@master# ansible-playbook /usr/share/ansible/openshift-ansible/playbooks/byo/config.yml


Once this completes successfully, run *oc get nodes* and you should see "Ready"

	
	root@master# oc get nodes
	NAME                      LABELS                                           STATUS
	ose3-master.example.com   kubernetes.io/hostname=ose3-master.example.com   Ready
	ose3-node1.example.com    kubernetes.io/hostname=ose3-node1.example.com    Ready
	ose3-node2.example.com    kubernetes.io/hostname=ose3-node2.example.com    Ready


Make sure you label your nodes if they aren't properly labeled by the ansible installer

	
	root@master# oc label node ose3-master.example.com region=infra zone=default
	root@master# oc label node ose3-node1.example.com region=primary zone=east
	root@master# oc label node ose3-node2.example.com region=primary zone=west


### Multi Master

Using the native MultiMaster config (with an external LB) looks like this

	
	# Create an OSEv3 group that contains the masters and nodes groups
	[OSEv3:children]
	masters
	nodes
	etcd
	
	# Set variables common for all OSEv3 hosts
	[OSEv3:vars]
	# SSH user, this user should allow ssh based auth without requiring a password
	ansible_ssh_user=chernand
	
	# If ansible_ssh_user is not root, ansible_sudo must be set to true
	ansible_sudo=true
	
	#product_type=openshift
	
	deployment_type=openshift-enterprise
	osm_default_subdomain=apps.mmchx.osecloud.com
	
	osm_cluster_network_cidr=10.1.0.0/16
	osm_host_subnet_length=8
	openshift_master_portal_net=172.30.0.0/16
	
	# uncomment the following to enable htpasswd authentication; defaults to DenyAllPasswordIdentityProvider
	openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/openshift/openshift-passwd'}]
	
	# Native high availbility cluster method with optional load balancer.
	# If no lb group is defined installer assumes that a load balancer has 
	# been preconfigured. For installation the value of
	# openshift_master_cluster_hostname must resolve to the load balancer
	# or to one or all of the masters defined in the inventory if no load
	# balancer is present.
	openshift_master_cluster_method=native
	openshift_master_cluster_hostname=mmchx.osecloud.com
	openshift_master_cluster_public_hostname=mmchx.osecloud.com
	openshift_master_public_api_url=https://mmchx.osecloud.com
	openshift_master_public_console_url=https://mmchx.osecloud.com/console
	
	# These ports reference the ports on the masters themselves 
	openshift_master_api_port=443
	openshift_master_console_port=443
	
	# override the default controller lease ttl 
	#osm_controller_lease_ttl=30
	
	# Multitenant instead of "flat" network
	os_sdn_network_plugin_name=redhat/openshift-ovs-multitenant
	
	# host group for masters
	[masters]
	ose-master01.c.jump-servers.internal openshift_public_hostname=53.137.155.104.bc.googleusercontent.com openshift_ip=10.128.0.3 openshift_public_ip=104.155.137.53 openshift_hostname=ose-master01.c.jump-servers.internal
	ose-master02.c.jump-servers.internal openshift_public_hostname=234.131.155.104.bc.googleusercontent.com openshift_ip=10.128.0.4 openshift_public_ip=104.155.131.234 openshift_hostname=ose-master02.c.jump-servers.internal
	ose-master03.c.jump-servers.internal openshift_public_hostname=5.233.223.199.bc.googleusercontent.com openshift_ip=10.128.0.5 openshift_public_ip=199.223.233.5 openshift_hostname=ose-master03.c.jump-servers.internal
	
	# host group for etcd
	[etcd]
	ose-master01.c.jump-servers.internal
	ose-master02.c.jump-servers.internal
	ose-master03.c.jump-servers.internal
	
	
	# host group for nodes, includes region info
	[nodes]
	ose-master01.c.jump-servers.internal openshift_public_hostname=53.137.155.104.bc.googleusercontent.com openshift_ip=10.128.0.3 openshift_public_ip=104.155.137.53 openshift_hostname=ose-master01.c.jump-servers.internal openshift_node_labels="{'region': 'primary', 'zone': 'master'}"
	ose-master02.c.jump-servers.internal openshift_public_hostname=234.131.155.104.bc.googleusercontent.com openshift_ip=10.128.0.4 openshift_public_ip=104.155.131.234 openshift_hostname=ose-master02.c.jump-servers.internal openshift_node_labels="{'region': 'primary', 'zone': 'master'}"
	ose-master03.c.jump-servers.internal openshift_public_hostname=5.233.223.199.bc.googleusercontent.com openshift_ip=10.128.0.5 openshift_public_ip=199.223.233.5 openshift_hostname=ose-master03.c.jump-servers.internal openshift_node_labels="{'region': 'primary', 'zone': 'master'}"
	ose-infra01.c.jump-servers.internal openshift_public_hostname=59.139.197.104.bc.googleusercontent.com openshift_ip=10.128.0.6 openshift_public_ip=104.197.139.59 openshift_hostname=ose-infra01.c.jump-servers.internal openshift_node_labels="{'region': 'primary', 'zone': 'infra'}"
	ose-infra02.c.jump-servers.internal openshift_public_hostname=157.157.155.104.bc.googleusercontent.com openshift_ip=10.128.0.7 openshift_public_ip=104.155.157.157 openshift_hostname=ose-infra02.c.jump-servers.internal openshift_node_labels="{'region': 'primary', 'zone': 'infra'}"
	ose-node01.c.jump-servers.internal openshift_public_hostname=46.180.197.104.bc.googleusercontent.com openshift_ip=10.128.0.8 openshift_public_ip=104.197.180.46 openshift_hostname=ose-node01.c.jump-servers.internal openshift_node_labels="{'region': 'primary', 'zone': 'node'}"
	##
	##


### HAProxy Config

If using HAProxy for your LB one possible config might look like this (YMMV)

	
	#---------------------------------------------------------------------
	# source balancing for OSE master
	#---------------------------------------------------------------------
	frontend master-ssl
	    bind 172.16.1.120:8443
	    mode tcp
	    default_backend master-backendssl
	
	frontend router-http
	    bind 172.16.1.120:80
	    default_backend router-backend-http
	
	frontend router-https
	    bind 172.16.1.120:443
	    mode tcp
	    default_backend router-backend-https
	
	backend master-backendssl
	    mode tcp
	    balance source
	    option httpchk get /healthz
	    http-check expect status 200
	    server master1-ssl 172.16.1.121:8443 check ssl verify none
	    server master2-ssl 172.16.1.122:8443 check ssl verify none
	    server master3-ssl 172.16.1.123:8443 check ssl verify none
	
	backend router-backend-http
	    balance roundrobin
	    mode http
	    option httpclose
	    option forwardfor
	    option httpchk get /healthz
	    http-check expect status 200
	    server router1 172.16.1.121:80 check port 1936
	    server router2 172.16.1.122:80 check port 1936
	
	
	backend router-backend-http
	    balance roundrobin
	    mode http
	    option httpclose
	    option forwardfor
	    option httpchk get /healthz
	    http-check expect status 200
	    server router1 172.16.1.121:80 check port 1936
	    server router2 172.16.1.122:80 check port 1936
	
	backend router-backend-https
	    mode tcp
	    balance source
	    option httpchk get /healthz
	    http-check expect status 200
	    server router1 172.16.1.121:443 check port 1936 ssl verify none
	    server router2 172.16.1.122:443 check port 1936 ssl verify none
	#---------------------------------------------------------------------
	## 



# Post Installation Steps

After you get OSE installed; you need to perform some post installation steps

## Schedule Master

Before you deploy a router/registry...you'll need to make the master scheduleable

	
	root@master# oadm manage-node ose3-master.example.com --schedulable
	

## Docker Registry

The registry stores docker images and metadata. If you simply deploy a pod with the registry, it will use an ephemeral volume that is destroyed once the pod exits. Any images anyone has built or pushed into the registry would disappear. That would be bad.

For now we will just show how to specify the directory and leave the NFS configuration as an exercise. On the master, as root...

	
	root@master# oadm registry \
	--config=/etc/origin/master/admin.kubeconfig \
	--credentials=/etc/origin/master/openshift-registry.kubeconfig \
	--service-account=registry \
	--images='openshift3/ose-${component}:${version}' \
	--selector="region=infra" \ 
	--mount-host=/registry



Wait a few moments and your registry will be up. Test with:

	
	root@master# curl -v $(oc get services | grep registry | awk '{print $4":"$5}/v2/' | sed 's,/[^/]\+$,/v2/,')


If you have a NFS server you'd like to use...

Deploy registry without the "--mount-host" option

	
	root@master# oadm registry \
	--config=/etc/origin/master/admin.kubeconfig \
	--credentials=/etc/origin/master/openshift-registry.kubeconfig \
	--service-account=registry \
	--images='openshift3/ose-${component}:${version}' \
	--selector="region=infra" 


Then specify backend storage

	
	root@master# oc volume deploymentconfigs/docker-registry --add --overwrite --name=registry-storage --mount-path=/registry --source='{"nfs": { "server": "`<fqdn>`", "path": "/path/to/export"}}`</code>`
	
	
	
	There are known issues when using multiple registry replicas with the same NFS volume. We recommend changing the docker-registry service’s sessionAffinity to ClientAPI like this:
	`<code>`
	root@master# oc get -o yaml svc docker-registry | \
	      sed 's/\(sessionAffinity:\s*\).*/\1ClientIP/' | \
	      oc replace -f -


### Connecting To Docker Registry

You can connect to the docker registry hosted by OpenShift. You can do this and do "pull" and "pushes" directly into the registry. Follow the steps below to get this behavior

#### Secure Registry

After you [deploy the registry](openshift_enterprise_3.x#docker_registry) find out the service IP:PORT mapping

	
	[root@ose3-master ~]# oc get se docker-registry
	NAME              LABELS                    SELECTOR                  IP(S)            PORT(S)
	docker-registry   docker-registry=default   docker-registry=default   172.30.209.118   5000/TCP


Create a server certificate for the registry service IP and the fqdn that's going to be your route (in this example it's *** docker-registry.cloudapps.example.com ***):

	
	[root@ose3-master ~]# CA=/etc/openshift/master
	[root@ose3-master ~]# oadm create-server-cert --signer-cert=$CA/ca.crt --signer-key=$CA/ca.key --signer-serial=$CA/ca.serial.txt --hostnames='docker-registry.cloudapps.example.com,172.30.209.118' --cert=registry.crt --key=registry.key


Create the secret for the registry certificates

	
	[root@ose3-master ~]# oc secrets new registry-secret registry.crt registry.key


Add the secret to the registry pod’s service account (i.e., the "registry" service account)

	
	[root@ose3-master ~]# oc secrets add serviceaccounts/registry secrets/registry-secret


Create the directory where the registry will mount the keys

	
	[root@ose3-master ~]# mkdir /registry-secrets
	[root@ose3-master ~]# cp registry.crt /registry-secrets
	[root@ose3-master ~]# cp registry.key /registry-secrets



Add the secret volume to the registry deployment configuration

	
	[root@ose3-master ~]# oc volume dc/docker-registry --add --type=secret --secret-name=registry-secret -m /registry-secrets 


Enable TLS by adding the following environment variables to the registry deployment configuration

	
	oc env dc/docker-registry REGISTRY_HTTP_TLS_CERTIFICATE=/registry-secrets/registry.crt  REGISTRY_HTTP_TLS_KEY=/registry-secrets/registry.key


Validate the registry is running in TLS mode. Wait until the // docker-registry // pod status changes to *Running* and verify the docker logs for the registry container. You should find an entry for *listening on :5000, tls*

	
	[root@ose3-master ~]# oc get pods
	NAME                      READY     STATUS    RESTARTS   AGE
	docker-registry-3-yqy8v   1/1       Running   0          25s
	router-1-vhjdc            1/1       Running   1          2d
	[root@ose3-master ~]# oc logs docker-registry-3-yqy8v | grep tls
	time="2015-08-27T16:34:56-04:00" level=info msg="listening on :5000, tls" instance.id=440700c4-16e2-4725-81c5-5835f72c7119 


Copy the CA certificate to the docker certificates directory. This must be done on all nodes in the cluster

	
	[root@ose3-master ~]# mkdir -p /etc/docker/certs.d/172.30.209.118:5000
	[root@ose3-master ~]# mkdir -p /etc/docker/certs.d/docker-registry.cloudapps.example.com:5000
	[root@ose3-master ~]# cp /etc/openshift/master/ca.crt /etc/docker/certs.d/172.30.209.118\:5000/
	[root@ose3-master ~]# cp /etc/openshift/master/ca.crt /etc/docker/certs.d/docker-registry.cloudapps.example.com\:5000/
	[root@ose3-master ~]# for i in ose3-node{1..2}.example.com; do ssh ${i} mkdir -p /etc/docker/certs.d/172.30.209.118\:5000; ssh ${i} mkdir -p /etc/docker/certs.d/docker-registry.cloudapps.example.com\:5000; scp /etc/openshift/master/ca.crt root@${i}:/etc/docker/certs.d/172.30.209.118\:5000/; scp /etc/openshift/master/ca.crt root@${i}:/etc/docker/certs.d/docker-registry.cloudapps.example.com\:5000/; done



#### Expose Registry

Now expose your registry

Create a route

	
	[root@ose3-master ~]# oc expose svc/docker-registry --hostname=docker-registry.cloudapps.example.com


Next edit the route and add the TLS termination to be "passthrough"...in the end it should look like this

	
	[root@ose3-master ~]# oc get route/docker-registry -o yaml 
	apiVersion: v1
	kind: Route
	metadata:
	  annotations:
	    openshift.io/host.generated: "false"
	  creationTimestamp: 2015-08-27T20:58:16Z
	  labels:
	    docker-registry: default
	  name: docker-registry
	  namespace: default
	  resourceVersion: "9557"
	  selfLink: /osapi/v1beta3/namespaces/default/routes/docker-registry
	  uid: 56a78ac4-4cfe-11e5-9ae1-525400baad4f
	spec:
	  host: docker-registry.cloudapps.example.com
	  tls:
	    termination: passthrough
	  to:
	    kind: Service
	    name: docker-registry
	status: {}


#### Connect to the Registry

Copy the CA cert to the client

	
	[root@ose3-master ~]# scp /etc/openshift/master/ca.crt 172.16.1.251:/tmp/


On the client, copy the cert into the created directory

	
	[christian@rhel7 ~]$ sudo mkdir /etc/docker/certs.d/docker-registry.cloudapps.example.com\:5000/
	[christian@rhel7 ~]$ sudo cp /tmp/ca.crt /etc/docker/certs.d/docker-registry.cloudapps.example.com\:5000/
	[christian@rhel7 ~]$ sudo cp -r /etc/docker/certs.d/docker-registry.cloudapps.example.com\:5000/ /etc/docker/certs.d/docker-registry.cloudapps.example.com
	[christian@rhel7 ~]$ sudo systemctl restart docker
	[christian@rhel7 ~]$ sudo systemctl restart docker


Obtain a key from oc (hey that rhymed!)

	
	[christian@rhel7 ~]$ oc whoami -t
	YMQeiPbrMNxgR9mWmSzr1utX7IIJWL-QSpnlBgK8XBU


Use this key to login

	
	[christian@rhel7 ~]$ docker login -u christian -e chernand@redhat.com -p YMQeiPbrMNxgR9mWmSzr1utX7IIJWL-QSpnlBgK8XBU docker-registry.cloudapps.example.com
	WARNING: login credentials saved in /home/christian/.docker/config.json
	Login Succeeded


Test it by pulling busybox to one of your projects

	
	[christian@rhel7 ~]$ oc get projects
	NAME      DISPLAY NAME        STATUS
	java      Java Applications   Active
	myphp     PHP Applicaitons    Active
	[christian@rhel7 ~]$ docker pull busybox
	[christian@rhel7 ~]$ docker tag busybox docker-registry.cloudapps.example.com/myphp/mybusybox
	[christian@rhel7 ~]$ docker push  docker-registry.cloudapps.example.com/myphp/mybusybox


On the master...verify that it's in the registry

	
	[root@ose3-master ~]# oc get is -n myphp

## Router

The OpenShift router is the ingress point for all traffic destined for services in your OpenShift installation.

First create the certificate that will be used for all default SSL connections

	
	root@master# CA=/etc/origin/master
	root@master# oadm ca create-server-cert --signer-cert=$CA/ca.crt --signer-key=$CA/ca.key --signer-serial=$CA/ca.serial.txt --hostnames='*.cloudapps.example.com' --cert=cloudapps.crt --key=cloudapps.key
	root@master# cat cloudapps.crt cloudapps.key $CA/ca.crt > cloudapps.router.pem


Now create the router

	
	root@master# oadm router --default-cert=cloudapps.router.pem --credentials='/etc/origin/master/openshift-router.kubeconfig' --selector='region=infra' --images='openshift3/ose-${component}:${version}' --service-account=router


## Host Path

If you are going to add "hostPath" then you might need to do the following

	
	oc edit scc privileged


And add under users

	
	- system:serviceaccount:default:registry
	- system:serviceaccount:default:docker


Maybe this will work too?

	
	oadm policy add-scc-to-user privileged -z registry
	oadm policy add-scc-to-user privileged -z router

## Image Streams

Now that you have your router and docker-registry up and running you can populate OSE with ImageStreams (canned versions of docker images supported by Red hat)

"One Shot" command

	
	root@master# cd /usr/share/openshift/
	root@master# find examples/ -type f -name '*.json' -exec oc create -f {} -n openshift \;


**__NOTE__** You might run into a bug where you see the "registry" address like so

	
	root@master# oc get is -n openshift
	NAME                                  DOCKER REPO                                                                  TAGS                         UPDATED
	jboss-amq-62                          registry.access.redhat.com/jboss-amq-6/amq62-openshift                       1.1,1.1-2,latest             26 minutes ago
	jboss-eap64-openshift                 registry.access.redhat.com/jboss-eap-6/eap64-openshift                       1.1,1.1-2,latest             26 minutes ago
	jboss-webserver30-tomcat7-openshift   registry.access.redhat.com/jboss-webserver-3/webserver30-tomcat7-openshift   latest,1.1,1.1-2             26 minutes ago
	jboss-webserver30-tomcat8-openshift   registry.access.redhat.com/jboss-webserver-3/webserver30-tomcat8-openshift   1.1-3,latest,1.1             26 minutes ago
	jenkins                               172.30.116.164:5000/openshift/jenkins                                        1,latest                     26 minutes ago
	mongodb                               172.30.116.164:5000/openshift/mongodb                                        latest,2.4,2.6               26 minutes ago
	mysql                                 172.30.116.164:5000/openshift/mysql                                          5.6,latest,5.5               26 minutes ago
	nodejs                                172.30.116.164:5000/openshift/nodejs                                         0.10,latest                  26 minutes ago
	perl                                  172.30.116.164:5000/openshift/perl                                           5.16,5.20,latest             26 minutes ago
	php                                   172.30.116.164:5000/openshift/php                                            5.6,latest,5.5               26 minutes ago
	postgresql                            172.30.116.164:5000/openshift/postgresql                                     9.4,latest,9.2               26 minutes ago
	python                                172.30.116.164:5000/openshift/python                                         3.4,latest,2.7 + 1 more...   26 minutes ago
	ruby                                  172.30.116.164:5000/openshift/ruby                                           2.2,2.0,latest               About a minute ago
	wildfly                               172.30.116.164:5000/openshift/wildfly                                        8.1,latest                   About a minute ago


Fix this by doing the following

	
	root@master# oc delete is --all -n openshift
	root@master# oc create -f https://raw.githubusercontent.com/rhtconsulting/rhc-ose/openshift-enterprise-3/provisioning/templates/image-streams-rhel7-ose3_0_2.json -n openshift
	root@master# oc create -n openshift -f /usr/share/openshift/examples/xpaas-streams/jboss-image-streams.json
	root@master# oc get is -n openshift
	NAME                                  DOCKER REPO                                                                  TAGS                                     UPDATED
	jboss-amq-62                          registry.access.redhat.com/jboss-amq-6/amq62-openshift                       latest,1.1,1.1-2                         1 seconds ago
	jboss-eap64-openshift                 registry.access.redhat.com/jboss-eap-6/eap64-openshift                       1.1,1.1-2,latest                         2 seconds ago
	jboss-webserver30-tomcat7-openshift   registry.access.redhat.com/jboss-webserver-3/webserver30-tomcat7-openshift   latest,1.1,1.1-2                         3 seconds ago
	jboss-webserver30-tomcat8-openshift   registry.access.redhat.com/jboss-webserver-3/webserver30-tomcat8-openshift   1.1,1.1-3,latest                         3 seconds ago
	jenkins                               registry.access.redhat.com/openshift3/jenkins-1-rhel7                        1.6-3,1.609-14,1 + 1 more...             10 seconds ago
	mongodb                               registry.access.redhat.com/openshift3/mongodb-24-rhel7                       v3.0.0.0,v3.0.1.0,v3.0.2.0 + 4 more...   11 seconds ago
	mysql                                 registry.access.redhat.com/openshift3/mysql-55-rhel7                         5.5-8,v3.0.1.0,v3.0.2.0 + 4 more...      14 seconds ago
	nodejs                                registry.access.redhat.com/openshift3/nodejs-010-rhel7                       v3.0.2.0,0.10,0.10-12 + 4 more...        20 seconds ago
	perl                                  registry.access.redhat.com/openshift3/perl-516-rhel7                         v3.0.0.0,v3.0.1.0,v3.0.2.0 + 4 more...   19 seconds ago
	php                                   registry.access.redhat.com/openshift3/php-55-rhel7                           v3.0.0.0,latest,v3.0.1.0 + 4 more...     17 seconds ago
	postgresql                            registry.access.redhat.com/openshift3/postgresql-92-rhel7                    latest,v3.0.0.0,v3.0.1.0 + 2 more...     13 seconds ago
	python                                registry.access.redhat.com/openshift3/python-33-rhel7                        3.1.0,3.3,3.3-13 + 4 more...             16 seconds ago
	ruby                                  registry.access.redhat.com/openshift3/ruby-20-rhel7                          v3.0.1.0,v3.0.2.0,2.0-12 + 4 more...     22 seconds ago


## LDAP Configuration

First (if using // ldaps //) you need to download the CA certificate (below example is using Red Hat IdM server)

	
	root@master# curl  http://ipa.example.com/ipa/config/ca.crt >> /etc/openshift/master/my-ldap-ca-bundle.crt


Make a backup copy of the config file

	
	root@master# cp /etc/openshift/master/master-config.yaml{,.bak}


Edit the // /etc/openshift/master/master-config.yaml // file with the following changes under the // identityProviders // section

	
	  identityProviders:
	  - name: "my_ldap_provider"
	    challenge: true
	    login: true
	    provider:
	      apiVersion: v1
	      kind: LDAPPasswordIdentityProvider
	      attributes:
	        id:
	        - dn
	        email:
	        - mail
	        name:
	        - cn
	        preferredUsername:
	        - uid
	      bindDN: "cn=directory manager"
	      bindPassword: "secret"
	      ca: my-ldap-ca-bundle.crt
	      insecure: false
	      url: "ldaps://ipa.example.com/cn=users,cn=accounts,dc=example,dc=com?uid"


Note you can customize what attributes it searches for. First non empty attribute returned is used.

Restart the openshift-master service

	
	systemctl restart atomic-openshift-master


### Active Directory

AD usually is using *sAMAccountName* as uid for login. Use the following ldapsearch to validate the informaiton

	
	ldapsearch -x -D "CN=xxx,OU=Service-Accounts,OU=DCS,DC=homeoffice,DC=example,DC=com" -W -H ldaps://ldaphost.example.com -b "ou=Users,dc=office,dc=example,DC=com" -s sub 'sAMAccountName=user1'


If the ldapsearch did not return any user, it means -D or -b may not be correct. Retry different *baseDN*. If there is too many entries returns, add filter to your search. Filter example is *(objectclass=people)* or *(objectclass=person)* if still having issues; increase logging as *OPTIONS=--loglevel=5* in // /etc/sysconfig/atomic-openshift-master //

If you see an error in // journalctl -u atomic-openshift-master//  there might be a conflict with the user identity when user trying to login (if you used *htpasswd* beforehand). Just do the following...

	
	oc get user oc delete user user1


Inspiration from :


*  [https://access.redhat.com/solutions/2016873](https///access.redhat.com/solutions/2016873)

*  [https://access.redhat.com/solutions/1978013](https///access.redhat.com/solutions/1978013)

The configuration in *master-config.yaml* Should look something like this:

	
	oauthConfig:
	  assetPublicURL: https://master.example.com:8443/console/
	  grantConfig:
	    method: auto
	  identityProviders:
	  - name: "OfficeAD"
	    challenge: true
	    login: true
	    provider:
	      apiVersion: v1
	      kind: LDAPPasswordIdentityProvider
	      attributes:
	        id:
	        - dn
	        email:
	        - mail
	        name:
	        - cn
	        preferredUsername:
	        - sAMAccountName
	      bindDN: "CN=LinuxSVC,OU=Service-Accounts,OU=DCS,DC=office,DC=example,DC=com"
	      bindPassword: "password"
	      ca: ad-ca.pem.crt
	      insecure: false
	      url: "ldaps://ad-server.example.com:636/CN=Users,DC=hoffice,DC=example,DC=com?sAMAccountName?sub"


If you need to look for a subclass...

	
	"ldaps://ad.corp.example.com:636/OU=Users,DC=corp,DC=example,DC=com?sAMAccountName?sub?(&(objectClass=person)"

## Users

###  Create User

The ansible scripts configured authentication using // htpasswd // so just create the users using the proper method

	
	root@host# htpasswd -b /etc/openshift/openshift-passwd demo demo


### Adding User to group

Currently, you can only add a user to a group by setting the "group" array to a group

	
	[root@ose3-master ~]# oc edit user/christian -o json
	{
	    "kind": "User",
	    "apiVersion": "v1",
	    "metadata": {
	        "name": "christian",
	        "selfLink": "/osapi/v1beta3/users/christian",
	        "uid": "a5c96638-4084-11e5-8a3c-fa163e2e3caf",
	        "resourceVersion": "1182",
	        "creationTimestamp": "2015-08-11T23:56:56Z"
	    },
	    "identities": [
	        "htpasswd_auth:christian"
	    ],
	    "groups": [
	        "mygroup"
	    ]
	}
	


# Misc

Misc info in no particular order

## Promotion

Rough notes taken from... https://blog.openshift.com/promoting-applications-across-environments/

### Create a new Project

Here are the commands used to create a new project with name “development” and providing “edit” access to developer and “view” access to the tester.

	
	oc new-project development —display-name="Development Project"
	oc policy add-role-to-user edit dev1
	oc policy add-role-to-user view test1

### Create a QA project

Commands needed to create a QA project and provide “edit’ access to the tester.

	
	oc new-project testing —display-name="QA Project"
	oc policy add-role-to-user edit test1


### Enable the test project to pull development images

Assigning the ''system:image-puller'' role to the service account “testing” which is the default service account for the testing project on the development project. By doing this, we are enabling the testing project to be able to pull images from the development project.

	
	oc policy add-role-to-group system:image-puller system:serviceaccounts:testing -n development

### Create an application in development

Switch over as developer and create an application in the development project.

	
	oc login -u dev1
	oc project development
	oc new-app --template=eap6-basic-sti -p APPLICATION_NAME=myapp,APPLICATION_HOSTNAME=myapp-dev.apps.demov3.osecloud.com,EAP_RELEASE=6.4,GIT_URI=https://github.com/VeerMuchandi/kitchensink.git,GIT_REF=,GIT_CONTEXT_DIR= -l name=myapp

### Identifying the image id

Finding the image stream name and identifying the full image id.

	
	oc get is
	oc describe is

The ''describe is'' command will show the full image id. You can copy that into clipboard. Use that to tag the specific image to promote.

	
	oc tag development/myapp:promote

### Deploy an application in the test project

Login as tester and deploy an application in the “testing” project.

	
	oc login -u test1
	oc project testing
	oc new-app development/myapp:promote

Note the service name and create a route.

	
	oc get svc
	oc expose svc


### Quick Notes

Add two users

	
	 htpasswd -b /etc/openshift/openshift-passwd qa qa
	 htpasswd -b /etc/openshift/openshift-passwd dev dev


As QA

	
	qa@host$ oc new-project qa
	qa@host$ oc new-app openshift/php~https://github.com/RedHatWorkshops/welcome-php.git
	qa@host$ oc get is


As DEV

	
	dev@host$ oc new-project dev


As Admin use the syntax

''oc policy add-role-to-group system:image-puller system:serviceaccounts:**$PROJECT_YOU_WANT_TO_PULL_TO** -n **$PROJECT_YOU_WANT_TO_PULL_FROM** ''

''oc policy add-role-to-user view **$USER_THAT_OWNS_THE_PROJECT_YOU_WANT_TO_PULL_TO** -n **$PROJECT_YOU_WANT_TO_PULL_FROM**''

	
	root@master# oc policy add-role-to-group system:image-puller system:serviceaccounts:dev -n qa
	root@master# oc policy add-role-to-user view dev -n qa


As DEV

	
	dev@host$ oc new-app qa/welcome-php

## Cluster Metrics

Quick and dirty notes; will clean up later

	
	oc project openshift-infra
	
	oc create -f - <<API
	apiVersion: v1
	kind: ServiceAccount
	metadata:
	  name: metrics-deployer
	secrets:
	- name: metrics-deployer
	API
	
	oadm policy add-role-to-user edit system:serviceaccount:openshift-infra:metrics-deployer
	
	oadm policy add-cluster-role-to-user cluster-reader system:serviceaccount:openshift-infra:heapster
	
	oc secrets new metrics-deployer nothing=/dev/null
	
	cp /usr/share/openshift/examples/infrastructure-templates/enterprise/metrics-deployer.yaml .
	
	oc process -f metrics-deployer.yaml -v IMAGE_PREFIX=openshift3/,IMAGE_VERSION=latest,HAWKULAR_METRICS_HOSTNAME=console.cloudapps.example.com,USE_PERSISTENT_STORAGE=false | oc create -f -
	
	[root@ose3-master ~]# grep -i metrics -B10 /etc/origin/master/master-config.yaml
	  masterPublicURL: https://ose3-master.example.com:8443
	  publicURL: https://ose3-master.example.com:8443/console/
	  servingInfo:
	    bindAddress: 0.0.0.0:8443
	    bindNetwork: tcp4
	    certFile: master.server.crt
	    clientCA: ""
	    keyFile: master.server.key
	    maxRequestsInFlight: 0
	    requestTimeoutSeconds: 0
	  metricsPublicURL: "https://console.cloudapps.example.com/hawkular/metrics" # add this to the /etc/origin/master/master-config.yaml file 
	
	systemctl restart atomic-openshift-master.service 
	
	# Visit https://console.cloudapps.example.com/hawkular/metrics in the browser to accept cert


"Clean UP" procedure

	
	oc project openshift-infra
	oc delete all --all
	oc delete templates --all
	oc delete secrets `oc get secrets | egrep 'metrics|hawk|heap' | awk '{print $1}'`
	oc delete sa hawkular cassandra heapster


## Logging

Quick notes

	
	oc create -n openshift -f \
	    /usr/share/openshift/examples/infrastructure-templates/enterprise/logging-deployer.yaml
	oadm new-project logging --node-selector=""
	oc project logging
	oc secrets new logging-deployer nothing=/dev/null
	oc create -f - <<API
	apiVersion: v1
	kind: ServiceAccount
	metadata:
	  name: logging-deployer
	secrets:
	- name: logging-deployer
	API
	oc policy add-role-to-user edit --serviceaccount logging-deployer
	oadm policy add-scc-to-user  \
	    privileged system:serviceaccount:logging:aggregated-logging-fluentd
	oadm policy add-cluster-role-to-user cluster-reader \
	    system:serviceaccount:logging:aggregated-logging-fluentd
	oc new-app logging-deployer-template \
	             --param KIBANA_HOSTNAME=kibana.cloudapps.example.com \
	             --param ES_CLUSTER_SIZE=1 \
	             --param PUBLIC_MASTER_URL=https://ose3node8.example.com:8443
	             --param MASTER_URL=https://ose3node8.example.com:8443
	oc new-app logging-support-template
	oc import-image logging-auth-proxy:3.2.0 \
	     --from registry.access.redhat.com/openshift3/logging-auth-proxy:3.2.0
	oc import-image logging-kibana:3.2.0 \
	     --from registry.access.redhat.com/openshift3/logging-kibana:3.2.0
	oc import-image logging-elasticsearch:3.2.0 \
	     --from registry.access.redhat.com/openshift3/logging-elasticsearch:3.2.0
	oc import-image logging-fluentd:3.2.0 \
	     --from registry.access.redhat.com/openshift3/logging-fluentd:3.2.0
	#update the logging-es pod name
	#oc patch dc/`<logging-es-mrwi1e61>` -p '{"spec":{"template":{"spec":{"nodeSelector":{"nodeLabel":"logging-es-node-1"}}}}}'
	oc new-app logging-es-template
	#loggingPublicURL: "https://kibana.cloudapps.example.com"
	oc scale dc/logging-fluentd --replicas=2

## Login

### User Login

To login

	
	user@host$ oc login https://ose3-master.example.com:8443 --insecure-skip-tls-verify --username=demo


### Admin Login

On the master

	
	root@master# oc login -u system:admin -n default

## Projects

To create a project as a user run the following command.

	
	user@host$ oc new-project demo --display-name="Demo Projects" --description="Demo projects go here"


If you're an OSE admin and want to create a project and assign a user to it with the *--admin=${user}* command.

	
	root@master# oadm new-project demo --display-name="OpenShift 3 Demo" --description="This is the first demo project with OpenShift v3" --admin=joe


## Create App

This is an example PHP-application you can use to test your OSEv3 environment.

Here is an example:

	
	user@host$ oc new-app openshift/php~https://github.com/christianh814/php-example-ose3


Things to keep in mind:

*  *ose new-app* Creates a new application on OSE3

*  *openshift/php* This tells OSEv3 to use the PHP image stream provided by OSE

*  Provide the git URL for the project
    * Syntax is "//**imagestream**~**souce**//"

Once you created the app, start your build

	
	user@host$ oc start-build php-example-ose3


View the build logs if you wish. Note the *-1* ...this is the build number. Find the build number with *oc get builds*

	
	user@host$ oc build-logs php-example-ose3-1


Once the build completes; create and add your route:

	
	user@host$ oc expose service php-example-ose3 --hostname=php-example.cloudapps.example.com



Scale up as you wish

	
	user@host$ oc scale --replicas=3 dc/php-example-ose3


If you'd like to add another route (aka "alias"); then you need to specify a new name for it

	
	user@host$ oc expose service php-example-ose3 --name=hello-openshift --hostname=hello-openshift.cloudapps.example.com


If you want to add SSL to your app. Edit the route and add "tls: termination: edge" it should look like this (other if you're serving SSL through the app itself. you may want to use "passthrough")

	
	user@host$ oc get route/jenkins -n demo -o yaml
	apiVersion: v1
	kind: Route
	metadata:
	  annotations:
	    openshift.io/host.generated: "true"
	  creationTimestamp: 2015-07-07T17:20:17Z
	  name: jenkins
	  namespace: demo
	  resourceVersion: "30360"
	  selfLink: /osapi/v1beta3/namespaces/demo/routes/jenkins
	  uid: 6fbd1876-24cc-11e5-accf-525400e68df5
	spec:
	  host: jenkins.demo.cloudapps.example.com
	  tls:
	    termination: edge
	  to:
	    kind: Service
	    name: jenkins
	status: {}
	


Note: To see what imageStreams are available to you...

	
	user@host$  oc get imageStreams -n openshift


Enter Container

Enter your container with the *oc rsh* command

	
	user@host$  oc rsh ${podID} 


Create an app with a 'Dockerfile' in github

	
	user@host$ oc new-app https://github.com/christianh814/ose3-ldap-auth --strategy=docker --name=auth -l appname=auth


Use Template for JBOSS

	
	user@host$ git clone https://github.com/openshift/openshift-ansible
	user@host$ cd openshift-ansible/roles/openshift_examples/files/examples/xpaas-templates/
	user@host$ oc process -f eap6-basic-sti.json -v APPLICATION_NAME=ks,APPLICATION_HOSTNAME=ks.demo.sbx.osecloud.com,GIT_URI=https://github.com/RedHatWorkshops/kitchensink,GIT_REF="",GIT_CONTEXT_DIR="" | oc create -f -


Custom Service

	
	user@host$ oc expose dc/basicauthurl --port=443 --generator=service/v1 -n auth

## Rolling Deployments

By default, when a new build is fired off it will stop the application while the new container is created. You can change the deployment time on an app

	
	user@host$ oc edit dc/php-example-ose3


Change the *Strategy* to *Rolling*

## Build Webhooks

You can trigger a build using the generic webhook (there is one for github too)

	
	curl -i -H "Accept: application/json" -H "X-HTTP-Method-Override: PUT" -X POST -k https://ose3-master.example.com:8443/osapi/v1beta3/namespaces/wiring/buildconfigs/ruby-example/webhooks/secret101/generic


## Run Dockerhub Images

In order to run Dockerhub images you need to lift the security in your cluster so that images are not forced to run as a pre-allocated UID, without granting everyone access to the privileged SCC, you can edit the restricted SCC and change the ***runAsUser*** strategy:

	
	root@master# oc edit scc restricted



...Change **//runAsUser//** Type to ***RunAsAny***.

	

**//__WARING:__//**This allows images to run as the root UID if no USER is specified in the Dockerfile.

Now you can pull docker images

	
	user@host$ oc new-app fedora/apache --name=apache
	user@host$ oc expose service apache


## SSH Key For Git

Create the secret first before using the SSH key to access the private repository:

	
	$ oc secrets new scmsecret ssh-privatekey=$HOME/.ssh/id_rsa


Add the secret to the builder service account:

	
	$ oc secrets add serviceaccount/builder secrets/scmsecret

Add a sourceSecret field into the source section inside the buildConfig and set it to the name of the secret that you created, in this case scmsecret:

	
	{
	  "apiVersion": "v1",
	  "kind": "BuildConfig",
	  "metadata": {
	    "name": "sample-build",
	  },
	  "parameters": {
	    "output": {
	      "to": {
	        "name": "sample-image"
	      }
	    },
	    "source": {
	      "git": {
	        "uri": "git@repository.com:user/app.git" 
	      },
	      "sourceSecret": {
	        "name": "scmsecret"
	      },
	      "type": "Git"
	    },
	    "strategy": {
	      "sourceStrategy": {
	        "from": {
	          "kind": "ImageStreamTag",
	          "name": "python-33-centos7:latest"
	        }
	      },
	      "type": "Source"
	    }
	  }

The URL of private repository is usually in the form // git@example.com:`<username>`/`<repository>` //

## Liveness Check for Apps

If A pod dies kubernetes will fire the pod back up.

But what if the pod is running but the application (pid) inside is hung or dead? Kubernetes needs a way to monitor the application.

This is done with a "health check" outlined [here](https///github.com/GoogleCloudPlatform/kubernetes/tree/master/docs/user-guide/liveness)

First edit the deploymentConfig

	
	user@host$ oc edit dc/myapp -o yaml


Inside "containers" and just after "image" add the following

	
	    livenessProbe:
	      httpGet:
	        path: /healthz
	        port: 8080
	      initialDelaySeconds: 15
	      timeoutSeconds: 1


In the end it should look something like this...

	
	apiVersion: v1
	kind: DeploymentConfig
	metadata:
	  creationTimestamp: 2015-07-30T16:15:16Z
	  labels:
	    appname: myapp
	  name: myapp
	  namespace: demo
	  resourceVersion: "255603"
	  selfLink: /osapi/v1beta3/namespaces/demo/deploymentconfigs/myapp
	  uid: 2a7f06f8-36d6-11e5-ba31-fa163e2e3caf
	spec:
	  replicas: 1
	  selector:
	    deploymentconfig: myapp
	  strategy:
	    resources: {}
	    type: Recreate
	  template:
	    metadata:
	      creationTimestamp: null
	      labels:
	        deploymentconfig: myapp
	    spec:
	      containers:
	      - env:
	        - name: PEARSON
	          value: value
	        image: 172.30.182.253:5000/demo/myapp@sha256:fec918b3e488a5233b2840e1c8db7d01ee9c2b9289ca0f69b45cfea955d629b2
	        imagePullPolicy: Always
	        livenessProbe:
	          httpGet:
	            path: /info.php
	            port: 8080
	          initialDelaySeconds: 15
	          timeoutSeconds: 1
	        name: myapp
	        ports:
	        - containerPort: 8080
	          name: myapp-tcp-8080
	          protocol: TCP
	        resources: {}
	        securityContext:
	          capabilities: {}
	          privileged: false
	        terminationMessagePath: /dev/termination-log
	      dnsPolicy: ClusterFirst
	      restartPolicy: Always
	  triggers:
	  - type: ConfigChange
	  - imageChangeParams:
	      automatic: true
	      containerNames:
	      - myapp
	      from:
	        kind: ImageStreamTag
	        name: myapp:latest
	      lastTriggeredImage: 172.30.182.253:5000/demo/myapp@sha256:fec918b3e488a5233b2840e1c8db7d01ee9c2b9289ca0f69b45cfea955d629b2
	    type: ImageChange
	status:
	  details:
	    causes:
	    - type: ConfigChange
	  latestVersion: 8
	



## REST API Notes

**//NOTE: These are QND Notes!//**

First get a token 

	
	curl -D - -u username:secret -k "https://ose3-master.sandbox.osecloud.com:8443/oauth/authorize?response_type=token&client_id=openshift-challenging-client" 2>&1 | grep -oP "access_token=\K[^&]*"


Use that token to list (GET) things

	
	curl -X GET -k "https://ose3-master.sandbox.osecloud.com:8443/api/v1/namespaces/demo/pods?access_token=0GGIgJ5JHG-TJ3k0pZfDuEKlvcjStCBXXOYEssmQ73U"


Some URLS are different for some reason

	
	 curl -X GET -k "https://ose3-master.sandbox.osecloud.com:8443/oapi/v1/namespaces/demo/builds?access_token=0GGIgJ5JHG-TJ3k0pZfDuEKlvcjStCBXXOYEssmQ73U" 


Do things with "POST"...fireoff a build example

	
	curl -X POST -k "https://ose3-master.sandbox.osecloud.com:8443/oapi/v1/namespaces/demo/buildconfigs/myapp/instantiate?access_token=SgS8mjOvsWVRd_grHFLn99sbnCT6TFeUxIAe3H4fgW0" -d '{"kind":"BuildRequest","apiVersion":"v1","metadata":{"name":"myapp","creationTimestamp":null}}'


To figure out the API is hit and miss...mostly do you "oc" commands with the "loglevel=8" to see what's happening and emulate that

	
	oc `<command>` --loglevel=8


## Node Selector

For some reason you now have to...

	
	oc edit namespace default


and add...

	
	openshift.io/node-selector: region=infra


In 3.2; you can do this from the command line

	
	root@master# oc annotate namespace default openshift.io/node-selector=region=infra

## NFS For Persistent Storage

You can provision your OpenShift cluster with persistent storage using NFS. The Kubernetes persistent volume framework allows administrators to provision a cluster with persistent storage and gives users a way to request those resources without having any knowledge of the underlying infrastructure.

### Adding Storage: Master

Example:

	
	{
	  "apiVersion": "v1",
	  "kind": "PersistentVolume",
	  "metadata": {
	    "name": "pv0001"
	  },
	  "spec": {
	    "capacity": {
	        "storage": "20Gi"
	        },
	    "accessModes": [ "ReadWriteMany" ],
	    "nfs": {
	        "path": "/var/export/vol1",
	        "server": "nfs.example.com"
	    }
	  }
	}


Create this object as the root (administrative) user

	
	root@master# oc create -f pv0001.json 
	persistentvolumes/pv0001


This defines a volume for OpenShift projects to use in deployments. The storage should correspond to how much is actually available (make each volume a separate filesystem if you want to enforce this limit). Take a look at it now:

	
	root@master# oc describe persistentvolumes/pv0001
	Name:		pv0001
	Labels:		`<none>`
	Status:		Available
	Claim:		
	Reclaim Policy:	%!d(api.PersistentVolumeReclaimPolicy=Retain)
	Message:	%!d(string=)



### Adding Storage: Client

Before you add the PV make sure you allow containers to mount NFS volumes

	
	root@master# setsebool -P virt_use_nfs=true
	root@node1#  setsebool -P virt_use_nfs=true
	root@node2#  setsebool -P virt_use_nfs=true


Now that the administrator has provided a PersistentVolume, any project can make a claim on that storage. We do this by creating a PersistentVolumeClaim that specifies what kind and how much storage is desired:

	
	{
	  "apiVersion": "v1",
	  "kind": "PersistentVolumeClaim",
	  "metadata": {
	    "name": "claim1"
	  },
	  "spec": {
	    "accessModes": [ "ReadWriteMany" ],
	    "resources": {
	      "requests": {
	        "storage": "20Gi"
	      }
	    }
	  }
	}


We can have alice do this in the project you created:

	
	user@host$ oc create -f pvclaim.json 
	persistentvolumeclaims/claim1


This claim will be bound to a suitable PersistentVolume (one that is big enough and allows the requested accessModes). The user does not have any real visibility into PersistentVolumes, including whether the backing storage is NFS or something else; they simply know when their claim has been filled ("bound" to a PersistentVolume).

	
	user@host$ oc get pvc
	NAME      LABELS    STATUS    VOLUME
	claim1    map[]     Bound     pv0001


Finally, we need to modify the DeploymentConfig to specify that this volume should be mounted 

	
	user@host$ oc edit dc/jenkins -o json


You'll notice in the // ** spec: volumes: ** // section that there is an *** "emptyDir" *** conifg there. Change it so you refer to your "claim1" you did above. Should look like this

	
	            "spec": {
	                "volumes": [
	                    {
	                        "name": "jenkins-16-centos7-volume-1",
	                        "persistentVolumeClaim": { "claimName": "claim1" }
	                    }
	                ],


Verify by scrolling down and seeing that the config in-fact does use *** jenkins-16-centos7-volume-1 *** 

	
	 "volumeMounts": [
	                            {
	                                "name": "jenkins-16-centos7-volume-1",
	                                "mountPath": "/var/lib/jenkins"
	                            }
	                        ],


**__OR__** you can do it from the cli

	
	oc volumes dc/gogs --add --claim-name=gogs-repos-claim --mount-path=/home/gogs/gogs-repositories -t persistentVolumeClaim 

