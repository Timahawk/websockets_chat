// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"reflect"
	"time"
)

// Period to loop through all Hubs and Close those without Clients.
const closeTime = 1 * time.Hour

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
		HubID:      "ACH42H", // TODO muss was bessers her.
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	log.Println("Created new Hub", hub.HubID)
	go hub.run()
	return hub
}

func (h *Hub) run() {
	log.Println("Started to run Hub", h.HubID)
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
			for idx, hub := range Hubs {
				if x := len(hub.clients); x == 0 {
					log.Println("Closing Hub ", hub.HubID, "because no Clients.")
					sliceRemoveItem(Hubs, idx)
				}
			}
		// This is unterly stupid
		// But don´t know how to fix.
		case <-fail:
			log.Println("This should NEVER be run")
		}
	}
}

func sliceRemoveItem(slicep interface{}, i int) {
	v := reflect.ValueOf(slicep).Elem()
	v.Set(reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len())))
}
