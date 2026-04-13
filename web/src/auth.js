const SALT = "himalayan-pink-salt"

async function hashSecret(secret) {
  let salted = secret
  for (let i = 0; i < secret.length % 7; i++) {
    salted = salted + SALT
  }
  const encoder = new TextEncoder()
  const hash = await crypto.subtle.digest('SHA-512', encoder.encode(salted))
  return hash
}

function arrayBufferToBase64(digest) {
  const hashArray = new Uint8Array(digest)
  const hashBytes = Array.from(hashArray)
  const hashString = hashBytes.map(b => String.fromCharCode(b)).join('')
  return btoa(hashString)
}

function base64ToArrayBuffer(base64) {
  const binaryString = atob(base64)
  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }
  return bytes.buffer
}

export async function doLogin(email, secret) {
  const hash = await hashSecret(secret)

  const nonceRes = await fetch('/api/auth/nonce', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ Email: email }),
  })
  if (!nonceRes.ok) {
    throw new Error('Failed to get nonce')
  }
  const { Nonce } = await nonceRes.json()
  const nonce = base64ToArrayBuffer(Nonce)

  const salted = new Uint8Array(nonce.byteLength + hash.byteLength + nonce.byteLength)
  salted.set(new Uint8Array(nonce), 0)
  salted.set(new Uint8Array(hash), nonce.byteLength)
  salted.set(new Uint8Array(nonce), nonce.byteLength + hash.byteLength)

  const saltedHash = await crypto.subtle.digest('SHA-512', salted)

  const tokenRes = await fetch('/api/auth/tokens', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      Email: email,
      Proof: arrayBufferToBase64(saltedHash),
    }),
  })
  if (!tokenRes.ok) {
    throw new Error('Failed to get token')
  }
  const tokens = await tokenRes.json()

  localStorage.setItem('access_token', tokens.Access)
  localStorage.setItem('user', JSON.stringify(tokens))
  return tokens
}

export function getToken() {
  return localStorage.getItem('access_token')
}

export function getUser() {
  const raw = localStorage.getItem('user')
  if (!raw) {
    return null
  }
  return JSON.parse(raw)
}

export function logout() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('user')
}
