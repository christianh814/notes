# Intro

Using Kubernetes you can run and manage docker containers from one system (a Kubernetes master) and deploy them to run on other systems (Kubernetes nodes). This procedure will help you:

* Set up three Red Hat Enterprise Linux 7  systems to use as a Kubernetes master and two Kubernetes nodes (also called minions).
* Configure networking between containers using Flannel.
* Create data files (in yaml format) that define services, pods and replication controllers.
* Configure a fourth RHEL 7 system as a docker registry server.

The services you start during this procedure include etcd, cadvisor, kube-apiserver, kubecontroller-manager, kube-proxy, kube-scheduler, and kubelet. To launch containers from yaml files and otherwise manage your Kubernetes environment, you run the kubectl command.

The result of this procedure is a Kubernetes environment that you can use as a foundation for deploying Docker containers across multiple host systems.

Assumptions:

We will be using 4 RHEL 7 servers.

* You are familiar with RHEL and RHEL Atomic
* You have the proper repos setup
* You are familiar with Docker
* You have a basic level of understanding of Kubernetes

# Docker


Docker is an open-source project that automates the deployment of applications inside software containers, by providing an additional layer of abstraction and automation of operating-system-level virtualization on Linux. Docker uses resource isolation features of the Linux kernel such as cgroups and kernel namespaces to allow independent "containers" to run within a single Linux instance, avoiding the overhead of starting virtual machines.

# Kubernetes


Kubernetes is a powerful system, developed by Google, for managing containerized applications in a clustered environment. It aims to provide better ways of managing related, distributed components across varied infrastructure. Kubernetes groups containers into a logical unit called "pods". A pod is a 1 or more containers that share common resources (such a network and storage).

# Installation

__Master__

I have the following installed as a RHEL 7 Host. This host will be the kuberentes master. I have the following IPaddr/DNS-names assigned

* rmaster.example.com - 172.16.1.241

Install kubernetes, etcd, and docker on the master

```
yum -y install kubernetes docker etcd
```

__Nodes__

I have the following installed as a RHEL 7 Host. These hosts will be running the docker containers. I have the following IPaddr/DNS-names assigned

* rnode1.example.com - 172.16.1.242
* rnode2.example.com - 172.16.1.243

Install kubernetes and docker on the nodes

```
yum -y install kubernetes docker
```

__Registry__

One host will be installed with RHEL 7 Server that will be running a docker registry.

*     rhub.example.com - 172.16.1.240

To install the registry; install the following packages

yum -y install docker-registry docker
systemctl enable docker-registry docker


__All Hosts__

Also make sure firewalld is disabled on all hosts

systemctl stop firewalld
systemctl disable firewalld

# Docker Registry Configuration


On the "hub" server; you will need to modify the systemd startup script. The best way to do this is copy it over to /etc/systemd/system directory (since files here "override" the files in /usr/lib/systemd/system)

cp /usr/lib/systemd/system/docker-registry.service /etc/systemd/system/

Edit the /etc/systemd/system/docker-registry.service file and add `--certfile` and `--keyfile` to the start up script.
```
root@host# cat /etc/systemd/system/docker-registry.service
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
```

Now generate the keyfile and the certificate

```
openssl genrsa -out /etc/pki/tls/private/dr.key 1024
openssl req -new -key /etc/pki/tls/private/dr.key -x509 -out /etc/pki/tls/certs/dr.crt
```

Copy the certificate locally and all nodes.

```
mkdir /etc/docker/certs.d/rhub.example.com:5000
for i in rmaster.example.com rnode{1..2}.example.com ; do ssh $i mkdir /etc/docker/certs.d/rhub.example.com:5000; done
cp /etc/pki/tls/certs/dr.crt /etc/docker/certs.d/rhub.example.com:5000/dr.crt
for i in rmaster.example.com rnode{1..2}.example.com ; do scp /etc/pki/tls/certs/dr.crt root@$i:/etc/docker/certs.d/rhub.example.com\:5000/ ; done
```

Now you can start and enable the service

```
systemctl enable docker-registry
systemctl start docker-registry
systemctl status docker-registry
```

For those that left the firewall running; add port 5000

```
firewall-cmd --zone=public --add-port=5000/tcp
firewall-cmd --permanent --zone=public --add-port=5000/tcp
```

At this point your docker-registry should be up and running. Next; you'll need to build the containers found on my github account.

```
git clone https://github.com/christianh814/dock-kube-example.git
```

In this repo, there are 2 directories webserver and dbserver. You can inspect the contents if you wish. Build the docker containers from the dockerfiles. Run docker images to verify they have been built.

```
cd dock-kube-example/webserver
docker build -t webserver .
cd dock-kube-example/
docker build -t dbserver .
docker images
```

In order to upload images to your registry you first must "tag" them referencing your registry URL:PORT. Then you can "push" these images to your registry.

```
docker tag dbserver rhub.example.com:5000/dbserver
docker tag webserver rhub.example.com:5000/webserver
docker push rhub.example.com:5000/webserver
docker push rhub.example.com:5000/dbserver
```

Lastly, pull the images
```
docker pull rhub.example.com:5000/webserver
docker pull rhub.example.com:5000/dbserver
```

# Configuring Kubernetes

__Kubernetes Master__

First, the etcd service needs to be configured in `/etc/etcd/etcd.conf` to listen on all interfaces to ports 2380 and 7001 and port 4001. The resulting uncommented lines in this file should appear as follows:

```
[root@rmaster src]# egrep -v '^#|^$' /etc/etcd/etcd.conf 
ETCD_NAME=default
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_PEER_URLS="http://0.0.0.0:2380,http://0.0.0.0:7001"
ETCD_LISTEN_CLIENT_URLS="http://0.0.0.0:4001"
```

Next Edit the `/etc/kubernetes/controller-manager` file.  Comment out the KUBELET_ADDRESSES line. While this variable can be used to identify each node to the master in the current release (the default identifies just the master as a node for all-in-one setups), this will not work in future versions. The better way is to create a json file and load the minion information using kubectl (as done a few steps ahead). So, comment out that variable in this file as follows:
```
[root@rmaster src]# grep KUBELET_ADDRESSES /etc/kubernetes/controller-manager 
#KUBELET_ADDRESSES="--machines=127.0.0.1"
```

In the `/etc/kubernetes/config` file; Change the KUBE_MASTER line to identify the location of your master server (it points to 127.0.0.1, by default). Leave other settings as they are. The file should look something like this

```
[root@rmaster src]# egrep -v '^#|^$' /etc/kubernetes/config 
KUBE_LOGTOSTDERR="--logtostderr=true"
KUBE_LOG_LEVEL="--v=0"
KUBE_ALLOW_PRIV="--allow_privileged=false"
KUBE_MASTER="--master=http://rmaster.example.com:8080"
```

In the `/etc/kubernetes/apiserver` file; add a new KUBE_ETCD_SERVERS line, then review and change other lines in the apiserver configuration file. Change KUBE_API_ADDRESS to listen on all network addresses(0.0.0.0), instead of just localhost. Set an address range on the KUBE_SERVICE_ADDRESS that Kubernetes can use to assign to services (see a description of this address below). Here are examples:. The file should look something like this.

```
[root@rmaster src]# egrep -v '^#|^$' /etc/kubernetes/apiserver
KUBE_API_ADDRESS="--address=0.0.0.0"
KUBE_ETCD_SERVERS="--etcd_servers=http://rmaster.example.com:4001"
KUBE_SERVICE_ADDRESSES="--portal_net=10.254.0.0/16"
KUBE_ADMISSION_CONTROL="--admission_control=NamespaceAutoProvision,LimitRanger,ResourceQuota"
KUBE_API_ARGS=""
```

Keep in mind what you are setting for the KUBE_SERVICE_ADDRESSES line...

* The address range is used by Kubernetes to assign to Kubernetes services.
* In the example just shown, the address range of 10.254.0.0/16 consumes a set of 10.254 subnets that can be assigned by Kubernetes as needed. For example, 10.254.1.X, 10.254.2.X and so on.
* Make sure this address range isn't used anywhere else in your environment.
* Each address range that is assigned is only used within a node and is not routable outside of that node.
* This address range must be different than the range used by flannel. (Flannel address ranges are assigned to pods.)

You need to start several services associated with a Kubernetes master. From the master, run the following command to start and enable Kubernetes systemd services on the master.

```
[root@rmaster src]# for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler docker; do 
systemctl restart $SERVICES
systemctl enable $SERVICES
systemctl status $SERVICES 
done
```

Lastly, define your nodes via the following JSON files

```
[root@rmaster src]# cat node1.json 
{
    "apiVersion": "v1beta3",
    "kind": "Node",
    "metadata": {
        "name": "rnode1.example.com",
        "labels": {"name" : "rnode1.example.com"}
    },
    "spec": {
        "externalID": "rnode1-external"
    }
}

[root@rmaster src]# cat node2.json 
{
    "apiVersion": "v1beta3",
    "kind": "Node",
    "metadata": {
        "name": "rnode2.example.com",
        "labels": {"name" : "rnode2.example.com"}
    },
    "spec": {
        "externalID": "rnode2-external"
     }
}
```

Load the configuration  using the **kubectl** command
```
[root@rmaster src]# kubectl create -f node1.json
[root@rmaster src]# kubectl create -f node2.json
[root@rmaster src]# kubectl get nodes
NAME                 LABELS        STATUS
rnode1.example.com   Schedulable   name=rnode1.example.com   Ready
rnode2.example.com   Schedulable   name=rnode2.example.com   Ready
```

__Kubernetes Nodes (minions)__

On each of the two Kubernetes nodes (rnode1.example.com and rnode2.example.com in this example), configure them to communicate with the master.

Due to a current bug in Kubernetes; you must create an auth file that is empty. The file is under `/var/lib/kubelet/auth` and it's a JSON file.

```
[root@rnodeX ~]# echo "{}" > /var/lib/kubelet/auth
```

In the `/etc/kubernetes/config` file; edit the KUBE_MASTER line to identify the location of your master (it is 127.0.0.1, by default). Leave other settings as they are.

```
[root@rnodeX ~]# egrep -v '^$|^#' /etc/kubernetes/config
KUBE_LOGTOSTDERR="--logtostderr=true"
KUBE_LOG_LEVEL="--v=0"
KUBE_ALLOW_PRIV="--allow_privileged=false"
KUBE_MASTER="--master=http://rmaster.example.com:8080"
```

In the `/etc/kubernetes/kubelet`  file on each node, modify KUBELET_ADDRESS (0.0.0.0 to listen on network interfaces), KUBELET_HOSTNAME (replace hostname_override with the hostname or IP address of the local system: node1.example.com or node2.example.com), and KUBELET_API_SERVER (set --api_servers=<nowiki>http://rmaster.example.com:8080 </nowiki> or other location of the master), the file should look something like this:

```
[root@rnodeX ~]# egrep -v '^#|^$' /etc/kubernetes/kubelet
KUBELET_ADDRESS="--address=127.0.0.1"
KUBELET_HOSTNAME="--hostname_override=rnodeX.example.com"
KUBELET_API_SERVER="--api_servers=http://rmaster.example.com:8080"
KUBELET_ARGS=""
```

In the `/etc/kubernetes/proxy` file, No settings are required in this file. If you have set KUBE_PROXY_ARGS, you can comment it out

```
[root@rnodeX ~]# grep KUBE_PROXY_ARGS /etc/kubernetes/proxy 
# KUBE_PROXY_ARGS=""
```

On each node, you need to start several services associated with a Kubernetes node
```
[root@rnodeX ~]# for SERVICES in docker kube-proxy.service kubelet.service; do 
    systemctl restart $SERVICES
    systemctl enable $SERVICES
    systemctl status $SERVICES 
done
```

From any node test to see if etcd is running.
```
[root@rnodeX ~]# curl -s -L http://rmaster.example.com:4001/version
```

# Setting up Flannel networking

The flannel package contains features that allow you to configure networking between the master and nodes in your Kubernetes cluster. You configure the flanneld service by creating and uploading a json configuration file with your network configuration to your etcd server (on the master). You then configure the flanneld systemd service on the master and each node to point to that etcd server and start the flanneld service.

Because the docker0 interface is probably already in place when you run this procedure, the IP address range assigned by flanneld to docker0 will not immediately take effect. You can either manually stop the docker0 interface and restart flanneld or simply reboot.

On ALL hosts install flannel
```
root@host# yum -y install flannel
```

On the master (rmaster.example.com) download the sample config file
```
[root@rmaster src]# cd /usr/local/src/
[root@rmaster src]# curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/flannel-config.json
[root@rmaster src]# cat flannel-config.json
{
  "Network": "10.20.0.0/16",
  "SubnetLen": 24,
  "Backend": {
    "Type": "vxlan",
    "VNI": 1
     }
}
```

Now load it to the master etcd service.
```
[root@rmaster src]# cd /usr/local/src/
[root@rmaster src]# etcdctl set coreos.com/network/config < flannel-config.json
[root@rmaster src]# etcdctl get coreos.com/network/config
{
  "Network": "10.20.0.0/16",
  "SubnetLen": 24,
  "Backend": {
    "Type": "vxlan",
    "VNI": 1
     }
}
```

On the master and both nodes, edit `/etc/sysconfig/flanneld` to insert the name or IP address of the system containing the etcd service (master) and set the network interface.The file should look something like this on all nodes.

```
[root@rnode1 ~]# egrep -v '^#|^$' /etc/sysconfig/flanneld
FLANNEL_ETCD="http://rmaster.example.com:4001"
FLANNEL_ETCD_KEY="/coreos.com/network"
FLANNEL_OPTIONS="eth0"
```

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

```
root@host# curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/db-service.yaml
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
```

The portalIP for the db-service must be in the range set in the /etc/kubernetes/apiserver file. This service in not directly accessed from the outside world. The webserver container will access it. The selector and labels name is set to db. To start that service, type the following on the master:

    kubectl create -f db-service.yaml

Download and inspect the db replication controller YAML file on the master. Note that I put the "replica" count to 1.

```
root@host# curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/db-rc.yaml
root@host# cat db-rc.yaml
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
            image: "rhub.example.com:5000/dbserver"
            ports:
              - containerPort: 3306
    labels:
      name: "db"
      selectorname: "db"
labels:
  name: "db-controller"
```

To start the replication controller run the kubectl command

```
kubectl create  -f db-rc.yaml
```

Download and inspect the webserver service YAML file

```
root@host# curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/webserver-service.yaml
root@host# cat webserver-service.yaml
kind: "Service"
id: "webserver-service"
apiVersion: "v1beta1"
port: 80
publicIPs:
  - 172.16.1.243
portalIP: "10.254.100.50"
selector:
  name: "webserver"
labels:
  name: "webserver"
```

The publicIPs value should align with an IP address associated with an external network interface on one of the nodes (on rnode2.example.com, the IP address for the eth0 interface is 172.16.1.243). That will make the service available to the outside world. The selector and labels name is set to webserver. To start that service, type the following on the master:

```
kubectl create -f webserver-service.yam
```

With the with the apiserver and the two services we just added running, run the following command to see those services:

```
kubectl get services
```

Next download and inspect the webserver replication controller

```
root@host# curl -s -O https://raw.githubusercontent.com/christianh814/dock-kube-example/master/yaml-json/webserver-rc.yaml
root@host# cat webserver-rc.yaml
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
            image: "rhub.example.com:5000/webserver"
            ports:
              - containerPort: 80
    labels:
      name: "webserver"
      selectorname: "webserver"
      uses: "db"
labels:
  name: "webserver-controller"
```

Based on the  yaml file, the replication controller will try to keep two pods labeled "webserver" running at all times. The pod definition is inside the replications controller yaml file, so no separate pod yaml file is needed. Any running pod with the webserver id will be taken by the replication controller as fulfilling the requiremen. Apply the configuration:

    kubectl create -f webserver-rc.yaml 

Inspect your configurations with the following commands:

```
kubectl get pods
kubectl get replicationControllers
kubectl get services
kubectl get nodes
```

Once you see that it's all clear; you can test your application

```
curl http://172.16.1.243
curl http://172.16.1.243/cgi-bin/action
```
