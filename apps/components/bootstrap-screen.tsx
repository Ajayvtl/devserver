'use client'

import { useEffect, useState } from 'react'

import { Progress } from './ui/progress'

const statuses = ['Checking installation...', 'Scanning configuration...', 'Preparing next step...']

interface Props {
  progress: number
  message: string
}

export function BootstrapScreen({ progress, message }: Props) {
  const [statusIndex, setStatusIndex] = useState(0)

  useEffect(() => {
    const timer = window.setInterval(() => {
      setStatusIndex((current) => (current + 1) % statuses.length)
    }, 1400)

    return () => window.clearInterval(timer)
  }, [])

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
        <p>{statuses[statusIndex]}</p>
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
