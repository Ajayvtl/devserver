import { getLoginData } from '@/lib/services/auth'

import { LoginForm } from '@/components/auth/login-form'

export default async function LoginPage() {
  const data = await getLoginData()
  return <LoginForm branding={data.branding} subtitle={data.subtitle} supportEmail={data.supportEmail} />
}
