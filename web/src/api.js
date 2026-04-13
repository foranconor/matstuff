import { getToken, logout } from './auth'

async function request(method, path, body) {
  const opts = { method, headers: {} }

  const token = getToken()
  if (token) {
    opts.headers['Authorization'] = 'Bearer ' + token
  }

  if (body !== undefined) {
    opts.headers['Content-Type'] = 'application/json'
    opts.body = JSON.stringify(body)
  }

  const r = await fetch('/api' + path, opts)

  if (r.status === 401 || r.status === 403) {
    logout()
    window.location.href = '/login'
    return
  }

  if (!r.ok) {
    throw new Error(await r.text())
  }

  if (r.status === 204) {
    return null
  }

  return r.json()
}

export const get  = (path)        => request('GET',    path)
export const put  = (path, body)  => request('PUT',    path, body)
export const post = (path, body)  => request('POST',   path, body)
export const del  = (path)        => request('DELETE', path)
