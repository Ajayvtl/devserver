import { requestOrFallback } from '../api/client'
import { knowledgeFallback } from './fallbacks'
import type { KnowledgeArticle } from '../types'

export async function getKnowledgeArticles(scope: 'setup' | 'login' | 'dashboard'): Promise<KnowledgeArticle[]> {
  return requestOrFallback<KnowledgeArticle[]>(`/api/knowledge?scope=${encodeURIComponent(scope)}`, knowledgeFallback(scope))
}
