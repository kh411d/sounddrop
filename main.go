package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/thejerf/suture"
	"github.com/tuarrep/sounddrop/service"
	"github.com/tuarrep/sounddrop/util"
	"os"
	"os/signal"
)

// Rev is set on build time and should contain the git commit
var rev = ""

func main() {
	util.InitLogger()

	log := util.GetContextLogger("main.go", "main")

	log.Info("Starting main process...")
	myID := util.GetMyID()

	log.Info("I'm known on mesh by: ", myID.String())

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	sb := util.GetServiceBag()
	sb.DeviceID = myID

	supervisor := suture.NewSimple("supervisor")

	messenger := &service.Messenger{Message: make(chan proto.Message)}
	supervisor.Add(messenger)

	server := &service.Server{Messenger: messenger}
	supervisor.Add(server)

	mesher := &service.Mesher{Messenger: messenger}
	supervisor.Add(mesher)

	if sb.Config.Streamer.AutoStart {
		streamer := &service.Streamer{Messenger: messenger}
		supervisor.Add(streamer)
	} //else {
	player := &service.Player{Messenger: messenger}
	supervisor.Add(player)
	//}

	supervisor.ServeBackground()

	log.Info("Main process started. Revision: ", rev)

	for {
		select {
		case <-stop:
			supervisor.Stop()
			log.Info("Clean exit. \n")
			os.Exit(0)
		}
	}
}
