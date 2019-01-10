# Kubernetes Installation

This guide will show you how to "easily" (as easy as possible) install a Multimaster system. You can also modify this to have a similer one contoller/one node setup (steps are almost exactly the same).

Also, this is here mainly for my notes; so I may be brief in some sections. I TRY to keep this up to date but you're better off looking at the [official documentation](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/)

This is broken off into sections

* [Prerequisites And Assumptions](#prerequisites-and-assumptions)
* [Infrastructure](#infrastructure)
* [Installing Kubeadm](#installing-kubeadm)
* [Initialize the Control Plane](#initialize-the-control-plane)
* [Bootstrap Remaining Controllers](#bootstrap-remaining-controllers)
* [Bootstrap Worker Nodes](#bootstrap-worker-nodes)
* [Smoke Test](#smoke-test)
* [Ingress](#ingress)
* [Misc](#miscellaneous-notes)

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

* Latest CentOS 7 (FULLY updated with `yum -y update` and `systemctl reboot`-ed)
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

# Installing Kubeadm

First, on ALL servers (from here on out when I say "ALL servers" I mean the 3 controllers and 3 workers), install the runtime. This can be `cri-o`, `containerd`, `docker`, or `rkt`. For ease of use; I installed `docker`

```
ansible all -m shell -a "yum -y install docker"
ansible all -m shell -a "systemctl enable --now docker"
```

Next, on ALL servers, install these packages:

* `kubeadm`: the command to bootstrap the cluster.
* `kubelet`: the component that runs on all of the machines in your cluster and does things like starting pods and containers.
* `kubectl`: the command line util to talk to your cluster.

You should also disable the firewall and SELinux at this point as well

```
ansible all -m shell -a "sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config"
ansible all -m shell -a "setenforce 0"
ansible all -m shell -a "systemctl stop firewalld"
ansible all -m shell -a "systemctl disable firewalld"
ansible all -m copy -a "src=files/k8s.repo dest=/etc/yum.repos.d/kubernetes.repo"
ansible all -m shell -a "yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes"
ansible all -m shell -a "systemctl enable --now kubelet"
ansible all -m copy -a "src=files/sysctl_k8s.conf dest=/etc/sysctl.d/k8s.conf"
ansible all -m shell -a "sysctl --system"
```

(Note: A copy of [k8s.repo](k8s-resources/k8s.repo) and [sysctl_k8s.conf](k8s-resources/sysctl_k8s.conf) are found in this repo)

If you run `systemctl status kubelet` you will see it error out becuase it's trying to connect to a control plane that's not there. This is a normal error.

At this point, "pre-pull" the required images on ONLY the controllers to run the control plane.

```
ansible controllers -m shell -a "kubeadm config images pull"
```

# Initialize the Control Plane

Pick one of the controllers (doesn't matter which one) and initialize the control plane with `kubeadm`. I used a config file since I was going through a loadbalancer. (Also note, if you're using Calico (like I was) you have to specify a `podSubnet` that's something other than `192.168.0.0/16`)

I did the following on my first controller (in my case my LB was `192.168.1.97`)

```
LBIP=192.168.1.97
cat <<EOF > kubeadm-config.yaml
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: stable
apiServer:
  certSANs:
  - ${LBIP}
controlPlaneEndpoint: ${LBIP}:6443
networking:
  podSubnet: 10.254.0.0/16
  serviceSubnet: 172.30.0.0/16
EOF
```

Use this `kubeadm-config.yaml` in install the first controller. On the controller you created the file, run the `kubeadm init` command...


```
kubeadm init --config=kubeadm-config.yaml
```

At this point, if you reload your LB's status page, you should see the first controller going green.

When the bootstrapping finishes; you should see a message like the following. SAVE THIS MESSAGE! You'll need this to join the other 2 controllers. 

```
kubeadm join <lb ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

Next, install a CNI compliant SDN. I used Calico since it was the easiest. First wget them

```
curl -O https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml
curl -O https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml
```

Edit the `calico.yaml` file to use the `podSubnet` you defined in the `kubeadm-config.yaml` file

```
sed -i 's/192\.168/10\.254/g' calico.yaml
```

Now load these into kubernetes

```
export KUBECONFIG=/etc/kubernetes/admin.conf
kubectl apply -f rbac-kdd.yaml
kubectl apply -f calico.yaml
```

Wait a few minutes and you should see CoreDNS and the Calico pods up

```
kubectl get pods --all-namespaces
NAMESPACE     NAME                                             READY   STATUS    RESTARTS   AGE
kube-system   calico-node-g6s98                                2/2     Running   0          33s
kube-system   coredns-86c58d9df4-gq6n8                         1/1     Running   0          15m
kube-system   coredns-86c58d9df4-pl6m6                         1/1     Running   0          15m
kube-system   etcd-dhcp-host-98.cloud.chx                      1/1     Running   0          14m
kube-system   kube-apiserver-dhcp-host-98.cloud.chx            1/1     Running   0          14m
kube-system   kube-controller-manager-dhcp-host-98.cloud.chx   1/1     Running   0          14m
kube-system   kube-proxy-6dn5z                                 1/1     Running   0          15m
kube-system   kube-scheduler-dhcp-host-98.cloud.chx            1/1     Running   0          13m
```
# Bootstrap Remaining Controllers

After the SDN is installed; this is the point you can boostrap the remaining controllers. Now copy the certificate files from the first control plane node to the rest of them

```
for host in 192.168.1.99 192.168.1.5 
do
    scp /etc/kubernetes/pki/ca.crt root@$host:
    scp /etc/kubernetes/pki/ca.key root@$host:
    scp /etc/kubernetes/pki/sa.key root@$host:
    scp /etc/kubernetes/pki/sa.pub root@$host:
    scp /etc/kubernetes/pki/front-proxy-ca.crt root@$host:
    scp /etc/kubernetes/pki/front-proxy-ca.key root@$host:
    scp /etc/kubernetes/pki/etcd/ca.crt root@$host:etcd-ca.crt
    scp /etc/kubernetes/pki/etcd/ca.key root@$host:etcd-ca.key
    scp /etc/kubernetes/admin.conf root@$host:
done
```

Now on the 2 remaining hosts run the follwing commands

```
mkdir -p /etc/kubernetes/pki/etcd
mv /root/ca.crt /etc/kubernetes/pki/
mv /root/ca.key /etc/kubernetes/pki/
mv /root/sa.pub /etc/kubernetes/pki/
mv /root/sa.key /etc/kubernetes/pki/
mv /root/front-proxy-ca.crt /etc/kubernetes/pki/
mv /root/front-proxy-ca.key /etc/kubernetes/pki/
mv /root/etcd-ca.crt /etc/kubernetes/pki/etcd/ca.crt
mv /root/etcd-ca.key /etc/kubernetes/pki/etcd/ca.key
mv /root/admin.conf /etc/kubernetes/admin.conf
```

Now finally join the remaining 2 controlplanes to the cluster

```
kubeadm join <lb ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash> --experimental-control-plane
```

Once that finishes; you can see them listed with `kubectl`

```
kubectl get nodes
NAME                     STATUS   ROLES    AGE     VERSION
dhcp-host-5.cloud.chx    Ready    master   3m34s   v1.13.1
dhcp-host-98.cloud.chx   Ready    master   34m     v1.13.1
dhcp-host-99.cloud.chx   Ready    master   4m52s   v1.13.1
```

# Bootstrap Worker Nodes

This part is the easiest; you just run the `kubeadm join ...` command WITHOUT the `--experimental-control-plane` flag


On the 3 workers run the following
```
kubeadm join <lb ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>
```

Once they have joined; get it with `kubectl get nodes`


```
[root@jumpbox k8s-with-kubeadm]# kubectl get nodes
NAME                     STATUS   ROLES    AGE     VERSION
dhcp-host-5.cloud.chx    Ready    master   9m32s   v1.13.1
dhcp-host-6.cloud.chx    Ready    <none>   53s     v1.13.1
dhcp-host-7.cloud.chx    Ready    <none>   53s     v1.13.1
dhcp-host-8.cloud.chx    Ready    <none>   54s     v1.13.1
dhcp-host-98.cloud.chx   Ready    master   40m     v1.13.1
dhcp-host-99.cloud.chx   Ready    master   10m     v1.13.1
```

# Smoke Test

Test your cluster by deploying a pod with a `NodePort` definition.

First create a namespace

```
kubectl create namespace test
```

Next, deploy a pod and a service to this namespace

```
cat <<EOF | kubectl apply -n test -f -
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: welcome-php
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: welcome-php
    spec:
      containers:
        - image: "quay.io/redhatworkshops/welcome-php:latest"
          imagePullPolicy: Always
          name: welcome-php
          ports:
            - containerPort: 8080
EOF
cat <<EOF | kubectl apply -n test -f -
apiVersion: v1
kind: Service
metadata:
  name: welcome-php
  labels:
    app: welcome-php
spec:
  type: NodePort
  ports:
    - port: 8080
      nodePort: 30000
      name: http
  selector:
    app: welcome-php
EOF
```

Check to see if your pods and svc is up

```
kubectl get pods -n test
NAME                           READY   STATUS    RESTARTS   AGE
welcome-php-77f6d8845b-mnv4r   1/1     Running   0          76s

kubectl get svc -n test
NAME          TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
welcome-php   NodePort   172.30.67.197   <none>        8080:30000/TCP   81s
```

Curling the svc address should get you that 200

```
curl -sI 172.30.67.197:8080
HTTP/1.1 200 OK
Date: Thu, 10 Jan 2019 01:36:08 GMT
Server: Apache/2.4.27 (Red Hat) OpenSSL/1.0.1e-fips
Content-Type: text/html; charset=UTF-8
```

If you visit the IP of ANY node in the cluster on port 30000, you should see the app come up.

Congrats! You have an HA k8s cluster up and running!


# Ingress

Right now, I have a cluster up and running and can deploy apps and use `nodePort` to expose these apps. In theory you can put a LB in front of it and you can serve your app.

You can have k8s "manage" your LB with an `ingress controller`. There are many out there but it seems like nginx is the most popular one. Let's get one up and running!

We will be using the `test` namespace we created during the [smoke test(#smoke-test). Clean it up if you wish

```
kubectl -n test delete all --all
```

First let's create some app yamls
```
cat > app-deployment.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: app1
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: app1
    spec:
      containers:
      - name: app1
        image: dockersamples/static-site
        env:
        - name: AUTHOR
          value: app1
        ports:
        - containerPort: 80
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: app2
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: app2
    spec:
      containers:
      - name: app2
        image: dockersamples/static-site
        env:
        - name: AUTHOR
          value: app2
        ports:
        - containerPort: 80
EOF
```

Next, some service yamls

```
cat > app-service.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: appsvc1
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: app1
---
apiVersion: v1
kind: Service
metadata:
  name: appsvc2
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: app2
EOF
```

Create these resources in the `test` namespace

```
kubectl create -n test -f app-deployment.yaml -f app-service.yaml
```

Next, create the nginx ingress controller. All resources for nginx has to be in it's own namespace; so let's create one

```
kubectl create namespace ingress
```

The first step is to create a default backend endpoint yaml. Default endpoint redirects all requests which are not defined by Ingress rules. Meaning when someone hits an endpoint that has not ingress rule.

```
cat > default-backend-deployment.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: default-backend
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: default-backend
    spec:
      terminationGracePeriodSeconds: 60
      containers:
      - name: default-backend
        image: gcr.io/google_containers/defaultbackend:1.0
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 30
          timeoutSeconds: 5
        ports:
        - containerPort: 8080
        resources:
          limits:
            cpu: 10m
            memory: 20Mi
          requests:
            cpu: 10m
            memory: 20Mi
EOF
```

Also, create a default backend service yaml

```
cat > default-backend-service.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: default-backend
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: default-backend
EOF
```

Now, create those resources in ingress namespace

```
kubectl create -n ingress -f default-backend-deployment.yaml -f default-backend-service.yaml
```

It's nice to have a status page for nginx...so let's create one

```
cat > nginx-ingress-controller-config-map.yaml <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-ingress-controller-conf
  labels:
    app: nginx-ingress-lb
data:
  enable-vts-status: 'true'
EOF
```

Now add it to k8s (in the ingress namespace)

```
kubectl create -n ingress -f nginx-ingress-controller-config-map.yaml
```

Now that you set up the side dishes...time for the meat and potatoes...the actual nginx ingress controller deplyment yaml :)

```
cat > nginx-ingress-controller-deployment.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-ingress-controller
spec:
  replicas: 1
  revisionHistoryLimit: 3
  template:
    metadata:
      labels:
        app: nginx-ingress-lb
    spec:
      terminationGracePeriodSeconds: 60
      serviceAccount: nginx
      containers:
        - name: nginx-ingress-controller
          image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.9.0
          imagePullPolicy: Always
          readinessProbe:
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
          livenessProbe:
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            initialDelaySeconds: 10
            timeoutSeconds: 5
          args:
            - /nginx-ingress-controller
            - --default-backend-service=\$(POD_NAMESPACE)/default-backend
            - --configmap=\$(POD_NAMESPACE)/nginx-ingress-controller-conf
            - --v=2
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          ports:
            - containerPort: 80
            - containerPort: 18080
EOF
```

Before we create this object, let's do some rbac stuff...

```
cat > nginx-ingress-controller-roles.yaml <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: nginx-role
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - endpoints
  - nodes
  - pods
  - secrets
  verbs:
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - get
  - list
  - update
  - watch
- apiGroups:
  - extensions
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
- apiGroups:
  - extensions
  resources:
  - ingresses/status
  verbs:
  - update
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: nginx-role
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-role
subjects:
- kind: ServiceAccount
  name: nginx
  namespace: ingress
EOF
```

Create the rbac object and the deployment object (in that order for cleanliness)

```
kubectl create -n ingress -f nginx-ingress-controller-roles.yaml -f nginx-ingress-controller-deployment.yaml
```

The controller is now set up but, check your work...it should look something like this

```
# kubectl get all -n ingress
NAME                                            READY   STATUS    RESTARTS   AGE
pod/default-backend-6f6db5f6cd-246xn            1/1     Running   0          8m13s
pod/default-backend-6f6db5f6cd-nmng9            1/1     Running   0          8m13s
pod/nginx-ingress-controller-854cff467c-h4nk2   1/1     Running   0          79s

NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
service/default-backend   ClusterIP   172.30.133.116   <none>        80/TCP    8m13s

NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/default-backend            2/2     2            2           8m13s
deployment.apps/nginx-ingress-controller   1/1     1            1           79s

NAME                                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/default-backend-6f6db5f6cd            2         2         2       8m13s
replicaset.apps/nginx-ingress-controller-854cff467c   1         1         1       79s
```

Next, define Ingress rules for load balancer status page, and your sample apps webpage

```
cat > nginx-ingress.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: nginx-ingress
spec:
  rules:
  - host: test.192.168.1.99.nip.io
    http:
      paths:
      - backend:
          serviceName: nginx-ingress
          servicePort: 18080
        path: /nginx_status
EOF
cat > app-ingress.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
  name: app-ingress
spec:
  rules:
  - host: test.192.168.1.99.nip.io
    http:
      paths:
      - backend:
          serviceName: appsvc1
          servicePort: 80
        path: /app1
      - backend:
          serviceName: appsvc2
          servicePort: 80
        path: /app2
EOF
```

Create both ingress rules

```
kubectl -n ingress create -f nginx-ingress.yaml
kubectl -n test create -f app-ingress.yaml
```

The last step is to expose `nginx-ingress-lb` deployment for external access. We will expose it with `NodePort`, but we could also use `ExternalIPs` here:

```
cat > nginx-ingress-controller-service.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
spec:
  type: NodePort
  ports:
    - port: 80
      nodePort: 30000
      name: http
    - port: 18080
      nodePort: 32000
      name: http-mgmt
  selector:
    app: nginx-ingress-lb
EOF
```

Apply this

```
kubectl create -n ingress -f nginx-ingress-controller-service.yaml
```

Verify these by accessing the following endpoinds

```
http://test.192.168.1.99.nip.io:30000/app1
http://test.192.168.1.99.nip.io:30000/app2
http://test.192.168.1.99.nip.io:32000/nginx_status
```

Any other endpoint results in the default 404 page

Note, once in place; all you need to expose your app is...

* An ingress object (i.e. the `app-ingress.yaml` file)
* A svc (i.e. `nginx-ingress-controller-service.yaml`)
  * If you've already exposed the port you wanted, you don't need this!
* The file will change (i.e. the `app-ingress.yaml` file), depending on the app and what path you want to route to what endpoint (see [using externalIPs](#using-externalips))

## Using ExternalIPs

To use an `externalIPs` object, the IP must already be assinged to the node (it doesn't create it for you). Since I'm testing this on VMs I'll just use the IP that's already on the node.

First get the IP from the ingress controller

```
[root@jumpbox test-yamls]# kubectl get pods -n ingress -o wide -l app=nginx-ingress-lb
NAME                                        READY   STATUS    RESTARTS   AGE   IP           NODE                    NOMINATED NODE   READINESS GATES
nginx-ingress-controller-854cff467c-h4nk2   1/1     Running   0          41m   10.254.3.5   dhcp-host-8.cloud.chx   <none>           <none>
[root@jumpbox test-yamls]# dig dhcp-host-8.cloud.chx +short
192.168.1.8
```

Looks like my ingress controller is running on host `dhcp-host-8.cloud.chx` with ip of `192.168.1.8`

Now I'm going to create my app's ingress object

```
cat > app-ingress-exip.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
  name: app-ingress-exip
spec:
  rules:
  - host: app.192.168.1.8.nip.io
    http:
      paths:
      - backend:
          serviceName: appsvc1
          servicePort: 80
        path: /
EOF
```

Create this in the `test` namespace

```
kubectl create -n test -f app-ingress-exip.yaml
```

Now create the ingress controller on port 80. Note the use of `externalIPs`

```
cat > ginx-ingress-controller-service-exip.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress-exip
spec:
  externalIPs:
  - 192.168.1.8
  ports:
    - port: 80
      name: http
  selector:
    app: nginx-ingress-lb
EOF
```

Now create this service in the `ingress` namespace

```
kubectl create -n ingress -f ginx-ingress-controller-service-exip.yaml
```

Test it out with curl (example below) or on your browser

```
# curl -s app.192.168.1.8.nip.io  | grep app1
<h1 id="toc_0">Hello app1!</h1>
```

Success!

**NOTE** Please note that you've effectively made `dhcp-host-8.cloud.chx` your "infrastrucure" node (if you're coming from "openshift" land)

# Miscellaneous Notes

Notes in no paticular order

## Get join token

First create the token

```
kubeadm token create
```

Next, get the hash from the ca cert

```
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```

This is all you need in order to join using the `kubeadm join --token <token> <lb>:<port> --discovery-token-ca-cert-hash sha256:<hash>` syntax. You can list "join tokens" with

```
kubeadm token list
```
## Create a Route

This QnD assumes the following

* You want to create a route for your app
* Your app is already up and running with a service
* You have installed an [ingress controller](#ingress)
* You have exported that ingress controller on a [node on port 80](#using-externalips)
* You have DNS in place that points to that ingress node (or you're using nip.io)

If ALL of the above are true...then feel free to proceed.

To create a "route" create an ingress object to tell the controller what to do when it recv a `HTTP_HOST` headder with what's in the `host:` section in your YAML. Example...

Create the YAML

```
cat > app2-ingress-exip.yaml <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
  name: app2-ingress-exip
spec:
  rules:
  - host: app2.192.168.1.8.nip.io
    http:
      paths:
      - backend:
          serviceName: appsvc2
          servicePort: 80
        path: /
EOF
```

Apply it to the app in your namespace

```
kubectl create -n test -f app2-ingress-exip.yaml
```
