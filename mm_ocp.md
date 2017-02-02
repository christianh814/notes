# OpenShift MM Config by NickNach
```
# Create an OSEv3 group that contains the masters and nodes groups
[OSEv3:children]
masters
nodes
new_nodes
etcd
lb

## Set variables common for all OSEv3 hosts
[OSEv3:vars]
ansible_ssh_user=root
deployment_type=openshift-enterprise

## If ansible_ssh_user is not root, ansible_sudo must be set to true
#ansible_sudo=true
#ansible_ssh_user=ec2-user
#ansible_become=yes

## authentication stuff
## htpasswd file auth
openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]
## ldap auth
#openshift_master_identity_providers=[{"name":"NNWIN","challenge":true,"login":true,"kind":"LDAPPasswordIdentityProvider","attributes":{"id":["dn"],"email":["mail"],"name":["cn"],"preferredUsername":["sAMAccountName"]},"bindDN":"CN=SVC-nn-ose,OU=SVC,OU=FNA,DC=nnwin,DC=ad,DC=nncorp,DC=com","bindPassword":"<REDACTED>","insecure":true,"url":"ldap://uswin.nicknach.com:389/DC=uswin,DC=ad,DC=nncorp,DC=com?sAMAccountName?sub"}]
#openshift_master_ldap_ca_file=/etc/ssl/certs/NNWINDC_Cert_Chain.pem

## service endpoints
openshift_hosted_metrics_public_url=https://hawkular-metrics.apps.ocp.nicknach.net/hawkular/metrics
openshift_master_logging_public_url=https://kibana.apps.ocp.nicknach.net/

##  cloud provider configs
##  AWS
#openshift_cloudprovider_kind=aws
#openshift_cloudprovider_aws_access_key=SKFHLJSDHFLJSDLKFHDSKJFDSFDSFSDFDSFDS
#openshift_cloudprovider_aws_secret_key=SKDJFHILDSHFUSWNFKJNSDKJNFLDSKJNFLKSDNFSLKDJNFSLKJFNDSf
##  Openstack
#openshift_cloudprovider_kind=openstack
#openshift_cloudprovider_openstack_auth_url=https://controller.home.nicknach.com:35357/v2.0
#openshift_cloudprovider_openstack_username=svc-openshift-np
#openshift_cloudprovider_openstack_password=KJHDSLEJFHLKJSDFSDFDSF
#openshift_cloudprovider_openstack_tenant_id=fasjfhasjdflakjsndflsakjndfsajdfs
#openshift_cloudprovider_openstack_tenant_name=nn-dev
#openshift_cloudprovider_openstack_region=RegionOne
#openshift_cloudprovider_openstack_lb_subnet_id=222222222-2222-461d-af28-322222222222

## domain stuff
osm_default_subdomain=apps.ocp.nicknach.net

## cluster stuff (uncomment for multi-master mode)
openshift_master_cluster_method=native
openshift_master_cluster_hostname=api.ocp.nicknach.net
openshift_master_cluster_public_hostname=console.ocp.nicknach.net

## network stuff
#os_sdn_network_plugin_name='redhat/openshift-ovs-multitenant'
#set these if you are behind a proxy
#openshift_http_proxy=
#openshift_https_proxy=
#openshift_no_proxy=.nicknach.net

#use these if there is a network IP conflict
#osm_cluster_network_cidr=10.128.0.0/14
#openshift_portal_net=172.30.0.0/16

#change api port
# If using the 443 below
# openshift_master_public_api_url=https://api.ocp.nicknach.net
# openshift_master_public_console_url=https://console.ocp.nicknach.net/console
# These ports reference the ports on the masters themselves 
# openshift_master_api_port=443
# openshift_master_console_port=443

## adjust max pods for scale testing
#openshift_node_kubelet_args={'max-pods': ['225'], 'image-gc-high-threshold': ['90'], 'image-gc-low-threshold': ['80']}
## adjust scheduler
#osm_controller_args={'node-monitor-period': ['2s'], 'node-monitor-grace-period': ['16s'], 'pod-eviction-timeout': ['30s']}

## start group defs
## load balancer
[lb]
lb.ocp.nicknach.net

## host group for etcd (uncomment for multi-master)
[etcd]
master01.ocp.nicknach.net
master02.ocp.nicknach.net
master03.ocp.nicknach.net

## host group for masters
[masters]
master01.ocp.nicknach.net
master02.ocp.nicknach.net
master03.ocp.nicknach.net

[nodes]
master01.ocp.nicknach.net openshift_node_labels="{'region': 'masters', 'zone': 'a'}"
master02.ocp.nicknach.net openshift_node_labels="{'region': 'masters', 'zone': 'b'}"
master03.ocp.nicknach.net openshift_node_labels="{'region': 'masters', 'zone': 'c'}"
infra01.ocp.nicknach.net openshift_node_labels="{'region': 'infra', 'zone': 'a'}"
infra02.ocp.nicknach.net openshift_node_labels="{'region': 'infra', 'zone': 'b'}"
infra03.ocp.nicknach.net openshift_node_labels="{'region': 'infra', 'zone': 'c'}"
node01.ocp.nicknach.net openshift_node_labels="{'region': 'primary', 'zone': 'a'}"
node02.ocp.nicknach.net openshift_node_labels="{'region': 'primary', 'zone': 'b'}"
node03.ocp.nicknach.net openshift_node_labels="{'region': 'primary', 'zone': 'c'}"

[new_nodes]
## hold for use when adding new nodes
```
