import { request, requestOrFallback, setApiToken } from '../api/client'
import { loginFallback } from './fallbacks'
import type { LoginData } from '../types'

export async function getLoginData(): Promise<LoginData> {
  return requestOrFallback<LoginData>('/api/auth/login', loginFallback, { auth: false })
}

export async function authenticate(email: string, password: string, remember: boolean): Promise<{ ok: boolean; message: string }> {
  const response = await request<{ token: string; role: string }>('/api/auth/login', {
    method: 'POST',
    body: { email, password, remember },
    auth: false,
  })

  setApiToken(response.token, remember)
  return { ok: true, message: `Signed in as ${response.role}.` }
}
