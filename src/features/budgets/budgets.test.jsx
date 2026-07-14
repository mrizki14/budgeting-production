import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import BudgetsPage from './BudgetsPage'
import { apiRequest } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiRequest: vi.fn(), apiDelete: vi.fn(), apiPost: vi.fn(), apiPut: vi.fn() }))

it('renders budgets with their related category and period', async () => {
  apiRequest.mockResolvedValue([{ id: 2, month: 7, year: 2026, limit_amount: 1500000, category: { name: 'Groceries', type: 'expense' } }])
  render(<MemoryRouter><BudgetsPage /></MemoryRouter>)

  expect((await screen.findAllByText('Groceries')).length).toBeGreaterThan(0)
  expect(screen.getByText('Period: July 2026')).toBeInTheDocument()
  expect(screen.getByText('Rp 1.500.000,00')).toBeInTheDocument()
})
