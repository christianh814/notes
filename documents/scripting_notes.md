# Scripting Notes

* [KSH Notes](#ksh-notes)
* [Misc Commands](#misc-commands)

## KSH Notes

You can test the "modulo" of a number in KSH. A modulo prints the remainder of a division calculation.

Example

```
root@host# echo $(( 9 % 3 ))
0
root@host# echo $(( 9 % 4 ))
1
```

In the first command the output is zero because 9 is evenly divisible by 3. In the second command, 9 is divisible by 4 with a remainder of 1.

The most practical example is to have a cron run every other week. You can see if the date is divisible by 2.

```
[[ $(( $(date +%U) % 2 )) -eq 0 ]] || exit
```

If the week of the year (from 00-53) is not divisible by 2; then exit.


In KSH you can get options using single letter or whole word

Example...

```
while getopts "f:(file)l:(logfile)" opts
do
case "${opts}" in
    f)
      file="${OPTARG}"
      ;;
    l)
      log="${OPTARG}"
      ;;
    \?)
      echo "$0 [ -f file | --file file] [ -l log | --log ]"
      exit
      ;;
esac
done

echo "My file is ${file} and my log is ${log}"
```

Can be called as

```
root@host# myscript -f file.txt -l log.txt
```

OR

```
root@host# myscript --file file.txt --log log.txt
```


Eval takes the line - evaluates it- then puts it back on the stack. This is handy when you want a variable within a variable.

```
#!/bin/ksh -fuhx
#ps -ef | grep -v grep | grep -c reddit
#
set -A tunnels reddit cci.key
reddit_cmd="ssh -i /usr/local/lib/mysql_reddit.key -CfNg 69.43.172.51 -l root -L 3306:69.43.172.51:3306"
cci_cmd="ssh -i /usr/local/lib/cci.key -CfNg 216.133.247.10 -l over4 -L 1433:192.180.1.7:1433"
#
for tunnel in ${tunnels[*]}
do
  if [ $(ps -ef | grep -v grep | grep -c ${tunnel}) -eq 0 ] ; then
    # Looks like the tunnel is not up let's start it
    logger -p local5.info -t SSHTUNNEL "Starting the tunnel for ${tunnel}"
    eval \${${tunnel%.*}_cmd}
  fi
done
#
##-30-
```

This will expand ${tunnel%.*} first - then drop it back in place with it's value; then run the command again with the expanded values. (it runs the line twice ; first with eval then without).

You need the backslash.


## Misc Commands

You can use `csplit` to split on a line number or a regex.

I built 17 ldiff files at 4over using
```
root@host# csplit ,4over.ldiff /-/+1 {*}
```

This means that "Split file using - as the delimiter and offset by one, which includes the - in the file ( /-/+1) and do it as many times as possible" basic syntax is: `csplit FILENAME /REGEX/OFFSET {INTEGER}`


Forward Trim

```
${tunnel%.*} <~ Min Forward trim using . as delimiter
${tunnel%%.*} <~ Max Forward trim using . as delimiter
```

Backward trim

```
${tunnel#.*} <~ Min Backward trim using . as delimiter
${tunnel##.*} <~ Max Backward trim using . as delimiter
```

You can use anything as a delimiter. Like . / or aything you want to trim on. * means "everything from here on out"



First, get file name without the path:
```
filename=$(basename "$fullfile")
extension="${filename##*.}"
filename="${filename%.*}"
```

Alternatively, you can focus on the last '/' of the path instead of the '.' which should work even if you have unpredictable file extensions:
```
filename="${fullfile##*/}"
```



To accurately determine the number of CPU sockets on a system without physically opening and manually inspecting it, one must rely on the system's DMI table (populated by the BIOS). Install the dmidecode package if necessary and then run the following as root.

```
root@host#  dmidecode -t4 | grep Socket.Designation: | wc -l
```

If all the CPU sockets on the system are filled, consulting the `/proc/cpuinfo` file or the lscpu command (RHEL6 and above) as a normal user will be sufficient.

```
user@host$ grep physical.id /proc/cpuinfo | sort -u | wc -l
user@host$ lscpu | grep -i "socket(s)"
```


You can run a process in "daemon" mode. This is expecially useful if you'd like a process to run on a port that is "privileged"(i.e. only root can do this) BUT still be run as an "unprivileged user". For example; like if you'd like apache to run on port 80...you have to run it in daemon mode.

Example:
```
root@host# . /etc/init.d/functions
root@host# daemon --user=myuser --pidfile=/path/to/pid/commands.pid "commands"
```

At 4over; we used mojolicious web server...I had to write a script to start it...

```
#!/bin/bash -fuh
#
if [ $(id -u) -ne 0 ]; then
  echo "Must be root to execute"
fi
. /etc/init.d/functions
export LD_LIBRARY_PATH=/usr/local/profundus/mojo_server/lib:/usr/local/profundus/mojo_server/script:/opt/perl5/perls/perl-5.14.4/lib
export PATH=/opt/perl5/perls/perl-5.14.4/bin:$PATH:/usr/local/profundus/mojo_server/script
daemon --user=mojo --pidfile=/usr/local/profundus/mojo_server/var/run/mojo.pid "/opt/perl5/perls/perl-5.14.4/bin/hypnotoad /usr/local/profundus/mojo_server/script/my4_over_app"
#
##-30-
```


You can get NIC information with ethtool

```
[chrish@storehost ~]$ sudo ethtool eth0
Settings for eth0:
	Supported ports: [ TP ]
	Supported link modes:   10baseT/Half 10baseT/Full 
	                        100baseT/Half 100baseT/Full 
	                        1000baseT/Full 
	Supported pause frame use: No
	Supports auto-negotiation: Yes
	Advertised link modes:  10baseT/Half 10baseT/Full 
	                        100baseT/Half 100baseT/Full 
	                        1000baseT/Full 
	Advertised pause frame use: No
	Advertised auto-negotiation: Yes
	Speed: 1000Mb/s
	Duplex: Full
	Port: Twisted Pair
	PHYAD: 0
	Transceiver: internal
	Auto-negotiation: on
	MDI-X: Unknown
	Supports Wake-on: d
	Wake-on: d
	Current message level: 0x00000007 (7)
			       drv probe link
	Link detected: yes
```

You can get statistics too
```
[chrish@storehost ~]$ sudo ethtool -S eth0
NIC statistics:
     rx_packets: 11378414
     tx_packets: 3434881
     rx_bytes: 5526628999
     tx_bytes: 5872736530
     rx_broadcast: 0
     tx_broadcast: 0
     rx_multicast: 0
     tx_multicast: 0
     rx_errors: 0
     tx_errors: 0
     tx_dropped: 0
     multicast: 0
     collisions: 0
     rx_length_errors: 0
     rx_over_errors: 0
     rx_crc_errors: 0
     rx_frame_errors: 0
     rx_no_buffer_count: 0
     rx_missed_errors: 0
     tx_aborted_errors: 0
     tx_carrier_errors: 0
     tx_fifo_errors: 0
     tx_heartbeat_errors: 0
     tx_window_errors: 0
     tx_abort_late_coll: 0
     tx_deferred_ok: 0
     tx_single_coll_ok: 0
     tx_multi_coll_ok: 0
     tx_timeout_count: 0
     tx_restart_queue: 18892
     rx_long_length_errors: 0
     rx_short_length_errors: 0
     rx_align_errors: 0
     tx_tcp_seg_good: 601131
     tx_tcp_seg_failed: 0
     rx_flow_control_xon: 0
     rx_flow_control_xoff: 0
     tx_flow_control_xon: 0
     tx_flow_control_xoff: 0
     rx_long_byte_count: 5526628999
     rx_csum_offload_good: 6157340
     rx_csum_offload_errors: 0
     alloc_rx_buff_failed: 0
     tx_smbus: 0
     rx_smbus: 0
     dropped_smbus: 0
```

Netstat provides good information too

```
[chrish@storehost ~]$ netstat --statistics
Ip:
    6157725 total packets received
    4 with invalid addresses
    0 forwarded
    0 incoming packets discarded
    6157721 incoming packets delivered
    3421001 requests sent out
Icmp:
    322010 ICMP messages received
    0 input ICMP message failed.
    ICMP input histogram:
        destination unreachable: 91
        echo requests: 321919
    336041 ICMP messages sent
    0 ICMP messages failed
    ICMP output histogram:
        destination unreachable: 14122
        echo replies: 321919
IcmpMsg:
        InType3: 91
        InType8: 321919
        OutType0: 321919
        OutType3: 14122
Tcp:
    1363 active connections openings
    38386 passive connection openings
    28 failed connection attempts
    175 connection resets received
    6 connections established
    5303877 segments received
    2986344 segments send out
    8883 segments retransmited
    0 bad segments received.
    353 resets sent
Udp:
    116050 packets received
    14122 packets to unknown port received.
    0 packet receive errors
    89685 packets sent
UdpLite:
TcpExt:
    8 invalid SYN cookies received
    20 packets pruned from receive queue because of socket buffer overrun
    782 TCP sockets finished time wait in fast timer
    1601 delayed acks sent
    Quick ack mode was activated 134 times
    76 packets directly queued to recvmsg prequeue.
    122384 packets directly received from backlog
    58628 packets directly received from prequeue
    3810483 packets header predicted
    132 packets header predicted and directly queued to user
    689136 acknowledgments not containing data received
    1051950 predicted acknowledgments
    2469 times recovered from packet loss due to SACK data
    17 congestion windows recovered after partial ack
    7972 TCP data loss events
    TCPLostRetransmit: 186
    71 timeouts after SACK recovery
    1 timeouts in loss state
    7280 fast retransmits
    701 forward retransmits
    402 retransmits in slow start
    345 other TCP timeouts
    80 sack retransmits failed
    6333 packets collapsed in receive queue due to low socket buffer
    21 DSACKs sent for old packets
    11 DSACKs received
    51 connections reset due to unexpected data
    171 connections reset due to early user close
    TCPDSACKIgnoredNoUndo: 6
    TCPSpuriousRTOs: 2
    TCPSackShifted: 20510
    TCPSackMerged: 12325
    TCPSackShiftFallback: 6874
    TCPBacklogDrop: 1
IpExt:
    InMcastPkts: 26543
    OutMcastPkts: 63
    InBcastPkts: 401610
    InOctets: 5081727296
    OutOctets: 5822099186
    InMcastOctets: 3187570
    OutMcastOctets: 5092
    InBcastOctets: 37182745
```
    
Pretty JSON output
```
curl -s https://registry.access.redhat.com/v1/repositories/jboss-datagrid-6/datagrid65-openshift/tags |  python -m json.tool
```

Sample case startup script in bash
```
#! /bin/sh
# /etc/init.d/blah
#

# Some things that run always
touch /var/lock/blah

# Carry out specific functions when asked to by the system
case "$1" in
  start)
    echo "Starting script blah "
    echo "Could do more here"
    ;;
  stop)
    echo "Stopping script blah"
    echo "Could do more here"
    ;;
  *)
    echo "Usage: /etc/init.d/blah {start|stop}"
    exit 1
    ;;
esac
```

Long options with bash
```
#!/bin/bash
#length=$*
#< /dev/urandom tr -dc '}?$_A-Z-a-z-0-9' | head -c ${length:-16} ; echo 
while :
do
  case "$1" in
    --length|-l)
      length=$2
      shift 2
      ;;
    --execute|-x)
      execute=$2
      shift 2
      ;;
    --)
      shift
      break
      ;;
    *)
      break
      ;;
  esac
done
echo "running $length and I'm doing it $execute times"
##
##
```

Taring up "sparse" file (i.e. VM discs)

```
####### Sparse tar
tar -cvzSpf fedora-unknown-1.qcow2.tgz fedora-unknown-1.qcow2

####### Sparse untar
tar -xzSpf fedora-unknown-1.qcow2.tgz
```

Edit "sleep" background on gnome3
```
xhost +SI:localuser:gdm
sudo -u gdm dbus-launch gnome-control-center
```

Pigz tar (use `-p 8` maybe?)
```
tar cf - /myarchive | pigz -9 -p 32 > archive.tar.gz
```

# GRUB 2

First See your grub entries

```
root@host# grep menuentry /boot/grub2/grub.cfg
```

Make note of what "position" the entry is. GRUB2 is 0 based. So if it's the 3rd entry; it'd be counted as 2.

Now open the `/etc/default/grub` file and make a change...

```
GRUB_DEFAULT=2
```

Also while I was here I edited the file to turn off apic (`apic=off`) in the `GRUB_CMDLINE_LINUX=` line.

Now after you made the changes; you need to recreate the grub file.

```
root@host# grub2-mkconfig -o /boot/grub2/grub.cfg
```

I ran into an issue after doing a "yum -y upgrade" on Fedora 18 where grub didn't "see" my Windows partition (and subsequently `grub2-mkconfig` gave errors when trying to rebuild the `grub.conf` file). It seems that the script `/etc/grub.d/30_os-prober` wasn't able to detect Windows when running the `grub2-mkconfig -o /boot/grub2/grub.cfg` command.

The errors looks something like this (it spit out quite a few of these lines).

```
ERROR: ddf1: seeking device "/dev/dm-6" to 18446744073709421056
```

I found something online (in an Ubuntu forum out of all places) that solved my issue.

First create a script called `/etc/grub.d/31_windows-probe` that looked like this

```
[root@fedora18 ~]# cat /etc/grub.d/31_windows-probe
#!/bin/sh -e
#
cat << EOF
menuentry "Microsoft Windows 7" {
  set root=(hd0,1)
  chainloader +1
}
EOF
#
##-30-
```

Then I made it executable

```
[root@fedora18 ~]# chmod +x /etc/grub.d/31_windows-probe
```

Now I was able to run the `grub2-mkconfig -o /boot/grub2/grub.cfg` and was able to see Windows 7 in the grub menu
