import { beforeEach, describe, expect, it, vi } from 'vitest'
import { apiRequest, setUnauthorizedHandler } from './api'
import { setAuth } from './auth-storage'

describe('apiRequest', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    setUnauthorizedHandler(() => {})
  })

  it('attaches bearer authentication and unwraps data', async () => {
    setAuth({ token: 'secret', user: { id: 1 } }, true)
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ data: { id: 7 } }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    })))

    const result = await apiRequest('/categories')

    expect(result).toEqual({ id: 7 })
    expect(fetch).toHaveBeenCalledWith('http://127.0.0.1:8080/api/categories', expect.objectContaining({
      headers: expect.objectContaining({ Authorization: 'Bearer secret' }),
    }))
  })

  it('normalizes API errors', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({
      message: 'Validation failed',
      errors: { email: ['email already exists'] },
    }), { status: 422, headers: { 'Content-Type': 'application/json' } })))

    await expect(apiRequest('/auth/register')).rejects.toMatchObject({
      status: 422,
      message: 'Validation failed',
      errors: { email: ['email already exists'] },
    })
  })

  it('notifies authentication when the API returns unauthorized', async () => {
    const unauthorized = vi.fn()
    setUnauthorizedHandler(unauthorized)
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('{}', { status: 401 })))

    await expect(apiRequest('/auth/me')).rejects.toMatchObject({ status: 401 })

    expect(unauthorized).toHaveBeenCalledOnce()
  })
})
