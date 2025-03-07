document.addEventListener('DOMContentLoaded', () => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    const messagesDiv = document.getElementById('messages');
    const chatForm = document.getElementById('chatForm');
    const messageInput = document.getElementById('messageInput');

    // Display incoming messages
    ws.onmessage = (event) => {
        const message = document.createElement('p');
        message.textContent = event.data;
        messagesDiv.appendChild(message);
        messagesDiv.scrollTop = messagesDiv.scrollHeight; // Auto-scroll to bottom
    };

    // Send messages when the form is submitted
    chatForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const message = messageInput.value.trim();
        if (message) {
            ws.send(message);
            messageInput.value = '';
        }
    });
});