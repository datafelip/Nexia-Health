
import { api } from './api'

import type { Profissional } from '../types/profissional'

export const profissionalService = {
  // Busca todos os profissionais cadastrados na API.
  async listar(): Promise<Profissional[]> {
    const response = await api.get<Profissional[]>(
      '/profissionais/'
    )

    return response.data
  },

  // Busca um profissional específico pelo ID.
  async buscarPorId(id: string): Promise<Profissional> {
    const response = await api.get<Profissional>(
      `/profissionais/${id}`
    )

    return response.data
  },
}
