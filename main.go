// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Hubs = map[string]*Hub{}
var logger = logrus.New()

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(logrus.DebugLevel)
	//
	logger.SetFormatter(
		&logrus.TextFormatter{TimestampFormat: "2006/01/02 - 15:04:05",
			FullTimestamp: true})
	// Initialize Random
	rand.Seed(time.Now().UnixNano())
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("public/templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "start.html", nil)
	})

	r.POST("/", func(c *gin.Context) {
		hub := newHub()
		Hubs[hub.HubID] = hub
		logger.Info("New Room created", hub.HubID)
		c.String(200, fmt.Sprintln("Room created", hub.HubID))
	})

	r.GET(":room/chat", func(c *gin.Context) {
		room := c.Param("room")
		_, err := getHub(room)
		if err != nil {
			logger.Warn("RoomID not available", err)
			c.String(404, "Room not found")
			return
		}
		c.HTML(200, "chat.html", nil)
	})

	r.GET(":room/ws", func(c *gin.Context) {
		room := c.Param("room")
		user := c.Query("user")

		// This is a replica with the Gin-Logger...
		logger.WithFields(logrus.Fields{
			"user": user,
			"room": room,
		}).Info()

		hub, err := getHub(room)
		if err != nil {
			logger.Info("RoomID not available In Websocket", err)
			c.String(404, "Room not found")
			return
		}
		// Handles the Websocket, for this particular requests.
		serveWs(hub, user, c.Writer, c.Request)
	})

	// Goroutine that checks if OpenHubs are connected to,
	// if not Hub is deleted.
	// TODO check if all depending goroutines are stopped/closed
	go CloseClientlessHubs(closeTime)

	r.Run()
}

// Simple Helper function to check if Hub exists.
func getHub(room string) (*Hub, error) {

	if hub, ok := Hubs[room]; ok {
		return hub, nil
	}
	return &Hub{}, errors.New(fmt.Sprintln("room not found for Room", room))
}
