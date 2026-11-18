
import { api } from './api'

import type { Fila } from '../types/fila'

export const filaService = {
  // Lista todas as entradas da fila.
  async listar(): Promise<Fila[]> {
    const response = await api.get<Fila[]>('/filas/')

    return response.data
  },

  // Busca uma entrada específica pelo identificador.
  async buscarPorId(id: string): Promise<Fila> {
    const response = await api.get<Fila>(
      `/filas/${id}`
    )

    return response.data
  },

  // Chama o próximo paciente de acordo com
  // a prioridade definida pelo backend.
  async chamarProximo(): Promise<Fila> {
    const response = await api.post<Fila>(
      '/filas/chamar-proximo'
    )

    return response.data
  },

  // Cancela uma entrada da fila.
  // O backend permite cancelar somente pacientes
  // com status Aguardando ou Chamado.
  async cancelar(id: string): Promise<void> {
    await api.patch(`/filas/${id}/cancelar`)
  },
}
