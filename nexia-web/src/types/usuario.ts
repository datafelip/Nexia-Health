
export type PerfilUsuario =
  | 'Administrador'
  | 'Recepcionista'
  | 'Enfermeiro'
  | 'Médico'
  | 'Gestor'

export type StatusUsuario =
  | 'Ativo'
  | 'Inativo'
  | 'Bloqueado'

export interface Usuario {
  id_usuario: string
  nome: string
  email: string
  perfil: PerfilUsuario
  status: StatusUsuario

  ultimo_login?: string | null

  created_at: string
  updated_at: string
}
