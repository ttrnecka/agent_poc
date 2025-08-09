import { getClientId } from './auth.js'
import ws from '@/services/websocket'
import { useApiStore } from '@/stores/apiStore';
import { v4 as uuidv4 } from 'uuid';

export const MESSAGE_TYPE = {
  'ONLINE': 1,
  'OFFLINE': 2,
  'REFRESH': 10,
  'RUN': 20,
  'RUNNING': 21,
  'FINISHED_OK': 22,
  'FINISHED_ERR': 23,
  'DATA': 24
};

export function sendMessage(type, destinationId, text, session = uuidv4()) {
    const clientId = getClientId(); 
    const msg = {
        Type: type,
        Source: clientId,
        Destination: destinationId,
        Text: text,
        Session: session
    };
    ws.getConn().send(JSON.stringify(msg));
    return session; // Return the session ID for tracking
  }

export function handleMessage(msg) {
  const apiStore = useApiStore();
  
  console.log('📥 Message:', msg);
  if (msg.Source != getClientId()) {
    if (msg.Type == MESSAGE_TYPE.ONLINE) {
      apiStore.updateCollectorStatus(msg.Source,"ONLINE")
    }
    if (msg.Type == MESSAGE_TYPE.OFFLINE) {
      apiStore.updateCollectorStatus(msg.Source,"OFFLINE")
    }
  }
}
