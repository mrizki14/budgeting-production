import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { AuthProvider, useAuth } from './AuthContext'
import { setAuth } from '../../lib/auth-storage'

function Probe() {
  const { ready, user, login, logout } = useAuth()
  return (
    <div>
      <span>{ready ? 'ready' : 'loading'}</span>
      <span>{user?.name ?? 'guest'}</span>
      <button onClick={() => login('token', { id: 1, name: 'Ibrahim' }, true)}>login</button>
      <button onClick={logout}>logout</button>
    </div>
  )
}

describe('AuthProvider', () => {
  it('starts ready for a guest and persists login preference', async () => {
    render(<AuthProvider><Probe /></AuthProvider>)

    expect(await screen.findByText('ready')).toBeInTheDocument()
    await userEvent.click(screen.getByText('login'))

    expect(screen.getByText('Ibrahim')).toBeInTheDocument()
    expect(localStorage.getItem('budgeting_auth')).toContain('token')
  })

  it('validates stored authentication and refreshes the user', async () => {
    setAuth({ token: 'stored', user: { id: 1, name: 'Old' } }, true)
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({
      data: { id: 1, name: 'Fresh' },
    }), { status: 200, headers: { 'Content-Type': 'application/json' } })))

    render(<AuthProvider><Probe /></AuthProvider>)

    await waitFor(() => expect(screen.getByText('Fresh')).toBeInTheDocument())
  })
})
