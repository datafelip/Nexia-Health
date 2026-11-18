
import { api } from './api'
import type { Paciente, PacientePayload } from '../types/paciente'


export const pacienteService = {
  async listar(): Promise<Paciente[]> {
    const response = await api.get<Paciente[]>(
      '/pacientes/'
    )

    return response.data
  },

  async buscarPorId(id: string): Promise<Paciente> {
    const response = await api.get<Paciente>(
      `/pacientes/${id}`
    )

    return response.data
  },

  async cadastrar(
    dados: PacientePayload
  ): Promise<Paciente> {
    const response = await api.post<Paciente>(
      '/pacientes/',
      dados
    )

    return response.data
  },

  async atualizar(
    id: string,
    dados: PacientePayload
  ): Promise<Paciente> {
    const response = await api.put<Paciente>(
      `/pacientes/${id}`,
      dados
    )

    return response.data
  },

  async excluir(id: string): Promise<void> {
    await api.delete(`/pacientes/${id}`)
  },
}
