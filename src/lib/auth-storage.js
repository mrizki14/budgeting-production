const AUTH_KEY = 'budgeting_auth'

export function setAuth(auth, remember) {
  clearAuth()
  const storage = remember ? localStorage : sessionStorage
  storage.setItem(AUTH_KEY, JSON.stringify(auth))
}

export function getAuth() {
  const raw = localStorage.getItem(AUTH_KEY) ?? sessionStorage.getItem(AUTH_KEY)
  if (!raw) return null

  try {
    return JSON.parse(raw)
  } catch {
    clearAuth()
    return null
  }
}

export function clearAuth() {
  localStorage.removeItem(AUTH_KEY)
  sessionStorage.removeItem(AUTH_KEY)
}
