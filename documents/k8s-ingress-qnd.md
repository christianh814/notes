# Ingress on K8S QnD

First find out what host you want to run on

```
[root@jumpbox shm]# kubectl get nodes -l kubernetes.io/hostname=dhcp-host-8.cloud.chx
NAME                    STATUS   ROLES    AGE    VERSION
dhcp-host-8.cloud.chx   Ready    <none>   3d4h   v1.13.1
```

Export the value from the `kubernetes.io/hostname` key. And the DNS of what your nginx route is going to be

```
export INFRAHOST=dhcp-host-8.cloud.chx
export NGINXROUTE=nginx.192.168.1.8.nip.io
```
