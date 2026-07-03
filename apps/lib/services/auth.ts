import { setApiToken, submitCommand } from '../api/client'
import { loginFallback } from './fallbacks'
import type { CommandRecord, LoginData } from '../types'

export async function getLoginData(): Promise<LoginData> {
  return loginFallback
}

export async function authenticate(email: string, password: string, remember: boolean): Promise<{ ok: boolean; message: string }> {
  const response = await submitCommand<CommandRecord>({
    capability: 'auth.login',
    parameters: { email, password, remember },
    name: 'Authenticate user',
  }, {
    auth: false,
  })

  const result = response.result as { token?: string; role?: string } | undefined
  if (!result?.token) {
    return { ok: false, message: 'Authentication did not return a token.' }
  }

  setApiToken(result.token, remember)
  return { ok: true, message: `Signed in as ${result.role ?? 'User'}.` }
}
