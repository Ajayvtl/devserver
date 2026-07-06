import { request, setApiToken, submitCommand } from '../api/client'
import type { CommandRecord, LoginData } from '../types'

export async function getLoginData(): Promise<LoginData> {
  return { branding: 'DevServer', subtitle: 'Sign in to continue', supportEmail: 'admin@localhost', knowledgeArticles: [] }
}

export async function authenticate(email: string, password: string, remember: boolean): Promise<{ ok: boolean; message: string }> {
  try {
    const result = await request<{ accessToken: string; refreshToken: string }>('/api/v1/auth/login', {
      method: 'POST',
      body: { provider: 'local', username: email, email, password, remember },
      auth: false,
    })

    if (!result?.accessToken) {
      return { ok: false, message: 'Authentication failed.' }
    }

    setApiToken(result.accessToken, remember)
    return { ok: true, message: `Signed in successfully.` }
  } catch (error: any) {
    return { ok: false, message: error.message ?? 'Authentication failed.' }
  }
}

export async function logout(): Promise<void> {
  // In WP-8.1/8.2 we didn't add a /api/v1/auth/logout endpoint, but we can clear the client token
  setApiToken(null)
}
