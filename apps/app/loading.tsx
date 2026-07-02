import { BootstrapScreen } from '@/components/bootstrap-screen'

export default function Loading() {
  return <BootstrapScreen progress={24} message="Preparing the DevServer experience..." />
}
