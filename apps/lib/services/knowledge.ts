import { knowledgeFallback } from './fallbacks'
import type { KnowledgeArticle } from '../types'

export async function getKnowledgeArticles(scope: 'setup' | 'login' | 'dashboard'): Promise<KnowledgeArticle[]> {
  return knowledgeFallback(scope)
}
