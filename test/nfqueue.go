package main

import (
	"github.com/iesreza/gutil/log"
	"github.com/iesreza/nfqueue"
	"os/exec"
)

type Handler struct{}

func main() {

	var cfg = &nfqueue.QueueConfig{
		MaxPackets: 32,
		QueueFlags: []nfqueue.QueueFlag{nfqueue.Conntrack},
		BufferSize: 16 * 1024 * 1024,
	}
	var i uint16
	for i = 0; i < 4; i++ {
		var handler Handler
		q := nfqueue.NewQueue(i, handler, cfg)
		q.Start()
	}

}

func (Handler) Handle(packet *nfqueue.Packet) {
	log.Notice("Packet on queue %d", packet.Q.ID)
	//fmt.Println(hex.Dump(packet.Buffer))
	packet.Accept()
}

func Run(command string, args ...string) string {

	c := exec.Command(command, args...)
	out, err := c.Output()
	if err != nil {
		log.Error("Unable to run %s %v", command, args)
		log.Error("%s", err)
	}
	return string(out)
}

func ipTables() {
	Run("iptables", "-A", "INPUT", "-j", "NFQUEUE", "--queue-balance", "0:3", "--queue-bypass")
}
