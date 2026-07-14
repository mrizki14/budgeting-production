import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AuthContext } from './AuthContext'
import AuthPage from './AuthPage'
import { apiPost } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiPost: vi.fn() }))

function renderPage(login = vi.fn()) {
  render(
    <AuthContext.Provider value={{ login }}>
      <MemoryRouter initialEntries={['/login']}>
        <Routes>
          <Route path="/login" element={<AuthPage />} />
          <Route path="/dashboard" element={<p>Dashboard destination</p>} />
        </Routes>
      </MemoryRouter>
    </AuthContext.Provider>,
  )
  return login
}

describe('AuthPage', () => {
  beforeEach(() => vi.clearAllMocks())

  it('switches between the Blade login and registration forms', async () => {
    renderPage()

    expect(screen.getByRole('button', { name: 'Login ke Dashboard' })).toBeInTheDocument()
    await userEvent.click(screen.getByRole('button', { name: 'Daftar' }))

    expect(screen.getByLabelText('Nama Lengkap')).toBeInTheDocument()
    expect(screen.getByLabelText('Konfirmasi Password')).toBeInTheDocument()
  })

  it('logs in with the remember preference and redirects', async () => {
    apiPost.mockResolvedValue({ token: 'token', user: { id: 1, name: 'Ibrahim' } })
    const login = renderPage()
    await userEvent.type(screen.getByLabelText('Email'), 'demo@example.com')
    await userEvent.type(screen.getByLabelText('Password'), 'password123')
    await userEvent.click(screen.getByLabelText('Ingat saya'))
    await userEvent.click(screen.getByRole('button', { name: 'Login ke Dashboard' }))

    await waitFor(() => expect(apiPost).toHaveBeenCalledWith('/auth/login', {
      email: 'demo@example.com', password: 'password123',
    }))
    expect(login).toHaveBeenCalledWith('token', { id: 1, name: 'Ibrahim' }, true)
    expect(screen.getByText('Dashboard destination')).toBeInTheDocument()
  })

  it('shows a confirmation error before registration submission', async () => {
    renderPage()
    await userEvent.click(screen.getByRole('button', { name: 'Daftar' }))
    await userEvent.type(screen.getByLabelText('Nama Lengkap'), 'Ibrahim')
    await userEvent.type(screen.getByLabelText('Email'), 'new@example.com')
    await userEvent.type(screen.getByLabelText('Password'), 'password123')
    await userEvent.type(screen.getByLabelText('Konfirmasi Password'), 'different')
    await userEvent.click(screen.getByRole('button', { name: 'Buat Akun' }))

    expect(screen.getByText('Konfirmasi password tidak cocok.')).toBeInTheDocument()
    expect(apiPost).not.toHaveBeenCalled()
  })
})
