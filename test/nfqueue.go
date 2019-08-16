package main

import (
	"github.com/iesreza/gutil/log"
	"github.com/iesreza/nfqueue"
	"os/exec"
	"time"
)

type Handler struct{}

func main() {

	var cfg = &nfqueue.QueueConfig{
		MaxPackets: 32,
		QueueFlags: []nfqueue.QueueFlag{nfqueue.Conntrack},
		BufferSize: 1600,
	}

	var i int
	var j uint16
	for i = 0; i < 4; i++ {

		go func() {
			j++
			log.Info("Queue %d created", j-1)
			var handler Handler
			q := nfqueue.NewQueue(j-1, handler, cfg)
			q.Start()

		}()

	}

	ipTables()
	for {
		time.Sleep(10 * time.Second)
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

	Run("iptables", "-F")
	Run("iptables", "-t", "filter", "-F")
	Run("iptables", "-t", "mangle", "-F")
	Run("iptables", "-t", "nat", "-F")

	//Enable forwarding
	Run("echo", "1 > /proc/sys/net/ipv4/ip_forward")
	Run("sysctl", "-w", "net.ipv4.conf.eth0.route_localnet=1")

	Run("iptables", "-t", "mangle", "-A", "INPUT", "-j", "NFQUEUE", "--queue-balance", "0:3", "--queue-cpu-fanout", "--queue-bypass")
	Run("iptables", "-t", "mangle", "-A", "OUTPUT", "-j", "NFQUEUE", "--queue-balance", "0:3", "--queue-cpu-fanout", "--queue-bypass")
	Run("iptables", "-t", "mangle", "-A", "FORWARD", "-j", "NFQUEUE", "--queue-balance", "0:3", "--queue-cpu-fanout", "--queue-bypass")
}
