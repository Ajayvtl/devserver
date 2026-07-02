'use client'

import type { ReactNode } from 'react'

interface TabItem {
  key: string
  label: string
  panel: ReactNode
}

interface Props {
  tabs: TabItem[]
  activeKey: string
  onChange: (key: string) => void
}

export function Tabs({ tabs, activeKey, onChange }: Props) {
  return (
    <div className="tabs">
      <div className="tabs__list" role="tablist" aria-label="Project sections">
        {tabs.map((tab) => {
          const active = tab.key === activeKey
          return (
            <button
              key={tab.key}
              type="button"
              role="tab"
              aria-selected={active}
              className={`tabs__tab${active ? ' is-active' : ''}`}
              onClick={() => onChange(tab.key)}
            >
              {tab.label}
            </button>
          )
        })}
      </div>

      <div className="tabs__panel" role="tabpanel">
        {tabs.find((tab) => tab.key === activeKey)?.panel}
      </div>
    </div>
  )
}
