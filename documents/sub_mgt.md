# Red Hat Subscription Management

There is a new way to "register" systems to the RHN called "Subscription Management"

Documentation [HERE](https://access.redhat.com/site/documentation/Red_Hat_Subscription_Management/)

But basically...

Register the System

```
subscription-manager register --username sysman@4over.com --password secret
```

Next list available channel pools
```
subscription-manager list --available
```

Attach by pool ID
```
subscription-manager attach --pool=8a85f9843aced34201334b05fbcc2556
```

List what channels you are consuming

```
subscription-manager list --consumed
```

Once you are connected to a pool; you can add custom channels. List them with...
```
subscription-manager repos --list
```

Since it's stdout; you can grep for what you want (or "less" it)
```
subscription-manager repos --list |  grep --color server-optional
  Repo ID:   rhel-6-server-optional-fastrack-rpms
  Repo ID:   rhel-6-server-optional-beta-debug-rpms
  Repo ID:   rhel-6-server-optional-fastrack-debug-rpms
  Repo ID:   rhel-6-server-optional-debug-rpms
  Repo ID:   rhel-6-server-optional-beta-source-rpms
  Repo ID:   rhel-6-server-optional-source-rpms
  Repo ID:   rhel-6-server-optional-beta-rpms
  Repo ID:   rhel-6-server-optional-rpms
  Repo ID:   rhel-6-server-optional-fastrack-source-rpms
```

Add other channels...

```
subscription-manager repos --enable rhel-6-server-optional-rpms
```

**__Note:__** You can "auto subscribe" servers to pools at register time
```
subscription-manager register --username sysman@4over.com --password secret --auto-attach
```

You can also do it on pool attach as well...
```
subscription-manager attach --auto  
```

# YUM

If running RHEL, you can only install security patches if you wish.

```
yum -y install yum-plugin-security
yum --security check-update
yum -y --security update
```

# Misc

The `/var/lib/rhsm/facts/facts.json` gets created at registration...so the chef recipe needs to look something like this maybe...

```
# Register with RHN if it's a redhat server
if node["platform"]
  execute "Register With RHN" do
    command "subscription-manager register --username sysman@4over.com --password secret --auto-attach && subscription-manager repos --enable rhel-6-server-optional-rpms && history -c"
    creates "/var/lib/rhsm/facts/facts.json"
    action :run
  end
end
```


# Repo Sync

You can create your own repo by using `reposync`

```
for repo in \
rhel-7-server-rpms \
rhel-7-server-extras-rpms \
rhel-7-fast-datapath-rpms \
rhel-7-server-ose-3.5-rpms
do
  reposync --gpgcheck -lm --repoid=${repo} --download_path=/path/to/repos
  createrepo -v </path/to/repos/>${repo} -o </path/to/repos/>${repo}
done</code>

If this is in an apache webserver you can host your own repo by creating `/etc/yum.repos.d/myrepo.repo` file with the following contents on the client
```
[myrepo]  
name=my repo
baseurl=http://myrepo.example.com/path/to/repos/
enabled=1
gpgcheck=0
```
