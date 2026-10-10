package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for UMKM PWA
	},
}

type Client struct {
	Conn       *websocket.Conn
	MerchantID string
	OutletID   string
	UserID     string
}

type Hub struct {
	sync.RWMutex
	Clients map[*Client]bool
}

var DefaultHub = &Hub{
	Clients: make(map[*Client]bool),
}

func (h *Hub) Register(c *Client) {
	h.Lock()
	h.Clients[c] = true
	h.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.Lock()
	if _, ok := h.Clients[c]; ok {
		delete(h.Clients, c)
		c.Conn.Close()
	}
	h.Unlock()
}

// BroadcastPaymentSuccess mengirim event sukses bayar ke kasir spesifik (MerchantID + OutletID)
func (h *Hub) BroadcastPaymentSuccess(merchantID string, outletID string, amount float64, message string) {
	h.RLock()
	defer h.RUnlock()

	event := map[string]interface{}{
		"type":    "PAYMENT_SUCCESS",
		"amount":  amount,
		"message": message,
		"time":    time.Now().Unix(),
	}
	eventJSON, _ := json.Marshal(event)

	for client := range h.Clients {
		if client.MerchantID == merchantID && (client.OutletID == outletID || client.OutletID == "" || outletID == "") {
			err := client.Conn.WriteMessage(websocket.TextMessage, eventJSON)
			if err != nil {
				log.Printf("WS error: %v", err)
				client.Conn.Close()
			}
		}
	}
}

var jwtSecret = []byte("pos_secret_key_2026")

func HandleWebSocket(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	merchantID := fmt.Sprintf("%v", claims["merchant_id"])
	outletID := fmt.Sprintf("%v", claims["outlet_id"])
	if outletID == "<nil>" {
		outletID = ""
	}
	userID := fmt.Sprintf("%v", claims["user_id"])

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WS Upgrade Error:", err)
		return
	}

	client := &Client{
		Conn:       conn,
		MerchantID: merchantID,
		OutletID:   outletID,
		UserID:     userID,
	}

	DefaultHub.Register(client)

	// Keep alive & wait for disconnect
	go func() {
		defer DefaultHub.Unregister(client)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
