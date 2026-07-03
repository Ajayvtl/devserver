'use client'

import { useState, type FormEvent } from 'react'
import { useRouter } from 'next/navigation'

import { authenticate } from '@/lib/services/auth'

import { useToast } from '../toast'
import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { ErrorState } from '../ui/error-state'
import { Input } from '../ui/input'
import { LearnCard } from '../ui/learn-card'

interface Props {
  branding: string
  subtitle: string
  supportEmail: string
}

export function LoginForm({ branding, subtitle, supportEmail }: Props) {
  const router = useRouter()
  const { push } = useToast()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [remember, setRemember] = useState(true)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setBusy(true)
    setError(null)

    const result = await authenticate(email, password, remember)
    if (!result.ok) {
      setError(result.message)
      setBusy(false)
      push({
        title: 'Sign in failed',
        message: result.message,
        tone: 'danger',
      })
      return
    }

    push({
      title: 'Welcome back',
      message: remember ? 'Session will be remembered on this device.' : 'Signed in for this session.',
      tone: 'success',
    })
    router.replace('/dashboard')
  }

  return (
    <div className="auth-layout">
      <Card className="auth-card">
        <div className="brand-lockup">
          <div className="brand-lockup__mark">D</div>
          <div>
            <div className="brand-lockup__title">{branding}</div>
            <div className="brand-lockup__subtitle">Sign in</div>
          </div>
        </div>

        <h1>Welcome back</h1>
        <p className="auth-card__subtitle">{subtitle}</p>

        {error ? <ErrorState title="Sign in failed" description={error} /> : null}

        <form onSubmit={submit} className="auth-form">
          <Input label="Email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} />
          <Input label="Password" type="password" value={password} onChange={(event) => setPassword(event.target.value)} />

          <label className="checkbox-row">
            <input type="checkbox" checked={remember} onChange={(event) => setRemember(event.target.checked)} />
            <span>Remember me</span>
          </label>

          <div className="auth-actions">
            <Button variant="ghost" href={`mailto:${supportEmail}`}>
              Forgot password?
            </Button>
            <Button type="submit" variant="primary" disabled={busy}>
              {busy ? 'Signing in...' : 'Sign in'}
            </Button>
          </div>
        </form>
      </Card>

      <aside className="auth-rail">
        <LearnCard
          title="Login help"
          articles={[
            {
              title: 'Secure access',
              question: 'Why sign in separately?',
              summary: 'Authentication is handled separately from setup so sessions and roles stay clear.',
              link: '#learn-access',
            },
            {
              title: 'Recovery',
              question: 'What if I lose access?',
              summary: 'The UI will later connect to a secure recovery flow through the service layer.',
              link: '#learn-recovery',
            },
          ]}
        />

        <Card>
          <Badge tone="info">Command-driven workflow</Badge>
          <h3 className="rail-title">Need a fresh install?</h3>
          <p className="rail-copy">If the platform is not set up yet, the bootstrap router will send you to the setup wizard first.</p>
          <Button href="/setup" variant="secondary">
            Open setup wizard
          </Button>
        </Card>
      </aside>
    </div>
  )
}
