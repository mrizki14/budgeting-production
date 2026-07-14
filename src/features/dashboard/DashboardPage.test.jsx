import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import { AuthContext } from '../auth/AuthContext'
import DashboardPage from './DashboardPage'
import { apiRequest } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiRequest: vi.fn() }))

it('renders the Blade dashboard totals and budget status', async () => {
  apiRequest.mockResolvedValue({ total_balance: 5000000, total_income: 8000000, total_expenses: 3000000, current_month_label: 'July 2026', expense_by_category: [], budget_usage: [{ category: 'Food', spent: 760000, limit: 1000000, percentage: 76, status: 'warn' }] })
  render(<AuthContext.Provider value={{ user: { name: 'Ibrahim' }, logout: vi.fn() }}><MemoryRouter><DashboardPage /></MemoryRouter></AuthContext.Provider>)
  expect(await screen.findByText('Rp 5.000.000,00')).toBeInTheDocument()
  expect(screen.getByText("Welcome back, Ibrahim! Here's your financial overview.")).toBeInTheDocument()
  expect(screen.getByText('76%')).toHaveClass('budget-badge-warn')
})
