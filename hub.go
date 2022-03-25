// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"
)

const (
	// Period to loop through all Hubs and Close those without Clients.
	closeTime = 5 * time.Minute
	// This are all the letters RandString can draw from.
	letterBytes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Hub ID
	HubID string

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	hub := &Hub{
		HubID:      RandString(6), // TODO muss was bessers her.
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	logger.Info("Created new Hub", hub.HubID)
	go hub.run()
	return hub
}

func (h *Hub) run() {
	logger.Info("Started to run Hub ", h.HubID)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {

				// I dont get why this is nessesarry
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}

				// client.send <- message
			}
		}
	}
}

func CloseClientlessHubs(closeTime time.Duration) {
	ticker := time.NewTicker(closeTime)
	defer func() {
		ticker.Stop()
	}()
	fail := make(chan bool)
	for {
		select {
		case <-ticker.C:
			for _, hub := range Hubs {
				if x := len(hub.clients); x == 0 {

					delete(Hubs, hub.HubID)

					logger.Info("Closed Hub ", hub.HubID, "because no Clients.")
				}
			}
		// This is utterly stupid
		// But don´t know how to fix.
		case <-fail:
			logger.Panic("This should NEVER be run")
		}
	}
}

// Generates Random String with Lenght n from Letters specified in letterBytes
//
// Stolen from Stack Overflow.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
