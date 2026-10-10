package handlers

import (
	"fmt"
	"net/http"
	"pos-backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

type WebhookPayload struct {
	MerchantID string  `json:"merchant_id"`
	OutletID   string  `json:"outlet_id"`
	Amount     float64 `json:"amount"`
	Message    string  `json:"message"`
}

func HandleQrisWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dalam kasus nyata, kita akan validasi signature/hmac dari provider di sini.
	
	// Broadcast event sukses ke frontend!
	msg := payload.Message
	if msg == "" {
		msg = fmt.Sprintf("Pembayaran QRIS Rp %.0f Berhasil!", payload.Amount)
	}

	websocket.DefaultHub.BroadcastPaymentSuccess(payload.MerchantID, payload.OutletID, payload.Amount, msg)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Notification broadcasted"})
}
