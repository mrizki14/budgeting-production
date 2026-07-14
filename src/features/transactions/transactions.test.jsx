import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import TransactionsPage from './TransactionsPage'
import { apiRequest } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiRequest: vi.fn(), apiDelete: vi.fn(), apiPost: vi.fn(), apiPut: vi.fn() }))

it('renders the Blade transaction table with signed amount', async () => {
  apiRequest.mockImplementation((path) => Promise.resolve(path === '/categories' ? [] : [{
    id: 3, type: 'expense', amount: 250000, date: '2026-07-14T00:00:00Z', description: 'Lunch', category: { name: 'Food' },
  }]))
  render(<MemoryRouter><TransactionsPage /></MemoryRouter>)

  expect(await screen.findByText('Food')).toBeInTheDocument()
  expect(screen.getByText('-Rp 250.000,00')).toHaveClass('text-rose-600')
  expect(screen.getByText('Jul 14, 2026')).toBeInTheDocument()
})
