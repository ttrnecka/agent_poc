import { v4 as uuidv4 } from 'uuid';

export function getClientId() {
  let clientId = localStorage.getItem("client_id");
  if (!clientId) {
    clientId = "web_client_" + uuidv4(); // Secure, random
    localStorage.setItem("client_id", clientId);
  }
  return clientId;
}