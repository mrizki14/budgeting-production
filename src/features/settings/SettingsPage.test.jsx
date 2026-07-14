import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import { AuthContext } from '../auth/AuthContext'
import SettingsPage from './SettingsPage'
import { apiPut } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiPut: vi.fn() }))

it('updates the profile and refreshes authenticated user data', async () => {
  apiPut.mockResolvedValue({ id: 1, name: 'Updated', email: 'new@example.com' })
  const updateUser = vi.fn()
  render(<AuthContext.Provider value={{ user: { id: 1, name: 'Ibrahim', email: 'old@example.com' }, updateUser }}><MemoryRouter><SettingsPage /></MemoryRouter></AuthContext.Provider>)
  const name = screen.getByLabelText('Account Name')
  await userEvent.clear(name); await userEvent.type(name, 'Updated')
  const email = screen.getByLabelText('Email')
  await userEvent.clear(email); await userEvent.type(email, 'new@example.com')
  await userEvent.click(screen.getByRole('button', { name: 'Save Profile' }))
  await waitFor(() => expect(updateUser).toHaveBeenCalledWith({ id: 1, name: 'Updated', email: 'new@example.com' }))
  expect(screen.getByText('Profil berhasil diperbarui.')).toBeInTheDocument()
})
