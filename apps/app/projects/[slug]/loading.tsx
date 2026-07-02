import { Card } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

export default function Loading() {
  return (
    <div className="page">
      <Card>
        <Skeleton lines={4} />
      </Card>
      <Card>
        <Skeleton lines={5} />
      </Card>
    </div>
  )
}
