import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import CategoriesPage from './CategoriesPage'
import { apiRequest } from '../../lib/api'

vi.mock('../../lib/api', () => ({ apiRequest: vi.fn(), apiDelete: vi.fn(), apiPost: vi.fn(), apiPut: vi.fn() }))

describe('CategoriesPage', () => {
  beforeEach(() => vi.clearAllMocks())

  it('renders category data using the Blade table', async () => {
    apiRequest.mockResolvedValue([{ id: 1, name: 'Salary', type: 'income', created_at: '2026-07-14T10:00:00Z' }])

    render(<MemoryRouter><CategoriesPage /></MemoryRouter>)

    expect(await screen.findByText('Salary')).toBeInTheDocument()
    expect(screen.getByText('Income')).toHaveClass('budget-badge-success')
    expect(screen.getByRole('link', { name: 'Add Category' })).toHaveAttribute('href', '/categories/create')
  })
})
