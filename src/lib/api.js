import { getAuth } from './auth-storage'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080/api'
let unauthorizedHandler = () => {}

export function setUnauthorizedHandler(handler) {
  unauthorizedHandler = handler
}

export function normalizeApiError(status, payload = {}) {
  return {
    status,
    message: payload.message ?? (status === 0 ? 'Tidak dapat terhubung ke server.' : 'Request failed'),
    errors: payload.errors ?? {},
  }
}

export async function apiRequest(path, options = {}) {
  const auth = getAuth()
  let response
  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(auth?.token ? { Authorization: `Bearer ${auth.token}` } : {}),
        ...options.headers,
      },
    })
  } catch {
    throw normalizeApiError(0)
  }

  const payload = await response.json().catch(() => ({}))
  if (!response.ok) {
    if (response.status === 401) unauthorizedHandler()
    throw normalizeApiError(response.status, payload)
  }

  return payload.data
}

export function apiPost(path, values) {
  return apiRequest(path, { method: 'POST', body: JSON.stringify(values) })
}

export function apiPut(path, values) {
  return apiRequest(path, { method: 'PUT', body: JSON.stringify(values) })
}

export function apiDelete(path) {
  return apiRequest(path, { method: 'DELETE' })
}
