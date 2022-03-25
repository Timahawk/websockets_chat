// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var Hubs = map[string]*Hub{}

func main() {

	/*
		hub := newHub()
		go hub.run()
	*/
	hub_a := &Hub{
		HubID:      "AAAAAA",
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool)}
	go hub_a.run()

	hub_b := &Hub{
		HubID:      "BBBBBB",
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool)}
	go hub_b.run()

	// Hubs = append(Hubs, hub_a, hub_b)
	Hubs[hub_a.HubID] = hub_a
	Hubs[hub_b.HubID] = hub_b

	r := gin.Default()
	r.LoadHTMLGlob("*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "start.html", nil)
	})

	r.POST("/", func(c *gin.Context) {
		hub := newHub()
		Hubs[hub.HubID] = hub
		log.Println("New Room created", hub.HubID)
		c.String(200, fmt.Sprintln("Room created", hub.HubID))
	})

	r.GET(":room/chat", func(c *gin.Context) {
		room := c.Param("room")
		_, err := getHub(room)
		if err != nil {
			log.Println("RoomID not available", err)
			c.String(404, "Room not found")
			return
		}
		c.HTML(200, "chat.html", nil)
	})

	r.GET(":room/ws", func(c *gin.Context) {
		room := c.Param("room")

		hub, err := getHub(room)
		if err != nil {
			log.Println("RoomID not available In Websocket", err)
			c.String(404, "Room not found")
			return
		}

		serveWs(hub, c.Writer, c.Request)
	})

	// Goroutine that checks if OpenHubs are connected to,
	// if not Hub is deleted.
	// TODO check if all depending goroutines are stopped/closed
	go CloseClientlessHubs(closeTime)

	// Initialize Random
	rand.Seed(time.Now().UnixNano()) //TODO -> move to init

	r.Run()
}

func getHub(room string) (*Hub, error) {

	if hub, ok := Hubs[room]; ok {
		return hub, nil
	}
	return &Hub{}, errors.New(fmt.Sprintln("room not found for Room", room))
}
