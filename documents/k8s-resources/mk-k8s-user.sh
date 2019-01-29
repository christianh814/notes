#!/bin/bash

if [[ "$#" -ne 2 ]]; then
	echo "Usage: $(basename $0) user group"
	exit 1
fi

USER=$1
GROUP=$2
CLUSTERENDPOINT=https://api.example.com
CLUSTERNAME=k8s.example.com
CACERT=cluster/tls/ca.crt
CAKEY=cluster/tls/ca.key
CLIENTCERTKEY=clients/$USER/$USER-key.pem
CLIENTCERTCSR=clients/$USER/$USER.csr

## Create the dir for the user with the artifacts
mkdir -p clients/$USER

## Create the certs and key
echo "{\"CN\":\"$USER\",\"hosts\":[\"\"],\"key\":{\"algo\":\"rsa\",\"size\":2048}}" | cfssl genkey  - | cfssljson -bare clients/$USER/$USER

## Submit them to K8S
cat <<EOF | kubectl create -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: user-request-$USER
spec:
  groups:
  - system:authenticated
  request: $(base64 -w0 < $CLIENTCERTCSR)
  usages:
  - digital signature
  - key encipherment
  - client auth
EOF

# Sleep here for testing
## Approve the request
sleep 2
kubectl certificate approve user-request-$USER

## Create a ns for the user
kubectl create ns sbx-$USER

## Bind the role
cat <<EOF | kubectl create -f -
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: $USER-edit-sbx
  namespace: sbx-$USER
subjects:
- kind: User
  name: $USER
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: edit
  apiGroup: rbac.authorization.k8s.io
EOF

## Get the user CRT
kubectl get csr user-request-$USER -o jsonpath='{.status.certificate}' | base64 -d > clients/$USER/$USER.crt
CLIENTCERTCRT=clients/$USER/$USER.crt

cat <<EOF > clients/$USER/kubeconfig
apiVersion: v1
kind: Config
preferences:
  colors: true
current-context: $CLUSTERNAME
clusters:
- name: $CLUSTERNAME
  cluster:
    server: $CLUSTERENDPOINT
    certificate-authority-data: $(base64 -w 0 < $CACERT)
contexts:
- context:
    cluster: $CLUSTERNAME
    user: $USER
    namespace: sbx-$USER
  name: $CLUSTERNAME
users:
- name: $USER
  user:
    client-certificate-data: $( base64 -w 0 < $CLIENTCERTCRT)
    client-key-data: $(base64 -w 0 < $CLIENTCERTKEY)
EOF

##
##
