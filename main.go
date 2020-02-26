package main

import (
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
	"strconv"
)

const (
	defaultServerPort = "5222"
	serverDomain      = iota
	correspondent
	login
	pass
	serverPort
	message
	correspondentIsRoom
)

func main() {
	// Find server port from action config or use default one
	var port string
	if os.Args[serverPort] == "" {
		port = defaultServerPort
	} else {
		port = os.Args[serverPort]
	}

	// Build client and connect to server
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: os.Args[serverDomain] + ":" + port,
		},
		Jid:          os.Args[login],
		Credential:   xmpp.Password(os.Args[pass]),
		StreamLogger: os.Stdout,
		Insecure:     false,
	}
	router := xmpp.NewRouter()
	client, err := xmpp.NewClient(config, router, errorHandler)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	// Check if we want to send to a chat room or a single user
	// Send presence to connect to chat room, if specified
	// Set the correspondentJid
	var correspondentJid *stanza.Jid
	isCorrespRoom, err := strconv.ParseBool(os.Args[correspondentIsRoom])
	if err != nil {
		panic("failed to determine if sending to a client or chat room : " + err.Error())
	}

	if isCorrespRoom {
		// Building Jid for the room.
		// Here we store the room name as the "node", the server domain as "domain" and the bot alias in the room as
		// the "resource". See XEP-0045
		correspondentJid, err = stanza.NewJid(os.Args[correspondent] + "@" + os.Args[serverDomain] + "/github_bot")
		if err != nil {
			panic(err)
		}
		// Sending room presence
		joinMUC(client, correspondentJid)
	} else {
		correspondentJid, err = stanza.NewJid(os.Args[correspondent])
		if err != nil {
			panic(err)
		}
	}

	// Send github message to recipient or chat room
	m := stanza.Message{Attrs: stanza.Attrs{To: correspondentJid.Bare(), Type: getMessageType(isCorrespRoom)}, Body: os.Args[message]}
	err = client.Send(m)
	if err != nil {
		panic(err)
	}

	// After sending the action message, let's disconnect from the chat room if we were connected to one.
	if isCorrespRoom {
		leaveMUC(client, correspondentJid)
	}
	// And disconnect from the server
	client.Disconnect()
}

// errorHandler is the client error handler
func errorHandler(err error) {
	fmt.Println(err.Error())
}

// joinMUC builds a presence stanza to request joining a chat room
func joinMUC(c xmpp.Sender, toJID *stanza.Jid) error {
	return c.Send(stanza.Presence{Attrs: stanza.Attrs{To: toJID.Full()},
		Extensions: []stanza.PresExtension{
			stanza.MucPresence{
				History: stanza.History{MaxStanzas: stanza.NewNullableInt(0)},
			}},
	})
}

// leaveMUC builds a presence stanza to request leaving a chat room
func leaveMUC(c xmpp.Sender, muc *stanza.Jid) {
	c.Send(stanza.Presence{Attrs: stanza.Attrs{
		To:   muc.Full(),
		Type: stanza.PresenceTypeUnavailable,
	}})
}

// getMessageType figures out the right message type for the github message, depending on what recipient is targeted.
// Message type is important since servers are allowed to ignore messages that do not have an appropriate type
func getMessageType(isCorrespRoom bool) stanza.StanzaType {
	if isCorrespRoom {
		return stanza.MessageTypeGroupchat
	}
	return stanza.MessageTypeChat
}
