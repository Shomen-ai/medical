import type { InjectionKey, ComputedRef } from 'vue'
import type { Specialty, Doctor, Service } from '~/types'

export const SpecialtiesKey: InjectionKey<ComputedRef<Specialty[]>> = Symbol('specialties')
export const DoctorsKey:     InjectionKey<ComputedRef<Doctor[]>>    = Symbol('doctors')
export const ServicesKey:    InjectionKey<ComputedRef<Service[]>>   = Symbol('services')
