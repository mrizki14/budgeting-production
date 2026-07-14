import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import App from './App'

describe('App', () => {
  it('renders the application entrypoint', () => {
    render(<App />)

    expect(screen.getByRole('heading', { name: 'Kelola budget pribadi dengan lebih rapi.' })).toBeInTheDocument()
  })
})
