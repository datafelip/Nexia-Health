
import { api } from './api'

import type { Especialidade } from '../types/especialidade'

export const especialidadeService = {
  async listar(): Promise<Especialidade[]> {
    const response = await api.get<Especialidade[]>(
      '/especialidades/'
    )

    return response.data
  },

  async buscarPorId(id: string): Promise<Especialidade> {
    const response = await api.get<Especialidade>(
      `/especialidades/${id}`
    )

    return response.data
  },
}
