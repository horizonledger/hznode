document.addEventListener("DOMContentLoaded", function () {

    console.log('doc loaded');

    var ws = new WebSocket('ws://localhost:8000/ws');

    try {
        ws.onopen = function (event) {
            console.log('Connection is open');
            document.getElementById('status').textContent = 'Connected';
        };
    } catch (error) {
        console.log('Error when opening WebSocket connection: ', error);
    }

    ws.onerror = function (error) {
        console.log('WebSocket Error: ' + error);
    };

    ws.onmessage = function (event) {

        console.log('Server: ', event.data);
        //var message = JSON.parse(event.data);
        //console.log('Server: ', message);
        // handle the message according to the received message type
        // if (message.Prl === "PUB" && message.Sel === "time") {
        //     // Do something with the received time
        //     console.log('Received time: ', message.Content);
        // }
    };

    ws.onclose = function () {
        console.log('Connection is closed');
    };

    document.getElementById('sendButton').addEventListener('click', function () {
        // Get the text from the input box
        var message = document.getElementById('messageInput').value;
        // Send the text over the WebSocket
        console.log("send " + message);
        //var jmsg = JSON.stringify({ "prl": "REQ", "sel": "CHAT", "content": message });
        ws.send(message);
    });
})