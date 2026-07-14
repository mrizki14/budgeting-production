import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import ReportsPage from './ReportsPage'
import { apiRequest } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiRequest: vi.fn() }))

it('renders report summary and category breakdowns', async () => {
  apiRequest.mockResolvedValue({ month: 7, year: 2026, total_income: 2000000, total_expenses: 500000, net_savings: 1500000, income_by_category: [{ category_name: 'Salary', total_amount: 2000000 }], expense_by_category: [] })
  render(<MemoryRouter><ReportsPage /></MemoryRouter>)
  expect(await screen.findByText('Rp 1.500.000,00')).toBeInTheDocument()
  expect(screen.getByText('Salary')).toBeInTheDocument()
  expect(screen.getAllByText('July 2026').length).toBeGreaterThan(0)
})

it('renders empty states when the API encodes empty breakdowns as null', async () => {
  apiRequest.mockResolvedValue({ month: 7, year: 2026, total_income: 0, total_expenses: 0, net_savings: 0, income_by_category: null, expense_by_category: null })

  render(<MemoryRouter><ReportsPage /></MemoryRouter>)

  expect(await screen.findByText('Belum ada pemasukan pada periode ini.')).toBeInTheDocument()
  expect(screen.getByText('Belum ada pengeluaran pada periode ini.')).toBeInTheDocument()
})
