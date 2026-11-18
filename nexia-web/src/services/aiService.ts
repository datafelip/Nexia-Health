
import { api } from './api'

// Dados enviados para a API Go.
export interface PrevisaoDemandaRequest {
  data: string
  feriado: 0 | 1
}

// Resposta produzida pelo modelo de Machine Learning.
export interface PrevisaoDemandaResponse {
  data: string
  atendimentos_previstos: number
  modelo: string
  origem_dados: string
  observacao: string
}

// Consulta o backend Go, que se comunica com o Nexia AI.
export async function preverDemanda(
  dados: PrevisaoDemandaRequest
): Promise<PrevisaoDemandaResponse> {
  const response = await api.post<PrevisaoDemandaResponse>(
    '/ai/prever-demanda',
    dados
  )

  return response.data
}
