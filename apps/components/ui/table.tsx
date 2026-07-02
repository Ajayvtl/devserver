'use client'

import type { ReactNode } from 'react'

interface Column {
  header: string
  align?: 'left' | 'right' | 'center'
}

interface Props {
  columns: Column[]
  rows: ReactNode[][]
  caption?: string
}

export function DataTable({ columns, rows, caption }: Props) {
  return (
    <div className="table-wrap">
      <table className="data-table">
        {caption ? <caption className="sr-only">{caption}</caption> : null}
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={column.header} style={{ textAlign: column.align ?? 'left' }}>
                {column.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, rowIndex) => (
            <tr key={rowIndex}>
              {row.map((cell, cellIndex) => (
                <td key={cellIndex}>{cell}</td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
