// static/ws.js

// Establish a WebSocket connection
const socket = new WebSocket("ws://localhost:8080/ws");

// Handle incoming messages
socket.onmessage = function(event) {
    const notification = event.data;
    displayNotification(notification);
};

// Display notification (this function can be customized)
function displayNotification(message) {
    const notificationArea = document.getElementById('notification-area');
    const notificationElement = document.createElement('div');
    notificationElement.className = 'notification';
    notificationElement.innerText = message;
    notificationArea.appendChild(notificationElement);

    // Optional: Remove notification after a few seconds
    setTimeout(() => {
        notificationArea.removeChild(notificationElement);
    }, 5000);
}