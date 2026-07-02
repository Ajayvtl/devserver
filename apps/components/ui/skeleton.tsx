interface Props {
  lines?: number
}

export function Skeleton({ lines = 3 }: Props) {
  return (
    <div className="skeleton-group" aria-hidden="true">
      {Array.from({ length: lines }).map((_, index) => (
        <div key={index} className="skeleton-line" />
      ))}
    </div>
  )
}
