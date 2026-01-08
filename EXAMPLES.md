# Kelutral Integration Examples

## Example 1: Unity C# Client

```csharp
using UnityEngine;
using WebSocketSharp;
using Newtonsoft.Json;
using System.Collections.Generic;

public class KelutralClient : MonoBehaviour
{
    private WebSocket ws;
    private long frameNumber = 0;
    private string playerId;

    [System.Serializable]
    public class FrameData
    {
        public string player_id;
        public long frame_num;
        public long timestamp;
        public Dictionary<string, object> actions;
    }

    void Start()
    {
        playerId = "player_" + SystemInfo.deviceUniqueIdentifier;
        Connect();
    }

    void Connect()
    {
        ws = new WebSocket("ws://localhost:8080/ws?player_id=" + playerId);
        
        ws.OnOpen += (sender, e) =>
        {
            Debug.Log("Connected to Kelutral server");
        };

        ws.OnMessage += (sender, e) =>
        {
            var data = JsonConvert.DeserializeObject<FrameData>(e.Data);
            HandleFrameUpdate(data);
        };

        ws.OnError += (sender, e) =>
        {
            Debug.LogError("WebSocket Error: " + e.Message);
        };

        ws.Connect();
    }

    void FixedUpdate()
    {
        // Send player state every frame (~60 FPS)
        if (ws != null && ws.ReadyState == WebSocketState.Open)
        {
            SendFrame(new Dictionary<string, object>
            {
                { "position_x", transform.position.x },
                { "position_y", transform.position.y },
                { "position_z", transform.position.z },
                { "rotation", transform.rotation.eulerAngles.y },
                { "velocity", GetComponent<Rigidbody>().velocity.magnitude }
            });
        }
    }

    void SendFrame(Dictionary<string, object> actions)
    {
        var frameData = new FrameData
        {
            player_id = playerId,
            frame_num = frameNumber++,
            actions = actions
        };

        string json = JsonConvert.SerializeObject(frameData);
        ws.Send(json);
    }

    void HandleFrameUpdate(FrameData data)
    {
        // Apply other players' state updates
        if (data.player_id != playerId)
        {
            // Update remote player positions, etc.
            Debug.Log($"Player {data.player_id} frame {data.frame_num}");
        }
    }

    void OnDestroy()
    {
        if (ws != null)
        {
            ws.Close();
        }
    }
}
```

## Example 2: Godot GDScript Client

```gdscript
extends Node

var ws = WebSocketClient.new()
var player_id = ""
var frame_number = 0

func _ready():
    player_id = "player_" + str(OS.get_unique_id())
    connect_to_server()

func connect_to_server():
    var url = "ws://localhost:8080/ws?player_id=" + player_id
    var err = ws.connect_to_url(url)
    
    if err != OK:
        print("Unable to connect")
        return
    
    ws.connect("connection_established", self, "_on_connection_established")
    ws.connect("data_received", self, "_on_data_received")
    ws.connect("connection_error", self, "_on_connection_error")

func _on_connection_established(protocol):
    print("Connected to Kelutral server")

func _on_data_received():
    var data = JSON.parse(ws.get_peer(1).get_packet().get_string_from_utf8())
    if data.error == OK:
        handle_frame_update(data.result)

func _on_connection_error():
    print("Connection error")

func _physics_process(delta):
    ws.poll()
    
    # Send frame data at 60 FPS
    if ws.get_connection_status() == WebSocketClient.CONNECTION_CONNECTED:
        send_frame({
            "position_x": global_position.x,
            "position_y": global_position.y,
            "velocity": velocity.length(),
            "action": current_action
        })

func send_frame(actions):
    var frame_data = {
        "player_id": player_id,
        "frame_num": frame_number,
        "actions": actions
    }
    frame_number += 1
    
    ws.get_peer(1).put_packet(JSON.print(frame_data).to_utf8())

func handle_frame_update(data):
    # Update game state based on received frame data
    if data.player_id != player_id:
        print("Received frame from ", data.player_id)
```

## Example 3: JavaScript/Phaser Client

```javascript
class GameScene extends Phaser.Scene {
    constructor() {
        super({ key: 'GameScene' });
        this.ws = null;
        this.playerId = 'player_' + Date.now();
        this.frameNumber = 0;
        this.remotePlayers = {};
    }

    create() {
        this.connectToServer();
        
        // Create player sprite
        this.player = this.add.circle(400, 300, 20, 0x00ff00);
        this.physics.add.existing(this.player);
        
        // Setup input
        this.cursors = this.input.keyboard.createCursorKeys();
    }

    connectToServer() {
        this.ws = new WebSocket(`ws://localhost:8080/ws?player_id=${this.playerId}`);
        
        this.ws.onopen = () => {
            console.log('Connected to Kelutral');
        };
        
        this.ws.onmessage = (event) => {
            const frameData = JSON.parse(event.data);
            this.handleFrameUpdate(frameData);
        };
    }

    update() {
        // Handle input
        const velocity = { x: 0, y: 0 };
        
        if (this.cursors.left.isDown) velocity.x = -200;
        if (this.cursors.right.isDown) velocity.x = 200;
        if (this.cursors.up.isDown) velocity.y = -200;
        if (this.cursors.down.isDown) velocity.y = 200;
        
        this.player.body.setVelocity(velocity.x, velocity.y);
        
        // Send frame every update (~60 FPS)
        this.sendFrame({
            x: this.player.x,
            y: this.player.y,
            velocity_x: velocity.x,
            velocity_y: velocity.y
        });
    }

    sendFrame(actions) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const frameData = {
                player_id: this.playerId,
                frame_num: this.frameNumber++,
                actions: actions
            };
            
            this.ws.send(JSON.stringify(frameData));
        }
    }

    handleFrameUpdate(frameData) {
        if (frameData.player_id === this.playerId) return;
        
        // Create or update remote player
        if (!this.remotePlayers[frameData.player_id]) {
            this.remotePlayers[frameData.player_id] = 
                this.add.circle(0, 0, 20, 0xff0000);
        }
        
        const remotePlayer = this.remotePlayers[frameData.player_id];
        const actions = frameData.actions;
        
        // Interpolate position
        remotePlayer.x = actions.x;
        remotePlayer.y = actions.y;
    }
}
```

## Example 4: Python Client (for testing/bots)

```python
import asyncio
import websockets
import json
import time

class KelutralClient:
    def __init__(self, server_url, player_id):
        self.server_url = server_url
        self.player_id = player_id
        self.frame_number = 0
        self.ws = None

    async def connect(self):
        uri = f"{self.server_url}?player_id={self.player_id}"
        async with websockets.connect(uri) as websocket:
            self.ws = websocket
            print(f"Connected as {self.player_id}")
            
            # Start receiving messages
            receive_task = asyncio.create_task(self.receive_messages())
            
            # Send frames at 60 FPS
            while True:
                await self.send_frame({
                    "position": {"x": 0, "y": 0},
                    "action": "idle"
                })
                await asyncio.sleep(1/60)  # 60 FPS

    async def send_frame(self, actions):
        if self.ws:
            frame_data = {
                "player_id": self.player_id,
                "frame_num": self.frame_number,
                "actions": actions
            }
            self.frame_number += 1
            
            await self.ws.send(json.dumps(frame_data))

    async def receive_messages(self):
        async for message in self.ws:
            data = json.loads(message)
            self.handle_frame_update(data)

    def handle_frame_update(self, frame_data):
        print(f"Received frame from {frame_data['player_id']}: {frame_data['frame_num']}")

# Usage
async def main():
    client = KelutralClient("ws://localhost:8080/ws", "bot_player_1")
    await client.connect()

if __name__ == "__main__":
    asyncio.run(main())
```

## Load Testing

You can use the Python client for load testing:

```python
import asyncio
import websockets
import json

async def simulate_player(player_id, duration=60):
    uri = f"ws://localhost:8080/ws?player_id={player_id}"
    async with websockets.connect(uri) as ws:
        frame_num = 0
        start_time = asyncio.get_event_loop().time()
        
        while asyncio.get_event_loop().time() - start_time < duration:
            await ws.send(json.dumps({
                "player_id": player_id,
                "frame_num": frame_num,
                "actions": {"x": frame_num % 100, "y": frame_num % 100}
            }))
            frame_num += 1
            await asyncio.sleep(1/60)

async def load_test(num_players=100):
    tasks = [simulate_player(f"player_{i}") for i in range(num_players)]
    await asyncio.gather(*tasks)

# Run: asyncio.run(load_test(100))
```

## Deployment Recommendations

### Production Configuration

```yaml
# docker-compose.prod.yml
services:
  kelutral:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
    restart: always
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '1'
          memory: 512M
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 3
```

### Nginx Load Balancer

```nginx
upstream kelutral_backend {
    least_conn;
    server kelutral1:8080;
    server kelutral2:8080;
    server kelutral3:8080;
}

server {
    listen 80;
    server_name game.example.com;

    location /ws {
        proxy_pass http://kelutral_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /health {
        proxy_pass http://kelutral_backend;
    }
}
```
