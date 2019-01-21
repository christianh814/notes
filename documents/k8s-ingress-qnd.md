# Ingress on K8S QnD

First find out what host you want to run on

```
[root@jumpbox shm]# kubectl get nodes -l kubernetes.io/hostname=dhcp-host-8.cloud.chx
NAME                    STATUS   ROLES    AGE    VERSION
dhcp-host-8.cloud.chx   Ready    <none>   3d4h   v1.13.1
```

Export the value from the `kubernetes.io/hostname` key. And the DNS of what your nginx route is going to be (**NOTE**, the DNS (a wildcard entry) must be pointed to the "infra" host)

```
export INFRAHOST=dhcp-host-8.cloud.chx
export NGINXROUTE=nginx.192.168.1.8.nip.io
export EXIP=$(dig ${INFRAHOST} +short)
```

Get the template and replace the two values

```
wget https://raw.githubusercontent.com/christianh814/notes/master/documents/k8s-resources/nginx-ingress-TEMPLATE.yaml
envsubst < nginx-ingress-TEMPLATE.yaml > nginx-ingress.yaml
```

Create the ingress controller

```
kubectl create namespace ingress
kubectl create  -n ingress -f nginx-ingress.yaml
```

Now you can create routes back at the [k8s document](k8s.md#ingress)
