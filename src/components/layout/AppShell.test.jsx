import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import { AuthContext } from '../../features/auth/AuthContext'
import AppShell from './AppShell'

describe('AppShell', () => {
  it('renders the same navigation and account controls as Blade', async () => {
    const logout = vi.fn()
    render(
      <AuthContext.Provider value={{ user: { name: 'Ibrahim', email: 'ibrahim@example.com' }, logout }}>
        <MemoryRouter initialEntries={['/dashboard']}>
          <Routes>
            <Route element={<AppShell />}>
              <Route path="/dashboard" element={<p>Page content</p>} />
            </Route>
          </Routes>
        </MemoryRouter>
      </AuthContext.Provider>,
    )

    for (const label of ['Dashboard', 'Categories', 'Transactions', 'Budgets', 'Reports', 'Settings']) {
      expect(screen.getByRole('link', { name: new RegExp(label) })).toBeInTheDocument()
    }
    expect(screen.getAllByText('Ibrahim').length).toBeGreaterThan(0)
    expect(screen.getByText('ibrahim@example.com')).toBeInTheDocument()

    await userEvent.click(screen.getAllByRole('button', { name: 'Logout' })[0])
    expect(logout).toHaveBeenCalledOnce()
  })
})
