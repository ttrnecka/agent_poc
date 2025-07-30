import { defineStore } from 'pinia';
import { reactive, toRefs } from 'vue';

export const useSessionStore = defineStore('session', () => {
  // reactive state object to hold sessionId -> data
  const sessions = reactive({});

  // action: handle incoming websocket data
  function saveSession(data) {
    console.log('📥 WebSocket message received:', data);
    if (data.Session) {
      sessions[data.Session] = data;
    }
  }

  // action: clear specific session data
  function clearSession(sessionId) {
    delete sessions[sessionId];
  }

  // expose reactive state and actions (toRefs to destructure if needed)
  return {
    sessions,
    saveSession,
    clearSession,
  };
});