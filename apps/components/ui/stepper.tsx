'use client'

interface Step {
  label: string
  active?: boolean
  done?: boolean
}

interface Props {
  steps: Step[]
}

export function Stepper({ steps }: Props) {
  return (
    <ol className="stepper">
      {steps.map((step) => (
        <li key={step.label} className={`stepper__step${step.active ? ' is-active' : ''}${step.done ? ' is-done' : ''}`}>
          <span className="stepper__dot" />
          <span className="stepper__label">{step.label}</span>
        </li>
      ))}
    </ol>
  )
}
