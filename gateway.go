package discordgo

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"time"
)

type GatewayConnection struct {
	Url                          string
	conn                         *websocket.Conn
	sessionID                    string
	lastSequence                 int
	heartbeatInterval            time.Duration
	lastHeartbeatAck             *time.Time
	eventHandler                 func(eventName string, eventData json.RawMessage)
	token                        string
	intents                      int
	active                       bool
	closed                       chan struct{}
	heartbeatChannel             chan int
	lastConnection               *time.Time
	reconnectionStartingWaitTime time.Duration
	logger                       *log.Logger
}

func ConnectToGateway(token string, intents int, handler func(eventName string, eventData json.RawMessage)) (*GatewayConnection, error) {
	logger := log.New(os.Stdout, "DiscordGo-Lite", 0)

	logger.Print("DiscordGo-Lite V0.1.4")

	var gatewayInfo = make(chan []byte)
	NewRestRequest().
		Method("GET").
		Token(token).
		Route("gateway").
		Route("bot").
		Callback(gatewayInfo).
		Enqueue()

	data := <-gatewayInfo

	var gatewayResponse = GetGatewayResponse{}
	err := json.Unmarshal(data, &gatewayResponse)

	if err != nil {
		return nil, err
	}

	c := GatewayConnection{
		Url:                          gatewayResponse.Url,
		token:                        token,
		intents:                      intents,
		eventHandler:                 handler,
		closed:                       make(chan struct{}),
		heartbeatChannel:             make(chan int),
		active:                       true,
		reconnectionStartingWaitTime: time.Second,
		logger:                       logger,
	}

	c.debug("Received Gateway URL: " + gatewayResponse.Url)

	go func() {
		defer close(c.closed)
		defer close(c.heartbeatChannel)
		for c.active {
			if c.lastConnection != nil {
				if time.Now().Sub(*c.lastConnection) < time.Minute { //if its already done a connection in the last minute
					c.debug("Delaying reconnection for " + c.reconnectionStartingWaitTime.String())
					time.Sleep(c.reconnectionStartingWaitTime) //wait for this duration
				} else { //if its been more than a minute since the last connection
					c.reconnectionStartingWaitTime = time.Second //reset the wait time
				}
			}
			now := time.Now()
			c.lastConnection = &now
			err := c.connect()
			if err != nil {
				c.error(&err)
			} else {
				c.listen()
			}
			c.reconnectionStartingWaitTime *= 2                     //double the wait time
			if c.reconnectionStartingWaitTime > (5 * time.Minute) { // if wait time longer than 5 minutes
				c.reconnectionStartingWaitTime = 5 * time.Minute //set it to 5 minutes
			}
		}
	}()

	return &c, nil

}

//opens a connection to the Gateway API and keeps it open with a heartbeat
func (c *GatewayConnection) connect() error {

	c.heartbeatChannel = make(chan int)

	c.debug("Connecting to Gateway")

	conn, _, err := websocket.DefaultDialer.Dial(c.Url, nil)
	if err != nil {
		return err
	}

	c.conn = conn

	response, err := c.read()

	if err != nil {
		return err
	}

	if response.Op != 10 {
		//return error expecting Hello payload
		return errors.New("expected Hello payload from server but received something else")
	}

	c.debug("Received Hello payload")

	var helloResponse HelloPayload
	err = json.Unmarshal(response.D, &helloResponse)

	c.heartbeatInterval = helloResponse.HeartbeatInterval
	c.lastHeartbeatAck = nil

	c.debug("Heartbeat interval: " + c.heartbeatInterval.String())

	if err != nil {
		return err
	}

	ticker := time.NewTicker(c.heartbeatInterval * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				if c.lastHeartbeatAck != nil {
					if time.Now().Sub(*c.lastHeartbeatAck) > (c.heartbeatInterval * time.Millisecond) {
						c.close(websocket.CloseAbnormalClosure)
						c.debug("No acknowledgement between heartbeats")
						return
					}
				}
				c.heartbeat()
			case <-c.heartbeatChannel:
				return
			}
		}
	}()

	var authError error
	if c.lastSequence != 0 {
		authError = c.resume()
	} else {
		authError = c.identify()
	}

	if authError != nil {
		return authError
	}

	return nil

}

func (c *GatewayConnection) Close() {
	c.active = false
	c.close(websocket.CloseNormalClosure)
}

func (c *GatewayConnection) close(closeCode int) {

	c.debug("Closing Gateway connection")

	close(c.heartbeatChannel)

	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, ""))

	if err != nil {
		err := c.conn.Close()
		c.error(&err)
	}

	select {
	case <-c.closed:
	case <-time.After(time.Second):
	}

	c.debug("Gateway connection closed")

}

func (c *GatewayConnection) identify() error {
	c.debug("Identifying")
	properties := IdentityConnectionProperties{
		Os:      "Go",
		Browser: "andrew-wilson/discord.go",
		Device:  "andrew-wilson/discord.go",
	}
	identity := IdentityPayload{
		Token:      c.token,
		Properties: properties,
		Intents:    c.intents,
	}

	raw, _ := json.Marshal(identity)

	payload := GatewayPayload{
		Op: 2,
		D:  raw,
	}

	return c.write(payload)

}

func (c *GatewayConnection) resume() error {
	c.debug("Resuming")
	resume := ResumePayload{
		Token:     c.token,
		SessionID: c.sessionID,
		Sequence:  c.lastSequence,
	}

	raw, _ := json.Marshal(resume)

	payload := GatewayPayload{
		Op: 6,
		D:  raw,
	}

	return c.write(payload)
}

func (c *GatewayConnection) heartbeat() {
	raw, _ := json.Marshal(c.lastSequence)
	payload := GatewayPayload{
		Op: 1,
		D:  raw,
	}
	err := c.write(payload)

	if err != nil {
		customErr := errors.New("failed to emit heartbeat")
		c.error(&customErr)
		c.error(&err)
		c.close(websocket.CloseAbnormalClosure)
	} else {
		c.debug("Heartbeat emitted")
	}
}

func (c *GatewayConnection) invalidateSession() {
	c.debug("Invalidating session. Will identify again for next connection")
	c.sessionID = ""
	c.lastSequence = 0
}

func (c *GatewayConnection) listen() {
	c.debug("Listening to Gateway")
	for {
		payload, err := c.read()

		if err != nil {
			return
		}

		switch payload.Op {
		case 0:
			c.debug("Received event dispatch")
			if payload.T == "READY" {
				var ready ReadyEvent
				err := json.Unmarshal(payload.D, &ready)

				if err != nil {
					c.error(&err)
					c.close(websocket.CloseAbnormalClosure)
				}

				c.sessionID = ready.SessionID
				c.debug("Received READY event with session ID: " + c.sessionID)
			}

			c.lastSequence = payload.S
			c.eventHandler(payload.T, payload.D)
		case 1:
			c.debug("Heartbeat requested from server")
			c.heartbeat()
		case 7:
			//should attempt to reconnect and resume immediately
			c.debug("Server requested a reconnection")
			c.close(websocket.CloseAbnormalClosure)
		case 9:
			//session is invalid. should reconnect and identify/resume accordingly
			c.debug("Server invalidated session")
			var resume bool
			err := json.Unmarshal(payload.D, &resume)

			c.error(&err)

			c.close(websocket.CloseAbnormalClosure)

			if !resume {
				c.invalidateSession()
			}

		case 11:
			//heartbeat ack  received.
			c.debug("Server acknowledged heartbeat")
			now := time.Now()
			c.lastHeartbeatAck = &now
		}
	}
}

func (c *GatewayConnection) write(payload GatewayPayload) error {
	err := c.conn.WriteJSON(payload)
	return err
}

func (c *GatewayConnection) read() (GatewayPayload, error) {
	var payload GatewayPayload
	err := c.conn.ReadJSON(&payload)
	return payload, err
}

func (c *GatewayConnection) debug(text interface{}) {
	c.logger.Print(text)
	//logp.Info("%v", text)
}

func (c *GatewayConnection) error(err *error) {
	if *err != nil {
		c.logger.Fatal(err)
		//logp.Err("%v", err)
	}
}
