
import { api } from './api'

import type { Setor } from '../types/setor'

export const setorService = {
  async listar(): Promise<Setor[]> {
    const response = await api.get<Setor[]>('/setores/')

    return response.data
  },

  async buscarPorId(id: string): Promise<Setor> {
    const response = await api.get<Setor>(
      `/setores/${id}`
    )

    return response.data
  },
}
