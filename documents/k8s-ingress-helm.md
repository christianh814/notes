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

```
helm install --name nginx-ingress stable/nginx-ingress --namespace ingress \
--set rbac.create=true --set controller.image.pullPolicy="Always" \
--set controller.nodeSelector.nginx="ingresshost" --set controller.stats.enabled=true \
--set controller.service.externalIPs={192.168.1.8} --set controller.service.type="ClusterIP"
```

^ Pay close attention to `controller.nodeSelector`...the syntax is `controller.nodeSelector.<key>="<value>"` ...note the use of the dot instead of an `=`. Also note that `controller.service.externalIPs` is an array



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

# TLS

The following sets up `nginx-ingress` with `cert-bot` for getting/updating ssl certs using [letsencrypt](https://letsencrypt.org/). Note the following two "gotchas"

* These steps worked on a cluster installed with [kops on aws](k8s-kops.md)
* These steps worked on an "internal" cluster with the "infra" node exposed to the internet on `80` and `443`

# Install Nginx Ingress

**LOCAL**

If you're installing nginx "internall" the [above installation with helm](#ingress-install-with-helm) should work just fine. Just pay attention to your `nodeSelector` and you're targeting your "exposed" node.

**CLOUD** 

After installing with [kops](k8s-kops.md), I did an `init` with [helm](k8s.md#helm). Then I did a simple install like this

```
helm install --name nginx-ingress stable/nginx-ingress \
--namespace ingress --set rbac.create=true --set controller.image.pullPolicy="Always"
```

The above automatically creates a Network `ELB` on `80` and `443` passing through to the NGINX ingress controller (that then sends that traffic to the pods). I took the `ELB` DNS name...

```
$ kubectl get svc nginx-ingress-controller -n ingress -o jsonpath='{.status.loadBalancer.ingress[0].hostname}{"\n"}'
```
Then I created a Route53 wildcard entry `*.apps.my.cluster.com` that CNAMEs to that ELB.

## Cert Manager

First thing to do, is use `helm` to install it

```
$ helm install stable/cert-manager
```

Create a  `cluster issuer` yaml. Replace the email with a valid email address

```
# cat <<EOF > clusterissuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    email: example@email.eg
    http01: {}
    privateKeySecretRef:
      key: ""
      name: letsencrypt-prod
    server: https://acme-v02.api.letsencrypt.org/directory
EOF
```

Create this in your env

```
kubectl create -f clusterissuer.yaml
```

Create your ingress config yaml

```
# cat <<EOF > welcome-php-ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    ingress.kubernetes.io/ssl-redirect: "true"
    kubernetes.io/tls-acme: "true"
    certmanager.k8s.io/cluster-issuer: letsencrypt-prod
    kubernetes.io/ingress.class: "nginx"
  name: welcome-php
spec:
  tls:
  - hosts:
    - welcome-php.8.8.8.8.nip.io
    secretName: welcome-php-tls
  rules:
  - host: welcome-php.8.8.8.8.nip.io
    http:
      paths:
      - backend:
          serviceName: welcome-php
          servicePort: 8080
        path: /
EOF
```

^ **Note** the annotations! Get the `cluster-issuer` name with `kops get cluster`

```
# kubectl create -n test -f welcome-php-ingress.yaml
```

Viola! You have a valid cert with letsencrypt (that will autorenew too!)
