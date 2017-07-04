package dockerVolumeRbd

import (
	"os"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

// TODO create a "state" interface and factory in order to use different state backends

const KEY_PREFIX = "/docker/volume/plugin/rbd/"

func (d *rbdDriver) setVolume(v *Volume) error {
	logrus.WithField("etcd.go", "setVolume").Debugf("%#v", v)

	err, kv := getConnection()
	if err != nil {
		return err
	}

	data, err := json.Marshal(v)
	if err != nil {
		logrus.WithField("etcd.go", "setVolume").Error(err)
		return err
	}
	_, err = kv.Put(context.Background(), getKeyFromName(v.Name), string(data))
	if err != nil {
		logrus.WithField("etcd.go", "setVolume").Error(err)
		panic(err)
	}

	return nil

}

func (d *rbdDriver) deleteVolume(name string) error {
	logrus.WithField("etcd.go", "deleteVolume").Debugf("volume name: %s", name)

	err, kv := getConnection()
	if err != nil {
		return err
	}
	_, err = kv.Delete(context.Background(), getKeyFromName(name))
	if err != nil {
		logrus.WithField("etcd.go", "deleteVolume").Error(err)
		return err
	}

	return nil
}

func (d *rbdDriver) getVolume(name string) (error, *Volume) {
	logrus.WithField("etcd.go", "getVolume").Debugf("volume name: %s", name)

	err, kv := getConnection()
	if err != nil {
		return err, nil
	}
	res, err := kv.Get(context.Background(), getKeyFromName(name))

	if err != nil {
		logrus.WithField("etcd.go", "getVolume").Error(err)
		return err, nil
	}

	v := Volume{}

	if res.Kvs != nil {
		if err := json.Unmarshal(res.Kvs[0].Value, &v); err != nil {
			logrus.WithField("etcd.go", "getVolume").Error(err)
			return err, nil
		}
	}

	return nil, &v
}

func (d *rbdDriver) getVolumes() (error, *map[string]*Volume) {
	logrus.WithField("etcd.go", "getVolumes").Debug("get list of volumes")

	err, kv := getConnection()
	if err != nil {
		return err, nil
	}
	res, err := kv.Get(context.Background(), KEY_PREFIX,client.WithPrefix())
	if err != nil {
		logrus.WithField("etcd.go", "getVolumes").Error(err)
		return err, nil
	}

	volumes := map[string]*Volume{}
	if res.Kvs != nil {
		for _, pair := range res.Kvs {

			v := Volume{}

			if err := json.Unmarshal(pair.Value, &v); err != nil {
				logrus.WithField("etcd.go", "getVolumes").Error(err)
				return err, nil
			}

			volumes[v.Name] = &v
		}
	}
	return nil, &volumes
}

func getConnection() (error, client.KV) {
	cfg := client.Config{
		Endpoints: []string{os.Getenv("ETCD_ENDPOINTS")},
		//Transport:               client.DefaultTransport,
		//HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		logrus.WithField("etcd.go", "getConnection").Error(err)
		return err, nil
	}
	kapi := client.NewKV(c)
	return nil, kapi

}

func getKeyFromName(name string) string {
	return KEY_PREFIX + name
}

