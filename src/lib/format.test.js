import { describe, expect, it } from 'vitest'
import { formatDate, formatMoney, monthName } from './format'

describe('display formatting', () => {
  it('formats rupiah like the Blade UI', () => {
    expect(formatMoney(1250000)).toBe('Rp 1.250.000,00')
  })

  it('formats API dates in the Blade table style', () => {
    expect(formatDate('2026-07-14T00:00:00Z')).toBe('Jul 14, 2026')
  })

  it('returns English month names used by the Blade UI', () => {
    expect(monthName(7)).toBe('July')
  })
})
