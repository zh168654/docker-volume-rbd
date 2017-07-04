package main

import (
    "os"
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
	"gitserver/iop/docker-volume-rbd/lib"
)

const socketAddress = "/run/docker/plugins/rbd.sock"



func main() {

	dockerVolumeRbdVersion := "1.0"
	err, rbdDriver := dockerVolumeRbd.NewDriver()

	if err != nil {
		logrus.Fatal(err)
	}

	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
		case "3":
			logrus.SetLevel(logrus.DebugLevel)
			break;
		case "2":
			logrus.SetLevel(logrus.InfoLevel)
			break;
		case "1":
			logrus.SetLevel(logrus.WarnLevel)
			break;
		default:
			logrus.SetLevel(logrus.ErrorLevel)
		}
	 file, err := os.OpenFile("/var/lib/docker-volume-rbd/rbd.log", os.O_CREATE|os.O_WRONLY, 0666)
	 if err == nil {
		 logrus.SetOutput(file)
	 } else {
		 logrus.SetOutput(os.Stdout)
	 }


	h := volume.NewHandler(rbdDriver)
	logrus.Infof("plugin(rbd) version(%s) started with log level(%s) attending socket(%s)", dockerVolumeRbdVersion, logLevel, socketAddress)
	logrus.Error(h.ServeUnix(socketAddress, 0))
}
