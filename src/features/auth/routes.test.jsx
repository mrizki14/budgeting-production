import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { AuthContext } from './AuthContext'
import GuestRoute from './GuestRoute'
import ProtectedRoute from './ProtectedRoute'

function renderRoute(value, element, initial = '/private', guardedPath = '/private', guardedText = 'Private') {
  return render(
    <AuthContext.Provider value={value}>
      <MemoryRouter initialEntries={[initial]}>
        <Routes>
          <Route element={element}>
            <Route path={guardedPath} element={<p>{guardedText}</p>} />
          </Route>
          <Route path="/dashboard" element={<p>Dashboard</p>} />
          {guardedPath !== '/login' && <Route path="/login" element={<p>Login</p>} />}
        </Routes>
      </MemoryRouter>
    </AuthContext.Provider>,
  )
}

describe('authentication routes', () => {
  it('redirects guests away from protected routes', () => {
    renderRoute({ ready: true, token: null }, <ProtectedRoute />)
    expect(screen.getByText('Login')).toBeInTheDocument()
  })

  it('renders protected routes for authenticated users', () => {
    renderRoute({ ready: true, token: 'token' }, <ProtectedRoute />)
    expect(screen.getByText('Private')).toBeInTheDocument()
  })

  it('redirects authenticated users away from login', () => {
    renderRoute({ ready: true, token: 'token' }, <GuestRoute />, '/login', '/login', 'Login')
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
  })
})
