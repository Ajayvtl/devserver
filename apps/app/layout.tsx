import type { Metadata } from 'next'
import { Manrope, Space_Grotesk } from 'next/font/google'
import type { ReactNode } from 'react'

import { Providers } from '@/components/providers'

import './globals.css'

const heading = Space_Grotesk({
  subsets: ['latin'],
  variable: '--font-heading',
})

const body = Manrope({
  subsets: ['latin'],
  variable: '--font-body',
})

export const metadata: Metadata = {
  title: 'DevServer',
  description: 'A production-ready control plane for server bootstrap, service management, and multi-target deployment.',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode
}>) {
  return (
    <html lang="en" className={`${heading.variable} ${body.variable}`}>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
