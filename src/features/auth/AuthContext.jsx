import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { apiRequest, setUnauthorizedHandler } from '../../lib/api'
import { clearAuth, getAuth, setAuth } from '../../lib/auth-storage'

export const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const stored = useMemo(() => getAuth(), [])
  const [token, setToken] = useState(stored?.token ?? null)
  const [user, setUser] = useState(stored?.user ?? null)
  const [ready, setReady] = useState(!stored?.token)

  function logout() {
    clearAuth()
    setToken(null)
    setUser(null)
    setReady(true)
  }

  useEffect(() => {
    setUnauthorizedHandler(logout)
    if (!stored?.token) return

    let active = true
    apiRequest('/auth/me')
      .then((freshUser) => {
        if (!active) return
        setUser(freshUser)
        setAuth({ ...stored, user: freshUser }, stored.remember !== false)
      })
      .catch(() => {
        if (active) logout()
      })
      .finally(() => {
        if (active) setReady(true)
      })

    return () => {
      active = false
      setUnauthorizedHandler(() => {})
    }
  }, [stored])

  function login(nextToken, nextUser, remember) {
    const nextAuth = { token: nextToken, user: nextUser, remember }
    setAuth(nextAuth, remember)
    setToken(nextToken)
    setUser(nextUser)
    setReady(true)
  }

  function updateUser(nextUser) {
    const current = getAuth()
    setUser(nextUser)
    if (current) setAuth({ ...current, user: nextUser }, current.remember !== false)
  }

  return (
    <AuthContext.Provider value={{ token, user, ready, login, logout, updateUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}
