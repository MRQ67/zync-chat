let ws = null;
const authDiv = document.getElementById('auth');
const chatDiv = document.getElementById('chat');
const messagesDiv = document.getElementById('messages');

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
        const message = document.createElement('p');
        message.textContent = event.data;
        const messagesDiv = document.getElementById('messages');
        messageDiv.classList.add('message', data.isSelf ? 'self' : 'other');
        messageDiv.innerHTML = `
            <div class="sender">${data.username}</div>
            <div>${data.message}</div>
            <div class="timestamp">${new Date().toLocaleTimeString()}</div>
        `;
        messagesDiv.appendChild(message);
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
        input.value = '';
    }
});