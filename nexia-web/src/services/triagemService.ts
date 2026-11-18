
import { api } from './api'

import type {
  Triagem,
  TriagemPayload,
  ReclassificarTriagemPayload,
  CancelarTriagemPayload,
} from '../types/triagem'

export const triagemService = {
  // Lista todas as triagens cadastradas.
  async listar(): Promise<Triagem[]> {
    const response = await api.get<Triagem[]>('/triagens/')
    return response.data
  },

  // Busca uma triagem pelo identificador.
  async buscarPorId(id: string): Promise<Triagem> {
    const response = await api.get<Triagem>(
      `/triagens/${id}`
    )

    return response.data
  },

  // Cadastra uma triagem e insere o paciente na fila.
  async cadastrar(dados: TriagemPayload): Promise<Triagem> {
    const response = await api.post<Triagem>(
      '/triagens/',
      dados
    )

    return response.data
  },

  // Mantido por compatibilidade com o frontend.
  // O backend bloqueia a edição direta de triagens.
  async atualizar(
    id: string,
    dados: TriagemPayload
  ): Promise<Triagem> {
    const response = await api.put<Triagem>(
      `/triagens/${id}`,
      dados
    )

    return response.data
  },

  // Reclassifica o risco e reorganiza a fila.
  async reclassificar(
    id: string,
    dados: ReclassificarTriagemPayload
  ): Promise<Triagem> {
    const response = await api.patch<Triagem>(
      `/triagens/${id}/reclassificar`,
      dados
    )

    return response.data
  },

  // Cancela a triagem, preservando o histórico clínico.
  // O backend valida o status e reorganiza a fila.
  async cancelar(
    id: string,
    dados: CancelarTriagemPayload
  ): Promise<Triagem> {
    const response = await api.patch<Triagem>(
      `/triagens/${id}/cancelar`,
      dados
    )

    return response.data
  },

  // Mantido por compatibilidade com o frontend.
  // O backend bloqueia a exclusão física de triagens.
  async excluir(id: string): Promise<void> {
    await api.delete(`/triagens/${id}`)
  },
}
