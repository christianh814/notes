# Kubernetes with Kops

With `kops` you can spin up a cluster in aws, gce, and digital ocean. These are high level notes.

## Prereqs

Keep these in mind

* I delegated a subdomain to route53
* I set up the `aws` cli on my local box
* I used an IAM account with pretty wide access
* These are highlevel; consult the [kops aws docs](https://github.com/kubernetes/kops/blob/master/docs/aws.md) and the [kops ha docs](https://github.com/kubernetes/kops/blob/master/docs/high_availability.md#advanced-example)
* Use `kops --help` ...it's well documented

## Kops Install

First download the binary

```
curl -LO https://github.com/kubernetes/kops/releases/download/$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
chmod +x kops-linux-amd64
sudo mv kops-linux-amd64 /usr/local/bin/kops
sudo chmod 555 /usr/local/bin/kops
```

You will need `kubectl` too

```
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
sudo mv kubectl /usr/local/bin/kubectl
sudo chmod 555 /usr/local/bin/kubectl
```

Having the `aws` cli is good too

```
sudo dnf -y install python2-boto3
curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
unzip awscli-bundle.zip
sudo ./awscli-bundle/install -i /usr/local/aws -b /bin/aws
rm -rf awscli-bundle*
```

Set `aws` up with your creds

```
mkdir ~/.aws
cat << EOF >>  ~/.aws/credentials
[default]
aws_access_key_id = YOURACCESSKEY
aws_secret_access_key = YOURSECRETACCESSKEY
EOF
```

Test it...this should return your identity

```
aws sts get-caller-identity
```

## Install On AWS

Before you install, create your bucket, and set some configs on your bucket storage

```
aws s3 mb s3://my-bucket-name --region us-east-1
aws s3api put-bucket-versioning --bucket my-bucket-name --versioning-configuration Status=Enabled
aws s3api put-bucket-encryption --bucket my-bucket-name \
--server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'
```

From a high level, first set your cluster config

```
kops create cluster \
    --node-count 3 \
    --zones us-west-2a,us-west-2b,us-west-2c \
    --master-zones us-west-2a,us-west-2b,us-west-2c \
    --dns-zone my.domain.com \
    --node-size t2.medium \
    --master-size t2.medium \
    --networking calico \
    --ssh-public-key ~/.ssh/id_rsa.pub \
    --state s3://my-bucket-name \
    --api-loadbalancer-type public \
    kube.my.domain.com
```

This sets up a config that you can edit further if you made a mistake above

```
kops edit cluster kube.my.domain.com
```

> You may `--master-count 3` if youre an in a region with only 2 AZ

**NOTE** You can change the image (by default Debian) by using `--image`...the below example images

* For Centos I used (didn't work): `--image "aws-marketplace/CentOS Linux 7 x86_64 HVM EBS ENA 1805_01-b7ee8a69-ee97-4a49-9e68-afaee216db2e-ami-77ec9308.4"`
* For Ubuntu I used: `--image "099720109477/ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-20180405"`

To "install" the cluster..."update" your aws account with the state saved in s3...

```
kops update cluster kube.my.domain.com --yes
```

This creates a `~/.kube/config` file for you...verify with...

```
$ kubectl get nodes
NAME                                          STATUS   ROLES    AGE   VERSION
ip-172-20-121-13.us-west-2.compute.internal   Ready    master   2h    v1.11.6
ip-172-20-126-35.us-west-2.compute.internal   Ready    node     2h    v1.11.6
ip-172-20-43-172.us-west-2.compute.internal   Ready    master   2h    v1.11.6
ip-172-20-46-252.us-west-2.compute.internal   Ready    node     2h    v1.11.6
ip-172-20-66-59.us-west-2.compute.internal    Ready    master   2h    v1.11.6
ip-172-20-95-20.us-west-2.compute.internal    Ready    node     2h    v1.11.6
```

Validating via `kops` lets you  see detailed info about your cluster

```
$ kops  validate cluster
Using cluster from kubectl context: kube.my.domain.com

Validating cluster kube.my.domain.com

INSTANCE GROUPS
NAME			ROLE	MACHINETYPE	MIN	MAX	SUBNETS
master-us-west-2a	Master	t2.medium	1	1	us-west-2a
master-us-west-2b	Master	t2.medium	1	1	us-west-2b
master-us-west-2c	Master	t2.medium	1	1	us-west-2c
nodes			Node	t2.medium	3	3	us-west-2a,us-west-2b,us-west-2c

NODE STATUS
NAME						ROLE	READY
ip-172-20-121-13.us-west-2.compute.internal	master	True
ip-172-20-126-35.us-west-2.compute.internal	node	True
ip-172-20-43-172.us-west-2.compute.internal	master	True
ip-172-20-46-252.us-west-2.compute.internal	node	True
ip-172-20-66-59.us-west-2.compute.internal	master	True
ip-172-20-95-20.us-west-2.compute.internal	node	True

Your cluster kube.my.domain.com is ready
```

## Delete cluster

If you want to delete your cluster, you can do it within the `kops` cli

First, do a "dry-run" to see what will be removed

```
kops delete cluster --name kube.my.domain.com
```

To actually delete it, run the same command but with `--yes`

```
kops delete cluster --name kube.my.domain.com --yes
```

If you have multiple clusters, get the name with...

```
kops get cluster
```

## Kops Export Kubecfg

To export the kubecfg file...

```
export KUBECONFIG=$HOME/test-kubeconfig.yaml
kops export kubecfg --name k8s.example.com --state s3://my-bucket-name
```

That'll save the state in the path that you set in the `KUBECONFIG` env variable
