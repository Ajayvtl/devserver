import { BootstrapScreen } from '@/components/bootstrap-screen'

export default function Loading() {
  return <BootstrapScreen progress={64} message="Preparing the setup wizard..." />
}
