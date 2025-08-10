import { sendMessage, handleMessage, MESSAGE_TYPE } from './messages.js';
import { getClientId } from './auth.js';
import { useSessionStore } from '@/stores/sessionStore.js';

let conn = null;
let reconnectAttempts = 0
let reconnectTimer = null
let isManuallyClosed = false

const MAX_RECONNECT_ATTEMPTS = 10
const RECONNECT_DELAY_BASE = 1000 // 1s, exponential backoff

const wsProtocol = location.protocol === 'https:' ? 'wss' : 'ws';

const wsUrl = wsProtocol + '://' + location.host + '/ws';

function connectWebSocket () {

  const sessionStore = useSessionStore();

  if (conn || isManuallyClosed) return;

  conn = new WebSocket(wsUrl);

  conn.addEventListener('message', (event) => {
    try {
      const data = JSON.parse(event.data);
      handleMessage(data);
      if (data.Destination == getClientId()) {
        sessionStore.saveSession(data);
      }
    } catch (err) {
      console.error("❌ Failed to parse message:", event.data, err);
    }
  });

  conn.onopen = function (evt) {
    console.log('WebSocket connected');
    reconnectAttempts = 0;
    sendMessage(MESSAGE_TYPE.COLLECTOR_ONLINE, null, "web_client connected", null);

  };

  conn.onclose = () => {
    console.log('WebSocket disconnected');
    conn = null;
    if (!isManuallyClosed) attemptReconnect()
  }

  conn.onerror = (err) => {
    console.error('WebSocket error', err);
  }
}

function getConn() {
  return conn;
}

const attemptReconnect = () => {
  if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
    console.error('Max reconnect attempts reached.')
    return
  }

  const delay = RECONNECT_DELAY_BASE * Math.pow(2, reconnectAttempts)
  console.log(`Attempting to reconnect in ${delay / 1000}s...`)

  reconnectAttempts++

  reconnectTimer = setTimeout(() => {
    connectWebSocket()
  }, delay)
}

const disconnectWebSocket = () => {
  isManuallyClosed = true
  clearTimeout(reconnectTimer)
  if (conn) {
    conn.close()
    conn = null
  }
}
export default { connectWebSocket , getConn, disconnectWebSocket }  