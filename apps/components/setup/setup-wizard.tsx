'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'

import type { SetupWizardData } from '@/lib/types'
import type { TaskEvent } from '@/lib/types'
import { completeSetupWizard } from '@/lib/services/setup'
import { connectSocket } from '@/lib/api/client'

import { useToast } from '../toast'
import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import { Input } from '../ui/input'
import { LearnCard } from '../ui/learn-card'
import { Progress } from '../ui/progress'
import { Stepper } from '../ui/stepper'

interface Props {
  initialData: SetupWizardData
}

interface FormState {
  licenseAccepted: boolean
  adminName: string
  adminEmail: string
  adminPassword: string
  provider: string
  installationType: string
}

const initialForm: FormState = {
  licenseAccepted: false,
  adminName: '',
  adminEmail: '',
  adminPassword: '',
  provider: 'None',
  installationType: 'quick',
}

export function SetupWizard({ initialData }: Props) {
  const router = useRouter()
  const { push } = useToast()
  const [stepIndex, setStepIndex] = useState(0)
  const [form, setForm] = useState<FormState>(initialForm)
  const [errors, setErrors] = useState<string[]>([])
  const [installProgress, setInstallProgress] = useState(0)
  const [installing, setInstalling] = useState(false)
  const [taskState, setTaskState] = useState<'idle' | 'running' | 'done'>('idle')
  const taskIdRef = useRef<string | null>(null)

  const currentStep = initialData.steps[stepIndex]
  const totalSteps = initialData.steps.length
  const percent = useMemo(() => Math.round(((stepIndex + 1) / totalSteps) * 100), [stepIndex, totalSteps])

  const learnArticles = useMemo(() => initialData.knowledgeArticles, [initialData.knowledgeArticles])

  useEffect(() => {
    const connection = connectSocket<TaskEvent>('/ws/tasks', {
      onMessage(event) {
        if (event.scope !== 'setup') {
          return
        }

        if (taskIdRef.current && event.taskId !== taskIdRef.current) {
          return
        }

        taskIdRef.current = event.taskId
        setInstallProgress(event.progress)
        setTaskState(event.state === 'Done' ? 'done' : 'running')

        if (event.state === 'Done') {
          setInstalling(false)
          setStepIndex(totalSteps - 1)
          push({
            title: 'Setup complete',
            message: 'DevServer is ready. Review the finish screen and continue.',
            tone: 'success',
          })
        }
      },
    })

    return () => connection.close()
  }, [push, totalSteps])

  const validate = () => {
    const nextErrors: string[] = []

    if (currentStep.key === 'license' && !form.licenseAccepted) {
      nextErrors.push('You must accept the license to continue.')
    }

    if (currentStep.key === 'admin') {
      if (form.adminName.trim().length < 2) {
        nextErrors.push('Administrator name is required.')
      }
      if (!form.adminEmail.includes('@')) {
        nextErrors.push('A valid administrator email is required.')
      }
      if (form.adminPassword.length < 8) {
        nextErrors.push('Administrator password must be at least 8 characters.')
      }
    }

    if (currentStep.key === 'install-type' && !form.installationType) {
      nextErrors.push('Choose an installation type.')
    }

    setErrors(nextErrors)
    return nextErrors.length === 0
  }

  const next = () => {
    if (!validate()) {
      push({
        title: 'Setup needs attention',
        message: 'Please fix the highlighted inputs before continuing.',
        tone: 'warning',
      })
      return
    }

    if (currentStep.key === 'installing') {
      return
    }

    if (currentStep.key === 'summary') {
      setInstalling(true)
      setTaskState('running')
      setStepIndex(stepIndex + 1)
      completeSetupWizard(form)
        .then((result) => {
          taskIdRef.current = result.taskId
          setInstallProgress(10)
        })
        .catch((error) => {
          setInstalling(false)
          setTaskState('idle')
          setErrors([error instanceof Error ? error.message : 'Unable to start setup.'])
        })
      return
    }

    setStepIndex((current) => Math.min(current + 1, totalSteps - 1))
  }

  const back = () => {
    if (stepIndex === 0) {
      return
    }

    setStepIndex((current) => Math.max(current - 1, 0))
  }

  const progressSteps = initialData.steps.map((step, index) => ({
    label: step.title,
    active: index === stepIndex,
    done: index < stepIndex,
  }))

  return (
    <div className="wizard-layout">
      <div className="wizard-main">
        <div className="brand-lockup brand-lockup--wizard">
          <div className="brand-lockup__mark">D</div>
          <div>
            <div className="brand-lockup__title">DevServer</div>
            <div className="brand-lockup__subtitle">Setup wizard</div>
          </div>
        </div>

        <div className="wizard-hero">
          <Badge tone="accent">First time setup</Badge>
          <h1>{currentStep.title}</h1>
          <p>{currentStep.subtitle}</p>
          <Progress value={percent} tone="accent" />
        </div>

        <Stepper steps={progressSteps} />

        {errors.length ? (
          <Card className="wizard-errors">
            {errors.map((error) => (
              <div key={error} className="wizard-error">
                {error}
              </div>
            ))}
          </Card>
        ) : null}

        <Card>
          {currentStep.key === 'welcome' ? (
            <div className="wizard-copy">
              <h2>Welcome to DevServer</h2>
              <p>
                We&apos;ll walk through license acceptance, system checks, administrator setup, AI
                providers, and installation preferences before provisioning your platform.
              </p>
              <div className="wizard-highlight">
                <strong>What you get</strong>
                <span>Bootstrap, services, deployment, monitoring, and a learning layer in one flow.</span>
              </div>
            </div>
          ) : null}

          {currentStep.key === 'license' ? (
            <div className="wizard-form">
              <label className="checkbox-card">
                <input
                  type="checkbox"
                  checked={form.licenseAccepted}
                  onChange={(event) => setForm((current) => ({ ...current, licenseAccepted: event.target.checked }))}
                />
                <span>
                  I accept the DevServer license and acknowledge that this setup will modify the server.
                </span>
              </label>
            </div>
          ) : null}

          {currentStep.key === 'scan' ? (
            <div className="wizard-scan">
              {initialData.checks.map((check) => (
                <div key={check.label} className="scan-row">
                  <div className={`scan-row__status scan-row__status--${check.status}`}>{check.status}</div>
                  <div>
                    <strong>{check.label}</strong>
                    <p>{check.detail}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : null}

          {currentStep.key === 'admin' ? (
            <div className="wizard-form wizard-form--grid">
              <Input
                label="Administrator name"
                value={form.adminName}
                onChange={(event) => setForm((current) => ({ ...current, adminName: event.target.value }))}
                placeholder="Ava Patel"
              />
              <Input
                label="Administrator email"
                type="email"
                value={form.adminEmail}
                onChange={(event) => setForm((current) => ({ ...current, adminEmail: event.target.value }))}
                placeholder="admin@company.com"
              />
              <Input
                label="Administrator password"
                type="password"
                value={form.adminPassword}
                onChange={(event) => setForm((current) => ({ ...current, adminPassword: event.target.value }))}
                placeholder="At least 12 characters"
              />
            </div>
          ) : null}

          {currentStep.key === 'ai' ? (
            <div className="choice-grid">
              {initialData.providers.map((provider) => (
                <button
                  key={provider}
                  type="button"
                  className={`choice-card${form.provider === provider ? ' is-selected' : ''}`}
                  onClick={() => setForm((current) => ({ ...current, provider }))}
                >
                  <strong>{provider}</strong>
                  <span>{provider === 'None' ? 'Skip AI for now' : 'Connect provider for recommendations and insights.'}</span>
                </button>
              ))}
            </div>
          ) : null}

          {currentStep.key === 'install-type' ? (
            <div className="choice-grid">
              {initialData.installationTypes.map((type) => (
                <button
                  key={type.key}
                  type="button"
                  className={`choice-card${form.installationType === type.key ? ' is-selected' : ''}`}
                  onClick={() => setForm((current) => ({ ...current, installationType: type.key }))}
                >
                  <strong>{type.label}</strong>
                  <span>{type.detail}</span>
                </button>
              ))}
            </div>
          ) : null}

          {currentStep.key === 'summary' ? (
            <div className="summary-grid">
              <div>
                <h3>Administrator</h3>
                <p>{form.adminName || 'Not set'}</p>
                <p>{form.adminEmail || 'Not set'}</p>
              </div>
              <div>
                <h3>AI Provider</h3>
                <p>{form.provider}</p>
              </div>
              <div>
                <h3>Installation</h3>
                <p>{form.installationType}</p>
              </div>
            </div>
          ) : null}

          {currentStep.key === 'installing' ? (
            <div className="installing-state">
              <h2>Installing</h2>
              <p>Applying configuration, registering services, and preparing your dashboard.</p>
              <Progress value={installProgress} tone={taskState === 'done' ? 'success' : 'accent'} />
              <div className="installing-state__meta">
                {installProgress}% complete {taskState === 'running' ? ' - provisioning tasks are running' : ''}
              </div>
            </div>
          ) : null}

          {currentStep.key === 'finished' ? (
            <div className="wizard-finished">
              <EmptyState
                title="Setup complete"
                description="DevServer is ready. You can now sign in and start managing projects."
                actionLabel="Open dashboard"
                actionHref="/dashboard"
              />
            </div>
          ) : null}
        </Card>

        <div className="wizard-actions">
          <Button variant="secondary" onClick={back} type="button" disabled={installing || stepIndex === 0}>
            Back
          </Button>
          <Button
            variant="primary"
            onClick={currentStep.key === 'finished' ? () => router.replace('/dashboard') : next}
            type="button"
            disabled={installing && currentStep.key !== 'finished'}
          >
            {currentStep.key === 'summary' ? 'Install now' : currentStep.key === 'finished' ? 'Open dashboard' : 'Continue'}
          </Button>
        </div>
      </div>

      <aside className="wizard-rail">
        <LearnCard title="Setup guidance" articles={learnArticles} />
        <Card>
          <div className="card__eyebrow">System</div>
          <h3>Ready checklist</h3>
          <div className="rail-checks">
            {initialData.checks.slice(0, 6).map((check) => (
              <div key={check.label} className="rail-check">
                <span>{check.label}</span>
                <strong>{check.status}</strong>
              </div>
            ))}
          </div>
        </Card>
      </aside>
    </div>
  )
}
