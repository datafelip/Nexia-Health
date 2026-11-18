
export type StatusProfissional =
  | 'Ativo'
  | 'Inativo'
  | 'Bloqueado'

export interface Profissional {
  id_profissional: string
  id_usuario: string
  id_especialidade: string
  id_setor: string

  cpf: string
  conselho_profissional: string
  numero_registro: string

  telefone?: string | null

  status: StatusProfissional

  created_at: string
  updated_at: string
}
