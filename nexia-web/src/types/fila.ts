
import type { ClassificacaoRisco } from './triagem'

export type StatusFila =
  | 'Aguardando'
  | 'Chamado'
  | 'Em Atendimento'
  | 'Finalizado'
  | 'Cancelado'

export interface Fila {
  id_fila: string
  id_paciente: string
  id_triagem: string
  posicao: number
  status: StatusFila
  data_entrada: string
  created_at: string
  updated_at: string
}

export interface FilaDetalhada {
  fila: Fila
  nomePaciente: string
  classificacaoRisco: ClassificacaoRisco | null
  dataTriagem: string
}
