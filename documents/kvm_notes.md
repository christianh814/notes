# KVM Notes

These are quick notes in no paticular order

# Nested KVM

In KVM - you can "pass along" CPU capabilities to the guest. This poses interesting scenarios where you can have a VM within a VM running a VM (sort of a "dream within a dream" type of thing)

First; check if your machine has Hardware Virt

```
root@host# egrep '(vmx|svm)' --color=always /proc/cpuinfo
```

If nothing is displayed, then your processor doesn't support hardware virtualization.

Now, check to see if nested KVM support is enabled on your machine by running:

```
root@host# cat /sys/module/kvm_intel/parameters/nested
```

If the answer is “N,” you can enable it by running:

```
root@host# echo "options kvm-intel nested=1" > /etc/modprobe.d/kvm-intel.conf
```

After adding that `kvm-intel.conf` file, reboot your machine, after which `cat /sys/module/kvm_intel/parameters/nested` should return `Y`

You can add virtulazation on the UI or on the command line like so

```
[root@dlp ~]# virsh edit www
# edit a virtual machine "www"
# add following lines

<cpu mode='custom' match='exact'>

# CPU model

  <model fallback='allow'>SandyBridge</model>

# CPU vendor

  <vendor>Intel</vendor>
  <feature policy='require' name='vmx'/>
</cpu>
```

# Export to VMware

I followed [this](https://blog.ktz.me/migrate-qcow2-images-from-kvm-to-vmware/) and [this](https://blog.ktz.me/gotchas-when-migrating-fedora-qcow2-images-to-vmware/) when migrating from Libvirtd to vSphere/ESXi.

In the end I still had to boot into "rescue" mode from grub and run the following in that shell

```
dracut --regenerate-all --force
```

Looks like [this](https://possiblelossofprecision.net/?p=2293) is also a good howto as well.
