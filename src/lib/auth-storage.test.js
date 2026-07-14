import { beforeEach, describe, expect, it } from 'vitest'
import { clearAuth, getAuth, setAuth } from './auth-storage'

describe('auth storage', () => {
  beforeEach(clearAuth)

  it('stores remembered authentication only in local storage', () => {
    setAuth({ token: 'abc', user: { id: 1 } }, true)

    expect(localStorage.getItem('budgeting_auth')).toContain('abc')
    expect(sessionStorage.getItem('budgeting_auth')).toBeNull()
    expect(getAuth().token).toBe('abc')
  })

  it('stores temporary authentication only in session storage', () => {
    setAuth({ token: 'xyz', user: { id: 2 } }, false)

    expect(sessionStorage.getItem('budgeting_auth')).toContain('xyz')
    expect(localStorage.getItem('budgeting_auth')).toBeNull()
  })

  it('clears both authentication stores', () => {
    localStorage.setItem('budgeting_auth', '{}')
    sessionStorage.setItem('budgeting_auth', '{}')

    clearAuth()

    expect(getAuth()).toBeNull()
  })
})
