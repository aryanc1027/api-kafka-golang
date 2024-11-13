package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    //"github.com/aryanc1027/api-kafka-golang/internal/kafka"
)

type MockProducer struct {
    mock.Mock
}

func (m *MockProducer) SendMessage(topic string, message []byte) error {
    args := m.Called(topic, message)
    return args.Error(0)
}

func (m *MockProducer) Close() error {
    args := m.Called()
    return args.Error(0)
}

type MockConsumer struct {
    mock.Mock
}

func (m *MockConsumer) ConsumeMessages(topic string, messages chan<- []byte, errors chan<- error) {
    m.Called(topic, messages, errors)
}

func (m *MockConsumer) Close() error {
    args := m.Called()
    return args.Error(0)
}

func TestStartStream(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockProducer := new(MockProducer)
    mockConsumer := new(MockConsumer)
    handler := NewHandler(mockProducer, mockConsumer)

    router := gin.Default()
    router.POST("/stream/start", handler.StartStream)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/stream/start", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)

    var response map[string]string
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Contains(t, response, "stream_id")
}

func TestSendData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	handler := NewHandler(mockProducer, mockConsumer)

	router := gin.Default()
	router.POST("/stream/:stream_id/send", handler.SendData)


	streamID := "test-stream"
	handler.streams[streamID] = &Stream{
		ID:       streamID,
		Messages: make(chan []byte),
		Errors:   make(chan error),
	}

	mockProducer.On("SendMessage", streamID, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	data := []byte(`{"test": "data"}`)
	req, _ := http.NewRequest("POST", "/stream/"+streamID+"/send", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Data sent successfully", response["status"])

	mockProducer.AssertExpectations(t)
	
}

func TestGetResults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)
	handler := NewHandler(mockProducer, mockConsumer)

	router := gin.Default()
	router.GET("/stream/:stream_id/results", handler.GetResults)


	streamID := "test-stream"
	handler.streams[streamID] = &Stream{
		ID:       streamID,
		Messages: make(chan []byte),
		Errors:   make(chan error),
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stream/"+streamID+"/results", nil)
	go router.ServeHTTP(w, req)


	assert.Equal(t, 500, w.Code)
}