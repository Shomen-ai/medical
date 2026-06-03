// Файл: types/index.ts
// Назначение: общие TypeScript-типы фронта — модели предметной области (Specialty, Doctor, Service, TimeSlot, Review) и мета-данные для UI.

// Медицинская специальность с длительностью базового слота приёма.
export interface Specialty {
  id: string
  name: string
  slot_duration_min: number
}

// Врач клиники с привязкой к специальности и данными для карточки.
export interface Doctor {
  id: string
  full_name: string
  specialty_id: string
  specialty_name: string
  bio: string
  education: string
  photo_url: string | null
  experience_years: number
  is_active: boolean
}

// Услуга клиники с ценой, длительностью и привязкой к специальности.
export interface Service {
  id: string
  name: string
  description: string
  price: number
  duration_min: number
  specialty_id: string
  is_active: boolean
}

// Временной слот записи (ISO 8601 начало/конец).
export interface TimeSlot {
  starts_at: string   // ISO 8601 datetime
  ends_at: string
}

// Отзыв пациента с локализованным текстом (ru/tk) и рейтингом.
export interface Review {
  author: string
  text: { ru: string; tk: string }
  rating: number
  date: string
}

// Декоративные мета-данные специальности для UI (иконка, цвет фона, краткое описание).
export interface SpecialtyMeta {
  icon: string
  description: string
  color: string
}

// Отзыв пациента из API (анонимный): рейтинг, текст, дата + имена врача/услуги для отображения.
export interface ReviewItem {
  id: string
  doctor_id: string
  service_id: string
  rating: number
  text: string
  created_at: string
  doctor_name?: string
  specialty_name?: string
  service_name?: string
  is_hidden?: boolean
}

// Завершённый визит пациента, по которому можно оставить отзыв (для формы).
export interface ReviewableAppt {
  appointment_id: string
  doctor_name: string
  service_name: string
  starts_at: string
}
