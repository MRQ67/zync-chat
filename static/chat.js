let ws = null;
const authDiv = document.getElementById('auth');
const chatDiv = document.getElementById('chat');
const messagesDiv = document.getElementById('messages');
// const Username = document.getElementById('username');
let currentUsername = null; // To track the logged-in user

// Handle registration
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('regUsername')?.value.trim() || "";
    const password = document.getElementById('regPassword')?.value.trim() || "";

    console.log("Sending registration request:", { username, password });

    if (!username || !password) {
        alert("Username and password cannot be empty.");
        return;
    }

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
        // const currentUsername = data.username; // Store the logged-in username
        localStorage.setItem('token', token); // Store token
        authDiv.style.display = 'none';
        chatDiv.style.display = 'block';
        document.getElementById('username').textContent = document.getElementById('loginUsername').value;
        connectWebSocket(token);
        e.target.reset();
    } else {
        alert('Login failed. Check your credentials.');
    }
});

// Connect to WebSocket with token
function connectWebSocket(token) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host; // e.g., localhost:8080 or zync-chat.onrender.com
    const wsUrl = `${protocol}//${host}/ws?token=${token}`;
    console.log("Attempting to connect to:", wsUrl); // Debug the URL
    ws = new WebSocket(wsUrl);
    ws.onopen = () => console.log("WebSocket connected");
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        const messageDiv = document.createElement('div');
        const isSelf = data.username === currentUsername;
        messageDiv.classList.add('message', isSelf ? 'self' : 'other');
        messageDiv.innerHTML = `    
            <div class="sender">${data.username}</div>
            <div>${data.message}</div>
            <div class="timestamp">${new Date(data.timestamp).toLocaleTimeString()}</div>
        `;
        messagesDiv.appendChild(messageDiv);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    };
    ws.onerror = (error) => console.error("WebSocket error:", error);
    ws.onclose = () => console.log("WebSocket closed");
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