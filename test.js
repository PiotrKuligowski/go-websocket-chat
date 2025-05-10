const ws = new WebSocket('ws://localhost:8000'); 

function sendMessage(message) {
  ws.send(message);
}

ws.onmessage = function(event) {
  console.log("Received message from GO server: ", event["data"]);
};

sendMessage("test"); 