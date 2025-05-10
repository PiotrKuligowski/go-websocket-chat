# Websocket chat

There are many ways to test this websocket server, one of them is to use node from command line or chrome devtools console with the following snippet.

Setup connection and sending function
```js
const ws = new WebSocket('ws://localhost:8000/ws'); 

function send(cid, pid, data) {
  ws.send(JSON.stringify({
    "ChatId": cid,
    "PlayerId": pid,
    "Data": data
  }));
}

ws.onmessage = function(event) {
   console.log("Received message from GO server: ", event.data);
};

```

To send message to global (or any) channel:
```js
send("global", "100", "This is client 1 speaking");
```

For 2 or more players to communicate, game client will send initial message to be added to channel:
```js
send("private", "101", "Client 101 connected!");
```
Next, players will send messages using same syntax.