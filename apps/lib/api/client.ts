export type ApiMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export interface ApiErrorShape {
  error?: string
  code?: string
  details?: unknown
}

export class ApiError extends Error {
  status: number
  code?: string
  details?: unknown

  constructor(message: string, status: number, code?: string, details?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
  }
}

export interface RequestOptions {
  method?: ApiMethod
  body?: unknown
  headers?: Record<string, string>
  auth?: boolean
  retries?: number
  signal?: AbortSignal
}

export interface WebSocketHandlers<T> {
  onMessage?: (value: T) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (error: Error) => void
  retry?: boolean
  reconnectDelayMs?: number
}

export interface CommandRequest {
  id?: string
  name?: string
  workspaceId?: string
  capability: string
  provider?: string
  target?: string
  parameters?: Record<string, unknown>
  metadata?: Record<string, string>
  retries?: number
  timeout?: number
}

const DEFAULT_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://127.0.0.1:8080'

let authToken: string | null = null

export function setApiToken(token: string | null, remember = true) {
  authToken = token
  if (typeof window !== 'undefined') {
    if (token && remember) {
      window.localStorage.setItem('devserver-token', token)
    } else {
      window.localStorage.removeItem('devserver-token')
    }
  }
}

export function getApiToken() {
  if (authToken) {
    return authToken
  }
  if (typeof window !== 'undefined') {
    const stored = window.localStorage.getItem('devserver-token')
    if (stored) {
      authToken = stored
      return stored
    }
  }
  return null
}

function baseUrl() {
  return DEFAULT_BASE_URL.replace(/\/$/, '')
}

function toUrl(path: string) {
  return `${baseUrl()}${path.startsWith('/') ? path : `/${path}`}`
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const method = options.method ?? 'GET'
  const retries = options.retries ?? (method === 'GET' ? 2 : 0)
  let attempt = 0

  while (true) {
    try {
      const response = await fetch(toUrl(path), {
        method,
        headers: {
          'Content-Type': 'application/json',
          ...(options.headers ?? {}),
          ...(options.auth !== false && getApiToken() ? { Authorization: `Bearer ${getApiToken()}` } : {}),
        },
        body: options.body === undefined ? undefined : JSON.stringify(options.body),
        signal: options.signal,
      })

      if (!response.ok) {
        const payload = (await safeJson(response)) as ApiErrorShape | undefined
        throw new ApiError(payload?.error ?? response.statusText, response.status, payload?.code, payload?.details)
      }

      if (response.status === 204) {
        return undefined as T
      }

      return (await safeJson(response)) as T
    } catch (error) {
      if (attempt >= retries) {
        if (error instanceof ApiError) {
          throw error
        }
        throw normalizeError(error)
      }

      attempt += 1
      await new Promise((resolve) => setTimeout(resolve, 200 * attempt))
    }
  }
}

export async function requestOrFallback<T>(path: string, fallback: T, options: RequestOptions = {}): Promise<T> {
  try {
    return await request<T>(path, options)
  } catch (error) {
    if (error instanceof ApiError && error.status !== 0) {
      throw error
    }
    return fallback
  }
}

export async function submitCommand<T>(command: CommandRequest, options: RequestOptions = {}): Promise<T> {
  return request<T>('/api/commands', {
    method: 'POST',
    auth: options.auth,
    headers: options.headers,
    signal: options.signal,
    body: command,
  })
}

async function safeJson(response: Response) {
  const contentType = response.headers.get('content-type') ?? ''
  if (contentType.includes('application/json')) {
    return response.json()
  }

  return response.text()
}

function normalizeError(error: unknown) {
  if (error instanceof ApiError) {
    return error
  }
  if (error instanceof Error) {
    return new ApiError(error.message, 0)
  }
  return new ApiError('Unknown API error', 0)
}

export interface SocketConnection<T> {
  close: () => void
}

export function connectSocket<T>(path: string, handlers: WebSocketHandlers<T> = {}): SocketConnection<T> {
  if (typeof window === 'undefined') {
    return { close: () => undefined }
  }

  let closed = false
  let socket: WebSocket | null = null
  let retries = 0
  let retryTimer: number | null = null

  const open = () => {
    const base = baseUrl().replace(/^http/, 'ws')
    socket = new WebSocket(`${base}${path.startsWith('/') ? path : `/${path}`}`)

    socket.onopen = () => {
      retries = 0
      handlers.onOpen?.()
    }

    socket.onmessage = (event) => {
      try {
        handlers.onMessage?.(JSON.parse(String(event.data)) as T)
      } catch (error) {
        handlers.onError?.(normalizeError(error))
      }
    }

    socket.onerror = () => {
      handlers.onError?.(new ApiError('WebSocket error', 0))
    }

    socket.onclose = () => {
      handlers.onClose?.()
      if (!closed && handlers.retry !== false) {
        retries += 1
        const delay = Math.min((handlers.reconnectDelayMs ?? 750) * retries, 5000)
        retryTimer = window.setTimeout(open, delay)
      }
    }
  }

  open()

  return {
    close() {
      closed = true
      if (retryTimer !== null) {
        window.clearTimeout(retryTimer)
        retryTimer = null
      }
      socket?.close()
    },
  }
}
