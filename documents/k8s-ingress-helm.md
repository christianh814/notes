# Ingress install with Helm

This is a QnD method but with Helm

You need to first decide where to put it, label your nodes something that you can use in your selector

```
# kubectl label node dhcp-host-8.cloud.chx nginx=ingresshost
node/dhcp-host-8.cloud.chx labeled

# kubectl get nodes -l nginx=ingresshost
NAME                    STATUS   ROLES    AGE   VERSION
dhcp-host-8.cloud.chx   Ready    <none>   10d   v1.13.1
```

The IP for this node is `192.168.1.8`

```
# dig dhcp-host-8.cloud.chx +short
192.168.1.8
```

I want to install this in the `ingress` namespace; so I'll create it

```
kubectl create namespace ingress
```

Using this information; I will deploy the `nginx-ingress` helm chart; giving it the name `nginx-ingress`. List of all options can be found on the [github page](https://github.com/helm/charts/tree/master/stable/nginx-ingress#configuration)

__For Internal clusters__
```
helm install --name nginx-ingress stable/nginx-ingress --namespace ingress \
--set rbac.create=true --set controller.image.pullPolicy="Always" \
--set controller.nodeSelector.nginx="ingresshost" --set controller.stats.enabled=true \
--set controller.service.externalIPs={192.168.1.8} --set controller.service.type="ClusterIP"
```

^ Pay close attention to `controller.nodeSelector`...the syntax is `controller.nodeSelector.<key>="<value>"` ...note the use of the dot instead of an `=`. Also note that `controller.service.externalIPs` is an array

__For cloud clusters__
```
helm install --name nginx-ingress stable/nginx-ingress --namespace ingress --set rbac.create=true --set controller.image.pullPolicy="Always" --set controller.nodeSelector.nginx="ingresshost" --set controller.stats.enabled=true --set controller.hostNetwork=true
```

Export the stats page if you wish (make sure the `svc` name and the `port`are right)

```
cat <<EOF | kubectl apply -n ingress -f -
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: nginx-ingress
spec:
  rules:
  - host: nginx.192.168.1.8.nip.io
    http:
      paths:
      - backend:
          serviceName: nginx-ingress-controller-stats
          servicePort: 18080
        path: /nginx_status
EOF
```

# Cert Manager for TLS

** WIP **
