package utils

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	ctrl "sigs.k8s.io/controller-runtime"
)

// // UpgradeToWebsocket attempts to upgrade connection from HTTP(S) to WS(S)
// func UpgradeToWebsocket(request *restful.Request, response *restful.Response) (*websocket.Conn, error) {
// 	var writer http.ResponseWriter = response
// 	log.Debug().Msg("Upgrading connection to websocket...")
// 	// Handles writing error to response
// 	upgrader := websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 4096,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}
// 	connection, err := upgrader.Upgrade(writer, request.Request, nil)
// 	return connection, err
// }

// WriteOnlyWebsocket discards text messages from the peer connection
func WriteOnlyWebsocket(connection *websocket.Conn, b *Broadcaster) {
	// The underlying connection is never closed so this cannot error
	subscriber, _ := b.Subscribe()
	go readControl(connection, b, subscriber)
	write(connection, subscriber)
}

// ping over the socket with a given deadline; if there's an error, close
func writePing(connection *websocket.Conn, deadline time.Time) {
	if err := connection.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
		ReportClosing(connection)
	}
}

// readControl will unsubscribe on connection failures
func readControl(connection *websocket.Conn, b *Broadcaster, s *Subscriber) {
	// Connection lifecycle handler
	connection.SetPongHandler(func(string) error {
		// Extend deadline to prevent expiration
		deadline := time.Now().Add(time.Second * 2)
		connection.SetReadDeadline(deadline)
		// Cut down on ping/pong traffic
		time.Sleep(time.Second)
		// Ellicit another ping
		writePing(connection, deadline)
		return nil
	})
	initialDeadline := time.Now().Add(time.Second)
	connection.SetReadDeadline(initialDeadline)
	// Kick off cycle
	writePing(connection, initialDeadline)
	for {
		// Connection has either decayed or close has been requested from server side
		if _, _, err := connection.ReadMessage(); err != nil {
			ctrl.Log.Error(err, "websocket connection to client lost")
			b.Unsubscribe(s)
			return
		}
	}
}

// ReportClosing sends close to client then closes connection
func ReportClosing(connection *websocket.Conn) {
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	connection.Close()
}

// Send data over the connection using the subscriber channel, if there's a failure we return
func write(connection *websocket.Conn, subscriber *Subscriber) {
	subChan := subscriber.SubChan()
	unsubChan := subscriber.UnsubChan()
	for {
		select {
		case socketData := <-subChan:
			if !websocketSend(connection, socketData) {
				return
			}
		case <-unsubChan:
			return
		}
	}
}

// Returns whether successful or not, closes connection on failures
func websocketSend(connection *websocket.Conn, data SocketData) bool {
	payload, err := json.Marshal(data)
	if err != nil {
		ctrl.Log.Error(err, "failed to marshal status")
		ReportClosing(connection)
		return false
	}
	if err := connection.WriteMessage(websocket.TextMessage, payload); err != nil {
		ctrl.Log.Error(err, "could not write the message to the websocket client connection")
		ReportClosing(connection)
		return false
	}
	return true
}
