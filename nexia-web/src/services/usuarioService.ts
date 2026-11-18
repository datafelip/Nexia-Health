
import { api } from './api'

import type { Usuario } from '../types/usuario'

export const usuarioService = {
  // Lista todos os usuários cadastrados.
  async listar(): Promise<Usuario[]> {
    const response = await api.get<Usuario[]>('/usuarios/')

    return response.data
  },

  // Busca um usuário pelo identificador.
  async buscarPorId(id: string): Promise<Usuario> {
    const response = await api.get<Usuario>(
      `/usuarios/${id}`
    )

    return response.data
  },
}
