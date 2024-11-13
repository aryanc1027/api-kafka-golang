package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/aryanc1027/api-kafka-golang/internal/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	producer kafka.Producer
	consumer kafka.Consumer
	streams  map[string]*Stream
}

type Stream struct {
	ID       string
	Messages chan []byte
	Errors   chan error
}


func NewHandler(producer kafka.Producer, consumer kafka.Consumer) *Handler {
	return &Handler{
		producer: producer,
		consumer: consumer,
		streams:  make(map[string]*Stream),
	}
}

// simple function to process data stream
func processData(data []byte) []byte {
	
	return bytes.ToUpper(data)
}

func (h *Handler) StartStream(c *gin.Context) {
	streamID := uuid.New().String()
	stream := &Stream{
		ID:       streamID,
		Messages: make(chan []byte),
		Errors:   make(chan error),
	}
	h.streams[streamID] = stream

	c.JSON(http.StatusOK, gin.H{"stream_id": streamID})
}

// SendData accepts data for a specific stream and sends it to Kafka.
func (h *Handler) SendData(c *gin.Context) {
	streamID := c.Param("stream_id")
	_, exists := h.streams[streamID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	

	var data json.RawMessage
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	processedData := processData(data)
	if err := h.producer.SendMessage(streamID, processedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data"})
		return
	}

	if err := h.producer.SendMessage(streamID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Data sent successfully"})
}

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

// GetResults sets up a WebSocket connection 
func (h *Handler) GetResults(c *gin.Context) {
	streamID := c.Param("stream_id")
	stream, exists := h.streams[streamID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer ws.Close()

	go h.consumer.ConsumeMessages(streamID, stream.Messages, stream.Errors)

	for {
		select {
		case msg := <-stream.Messages:
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case err := <-stream.Errors:
			if err := ws.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
				return
			}
			return 
		}
	}
}