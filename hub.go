package main

type hub struct {
	clients          map[*client]bool
	monitoringClient *client
	register         chan *client
	unregister       chan *client
	content          string
}

var connectedCount int

var h = hub{
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
	content:    "",
}

func (h *hub) setMonitoringClient(c *client) {
	h.monitoringClient = c
}

func (h *hub) run() {
	for {
		log.Debug("run")
		select {
		case c := <-h.register:
			h.clients[c] = true
			log.Debugf("[%s] register new client", c.ws.RemoteAddr())
			connectedCount++
			log.Infof("connected %d", connectedCount)
			break

		case c := <-h.unregister:
			log.Infof("[%s] unregister new client", c.ws.RemoteAddr())
			_, ok := h.clients[c]
			if ok {
				//		c.close <- true
				delete(h.clients, c)
				close(c.outcomingDataCh)
				connectedCount--
			}
			break
		}
		log.Debug("run close")
	}
}
