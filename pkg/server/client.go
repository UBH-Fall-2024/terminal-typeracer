// Server-side networking information.
//
// Contains all the information for lobbies, clients, etc

package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	"github.com/charmbracelet/log"
)

type ClientID uint

// The server side representation of a Connected client
type Client struct {
	Conn  net.Conn
	Id    ClientID
	Name  string
	enc   gob.Encoder
	dec   gob.Decoder
	Lobby *Lobby
}

func (c *Client) Start() {
	c.enc = *gob.NewEncoder(c.Conn)
	c.dec = *gob.NewDecoder(c.Conn)

	for {
		// Read a message from the tcp connection
		var message network.Message
		if err := c.dec.Decode(&message); err != nil {
			// TODO better error handling, don't just kill the client maybe
			log.Error("Could not read message from client")
			c.Disconnect()
			break
		}
		// TODO perhaps make handling messages concurrent? idk if its necessary
		// or could cause race conditions
		c.handleMessage(message)

	}
}

func (c *Client) SendMessage(message *network.Message) error {
	return c.enc.Encode(message)
}

// Used if the client was very naughty (Sending invalid requests and stuff!!)
func (c *Client) SendError(err string) error {
	message := network.Message{
		Header: uint8(network.Error),
		Data:   err,
	}
	return c.SendMessage(&message)
}

func (c *Client) handleMessage(message network.Message) {
	req := network.Request(message.Header)

	if c.Lobby == nil && req == network.Connect {
		// Make the length of the request only 16 chars long (client names
		// should not be over 16 chars)
		clientName := message.Data
		clientName = clientName[:16]

		server.Clients = append(server.Clients, c)
		return
	}

	// YOu need to be in a lobby
	if c.Lobby != nil {
		c.SendError("Why aren't you in a lobby wtf")
		return
	}
	c.Lobby.commands <- func() {
		switch req {
		case network.JoinLobby:
			if c.Lobby != nil {
				c.SendError("You're already in a lobby")
				return
			}

			lobby := c.Lobby

			// Message to be sent to all clients already connected to the lobby
			joinedMsg := network.Message{
				Header: uint8(network.JoinedLobby),
				Data:   fmt.Sprintf("%d,%s", c.Id, c.Name),
			}

			// Let the client trying to join (c) know who is in the lobby they're
			// joining, and tell everyone in the lobby currently that someone new
			// has joined
			for _, client := range lobby.Clients {
				client.SendMessage(&joinedMsg)
				// Message to be sent to the client that joined to let them know
				// what clients are connected
				c.SendMessage(&network.Message{
					Header: uint8(network.JoinedLobby),
					Data:   fmt.Sprintf("%d,%s", client.Id, client.Name),
				})
			}

			if len(lobby.Clients) >= LOBBY_SIZE {
				lobby.Start()
			}
		case network.Progress:
			c.Lobby.SendMessage(&network.Message{
				Header: uint8(network.ProgUpdate),
				Data:   fmt.Sprintf("%d,%s", c.Id, message.Data),
			})

			res, err := strconv.Atoi(message.Data)
			if err != nil {
				c.SendError("Please give me correct data man")
				return
			}
			if res == 100 {
				c.Lobby.finished <- struct{}{}
			}

		default:
			c.SendError("Wtf man, not allowed")
		}
	}
}

func (c *Client) Disconnect() {
	lobby := c.Lobby
	c.Lobby = nil

	if err := c.SendError("Disconnected from server"); err != nil {
		log.Print("Could not tell the client they suck")
	}

	for i, client := range lobby.Clients {
		// Check if the client being iterated over has the same address as c
		if client == c {
			// remove the client from the list
			lobby.Clients[i] = lobby.Clients[len(lobby.Clients)-1]
			lobby.Clients = lobby.Clients[:len(lobby.Clients)-1]
			break
		}
	}

	message := network.Message{
		Header: uint8(network.LeftLobby),
		Data:   string(c.Id),
	}

	// Doing double for loop because i am scared removing from list will mess
	// with iteration
	for _, client := range lobby.Clients {
		client.SendMessage(&message)
	}
}
