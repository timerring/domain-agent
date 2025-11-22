package agent

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要严格检查
	},
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket client connected")

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// 处理消息并流式返回
		response := map[string]interface{}{
			"type":    "message",
			"content": "收到消息: " + msg["message"].(string),
		}

		if err := conn.WriteJSON(response); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}
