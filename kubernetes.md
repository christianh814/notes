# Deprecated

This information is deprecated. I'm keeping this here for historical information only.

See newer information [ HERE ]( kubernetes_v2 )

# Intro

Using Kubernetes you can run and manage docker containers from one system (a Kubernetes master) and deploy them to run on other systems (Kubernetes nodes). This procedure will help you:

*     Set up three Red Hat Enterprise Linux 7 or RHEL Atomic systems to use as a Kubernetes master and two Kubernetes nodes (also called minions).
*     Configure networking between containers using Flannel.
*     Create data files (in yaml format) that define services, pods and replication controllers.

The services you start during this procedure include etcd, cadvisor, kube-apiserver, kubecontroller-manager, kube-proxy, kube-scheduler, and kubelet. To launch containers from yaml files and otherwise manage your Kubernetes environment, you run the kubectl command.

The result of this procedure is a Kubernetes environment that you can use as a foundation for deploying Docker containers across multiple host systems.

Assumptions:

We will be using 2 RHEL 7 Atomic Host with one RHEL 7 server.

*     You are familiar with RHEL and RHEL Atomic
*     You have the proper repos setup
*     You are familiar with Docker
*     You have a basic level of understanding of Kubernetes

# Docker


Docker is an open-source project that automates the deployment of applications inside software containers, by providing an additional layer of abstraction and automation of operating-system-level virtualization on Linux. Docker uses resource isolation features of the Linux kernel such as cgroups and kernel namespaces to allow independent "containers" to run within a single Linux instance, avoiding the overhead of starting virtual machines.

# Kubernetes


Kubernetes is a powerful system, developed by Google, for managing containerized applications in a clustered environment. It aims to provide better ways of managing related, distributed components across varied infrastructure. Kubernetes groups containers into a logical unit called "pods". A pod is a 1 or more containers that share common resources (such a network and storage).

# Installation

__Nodes__

I have the following installed as a RHEL 7 Atomic Host. These hosts will be running the docker containers. Also, one (atomic01) will be the Kubernetes master. Atomic comes with everything you need to run containers by default so there is nothing to install.

*     atomic01 - 172.16.1.121
*     atomic02 - 172.16.1.122
*     atomic03 - 172.16.1.123

__Registry__

One host will be installed with RHEL 7 Server that will be running a docker registry.

*     hub - 172.16.1.224

To install the registry; install the following packages

    yum -y install docker-registry
    systemctl enable docker-registry

Also make sure firewalld is disabled

    systemctl stop firewalld
    systemctl disable firewalld

# Docker Registry Configuration


On the "hub" server; you will need to modify the systemd startup script. The best way to do this is copy it over to /etc/systemd/system directory (since files here "override" the files in /usr/lib/systemd/system)

    cp /usr/lib/systemd/system/docker-registry.service /etc/systemd/system/

Edit the /etc/systemd/system/docker-registry.service file and add *--certfile* and *--keyfile* to the start up script.

	
	    cat /etc/systemd/system/docker-registry.service
	    [Unit]
	    Description=Registry server for Docker
	
	    [Service]
	    Type=simple
	    Environment=DOCKER_REGISTRY_CONFIG=/etc/docker-registry.yml
	    EnvironmentFile=-/etc/sysconfig/docker-registry
	    WorkingDirectory=/usr/lib/python2.7/site-packages/docker-registry
	    ExecStart=/usr/bin/gunicorn --certfile /etc/pki/tls/certs/dr.crt --keyfile /etc/pki/tls/private/dr.key --access-logfile - --debug --max-requests 100 --graceful-timeout 3600 -t 3600 -k gevent -b ${REGISTRY_ADDRESS}:${REGISTRY_PORT} -w $GUNICORN_WORKERS docker_registry.wsgi:application
	    Restart=on-failure
	
	    [Install]
	    WantedBy=multi-user.target


Now generate the keyfile and the certificate

    openssl genrsa -out /etc/pki/tls/private/dr.key 1024
    openssl req -new -key /etc/pki/tls/private/dr.key -x509 -out /etc/pki/tls/certs/dr.crt

Copy the certificate locally and all nodes.

    mkdir /etc/docker/certs.d/hub.example.net:5000
    for i in atomic0{1..2}.example.net ; do ssh $i mkdir /etc/docker/certs.d/hub.example.net:5000; done
    cp /etc/pki/tls/certs/dr.crt /etc/docker/certs.d/hub.example.net:5000/dr.crt
    for i in atomic0{1..2}.example.net ; do scp /etc/pki/tls/certs/dr.crt root@$i:/etc/docker/certs.d/hub.example.net\:5000/ ; done

Now you can start and enable the service

    systemctl enable docker-registry
    systemctl start docker-registry
    systemctl status docker-registry

For those that left the firewall running; add port 5000

    firewall-cmd --zone=public --add-port=5000/tcp
    firewall-cmd --permanent --zone=public --add-port=5000/tcp

At this point your docker-registry should be up and running. Next; you'll need to build the containers found on my github account.

    git clone https://github.com/christianh814/dock-kube-example.git

In this repo, there are 2 directories webserver and dbserver. You can inspect the contents if you wish. Build the docker containers from the dockerfiles. Run docker images to verify they have been built.

    cd dock-kube-example/webserver
    docker build -t webserver .
    cd dock-kube-example/
    docker build -t dbserver .
    docker images

In order to upload images to your registry you first must "tag" them referencing your registry URL:PORT. Then you can "push" these images to your registry.

    docker tag dbserver hub.example.net:5000/dbserver
    docker tag webserver hub.example.net:5000/webserver
    docker push hub.example.net:5000/webserver
    docker push hub.example.net:5000/dbserver

Lastly, pull the images from your atomic hosts

    docker pull hub.example.net:5000/webserver
    docker pull hub.example.net:5000/dbserver

 
# Configuring Kubernetes

__Kubernetes Master__

On the master, you need to configure what your nodes are, the etcd server, and API server. For this document we will be using atomic01 as our master.

First Edit the /etc/kubernetes/controller-manager file. In the KUBELET_ADDRESSES line add your nodes using comma to separate them.

    egrep -v '^#|^$' /etc/kubernetes/controller-manager
    KUBELET_ADDRESSES="--machines=atomic02.example.net,atomic03.example.net"
    KUBE_CONTROLLER_MANAGER_ARGS=""

In the /etc/kubernetes/config file; comment out the KUBE_ETCD_SERVERS line and add a KUBE_MASTER line. Leave other settings in tact. The two changes should look like this

    egrep 'KUBE_ETCD_SERVERS|KUBE_MASTER' /etc/kubernetes/config
    # KUBE_ETCD_SERVERS="--etcd_servers=http://127.0.0.1:4001"
    KUBE_MASTER="--master=http://atomic01.example.net:8080"

In the /etc/kubernetes/apiserver file; add a new KUBE_ETCD_SERVERS line, change the KUBE_API_ADDRESS to listen on all network addresses, change KUBE_MASTER to the master fqdn, and set an address range on the KUBE_SERVICE_ADDRESS that Kubernetes can use to assign to services. The file should look something like this.

	
	    egrep -v '^#|^$' /etc/kubernetes/apiserver
	    KUBE_ETCD_SERVERS="--etcd_servers=http://atomic01.example.net:4001"
	    KUBE_API_ADDRESS="--address=0.0.0.0"
	    KUBE_API_PORT="--port=8080"
	    KUBE_MASTER="--master=127.0.0.1:8080"
	    KUBELET_PORT="--kubelet_port=10250"
	    KUBE_SERVICE_ADDRESSES="--portal_net=10.254.0.0/16"
	    KUBE_API_ARGS=""

Keep in mind what you are setting for the KUBE_SERVICE_ADDRESSES line...

    The address range is used by Kubernetes to assign to Kubernetes services.
    In the example just shown, the address range of 10.254.0.0/16 consumes a set of 10.254 subnets that can be assigned by Kubernetes as needed. For example, 10.254.1.X, 10.254.2.X and so on.
    Make sure this address range isn't used anywhere else in your environment.
    Each address range that is assigned is only used within a node and is not routable outside of that node.
    This address range must be different than the range used by flannel. (Flannel address ranges are assigned to pods.)

You need to start several services associated with a Kubernetes master. From the master, run the following command to start and enable Kubernetes systemd services on the master.

	
	# for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler; do 
	    systemctl restart $SERVICES
	    systemctl enable $SERVICES
	    systemctl status $SERVICES 
	done

__Kubernetes Nodes (minions)__

On each of the two Kubernetes nodes (atomic02 and atomic03 in this example), configure them to communicate with the master.

Due to a current bug in Kubernetes; you must create an auth file that is empty. The file is under /var/lib/kubelet/auth and it's a JSON file.

    echo "{}" > /var/lib/kubelet/auth

In the /etc/kubernetes/config file; comment out the KUBE_ETCD_SERVERS line and add a KUBE_MASTER line. Leave other settings in tact. The two changes should look like this

    egrep 'KUBE_ETCD_SERVERS|KUBE_MASTER' /etc/kubernetes/config
    # KUBE_ETCD_SERVERS="--etcd_servers=http://127.0.0.1:4001"
    KUBE_MASTER="--master=http://atomic01.example.net:8080"

In the /etc/kubernetes/kubelet  file on each node, modify KUBELET_ADDRESS to listen on network interfaces, on KUBELET_HOSTNAME replace hostname_override with the hostname or IP address of the local system: atomic02.example.com or atomic03.example.com, and KUBELET_ARGS is set --auth_path=/var/lib/kubelet/auth. Set KUBE_ETCD_SERVERS to --api_servers=http://atomic01.example.net:8080. File should look something like this:

	
	egrep -v '^#|^$' /etc/kubernetes/kubelet
	KUBELET_ADDRESS="--address=0.0.0.0"
	KUBELET_PORT="--port=10250"
	KUBELET_HOSTNAME="--hostname_override=atomic0X.example.net"
	KUBELET_ARGS="--auth_path=/var/lib/kubelet/auth"
	KUBE_ETCD_SERVERS="--api_servers=http://atomic01.example.net:8080"

On each node, add the following line to the /etc/kubernetes/proxy file

    KUBE_PROXY_ARGS="--master=http://atomic01.example.net:8080"

On each node, you need to start several services associated with a Kubernetes node

    for cmd in enable restart start status; do systemctl ${cmd} docker kube-proxy.service kubelet.service ; done

From any node test to see if etcd is running.

    curl -s -L http://atomic01.example.net:4001/version

# Setting up Flannel networking


The flannel package contains features that allow you to configure networking between the master and nodes in your Kubernetes cluster. You configure the flanneld service by creating and uploading a json configuration file with your network configuration to your etcd server (on the master). You then configure the flanneld systemd service on the master and each node to point to that etcd server and start the flanneld service.

Because the docker0 interface is probably already in place when you run this procedure, the IP address range assigned by flanneld to docker0 will not immediately take effect. You can either manually stop the docker0 interface and restart flanneld or simply reboot.

On the master (atomic01) download the sample config file

    cd /usr/local/src/
    curl -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/flannel-config.json
    cat flannel-config.json
    {
    "Network": "10.20.0.0/16",
    "SubnetLen": 24,
    "Backend": {
    "Type": "vxlan",
    "VNI": 1
    }
    }

Now load it to the master etcd service.

    cd /usr/local/src/
    curl -L http://atomic01.example.net:4001/v2/keys/coreos.com/network/config -XPUT --data-urlencode value@flannel-config.json

Check to see if it was uploaded successfully.

    curl -s -L http://atomic01.example.net:4001/v2/keys/coreos.com/network/config | python -mjson.tool

On the master and both nodes, edit /etc/sysconfig/flanneld to insert the name or IP address of the system containing the etcd service (master) and set the network interface.The file should look something like this on all nodes.

    egrep -v '^#|^$' /etc/sysconfig/flanneld
    FLANNEL_ETCD="http://atomic01.example.net:4001"
    FLANNEL_ETCD_KEY="/coreos.com/network"
    FLANNEL_OPTIONS="eth0"

Start flanneld on master and nodes: First on the master, then on the two nodes:

    for cmd in restart enable status; do systemctl ${cmd} flanneld ; done

For the docker systemd service to pick up the flannel changes, and to make sure the network interfaces all come up properly, reboot all nodes

    systemctl reboot

# Launching Services, Replication Controllers, and Container Pods with Kubernetes


With the Kubernetes cluster in place, you can now create the yaml files needed to set up Kubernetes services, define replication controllers and launch pods of containers. Using the two containers described earlier in this document (Web and DB), you will create the following types of Kubernetes objects:

*     Services: Creating a Kubernetes service lets you assign a specific IP address and port number to a label. Because pods and IP addresses can come and go with Kubernetes, that label can be used within a pod to find the location of the services it needs.
*     Replication Controllers: By defining a replication controller, you can set not only which pods to start, but how many replicas of each pod should start. If a pod stops, the replication controller starts another to replace it.
*     Pods: A pod loads one or more containers, along with options associated with running the containers.

*  
In this example, we will use the YAML files found on my github repo: https://github.com/christianh814/dock-kube-example.git to set up the services. 

Download and inspect the db services YAML file on the master

    curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/db-service.yaml
    cat db-service.yaml
    id: "db-service"
    kind: "Service"
    apiVersion: "v1beta1"
    port: 3306
    portalIP: "10.254.100.1"
    selector:
      name: "db"
    labels:
      name: "db"

The portalIP for the db-service must be in the range set in the /etc/kubernetes/apiserver file. This service in not directly accessed from the outside world. The webserver container will access it. The selector and labels name is set to db. To start that service, type the following on the master:

    kubectl create -f db-service.yaml

Download and inspect the db replication controller YAML file on the master. Note that I put the "replica" count to 1.

    curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/db-rc.yaml
    cat db-rc.yaml
    id: "db-controller"
    kind: "ReplicationController"
    apiVersion: "v1beta1"
    desiredState:
      replicas: 1
      replicaSelector:
        selectorname: "db"
      podTemplate:
        desiredState:
          manifest:
            version: "v1beta1"
            id: "db-controller"
            containers:
              - name: "db"
                image: "hub.example.net:5000/dbserver"
                ports:
                  - containerPort: 3306
        labels:
          name: "db"
          selectorname: "db"
    labels:
      name: "db-controller"

To start the replication controller run the kubectl command

    kubectl create  -f db-rc.yaml

Download and inspect the webserver service YAML file

    curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/webserver-service.yaml
    cat webserver-service.yaml
    kind: "Service"
    id: "webserver-service"
    apiVersion: "v1beta1"
    port: 80
    publicIPs:
      - 172.16.1.122
    portalIP: "10.254.100.50"
    selector:
      name: "webserver"
    labels:
      name: "webserver"

The publicIPs value should align with an IP address associated with an external network interface on one of the nodes (on atomic02, the IP address for the eth0 interface is 172.16.1.122). That will make the service available to the outside world. The selector and labels name is set to webserver. To start that service, type the following on the master:

    kubectl create -f webserver-service.yam

With the with the apiserver and the two services we just added running, run the following command to see those services:

    kubectl get services

Next download and inspect the webserver replication controller

    curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/webserver-rc.yaml
    cat webserver-rc.yaml
    id: "webserver-controller"
    kind: "ReplicationController"
    apiVersion: "v1beta1"
    desiredState:
      replicas: 2
      replicaSelector:
        selectorname: "webserver"
      podTemplate:
        desiredState:
          manifest:
            version: "v1beta1"
            id: "webserver-controller"
            containers:
              - name: "apache-frontend"
                image: "hub.example.net:5000/webserver"
                ports:
                  - containerPort: 80
        labels:
          name: "webserver"
          selectorname: "webserver"
          uses: "db"
    labels:
      name: "webserver-controller"

Based on the  yaml file, the replication controller will try to keep two pods labeled "webserver" running at all times. The pod definition is inside the replications controller yaml file, so no separate pod yaml file is needed. Any running pod with the webserver id will be taken by the replication controller as fulfilling the requiremen. Apply the configuration:

    kubectl create -f webserver-rc.yaml 

Inspect your configurations with the following commands:

    kubectl get pods
    kubectl get replicationControllers
    kubectl get services
    kubectl get nodes

Once you see that it's all clear; you can test your application

    curl http://172.16.1.122
    curl http://172.16.1.122/cgi-bin/action
