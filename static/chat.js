let ws = null;
const authDiv = document.getElementById('auth');
const chatDiv = document.getElementById('chat');
const messagesDiv = document.getElementById('messages');
let currentUsername = null; // To track the logged-in user

// Handle registration
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('regUsername').value;
    const password = document.getElementById('regPassword').value;

    const response = await fetch('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        alert('Registration successful, please log in.');
        e.target.reset();
    } else {
        alert('Registration failed. Username may already be taken.');
    }
});

// Handle login
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    const response = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    if (response.ok) {
        const data = await response.json();
        const token = data.token;
        currentUsername = username; // Store the logged-in username
        localStorage.setItem('token', token); // Store token
        authDiv.style.display = 'none';
        chatDiv.style.display = 'block';
        connectWebSocket(token);
        e.target.reset();
    } else {
        alert('Login failed. Check your credentials.');
    }
});

// Connect to WebSocket with token
function connectWebSocket(token) {
    ws = new WebSocket(`wss://${window.location.host}/ws?token=${token}`);
    ws.onopen = () => {
        console.log("WebSocket connected");
    };
    ws.onmessage = (event) => {
        // Parse the incoming message (assuming JSON format from the server)
        const data = JSON.parse(event.data);
        const messageDiv = document.createElement('div'); // Use div for better styling control
        const isSelf = data.username === currentUsername; // Check if the message is from the current user
        messageDiv.classList.add('message', isSelf ? 'self' : 'other');
        messageDiv.innerHTML = `
            <div class="sender">${data.username}</div>
            <div>${data.message}</div>
            <div class="timestamp">${new Date(data.timestamp).toLocaleTimeString()}</div>
        `;
        messagesDiv.appendChild(messageDiv);
        messagesDiv.scrollTop = messagesDiv.scrollHeight; // Auto-scroll to the latest message
    };
    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };
    ws.onclose = () => {
        console.log("WebSocket closed");
    };
}

// Handle message sending
document.getElementById('messageForm').addEventListener('submit', (e) => {
    e.preventDefault();
    const input = document.getElementById('messageInput');
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(input.value);
        input.value = ''; // Clear input after sending
        console.log("Message sent"); // Optional feedback
    } else {
        console.error("WebSocket is not open. Cannot send message.");
    }
});