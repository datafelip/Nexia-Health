
import { api } from './api'

import type {
  Atendimento,
  AtendimentoPayload,
} from '../types/atendimento'

export interface AtendimentoFiltros {
  id_paciente?: string
  data_inicio?: string
  data_fim?: string
}

export interface FinalizarAtendimentoPayload {
  data_fim: string
  diagnostico?: string | null
  observacoes?: string | null
}

export const atendimentoService = {
  // Lista os atendimentos, permitindo filtros opcionais.
  async listar(
    filtros?: AtendimentoFiltros
  ): Promise<Atendimento[]> {
    const response = await api.get<Atendimento[]>(
      '/atendimentos/',
      {
        params: filtros,
      }
    )

    return response.data
  },

  // Busca um atendimento pelo identificador.
  async buscarPorId(id: string): Promise<Atendimento> {
    const response = await api.get<Atendimento>(
      `/atendimentos/${id}`
    )

    return response.data
  },

  // Inicia um atendimento.
  // O backend define data_inicio e deve sincronizar
  // o status da fila por meio de CreateComFila.
  async iniciar(
    dados: AtendimentoPayload
  ): Promise<Atendimento> {
    const response = await api.post<Atendimento>(
      '/atendimentos/',
      dados
    )

    return response.data
  },

  // Finaliza um atendimento.
  // O backend exige data_fim e utiliza FinalizarComFila.
  async finalizar(
    id: string,
    dados: FinalizarAtendimentoPayload
  ): Promise<Atendimento> {
    const response = await api.put<Atendimento>(
      `/atendimentos/${id}`,
      dados
    )

    return response.data
  },
}
