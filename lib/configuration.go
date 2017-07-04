package dockerVolumeRbd

import (
	"os"
	"flag"
	"text/template"
	"github.com/Sirupsen/logrus"
	"strings"
	"github.com/larspensjo/config"

)

var  (
	configFile = flag.String("configfile", "/var/lib/docker-volume-rbd/rbd-volume-plugin.ini", "rbd volume configuration file")
)


// Configure Ceph
// get conf files
// create the ceph.conf
// and the ceph.keyring used to authenticate with cephx
//
func (d *rbdDriver) configure() error {

	//var err error

	// set default confs:
	d.conf["cluster"] = "ceph"
	d.conf["device_map_root"] = "/dev/rbd"

	d.loadRbdConfigVars();
	//var err error
	//err = createConf("templates/ceph.conf.tmpl", "/etc/ceph/ceph.conf", d.conf);
	//if err != nil {
	//	return err
	//}
	//
	//err = createConf("templates/ceph.keyring.tmpl", "/etc/ceph/ceph.keyring", d.conf);
	//if err != nil {
	//	return err
	//}

	return nil
}


// Get only the env vars starting by RBD_CONF_*
// i.e. RBD_CONF_GLOBAL_MON_HOST is saved in d.conf[global_mon_host]
//
func (d *rbdDriver) loadEnvironmentRbdConfigVars() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		if (strings.HasPrefix(pair[0], "RBD_CONF_")) {
			configPair := strings.Split(pair[0], "RBD_CONF_")
			d.conf[strings.ToLower(configPair[1])] = pair[1]
		}
	}

}
// Get only the env vars starting by RBD_CONF_*
// i.e. RBD_CONF_GLOBAL_MON_HOST is saved in d.conf[global_mon_host]
//
func (d *rbdDriver) loadRbdConfigVars() {

	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		logrus.Error("read default config file failed")
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("default") {
		section, err := cfg.SectionOptions("default")
		if err == nil {
			for _, v := range section {
				options, _ := cfg.String("default", v)
				logrus.Info("add env "+v+"="+options)
				os.Setenv(v,options)
				if (strings.HasPrefix(v, "RBD_CONF_")) {
					configPair := strings.Split(v, "RBD_CONF_")
					d.conf[strings.ToLower(configPair[1])] = options
				}

			}
		}
	}

}



func createConf(templateFile string, outputFile string, config map[string]string) error {

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		logrus.WithField("utils", "createConf").Error(err)
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		logrus.WithField("utils", "createConf").Error(err)
		return err
	}

	err = t.Execute(f, config)
	if err != nil {
		logrus.WithField("utils", "createConf").Error(err)
		return err
	}

	f.Close()

	return nil
}

