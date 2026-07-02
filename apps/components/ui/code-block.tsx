'use client'

interface Props {
  title: string
  children: string
}

export function CodeBlock({ title, children }: Props) {
  return (
    <div className="code-block">
      <div className="code-block__title">{title}</div>
      <pre>
        <code>{children}</code>
      </pre>
    </div>
  )
}
