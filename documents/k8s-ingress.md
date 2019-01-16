# Ingress

Right now, I have a cluster up and running and can deploy apps and use `nodePort` to expose these apps. In theory you can put a LB in front of it and you can serve your app.

You can have k8s "manage" your LB with an `ingress controller`. There are many out there but it seems like nginx is the most popular one. Let's get one up and running!

(If you just need to "get one up and running" [look here](k8s-ingress-qnd.md))

We will be using the `test` namespace we created during the [smoke test](#smoke-test). Clean it up if you wish

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

[Back to K8S notes](k8s.md#ingress)
