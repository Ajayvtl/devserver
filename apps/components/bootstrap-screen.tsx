'use client'

import { Progress } from './ui/progress'

interface Props {
  progress: number
  message: string
}

export function BootstrapScreen({ progress, message }: Props) {
  const status = progress >= 100 ? 'Ready to continue...' : progress >= 60 ? 'Finalizing runtime state...' : 'Checking installation...'

  return (
    <div className="bootstrap-screen">
      <div className="bootstrap-screen__card">
        <div className="brand-lockup">
          <div className="brand-lockup__mark">D</div>
          <div>
            <div className="brand-lockup__title">DevServer</div>
            <div className="brand-lockup__subtitle">Infrastructure platform</div>
          </div>
        </div>

        <h1>Loading...</h1>
        <p>{status}</p>
        <p className="bootstrap-screen__message">{message}</p>
        <Progress value={progress} tone="accent" />

        <div className="bootstrap-screen__bars" aria-hidden="true">
          <span />
          <span />
          <span />
          <span />
        </div>
      </div>
    </div>
  )
}
