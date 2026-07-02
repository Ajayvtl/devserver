'use client'

import type { KnowledgeArticle } from '@/lib/types'

interface Props {
  title: string
  articles: KnowledgeArticle[]
}

export function LearnCard({ title, articles }: Props) {
  return (
    <section className="learn-card">
      <div className="learn-card__eyebrow">Learn</div>
      <h3>{title}</h3>
      <div className="learn-card__list">
        {articles.map((article) => (
          <article key={article.title} className="learn-card__item">
            <strong>{article.title}</strong>
            <span>{article.question}</span>
            <p>{article.summary}</p>
          </article>
        ))}
      </div>
    </section>
  )
}
