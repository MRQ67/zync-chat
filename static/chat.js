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
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?token=${token}`;
    ws = new WebSocket(wsUrl);
    ws.onopen = () => console.log("WebSocket connected");
    ws.onmessage = handleMessage;
    ws.onerror = () => console.error("WebSocket error");
    ws.onclose = () => console.log("WebSocket closed");
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