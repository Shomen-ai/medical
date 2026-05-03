export interface Specialty {
  id: string
  name: string
  slot_duration_min: number
}

export interface Doctor {
  id: string
  full_name: string
  specialty_id: string
  specialty_name: string
  bio: string
  photo_url: string | null
  experience_years: number
  is_active: boolean
}

export interface Service {
  id: string
  name: string
  description: string
  price: number
  duration_min: number
  specialty_id: string
  is_active: boolean
}

export interface TimeSlot {
  starts_at: string   // ISO 8601 datetime
  ends_at: string
}

export interface Review {
  author: string
  text: string
  rating: number
  date: string
}

export interface SpecialtyMeta {
  icon: string
  description: string
  color: string
}
