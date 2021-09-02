package main

type channel struct {
	name    string
	clients map[*client]bool
}

func newChannel(name string) *channel {
	return &channel{
		name:    name,
		clients: make(map[*client]bool),
	}
}

func (c *channel) broadcast(s string, m []byte) {

	player := newPlayer(s)

	if player != nil {
		player.read()

		switch chname := c.name; chname {
		case "#EMOJI":
			if player.Exist {
				msg := append([]byte("$"), []byte(player.ShoutColor)...)
				msg = append(msg, []byte(s)...)
				msg = append(msg, ": "...)
				msg = append(msg, m...)
				msg = append(msg, '\n')

				for cl := range c.clients {
					cl.conn.Write(msg)
				}
			}

		case "#STAFF":
			if player.Exist {
				if player.Authority == "99" || player.Authority == "100" {

					msg := append([]byte("%"), []byte(player.ShoutColor)...)
					msg = append(msg, []byte(s)...)
					msg = append(msg, ";STAFF];"...)
					msg = append(msg, m...)
					msg = append(msg, '\n')

					for cl := range c.clients {
						cl.conn.Write(msg)
					}
				}
			}

		case "#BZN":
			if player.Exist {
				if player.ShoutCount > 0 {
					defer player.decshout()

					msg := append([]byte("&"), []byte(player.ShoutColor)...)
					msg = append(msg, []byte(s)...)
					msg = append(msg, ": "...)
					msg = append(msg, m...)
					msg = append(msg, '\n')

					for cl := range c.clients {
						cl.conn.Write(msg)
					}

				} else if player.ShoutCount == 0 {
					msgx := append([]byte("&"), []byte("08")...)
					msgx = append(msgx, []byte("@Sistema")...)
					msgx = append(msgx, "] "...)
					msgx = append(msgx, []byte("Voce nao possui buzinas.")...)
					msgx = append(msgx, '\n')

					for cl := range c.clients {
						if cl.username == s {
							cl.conn.Write(msgx)
						}
					}
				}
			}
		default:
			msg := append([]byte(s), ": "...)
			msg = append(msg, m...)
			msg = append(msg, '\n')

			for cl := range c.clients {
				cl.conn.Write(msg)
			}
		}
	}
}
