
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import {
  Activity,
  AlertCircle,
  CalendarClock,
  CheckCircle2,
  ClipboardList,
  Clock3,
  FileText,
  History,
  Loader2,
  Play,
  RefreshCw,
  Search,
  ShieldCheck,
  Stethoscope,
  UserRound,
  Users,
  X,
} from 'lucide-react'

import { atendimentoService } from '../services/atendimentoService'
import { filaService } from '../services/filaService'
import { pacienteService } from '../services/pacienteService'
import { profissionalService } from '../services/profissionalService'
import { usuarioService } from '../services/usuarioService'
import { triagemService } from '../services/triagemService'
import { especialidadeService } from '../services/especialidadeService'
import { setorService } from '../services/setorService'

import type { Atendimento } from '../types/atendimento'
import type { Fila } from '../types/fila'
import type { Paciente } from '../types/paciente'
import type { Profissional } from '../types/profissional'
import type { Triagem } from '../types/triagem'
import type { Especialidade } from '../types/especialidade'
import type { Setor } from '../types/setor'

type UsuarioResumo = {
  id_usuario: string
  nome?: string
  nome_completo?: string
}

type Aba = 'chamados' | 'andamento' | 'historico'

type Modal =
  | { tipo: 'iniciar'; fila: Fila }
  | { tipo: 'finalizar'; atendimento: Atendimento }
  | null

type Dados = {
  atendimentos: Atendimento[]
  filas: Fila[]
  pacientes: Paciente[]
  profissionais: Profissional[]
  usuarios: UsuarioResumo[]
  triagens: Triagem[]
  especialidades: Especialidade[]
  setores: Setor[]
}

const dadosVazios: Dados = {
  atendimentos: [],
  filas: [],
  pacientes: [],
  profissionais: [],
  usuarios: [],
  triagens: [],
  especialidades: [],
  setores: [],
}

const nomesTeste = new Set([
  'paciente teste',
  'paciente teste 2',
])

function pacienteVisivel(nome: string): boolean {
  return !nomesTeste.has(
    nome.trim().toLocaleLowerCase('pt-BR')
  )
}

function formatarData(valor?: string | null): string {
  if (!valor) return '—'

  const data = new Date(valor)

  if (Number.isNaN(data.getTime())) {
    return '—'
  }

  return data.toLocaleString('pt-BR', {
    dateStyle: 'short',
    timeStyle: 'short',
  })
}

function obterIniciais(nome: string): string {
  return nome
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((parte) => parte.charAt(0).toUpperCase())
    .join('')
}

function mensagemErro(erro: unknown): string {
  if (
    typeof erro === 'object' &&
    erro !== null &&
    'response' in erro
  ) {
    const resposta = (
      erro as {
        response?: {
          data?: {
            error?: string
            message?: string
          }
        }
      }
    ).response

    return (
      resposta?.data?.error ??
      resposta?.data?.message ??
      'Não foi possível concluir a operação.'
    )
  }

  return erro instanceof Error
    ? erro.message
    : 'Ocorreu um erro inesperado.'
}

function BadgeRisco({
  risco,
}: {
  risco?: string | null
}) {
  const estilo =
    risco === 'Vermelho'
      ? 'border-red-200 bg-red-50 text-red-700'
      : risco === 'Amarelo'
        ? 'border-amber-200 bg-amber-50 text-amber-700'
        : risco === 'Verde'
          ? 'border-emerald-200 bg-emerald-50 text-emerald-700'
          : 'border-slate-200 bg-slate-50 text-slate-600'

  const ponto =
    risco === 'Vermelho'
      ? 'bg-red-500'
      : risco === 'Amarelo'
        ? 'bg-amber-400'
        : risco === 'Verde'
          ? 'bg-emerald-500'
          : 'bg-slate-400'

  return (
    <span
      className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-bold ${estilo}`}
    >
      <span
        className={`h-2 w-2 rounded-full ${ponto}`}
      />

      {risco ?? 'Risco não identificado'}
    </span>
  )
}

function AvatarPaciente({
  nome,
}: {
  nome: string
}) {
  return (
    <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-teal-50 text-xs font-bold text-teal-700">
      {obterIniciais(nome)}
    </div>
  )
}

function EstadoVazio({
  texto,
}: {
  texto: string
}) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 px-5 py-16 text-center">
      <div className="rounded-2xl bg-slate-50 p-4 text-slate-400">
        <ClipboardList size={30} />
      </div>

      <p className="text-sm font-semibold text-slate-700">
        {texto}
      </p>

      <p className="text-xs text-slate-400">
        Os registros aparecerão aqui quando estiverem disponíveis.
      </p>
    </div>
  )
}

const inputClassName =
  'w-full rounded-xl border border-slate-200 ' +
  'bg-slate-50/70 px-4 py-3 text-sm text-slate-800 ' +
  'outline-none transition placeholder:text-slate-400 ' +
  'focus:border-teal-500 focus:bg-white ' +
  'focus:ring-2 focus:ring-teal-500/10 ' +
  'disabled:cursor-not-allowed disabled:opacity-60'

export default function Atendimentos() {
  const [dados, setDados] = useState<Dados>(dadosVazios)

  const [carregando, setCarregando] = useState(true)
  const [salvando, setSalvando] = useState(false)

  const [erro, setErro] = useState('')
  const [sucesso, setSucesso] = useState('')

  const [aba, setAba] = useState<Aba>('chamados')
  const [busca, setBusca] = useState('')

  const [modal, setModal] = useState<Modal>(null)

  const [profissionalID, setProfissionalID] = useState('')
  const [diagnostico, setDiagnostico] = useState('')
  const [observacoes, setObservacoes] = useState('')

  const ultimaConsulta = useRef(0)
  const operacaoEmAndamento = useRef(false)

  const carregar = useCallback(async () => {
    const consultaAtual = ++ultimaConsulta.current

    setCarregando(true)
    setErro('')

    try {
      const [
        atendimentos,
        filas,
        pacientes,
        profissionais,
        usuarios,
        triagens,
        especialidades,
        setores,
      ] = await Promise.all([
        atendimentoService.listar(),
        filaService.listar(),
        pacienteService.listar(),
        profissionalService.listar(),
        usuarioService.listar(),
        triagemService.listar(),
        especialidadeService.listar(),
        setorService.listar(),
      ])

      if (consultaAtual !== ultimaConsulta.current) {
        return
      }

      setDados({
        atendimentos: atendimentos ?? [],
        filas: filas ?? [],
        pacientes: pacientes ?? [],
        profissionais: profissionais ?? [],
        usuarios: (usuarios ?? []) as UsuarioResumo[],
        triagens: triagens ?? [],
        especialidades: especialidades ?? [],
        setores: setores ?? [],
      })
    } catch (error) {
      if (consultaAtual === ultimaConsulta.current) {
        setErro(mensagemErro(error))
      }
    } finally {
      if (consultaAtual === ultimaConsulta.current) {
        setCarregando(false)
      }
    }
  }, [])

  useEffect(() => {
    void carregar()

    return () => {
      ultimaConsulta.current += 1
    }
  }, [carregar])

  const pacientesPorID = useMemo(
    () =>
      new Map(
        dados.pacientes.map((paciente) => [
          paciente.id,
          paciente,
        ])
      ),
    [dados.pacientes]
  )

  const triagensPorID = useMemo(
    () =>
      new Map(
        dados.triagens.map((triagem) => [
          triagem.id_triagem,
          triagem,
        ])
      ),
    [dados.triagens]
  )

  const profissionaisPorID = useMemo(
    () =>
      new Map(
        dados.profissionais.map((profissional) => [
          profissional.id_profissional,
          profissional,
        ])
      ),
    [dados.profissionais]
  )

  const usuariosPorID = useMemo(
    () =>
      new Map(
        dados.usuarios.map((usuario) => [
          usuario.id_usuario,
          usuario,
        ])
      ),
    [dados.usuarios]
  )

  const especialidadesPorID = useMemo(
    () =>
      new Map(
        dados.especialidades.map((especialidade) => [
          especialidade.id_especialidade,
          especialidade,
        ])
      ),
    [dados.especialidades]
  )

  const setoresPorID = useMemo(
    () =>
      new Map(
        dados.setores.map((setor) => [
          setor.id_setor,
          setor,
        ])
      ),
    [dados.setores]
  )

  function nomePaciente(id: string): string {
    return (
      pacientesPorID.get(id)?.nome ??
      'Paciente não identificado'
    )
  }

  function nomeProfissional(id: string): string {
    const profissional = profissionaisPorID.get(id)

    if (!profissional) {
      return 'Profissional não identificado'
    }

    const usuario = usuariosPorID.get(
      profissional.id_usuario
    )

    return (
      usuario?.nome ??
      usuario?.nome_completo ??
      `${profissional.conselho_profissional} ${profissional.numero_registro}`
    )
  }

  const profissionaisAtivos = useMemo(
    () =>
      dados.profissionais.filter(
        (profissional) =>
          profissional.status === 'Ativo' &&
          especialidadesPorID.has(
            profissional.id_especialidade
          ) &&
          setoresPorID.has(profissional.id_setor)
      ),
    [
      dados.profissionais,
      especialidadesPorID,
      setoresPorID,
    ]
  )

  const profissionalSelecionado =
    profissionaisPorID.get(profissionalID)

  const chamados = dados.filas.filter(
    (fila) =>
      fila.status === 'Chamado' &&
      pacienteVisivel(nomePaciente(fila.id_paciente))
  )

  const emAndamento = dados.atendimentos.filter(
    (atendimento) =>
      !atendimento.data_fim &&
      pacienteVisivel(
        nomePaciente(atendimento.id_paciente)
      )
  )

  const historico = dados.atendimentos.filter(
    (atendimento) =>
      Boolean(atendimento.data_fim) &&
      pacienteVisivel(
        nomePaciente(atendimento.id_paciente)
      )
  )

  const termoBusca = busca
    .trim()
    .toLocaleLowerCase('pt-BR')

  const chamadosFiltrados = chamados.filter((fila) =>
    nomePaciente(fila.id_paciente)
      .toLocaleLowerCase('pt-BR')
      .includes(termoBusca)
  )

  const andamentoFiltrado = emAndamento.filter(
    (atendimento) =>
      nomePaciente(atendimento.id_paciente)
        .toLocaleLowerCase('pt-BR')
        .includes(termoBusca)
  )

  const historicoFiltrado = historico.filter(
    (atendimento) =>
      nomePaciente(atendimento.id_paciente)
        .toLocaleLowerCase('pt-BR')
        .includes(termoBusca)
  )

  const indicadores = [
    {
      titulo: 'Pacientes chamados',
      valor: chamados.length,
      descricao: 'Aguardando início do atendimento',
      icone: Clock3,
      corIcone: 'text-blue-600',
      corFundo: 'bg-blue-50',
      corBarra: 'bg-blue-500',
    },
    {
      titulo: 'Em atendimento',
      valor: emAndamento.length,
      descricao: 'Atendimentos em andamento',
      icone: Activity,
      corIcone: 'text-violet-600',
      corFundo: 'bg-violet-50',
      corBarra: 'bg-violet-500',
    },
    {
      titulo: 'Finalizados',
      valor: historico.length,
      descricao: 'Atendimentos concluídos',
      icone: CheckCircle2,
      corIcone: 'text-emerald-600',
      corFundo: 'bg-emerald-50',
      corBarra: 'bg-emerald-500',
    },
  ]

  function abrirInicio(fila: Fila) {
    setProfissionalID('')
    setDiagnostico('')
    setObservacoes('')
    setErro('')
    setSucesso('')
    setModal({ tipo: 'iniciar', fila })
  }

  function abrirFinalizacao(atendimento: Atendimento) {
    setDiagnostico(atendimento.diagnostico ?? '')
    setObservacoes(atendimento.observacoes ?? '')
    setErro('')
    setSucesso('')

    setModal({
      tipo: 'finalizar',
      atendimento,
    })
  }

  function fecharModal() {
    if (salvando || operacaoEmAndamento.current) {
      return
    }

    setModal(null)
    setProfissionalID('')
    setDiagnostico('')
    setObservacoes('')
    setErro('')
  }

  async function iniciarAtendimento() {
    if (
      !modal ||
      modal.tipo !== 'iniciar' ||
      operacaoEmAndamento.current
    ) {
      return
    }

    const profissional =
      profissionaisPorID.get(profissionalID)

    if (
      !profissional ||
      profissional.status !== 'Ativo'
    ) {
      setErro('Selecione um profissional ativo.')
      return
    }

    if (
      !especialidadesPorID.has(
        profissional.id_especialidade
      ) ||
      !setoresPorID.has(profissional.id_setor)
    ) {
      setErro(
        'O profissional não possui especialidade ou setor válido.'
      )
      return
    }

    const fila = modal.fila
    const paciente = nomePaciente(fila.id_paciente)

    operacaoEmAndamento.current = true
    setSalvando(true)
    setErro('')
    setSucesso('')

    try {
      await atendimentoService.iniciar({
        id_paciente: fila.id_paciente,
        id_triagem: fila.id_triagem,
        id_profissional: profissional.id_profissional,
        id_especialidade: profissional.id_especialidade,
        id_setor: profissional.id_setor,
      })

      setModal(null)
      setAba('andamento')
      setBusca('')

      await carregar()

      setSucesso(
        `Atendimento de ${paciente} iniciado com sucesso!`
      )
    } catch (error) {
      setErro(mensagemErro(error))
    } finally {
      operacaoEmAndamento.current = false
      setSalvando(false)
    }
  }

  async function finalizarAtendimento() {
    if (
      !modal ||
      modal.tipo !== 'finalizar' ||
      operacaoEmAndamento.current
    ) {
      return
    }

    const atendimento = modal.atendimento
    const paciente = nomePaciente(
      atendimento.id_paciente
    )

    operacaoEmAndamento.current = true
    setSalvando(true)
    setErro('')
    setSucesso('')

    try {
      await atendimentoService.finalizar(
        atendimento.id_atendimento,
        {
          data_fim: new Date().toISOString(),
          diagnostico: diagnostico.trim() || null,
          observacoes: observacoes.trim() || null,
        }
      )

      setModal(null)
      setAba('historico')
      setBusca('')

      await carregar()

      setSucesso(
        `Atendimento de ${paciente} finalizado com sucesso!`
      )
    } catch (error) {
      setErro(mensagemErro(error))
    } finally {
      operacaoEmAndamento.current = false
      setSalvando(false)
    }
  }

  const quantidadeExibida =
    aba === 'chamados'
      ? chamadosFiltrados.length
      : aba === 'andamento'
        ? andamentoFiltrado.length
        : historicoFiltrado.length

  return (
    <div className="space-y-7 pb-4">
      {/* Cabeçalho premium */}
      <section className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-6 shadow-sm sm:p-7">
        <div className="pointer-events-none absolute -right-16 -top-20 h-56 w-56 rounded-full bg-teal-50 blur-3xl" />

        <div className="relative flex flex-col justify-between gap-5 md:flex-row md:items-center">
          <div className="flex items-start gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-teal-600 text-white shadow-lg shadow-teal-100">
              <Stethoscope size={25} />
            </div>

            <div>
              <div className="mb-2 flex flex-wrap items-center gap-2">
                <span className="rounded-full border border-teal-100 bg-teal-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-teal-700">
                  Gestão assistencial
                </span>

                <span className="inline-flex items-center gap-1.5 text-xs font-medium text-slate-500">
                  <span className="h-2 w-2 rounded-full bg-teal-500" />
                  Fluxo clínico
                </span>
              </div>

              <h2 className="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
                Atendimentos
              </h2>

              <p className="mt-2 max-w-xl text-sm leading-relaxed text-slate-500">
                Inicie, acompanhe e finalize atendimentos,
                mantendo o histórico clínico organizado.
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={() => void carregar()}
            disabled={carregando || salvando}
            className="inline-flex shrink-0 items-center justify-center gap-2 self-start rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 transition hover:border-teal-200 hover:bg-teal-50 hover:text-teal-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <RefreshCw
              size={18}
              className={carregando ? 'animate-spin' : ''}
            />

            Atualizar
          </button>
        </div>
      </section>

      {/* Mensagens */}
      {sucesso && (
        <div
          role="status"
          className="flex items-center justify-between gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm font-medium text-emerald-700"
        >
          <div className="flex items-center gap-2">
            <CheckCircle2
              size={19}
              className="shrink-0"
            />

            <span>{sucesso}</span>
          </div>

          <button
            type="button"
            onClick={() => setSucesso('')}
            aria-label="Fechar mensagem"
            className="rounded-lg p-1 transition hover:bg-emerald-100"
          >
            <X size={17} />
          </button>
        </div>
      )}

      {erro && !modal && (
        <div
          role="alert"
          className="flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
        >
          <AlertCircle
            size={19}
            className="shrink-0"
          />

          <span>{erro}</span>
        </div>
      )}

      {/* Indicadores */}
      <section>
        <div className="mb-4">
          <h3 className="text-base font-bold text-slate-900">
            Visão operacional
          </h3>

          <p className="mt-1 text-xs text-slate-500">
            Acompanhamento do fluxo de atendimentos
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-3">
          {indicadores.map((indicador) => {
            const Icone = indicador.icone

            return (
              <div
                key={indicador.titulo}
                className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
              >
                <div
                  className={`absolute inset-x-0 top-0 h-1 ${indicador.corBarra}`}
                />

                <div className="flex items-start justify-between gap-3">
                  <div
                    className={`flex h-12 w-12 items-center justify-center rounded-2xl ${indicador.corFundo} ${indicador.corIcone}`}
                  >
                    <Icone size={23} />
                  </div>

                  <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
                    Atendimento
                  </span>
                </div>

                <p className="mt-5 text-sm font-medium text-slate-500">
                  {indicador.titulo}
                </p>

                <p className="mt-2 text-4xl font-bold tracking-tight text-slate-900">
                  {carregando &&
                  dados.pacientes.length === 0
                    ? '—'
                    : indicador.valor}
                </p>

                <p className="mt-2 text-xs text-slate-400">
                  {indicador.descricao}
                </p>
              </div>
            )
          })}
        </div>
      </section>

      {/* Painel principal */}
      <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        {/* Abas */}
        <div className="flex gap-6 overflow-x-auto border-b border-slate-100 px-5 sm:px-6">
          {(
            [
              {
                valor: 'chamados',
                rotulo: 'Chamados',
                total: chamados.length,
                icone: Users,
              },
              {
                valor: 'andamento',
                rotulo: 'Em atendimento',
                total: emAndamento.length,
                icone: Activity,
              },
              {
                valor: 'historico',
                rotulo: 'Histórico',
                total: historico.length,
                icone: History,
              },
            ] as const
          ).map((item) => {
            const Icone = item.icone
            const selecionada = aba === item.valor

            return (
              <button
                key={item.valor}
                type="button"
                onClick={() => {
                  setAba(item.valor)
                  setBusca('')
                }}
                aria-pressed={selecionada}
                className={`inline-flex shrink-0 items-center gap-2 border-b-2 py-5 text-sm font-semibold transition ${
                  selecionada
                    ? 'border-teal-600 text-teal-700'
                    : 'border-transparent text-slate-500 hover:text-slate-800'
                }`}
              >
                <Icone size={18} />

                {item.rotulo}

                <span
                  className={`rounded-full px-2.5 py-0.5 text-xs ${
                    selecionada
                      ? 'bg-teal-50 text-teal-700'
                      : 'bg-slate-100 text-slate-600'
                  }`}
                >
                  {item.total}
                </span>
              </button>
            )
          })}
        </div>

        {/* Pesquisa */}
        <div className="flex flex-col justify-between gap-4 border-b border-slate-100 p-5 sm:p-6 lg:flex-row lg:items-center">
          <div>
            <h3 className="flex items-center gap-2 text-base font-bold text-slate-900">
              <ClipboardList
                size={19}
                className="text-teal-600"
              />

              {aba === 'chamados'
                ? 'Pacientes chamados'
                : aba === 'andamento'
                  ? 'Atendimentos em andamento'
                  : 'Histórico de atendimentos'}
            </h3>

            <p className="mt-1 text-xs text-slate-500">
              {quantidadeExibida} registro(s) encontrado(s)
            </p>
          </div>

          <div className="relative w-full lg:w-80">
            <Search
              size={18}
              className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400"
            />

            <input
              type="search"
              value={busca}
              onChange={(event) =>
                setBusca(event.target.value)
              }
              placeholder="Buscar paciente..."
              aria-label="Buscar paciente"
              className="w-full rounded-xl border border-slate-200 bg-slate-50 py-2.5 pl-10 pr-4 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-teal-500 focus:bg-white focus:ring-2 focus:ring-teal-500/10"
            />
          </div>
        </div>

        {/* Conteúdo das abas */}
        {carregando ? (
          <div
            role="status"
            className="flex flex-col items-center justify-center gap-3 px-5 py-16"
          >
            <Loader2
              size={29}
              className="animate-spin text-teal-600"
            />

            <p className="text-sm text-slate-500">
              Carregando atendimentos...
            </p>
          </div>
        ) : aba === 'chamados' ? (
          chamadosFiltrados.length === 0 ? (
            <EstadoVazio texto="Nenhum paciente chamado encontrado." />
          ) : (
            <div className="divide-y divide-slate-100">
              {chamadosFiltrados.map((fila) => {
                const triagem = triagensPorID.get(
                  fila.id_triagem
                )

                const paciente = nomePaciente(
                  fila.id_paciente
                )

                return (
                  <div
                    key={fila.id_fila}
                    className="flex flex-col justify-between gap-5 p-5 transition hover:bg-teal-50/30 sm:p-6 lg:flex-row lg:items-center"
                  >
                    <div className="flex items-start gap-4">
                      <AvatarPaciente nome={paciente} />

                      <div className="min-w-0">
                        <h4 className="text-sm font-bold text-slate-900">
                          {paciente}
                        </h4>

                        <div className="mt-2 flex flex-wrap items-center gap-2">
                          <BadgeRisco
                            risco={
                              triagem?.classificacao_risco
                            }
                          />

                          <span className="inline-flex items-center gap-1.5 rounded-full border border-blue-200 bg-blue-50 px-3 py-1.5 text-xs font-semibold text-blue-700">
                            <span className="h-2 w-2 rounded-full bg-blue-500" />
                            Chamado
                          </span>
                        </div>

                        <p className="mt-3 flex items-center gap-1.5 text-xs text-slate-500">
                          <CalendarClock size={14} />

                          Entrada na fila:{' '}
                          {formatarData(fila.data_entrada)}
                        </p>
                      </div>
                    </div>

                    <button
                      type="button"
                      onClick={() => abrirInicio(fila)}
                      disabled={salvando}
                      className="inline-flex shrink-0 items-center justify-center gap-2 rounded-xl bg-teal-600 px-5 py-3 text-sm font-semibold text-white shadow-sm shadow-teal-100 transition hover:bg-teal-700 disabled:opacity-50"
                    >
                      <Play size={17} />
                      Iniciar atendimento
                    </button>
                  </div>
                )
              })}
            </div>
          )
        ) : aba === 'andamento' ? (
          andamentoFiltrado.length === 0 ? (
            <EstadoVazio texto="Nenhum atendimento em andamento." />
          ) : (
            <div className="divide-y divide-slate-100">
              {andamentoFiltrado.map((atendimento) => {
                const paciente = nomePaciente(
                  atendimento.id_paciente
                )

                return (
                  <div
                    key={atendimento.id_atendimento}
                    className="flex flex-col justify-between gap-5 p-5 transition hover:bg-teal-50/30 sm:p-6 lg:flex-row lg:items-center"
                  >
                    <div className="flex items-start gap-4">
                      <AvatarPaciente nome={paciente} />

                      <div>
                        <h4 className="text-sm font-bold text-slate-900">
                          {paciente}
                        </h4>

                        <p className="mt-1 flex items-center gap-1.5 text-sm text-slate-600">
                          <Stethoscope
                            size={15}
                            className="text-slate-400"
                          />

                          {nomeProfissional(
                            atendimento.id_profissional
                          )}
                        </p>

                        <div className="mt-3 flex flex-wrap items-center gap-3">
                          <span className="inline-flex items-center gap-1.5 rounded-full border border-violet-200 bg-violet-50 px-3 py-1.5 text-xs font-semibold text-violet-700">
                            <span className="h-2 w-2 rounded-full bg-violet-500" />
                            Em atendimento
                          </span>

                          <span className="inline-flex items-center gap-1.5 text-xs text-slate-500">
                            <Clock3 size={14} />

                            Iniciado em:{' '}
                            {formatarData(
                              atendimento.data_inicio
                            )}
                          </span>
                        </div>
                      </div>
                    </div>

                    <button
                      type="button"
                      onClick={() =>
                        abrirFinalizacao(atendimento)
                      }
                      disabled={salvando}
                      className="inline-flex shrink-0 items-center justify-center gap-2 rounded-xl bg-emerald-600 px-5 py-3 text-sm font-semibold text-white transition hover:bg-emerald-700 disabled:opacity-50"
                    >
                      <CheckCircle2 size={17} />
                      Finalizar atendimento
                    </button>
                  </div>
                )
              })}
            </div>
          )
        ) : historicoFiltrado.length === 0 ? (
          <EstadoVazio texto="Nenhum atendimento finalizado encontrado." />
        ) : (
          <div className="divide-y divide-slate-100">
            {historicoFiltrado.map((atendimento) => {
              const paciente = nomePaciente(
                atendimento.id_paciente
              )

              return (
                <div
                  key={atendimento.id_atendimento}
                  className="p-5 transition hover:bg-slate-50/60 sm:p-6"
                >
                  <div className="flex flex-wrap items-start justify-between gap-4">
                    <div className="flex items-start gap-3">
                      <AvatarPaciente nome={paciente} />

                      <div>
                        <h4 className="text-sm font-bold text-slate-900">
                          {paciente}
                        </h4>

                        <p className="mt-1 flex items-center gap-1.5 text-sm text-slate-600">
                          <Stethoscope
                            size={15}
                            className="text-slate-400"
                          />

                          {nomeProfissional(
                            atendimento.id_profissional
                          )}
                        </p>
                      </div>
                    </div>

                    <span className="inline-flex items-center gap-1.5 rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1.5 text-xs font-semibold text-emerald-700">
                      <CheckCircle2 size={14} />
                      Finalizado
                    </span>
                  </div>

                  <div className="mt-5 grid gap-3 sm:grid-cols-2">
                    <div className="rounded-xl border border-slate-100 bg-slate-50/70 p-4">
                      <p className="flex items-center gap-1.5 text-xs font-medium text-slate-500">
                        <CalendarClock size={14} />
                        Início
                      </p>

                      <p className="mt-2 text-sm font-semibold text-slate-800">
                        {formatarData(
                          atendimento.data_inicio
                        )}
                      </p>
                    </div>

                    <div className="rounded-xl border border-slate-100 bg-slate-50/70 p-4">
                      <p className="flex items-center gap-1.5 text-xs font-medium text-slate-500">
                        <CheckCircle2 size={14} />
                        Finalização
                      </p>

                      <p className="mt-2 text-sm font-semibold text-slate-800">
                        {formatarData(
                          atendimento.data_fim
                        )}
                      </p>
                    </div>
                  </div>

                  {(atendimento.diagnostico ||
                    atendimento.observacoes) && (
                    <div className="mt-4 space-y-3 rounded-xl border border-slate-100 bg-white p-4">
                      {atendimento.diagnostico && (
                        <div>
                          <p className="flex items-center gap-2 text-xs font-bold text-slate-700">
                            <FileText
                              size={15}
                              className="text-teal-600"
                            />

                            Diagnóstico
                          </p>

                          <p className="mt-1 whitespace-pre-wrap text-sm leading-relaxed text-slate-600">
                            {atendimento.diagnostico}
                          </p>
                        </div>
                      )}

                      {atendimento.observacoes && (
                        <div>
                          <p className="flex items-center gap-2 text-xs font-bold text-slate-700">
                            <ClipboardList
                              size={15}
                              className="text-teal-600"
                            />

                            Observações
                          </p>

                          <p className="mt-1 whitespace-pre-wrap text-sm leading-relaxed text-slate-600">
                            {atendimento.observacoes}
                          </p>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              )
            })}
          </div>
        )}

        {!carregando && (
          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-slate-100 bg-slate-50/50 px-5 py-4 sm:px-6">
            <p className="text-xs text-slate-500">
              Exibindo{' '}
              <strong className="text-slate-700">
                {quantidadeExibida}
              </strong>{' '}
              registro(s)
            </p>

            <span className="inline-flex items-center gap-1.5 text-xs font-medium text-teal-700">
              <ShieldCheck size={14} />
              Gestão assistencial
            </span>
          </div>
        )}
      </section>

      {/* Modal premium */}
      {modal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 p-4 backdrop-blur-sm">
          <div
            role="dialog"
            aria-modal="true"
            aria-labelledby="titulo-modal-atendimento"
            className="flex max-h-[90vh] w-full max-w-xl flex-col overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-2xl"
          >
            {/* Cabeçalho do modal */}
            <div className="flex shrink-0 items-start justify-between gap-4 border-b border-slate-100 bg-gradient-to-r from-white to-teal-50/60 p-5 sm:p-6">
              <div className="flex items-start gap-3">
                <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-teal-600 text-white">
                  {modal.tipo === 'iniciar' ? (
                    <Play size={21} />
                  ) : (
                    <CheckCircle2 size={21} />
                  )}
                </div>

                <div>
                  <span className="text-[10px] font-bold uppercase tracking-wider text-teal-700">
                    Gestão clínica
                  </span>

                  <h2
                    id="titulo-modal-atendimento"
                    className="mt-1 text-lg font-bold text-slate-900"
                  >
                    {modal.tipo === 'iniciar'
                      ? 'Iniciar atendimento'
                      : 'Finalizar atendimento'}
                  </h2>

                  <p className="mt-1 text-xs text-slate-500">
                    {modal.tipo === 'iniciar'
                      ? 'Defina o profissional responsável.'
                      : 'Registre o encerramento clínico.'}
                  </p>
                </div>
              </div>

              <button
                type="button"
                onClick={fecharModal}
                disabled={salvando}
                aria-label="Fechar formulário"
                className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-50 disabled:opacity-50"
              >
                <X size={18} />
              </button>
            </div>

            {/* Corpo do modal */}
            <div className="space-y-6 overflow-y-auto p-5 sm:p-6">
              <div className="flex items-center gap-3 rounded-xl border border-teal-100 bg-teal-50/60 p-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-white text-teal-700">
                  <UserRound size={19} />
                </div>

                <div>
                  <p className="text-[11px] font-semibold uppercase tracking-wider text-teal-700">
                    Paciente
                  </p>

                  <p className="mt-1 text-sm font-bold text-slate-900">
                    {nomePaciente(
                      modal.tipo === 'iniciar'
                        ? modal.fila.id_paciente
                        : modal.atendimento.id_paciente
                    )}
                  </p>
                </div>
              </div>

              {modal.tipo === 'iniciar' ? (
                <section className="space-y-5">
                  <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                      <Stethoscope size={20} />
                    </div>

                    <div>
                      <h3 className="text-sm font-bold text-slate-900">
                        Profissional responsável
                      </h3>

                      <p className="mt-1 text-xs text-slate-500">
                        Selecione quem realizará o atendimento.
                      </p>
                    </div>
                  </div>

                  <div>
                    <label
                      htmlFor="profissional-atendimento"
                      className="mb-2 block text-sm font-semibold text-slate-700"
                    >
                      Profissional{' '}
                      <span className="text-red-500">*</span>
                    </label>

                    <select
                      id="profissional-atendimento"
                      value={profissionalID}
                      onChange={(event) =>
                        setProfissionalID(event.target.value)
                      }
                      disabled={salvando}
                      className={inputClassName}
                    >
                      <option value="">
                        Selecione um profissional
                      </option>

                      {profissionaisAtivos.map(
                        (profissional) => (
                          <option
                            key={profissional.id_profissional}
                            value={profissional.id_profissional}
                          >
                            {nomeProfissional(
                              profissional.id_profissional
                            )}
                          </option>
                        )
                      )}
                    </select>
                  </div>

                  {profissionalSelecionado && (
                    <div className="grid gap-3 sm:grid-cols-2">
                      <div className="rounded-xl border border-slate-100 bg-slate-50 p-4">
                        <p className="text-xs font-medium text-slate-500">
                          Especialidade
                        </p>

                        <p className="mt-2 text-sm font-bold text-slate-800">
                          {especialidadesPorID.get(
                            profissionalSelecionado.id_especialidade
                          )?.nome ?? 'Não identificada'}
                        </p>
                      </div>

                      <div className="rounded-xl border border-slate-100 bg-slate-50 p-4">
                        <p className="text-xs font-medium text-slate-500">
                          Setor
                        </p>

                        <p className="mt-2 text-sm font-bold text-slate-800">
                          {setoresPorID.get(
                            profissionalSelecionado.id_setor
                          )?.nome ?? 'Não identificado'}
                        </p>
                      </div>
                    </div>
                  )}

                  {profissionaisAtivos.length === 0 && (
                    <div className="flex items-start gap-2 rounded-xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
                      <AlertCircle
                        size={18}
                        className="shrink-0"
                      />

                      Não há profissionais ativos com
                      especialidade e setor cadastrados.
                    </div>
                  )}

                  <div className="flex items-start gap-2 rounded-xl border border-blue-100 bg-blue-50 p-4 text-xs leading-relaxed text-blue-800">
                    <ShieldCheck
                      size={17}
                      className="shrink-0"
                    />

                    O atendimento será iniciado utilizando a
                    especialidade e o setor vinculados ao
                    profissional selecionado.
                  </div>
                </section>
              ) : (
                <section className="space-y-5">
                  <div className="flex items-center gap-3 border-b border-slate-100 pb-4">
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                      <FileText size={20} />
                    </div>

                    <div>
                      <h3 className="text-sm font-bold text-slate-900">
                        Registro clínico
                      </h3>

                      <p className="mt-1 text-xs text-slate-500">
                        Informações de encerramento do atendimento.
                      </p>
                    </div>
                  </div>

                  <div>
                    <label
                      htmlFor="diagnostico-atendimento"
                      className="mb-2 block text-sm font-semibold text-slate-700"
                    >
                      Diagnóstico
                    </label>

                    <textarea
                      id="diagnostico-atendimento"
                      value={diagnostico}
                      onChange={(event) =>
                        setDiagnostico(event.target.value)
                      }
                      disabled={salvando}
                      rows={3}
                      placeholder="Registre o diagnóstico..."
                      className={`${inputClassName} resize-y`}
                    />
                  </div>

                  <div>
                    <label
                      htmlFor="observacoes-atendimento"
                      className="mb-2 block text-sm font-semibold text-slate-700"
                    >
                      Observações
                    </label>

                    <textarea
                      id="observacoes-atendimento"
                      value={observacoes}
                      onChange={(event) =>
                        setObservacoes(event.target.value)
                      }
                      disabled={salvando}
                      rows={4}
                      placeholder="Registre as observações sobre o atendimento..."
                      className={`${inputClassName} resize-y`}
                    />
                  </div>

                  <div className="flex items-start gap-2 rounded-xl border border-amber-200 bg-amber-50 p-4 text-xs leading-relaxed text-amber-800">
                    <AlertCircle
                      size={17}
                      className="shrink-0"
                    />

                    Ao confirmar, o atendimento será encerrado
                    e a entrada correspondente na fila será
                    marcada como Finalizado.
                  </div>
                </section>
              )}

              {erro && (
                <div
                  role="alert"
                  className="flex items-start gap-2 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
                >
                  <AlertCircle
                    size={18}
                    className="shrink-0"
                  />

                  <span>{erro}</span>
                </div>
              )}
            </div>

            {/* Rodapé do modal */}
            <div className="flex shrink-0 flex-col-reverse gap-3 border-t border-slate-100 bg-slate-50/70 p-5 sm:flex-row sm:justify-end sm:p-6">
              <button
                type="button"
                onClick={fecharModal}
                disabled={salvando}
                className="inline-flex items-center justify-center rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 disabled:opacity-50"
              >
                Cancelar
              </button>

              <button
                type="button"
                onClick={() =>
                  void (
                    modal.tipo === 'iniciar'
                      ? iniciarAtendimento()
                      : finalizarAtendimento()
                  )
                }
                disabled={
                  salvando ||
                  (modal.tipo === 'iniciar' &&
                    !profissionalID)
                }
                className={`inline-flex items-center justify-center gap-2 rounded-xl px-5 py-3 text-sm font-semibold text-white transition disabled:cursor-not-allowed disabled:opacity-50 ${
                  modal.tipo === 'iniciar'
                    ? 'bg-teal-600 hover:bg-teal-700'
                    : 'bg-emerald-600 hover:bg-emerald-700'
                }`}
              >
                {salvando ? (
                  <Loader2
                    size={17}
                    className="animate-spin"
                  />
                ) : modal.tipo === 'iniciar' ? (
                  <Play size={17} />
                ) : (
                  <CheckCircle2 size={17} />
                )}

                {salvando
                  ? 'Processando...'
                  : modal.tipo === 'iniciar'
                    ? 'Iniciar atendimento'
                    : 'Confirmar finalização'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
