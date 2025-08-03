import { v4 as uuidv4 } from 'uuid';
import axios from 'axios';

export function getClientId() {
  let clientId = localStorage.getItem("client_id");
  if (!clientId) {
    clientId = "web_client_" + uuidv4(); // Secure, random
    localStorage.setItem("client_id", clientId);
  }
  return clientId;
}

export async function getUser(opts = {}) {
  return await axios.get("/api/user", opts)
}

export async function logOut(opts = {}) {
  return await axios.get("/api/logout", opts)
}

export async function logIn(username,password,opts = {}) {
    const formData = new FormData()
    formData.append('username', username)
    formData.append('password', password)
    return await axios.post('/api/login',formData)
}