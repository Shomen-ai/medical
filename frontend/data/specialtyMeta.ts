import type { SpecialtyMeta } from '~/types'

export const SPECIALTY_META: Record<string, SpecialtyMeta> = {
  'Косметология': {
    icon: '✨',
    description: 'Аппаратная косметология, инъекции, пилинги и комплексный уход за кожей.',
    color: 'bg-pink-50',
  },
  'Дерматология': {
    icon: '🔬',
    description: 'Диагностика и лечение акне, дерматитов, розацеа. Удаление новообразований.',
    color: 'bg-blue-50',
  },
  'Трихология': {
    icon: '💆',
    description: 'Лечение выпадения волос, заболеваний кожи головы, PRP-терапия.',
    color: 'bg-green-50',
  },
  'Эстетическая медицина': {
    icon: '💎',
    description: 'Нитевой лифтинг, объёмное моделирование, коррекция возрастных изменений.',
    color: 'bg-purple-50',
  },
}

export const getSpecialtyMeta = (name: string): SpecialtyMeta =>
  SPECIALTY_META[name] ?? { icon: '🏥', description: '', color: 'bg-gray-50' }
