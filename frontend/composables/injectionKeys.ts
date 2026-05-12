// Файл: composables/injectionKeys.ts
// Назначение: типизированные ключи Symbol для provide/inject — позволяют дочерним компонентам получать списки специальностей, врачей и услуг из корневого app.vue.
import type { InjectionKey, ComputedRef } from 'vue'
import type { Specialty, Doctor, Service } from '~/types'

// Ключ для inject списка специальностей.
export const SpecialtiesKey: InjectionKey<ComputedRef<Specialty[]>> = Symbol('specialties')
// Ключ для inject списка врачей.
export const DoctorsKey:     InjectionKey<ComputedRef<Doctor[]>>    = Symbol('doctors')
// Ключ для inject списка услуг.
export const ServicesKey:    InjectionKey<ComputedRef<Service[]>>   = Symbol('services')
