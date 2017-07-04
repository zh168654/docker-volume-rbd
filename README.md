# Docker volume plugin for RBD

## Requirements

1. Docker >=1.13.1 (recommended)
2. Ceph cluster
3. Etcd. We need a KV store to persist state.






### 2 - Install the plugin

```bash
1.yum install librados2-devel librbd1-devel ceph-common xfsprogs -y

2.mkdir -p /var/lib/docker-volume-rbd/volumes

3.cp rbd-volume-plugin.ini  /var/lib/docker-volume-rbd/

```
LOG_LEVEL=[0:ErrorLevel; 1:WarnLevel; 2:InfoLevel; 3:DebugLevel] defaults to 0

LOG_LEVEL=3

ETCD_ENDPOINTS=10.110.17.121:2379

RBD_CONF_MAP_DEVICE_ROOT="/dev/rbd"

RBD_CONF_CLUSTER=ceph

RBD_CONF_POOL=hikube

RBD_CONF_KEYRING_USER=client.admin

RBD_CONF_KEYRING_KEY="AQAzZ0JZMSqpLxAAF1yajF3ud0X5mKsvp1+I6g=="

RBD_CONF_KEYRING_CAPS_MDS="allow *"

RBD_CONF_KEYRING_CAPS_MON="allow *"

RBD_CONF_KEYRING_CAPS_OSD="allow *"

RBD_CONF_GLOBAL_FSID="52ecee0d-d728-4cd7-9703-f94f2d0c0fc7"

RBD_CONF_GLOBAL_MON_INITIAL_MEMBERS="node-6.k8s.test, node-7.k8s.test, node-8.k8s.test"

RBD_CONF_GLOBAL_MON_HOST="10.110.20.6,10.110.20.7,10.110.20.8"

RBD_CONF_GLOBAL_AUTH_CLUSTER_REQUIRED=cephx

RBD_CONF_GLOBAL_AUTH_SERVICE_REQUIRED=cephx

RBD_CONF_GLOBAL_AUTH_CLIENT_REQUIRED=cephx

RBD_CONF_GLOBAL_OSD_POOL_DEFAULT_SIZE=3

RBD_CONF_GLOBAL_PUBLIC_NETWORK="0.0.0.0/0"

RBD_CONF_CLIENT_RBD_DEFAULT_FEATURES=2

RBD_CONF_MDS_SESSION_TIMEOUT=120

RBD_CONF_MDS_SESSION_AUTOCLOSE=600
```
4.generate file: /etc/ceph/ceph.conf and /etc/ceph/ceph.keyring
``` /etc/ceph/ceph.conf
[global]

fsid = "52ecee0d-d728-4cd7-9703-f94f2d0c0fc7"
mon_initial_members = "node-6.k8s.test, node-7.k8s.test, node-8.k8s.test"
mon_host = "10.110.20.6,10.110.20.7,10.110.20.8"
auth_cluster_required = cephx
auth_service_required = cephx
auth_client_required = cephx

# Write an object n times.
osd pool default size = 2

#All clusters have a front-side public network.
#If you have two NICs, you can configure a back side cluster
#network for OSD object replication, heart beats, backfilling,
#recovery, etc.
public network = "0.0.0.0/0"



[client]
rbd_default_features = 2
keyring = /etc/ceph/ceph.keyring



[mds]
mds_session_timeout = 120
mds_session_autoclose = 600
```


```/etc/ceph/ceph.keyring
[client.admin]
        key = "AQAzZ0JZMSqpLxAAF1yajF3ud0X5mKsvp1+I6g=="
        caps mds = ""allow *""
        caps mon = ""allow *""
        caps osd = ""allow *""

```

4.run: ./docker-volume-rbd
```

### 3 - Create and use a volume

#### Available volume driver options:

```conf
fstype: optional, defauls to xfs
size: optional, defaults to 512 (512MB)
order: optional, defaults to 22 (4KB Objects)
```

#### 3.A - Create a volume:


```sh
docker volume create -d wetopi/rbd  -o size=206 my_rbd_volume

## Known problems:

1. **WHEN** node restart **THEN** rbd plugin breaks: `(modprobe: ERROR: could not insert 'rbd': Operation not permitted //rbd: failed to load rbd kernel module (1) // rbd: sysfs write failed // In some cases useful info is found in syslog - try "dmesg | tail" or so. // rbd: map failed: (2) No such file or directory`
 **SOLUTION** load the module in your hosts: `modprobe rbd` **THEN** plugin works (our container plugin then finds its rbd module on kernel)

```bash
###Todo
###curl -s curl http://localhost:8500/v1/kv/docker/volume/rbd/my_rbd_volume?raw
```

### Use curl to debug plugin socket issues.

To verify if the plugin API socket that the docker daemon communicates with is responsive, use curl. In this example, we will make API calls from the docker host to volume and network plugins using curl to ensure that the plugin is listening on the said socket. For a well functioning plugin, these basic requests should work. Note that plugin sockets are available on the host under /var/run/docker/plugins/<pluginID>

```bash
curl -H "Content-Type: application/json" -XPOST -d '{}' --unix-socket /var/run/docker/plugins/546ac5b9043ce0f49552b14e9fb73dc78f1028d2da7e894ab599e6546566c0df/rbd.sock http:/VolumeDriver.List

{"Mountpoint":"","Err":"","Volumes":[{"Name":"rbd_test","Mountpoint":"","Status":null},{"Name":"demo_test","Mountpoint":"/mnt/volumes/demo_test","Status":null}],"Volume":null,"Capabilities":{"Scope":""}}
```
