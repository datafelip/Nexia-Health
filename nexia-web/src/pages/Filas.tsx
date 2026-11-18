
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import axios from 'axios'

import {
  AlertCircle,
  CheckCircle2,
  Clock3,
  History,
  ListOrdered,
  Loader2,
  Megaphone,
  RefreshCw,
  Search,
  ShieldCheck,
  Stethoscope,
  UserRound,
  Users,
  X,
  XCircle,
} from 'lucide-react'

import { filaService } from '../services/filaService'
import { pacienteService } from '../services/pacienteService'
import { triagemService } from '../services/triagemService'

import type { Fila, FilaDetalhada, StatusFila } from '../types/fila'
import type { ClassificacaoRisco } from '../types/triagem'

type Aba = 'ativa' | 'historico'
type FiltroStatus = 'Todos' | StatusFila

const nomesTeste = new Set([
  'paciente teste',
  'paciente teste 2',
])

const estilosRisco: Record<ClassificacaoRisco, string> = {
  Vermelho: 'border-red-200 bg-red-50 text-red-700',
  Amarelo: 'border-amber-200 bg-amber-50 text-amber-700',
  Verde: 'border-emerald-200 bg-emerald-50 text-emerald-700',
}

const pontosRisco: Record<ClassificacaoRisco, string> = {
  Vermelho: 'bg-red-500',
  Amarelo: 'bg-amber-400',
  Verde: 'bg-emerald-500',
}

const estilosStatus: Record<StatusFila, string> = {
  Aguardando: 'border-amber-200 bg-amber-50 text-amber-700',
  Chamado: 'border-blue-200 bg-blue-50 text-blue-700',
  'Em Atendimento':
    'border-violet-200 bg-violet-50 text-violet-700',
  Finalizado:
    'border-emerald-200 bg-emerald-50 text-emerald-700',
  Cancelado: 'border-slate-200 bg-slate-100 text-slate-600',
}

const pontosStatus: Record<StatusFila, string> = {
  Aguardando: 'bg-amber-400',
  Chamado: 'bg-blue-500',
  'Em Atendimento': 'bg-violet-500',
  Finalizado: 'bg-emerald-500',
  Cancelado: 'bg-slate-400',
}

const filtrosAtivos: FiltroStatus[] = [
  'Todos',
  'Aguardando',
  'Chamado',
  'Em Atendimento',
]

const filtrosHistorico: FiltroStatus[] = [
  'Todos',
  'Finalizado',
  'Cancelado',
]

function ehStatusAtivo(status: StatusFila): boolean {
  return (
    status === 'Aguardando' ||
    status === 'Chamado' ||
    status === 'Em Atendimento'
  )
}

function ehClassificacaoValida(
  valor: unknown
): valor is ClassificacaoRisco {
  return (
    valor === 'Vermelho' ||
    valor === 'Amarelo' ||
    valor === 'Verde'
  )
}

function ehPacienteTeste(nome: string): boolean {
  return nomesTeste.has(
    nome.trim().toLocaleLowerCase('pt-BR')
  )
}

function formatarData(data: string): string {
  if (!data) return '—'

  const valor = new Date(data)

  if (Number.isNaN(valor.getTime())) {
    return '—'
  }

  return valor.toLocaleString('pt-BR', {
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

function obterMensagemErro(erro: unknown): string {
  if (axios.isAxiosError(erro)) {
    const dados = erro.response?.data

    if (typeof dados?.error === 'string') {
      return dados.error
    }

    if (typeof dados?.message === 'string') {
      return dados.message
    }

    if (!erro.response) {
      return 'Não foi possível conectar à API.'
    }

    if (erro.response.status === 409) {
      return 'A operação não pôde ser realizada devido ao estado atual da fila.'
    }
  }

  return 'Ocorreu um erro inesperado. Tente novamente.'
}

function BadgeRisco({
  risco,
}: {
  risco: ClassificacaoRisco | null
}) {
  if (!risco) {
    return (
      <span className="inline-flex items-center rounded-full border border-slate-200 bg-slate-50 px-3 py-1.5 text-xs font-semibold text-slate-600">
        Não identificada
      </span>
    )
  }

  return (
    <span
      className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-bold ${estilosRisco[risco]}`}
    >
      <span
        className={`h-2 w-2 rounded-full ${pontosRisco[risco]}`}
      />
      {risco}
    </span>
  )
}

function BadgeStatus({ status }: { status: StatusFila }) {
  return (
    <span
      className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-semibold ${
        estilosStatus[status] ??
        'border-slate-200 bg-slate-50 text-slate-700'
      }`}
    >
      <span
        className={`h-2 w-2 rounded-full ${
          pontosStatus[status] ?? 'bg-slate-400'
        }`}
      />
      {status}
    </span>
  )
}

export default function Filas() {
  const [filas, setFilas] = useState<FilaDetalhada[]>([])
  const [aba, setAba] = useState<Aba>('ativa')
  const [pesquisa, setPesquisa] = useState('')
  const [filtroStatus, setFiltroStatus] =
    useState<FiltroStatus>('Todos')

  const [carregando, setCarregando] = useState(true)
  const [processando, setProcessando] = useState(false)

  const [erro, setErro] = useState<string | null>(null)
  const [sucesso, setSucesso] = useState<string | null>(null)

  const [filaParaCancelar, setFilaParaCancelar] =
    useState<FilaDetalhada | null>(null)

  // Impede respostas antigas de sobrescreverem consultas recentes.
  const ultimaConsulta = useRef(0)

  // Proteção imediata contra cliques duplicados.
  const operacaoEmAndamento = useRef(false)

  const carregarFilas = useCallback(async () => {
    const numeroConsulta = ++ultimaConsulta.current

    try {
      setCarregando(true)
      setErro(null)

      const [
        filasResponse,
        pacientesResponse,
        triagensResponse,
      ] = await Promise.all([
        filaService.listar(),
        pacienteService.listar(),
        triagemService.listar(),
      ])

      if (numeroConsulta !== ultimaConsulta.current) {
        return
      }

      const pacientesPorId = new Map(
        (pacientesResponse ?? []).map((paciente) => [
          paciente.id,
          paciente.nome,
        ])
      )

      const triagensPorId = new Map(
        (triagensResponse ?? []).map((triagem) => [
          triagem.id_triagem,
          triagem,
        ])
      )

      const filasDetalhadas: FilaDetalhada[] = (
        filasResponse ?? []
      )
        .filter((fila: Fila) => {
          const nome = pacientesPorId.get(fila.id_paciente)

          // Mantém a regra existente para pacientes de teste.
          return !nome || !ehPacienteTeste(nome)
        })
        .map((fila: Fila) => {
          const triagem = triagensPorId.get(fila.id_triagem)
          const risco = triagem?.classificacao_risco

          return {
            fila,
            nomePaciente:
              pacientesPorId.get(fila.id_paciente) ??
              'Paciente não identificado',
            classificacaoRisco:
              ehClassificacaoValida(risco) ? risco : null,
            dataTriagem: triagem?.data_triagem ?? '',
          }
        })

      filasDetalhadas.sort((a, b) => {
        const ativoA = ehStatusAtivo(a.fila.status)
        const ativoB = ehStatusAtivo(b.fila.status)

        if (ativoA !== ativoB) {
          return ativoA ? -1 : 1
        }

        if (
          a.fila.status === 'Aguardando' &&
          b.fila.status === 'Aguardando'
        ) {
          return a.fila.posicao - b.fila.posicao
        }

        return (
          new Date(b.fila.data_entrada).getTime() -
          new Date(a.fila.data_entrada).getTime()
        )
      })

      setFilas(filasDetalhadas)
    } catch (erroCapturado) {
      if (numeroConsulta === ultimaConsulta.current) {
        setErro(obterMensagemErro(erroCapturado))
      }
    } finally {
      if (numeroConsulta === ultimaConsulta.current) {
        setCarregando(false)
      }
    }
  }, [])

  useEffect(() => {
    void carregarFilas()

    return () => {
      ultimaConsulta.current += 1
    }
  }, [carregarFilas])

  const filasAtivas = useMemo(
    () =>
      filas.filter(({ fila }) => ehStatusAtivo(fila.status)),
    [filas]
  )

  const filasHistorico = useMemo(
    () =>
      filas.filter(({ fila }) => !ehStatusAtivo(fila.status)),
    [filas]
  )

  const filasDaAba =
    aba === 'ativa' ? filasAtivas : filasHistorico

  const filasFiltradas = useMemo(() => {
    const termo = pesquisa.trim().toLocaleLowerCase('pt-BR')

    return filasDaAba.filter(({ fila, nomePaciente }) => {
      const correspondeStatus =
        filtroStatus === 'Todos' ||
        fila.status === filtroStatus

      const correspondePesquisa =
        !termo ||
        nomePaciente.toLocaleLowerCase('pt-BR').includes(termo) ||
        fila.status.toLocaleLowerCase('pt-BR').includes(termo)

      return correspondeStatus && correspondePesquisa
    })
  }, [filasDaAba, pesquisa, filtroStatus])

  const totalAguardando = filasAtivas.filter(
    ({ fila }) => fila.status === 'Aguardando'
  ).length

  const totalChamados = filasAtivas.filter(
    ({ fila }) => fila.status === 'Chamado'
  ).length

  const totalEmAtendimento = filasAtivas.filter(
    ({ fila }) => fila.status === 'Em Atendimento'
  ).length

  const indicadores = [
    {
      titulo: 'Aguardando',
      valor: totalAguardando,
      descricao: 'Pacientes na fila de espera',
      icone: Clock3,
      corIcone: 'text-amber-600',
      corFundo: 'bg-amber-50',
      corBarra: 'bg-amber-400',
    },
    {
      titulo: 'Chamados',
      valor: totalChamados,
      descricao: 'Aguardando início do atendimento',
      icone: Megaphone,
      corIcone: 'text-blue-600',
      corFundo: 'bg-blue-50',
      corBarra: 'bg-blue-500',
    },
    {
      titulo: 'Em atendimento',
      valor: totalEmAtendimento,
      descricao: 'Atendimentos em andamento',
      icone: Stethoscope,
      corIcone: 'text-violet-600',
      corFundo: 'bg-violet-50',
      corBarra: 'bg-violet-500',
    },
  ]

  function mudarAba(novaAba: Aba) {
    setAba(novaAba)
    setFiltroStatus('Todos')
    setPesquisa('')
  }

  async function handleChamarProximo() {
    if (
      operacaoEmAndamento.current ||
      carregando ||
      totalAguardando === 0
    ) {
      return
    }

    operacaoEmAndamento.current = true

    try {
      setProcessando(true)
      setErro(null)
      setSucesso(null)

      // O backend decide qual paciente deve ser chamado.
      const filaChamada = await filaService.chamarProximo()

      const paciente = filas.find(
        ({ fila }) => fila.id_fila === filaChamada.id_fila
      )

      setSucesso(
        paciente
          ? `Paciente ${paciente.nomePaciente} chamado com sucesso!`
          : 'Próximo paciente chamado com sucesso!'
      )

      await carregarFilas()
    } catch (erroCapturado) {
      setErro(obterMensagemErro(erroCapturado))
    } finally {
      operacaoEmAndamento.current = false
      setProcessando(false)
    }
  }

  async function handleCancelar() {
    if (!filaParaCancelar || operacaoEmAndamento.current) {
      return
    }

    const filaSelecionada = filaParaCancelar

    operacaoEmAndamento.current = true

    try {
      setProcessando(true)
      setErro(null)
      setSucesso(null)

      await filaService.cancelar(filaSelecionada.fila.id_fila)

      setFilaParaCancelar(null)

      setSucesso(
        `Entrada de ${filaSelecionada.nomePaciente} cancelada com sucesso.`
      )

      await carregarFilas()
    } catch (erroCapturado) {
      setErro(obterMensagemErro(erroCapturado))
    } finally {
      operacaoEmAndamento.current = false
      setProcessando(false)
    }
  }

  const filtrosDisponiveis =
    aba === 'ativa' ? filtrosAtivos : filtrosHistorico

  return (
    <div className="space-y-7 pb-4">
      {/* Cabeçalho Medical Premium */}
      <section className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-6 shadow-sm sm:p-7">
        <div className="pointer-events-none absolute -right-16 -top-20 h-56 w-56 rounded-full bg-teal-50 blur-3xl" />

        <div className="relative flex flex-col justify-between gap-5 lg:flex-row lg:items-center">
          <div className="flex items-start gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-teal-600 text-white shadow-lg shadow-teal-100">
              <ListOrdered size={24} />
            </div>

            <div>
              <div className="mb-2 flex flex-wrap items-center gap-2">
                <span className="rounded-full border border-teal-100 bg-teal-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-teal-700">
                  Gestão assistencial
                </span>

                <span className="inline-flex items-center gap-1.5 text-xs font-medium text-slate-500">
                  <span className="h-2 w-2 rounded-full bg-teal-500" />
                  Fluxo de pacientes
                </span>
              </div>

              <h2 className="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
                Fila de atendimento
              </h2>

              <p className="mt-2 max-w-xl text-sm leading-relaxed text-slate-500">
                Acompanhe as prioridades, monitore os status e
                gerencie o fluxo de pacientes da unidade.
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={() => void handleChamarProximo()}
            disabled={
              processando ||
              carregando ||
              totalAguardando === 0
            }
            className="inline-flex shrink-0 items-center justify-center gap-2 self-start rounded-xl bg-teal-600 px-5 py-3 text-sm font-semibold text-white shadow-sm shadow-teal-100 transition hover:bg-teal-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {processando ? (
              <Loader2
                size={18}
                className="animate-spin"
              />
            ) : (
              <Megaphone size={18} />
            )}

            {processando
              ? 'Processando...'
              : 'Chamar próximo'}
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
            onClick={() => setSucesso(null)}
            aria-label="Fechar mensagem"
            className="rounded-lg p-1 transition hover:bg-emerald-100"
          >
            <X size={17} />
          </button>
        </div>
      )}

      {erro && (
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

      {/* Indicadores operacionais */}
      <section>
        <div className="mb-4">
          <h3 className="text-base font-bold text-slate-900">
            Visão operacional
          </h3>

          <p className="mt-1 text-xs text-slate-500">
            Situação atual do fluxo de atendimento
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-3">
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
                    Fila
                  </span>
                </div>

                <p className="mt-5 text-sm font-medium text-slate-500">
                  {indicador.titulo}
                </p>

                <p className="mt-2 text-4xl font-bold tracking-tight text-slate-900">
                  {carregando && filas.length === 0
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

      {/* Painel da fila */}
      <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        {/* Abas */}
        <div className="flex gap-6 overflow-x-auto border-b border-slate-100 px-5 sm:px-6">
          <button
            type="button"
            onClick={() => mudarAba('ativa')}
            aria-pressed={aba === 'ativa'}
            className={`inline-flex shrink-0 items-center gap-2 border-b-2 py-5 text-sm font-semibold transition ${
              aba === 'ativa'
                ? 'border-teal-600 text-teal-700'
                : 'border-transparent text-slate-500 hover:text-slate-800'
            }`}
          >
            <Users size={18} />
            Fila ativa

            <span
              className={`rounded-full px-2.5 py-0.5 text-xs ${
                aba === 'ativa'
                  ? 'bg-teal-50 text-teal-700'
                  : 'bg-slate-100 text-slate-600'
              }`}
            >
              {filasAtivas.length}
            </span>
          </button>

          <button
            type="button"
            onClick={() => mudarAba('historico')}
            aria-pressed={aba === 'historico'}
            className={`inline-flex shrink-0 items-center gap-2 border-b-2 py-5 text-sm font-semibold transition ${
              aba === 'historico'
                ? 'border-teal-600 text-teal-700'
                : 'border-transparent text-slate-500 hover:text-slate-800'
            }`}
          >
            <History size={18} />
            Histórico

            <span
              className={`rounded-full px-2.5 py-0.5 text-xs ${
                aba === 'historico'
                  ? 'bg-teal-50 text-teal-700'
                  : 'bg-slate-100 text-slate-600'
              }`}
            >
              {filasHistorico.length}
            </span>
          </button>
        </div>

        {/* Filtros */}
        <div className="flex flex-col justify-between gap-4 border-b border-slate-100 p-5 lg:flex-row lg:items-center sm:p-6">
          <div>
            <h3 className="flex items-center gap-2 text-base font-bold text-slate-900">
              <ListOrdered
                size={19}
                className="text-teal-600"
              />

              {aba === 'ativa'
                ? 'Pacientes em fluxo ativo'
                : 'Histórico da fila'}
            </h3>

            <p className="mt-1 text-xs text-slate-500">
              {filasFiltradas.length} registro(s) encontrado(s)
            </p>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative w-full sm:w-64">
              <Search
                size={17}
                className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400"
              />

              <input
                type="search"
                value={pesquisa}
                onChange={(event) =>
                  setPesquisa(event.target.value)
                }
                placeholder="Buscar paciente..."
                aria-label="Buscar paciente"
                className="w-full rounded-xl border border-slate-200 bg-slate-50 py-2.5 pl-10 pr-4 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-teal-500 focus:bg-white focus:ring-2 focus:ring-teal-500/10"
              />
            </div>

            <select
              value={filtroStatus}
              onChange={(event) =>
                setFiltroStatus(
                  event.target.value as FiltroStatus
                )
              }
              aria-label="Filtrar por status"
              className="rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-sm text-slate-700 outline-none transition focus:border-teal-500 focus:ring-2 focus:ring-teal-500/10"
            >
              {filtrosDisponiveis.map((filtro) => (
                <option key={filtro} value={filtro}>
                  {filtro === 'Todos'
                    ? 'Todos os status'
                    : filtro}
                </option>
              ))}
            </select>

            <button
              type="button"
              onClick={() => void carregarFilas()}
              disabled={carregando || processando}
              title="Atualizar fila"
              aria-label="Atualizar fila"
              className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold text-slate-700 transition hover:border-teal-200 hover:bg-teal-50 hover:text-teal-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <RefreshCw
                size={17}
                className={carregando ? 'animate-spin' : ''}
              />
              Atualizar
            </button>
          </div>
        </div>

        {/* Conteúdo */}
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
              Carregando fila...
            </p>
          </div>
        ) : filasFiltradas.length === 0 ? (
          <div className="flex flex-col items-center gap-3 px-5 py-16 text-center">
            <div className="rounded-2xl bg-slate-50 p-4 text-slate-400">
              <ListOrdered size={29} />
            </div>

            <p className="text-sm font-semibold text-slate-700">
              Nenhum registro encontrado
            </p>

            <p className="text-xs text-slate-500">
              Não existem registros para os filtros selecionados.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[850px] text-left">
              <thead className="bg-slate-50/80">
                <tr className="text-[11px] font-bold uppercase tracking-wider text-slate-500">
                  <th scope="col" className="px-5 py-4">
                    Posição
                  </th>
                  <th scope="col" className="px-5 py-4">
                    Paciente
                  </th>
                  <th scope="col" className="px-5 py-4">
                    Risco
                  </th>
                  <th scope="col" className="px-5 py-4">
                    Status
                  </th>
                  <th scope="col" className="px-5 py-4">
                    Entrada
                  </th>

                  {aba === 'ativa' && (
                    <th
                      scope="col"
                      className="px-5 py-4 text-right"
                    >
                      Ações
                    </th>
                  )}
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {filasFiltradas.map((item) => {
                  const {
                    fila,
                    nomePaciente,
                    classificacaoRisco,
                  } = item

                  const podeCancelar =
                    fila.status === 'Aguardando' ||
                    fila.status === 'Chamado'

                  return (
                    <tr
                      key={fila.id_fila}
                      className="transition-colors hover:bg-teal-50/30"
                    >
                      <td className="px-5 py-4">
                        {fila.status === 'Aguardando' &&
                        fila.posicao > 0 ? (
                          <span className="inline-flex h-10 min-w-10 items-center justify-center rounded-xl border border-slate-200 bg-slate-50 px-2 text-sm font-bold text-slate-800">
                            {fila.posicao}º
                          </span>
                        ) : (
                          <span className="text-sm text-slate-400">
                            —
                          </span>
                        )}
                      </td>

                      <td className="px-5 py-4">
                        <div className="flex items-center gap-3">
                          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-teal-50 text-xs font-bold text-teal-700">
                            {obterIniciais(nomePaciente)}
                          </div>

                          <div>
                            <p className="font-semibold text-slate-800">
                              {nomePaciente}
                            </p>

                            <p className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
                              <UserRound size={12} />
                              Paciente
                            </p>
                          </div>
                        </div>
                      </td>

                      <td className="px-5 py-4">
                        <BadgeRisco
                          risco={classificacaoRisco}
                        />
                      </td>

                      <td className="px-5 py-4">
                        <BadgeStatus status={fila.status} />
                      </td>

                      <td className="px-5 py-4">
                        <div className="flex items-center gap-1.5 text-sm text-slate-600">
                          <Clock3
                            size={15}
                            className="shrink-0 text-slate-400"
                          />
                          {formatarData(fila.data_entrada)}
                        </div>
                      </td>

                      {aba === 'ativa' && (
                        <td className="px-5 py-4 text-right">
                          {podeCancelar && (
                            <button
                              type="button"
                              onClick={() =>
                                setFilaParaCancelar(item)
                              }
                              disabled={processando}
                              className="inline-flex items-center gap-1.5 rounded-xl border border-transparent px-3 py-2 text-xs font-semibold text-red-600 transition hover:border-red-100 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-50"
                            >
                              <XCircle size={16} />
                              Cancelar
                            </button>
                          )}
                        </td>
                      )}
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        )}

        {!carregando && (
          <div className="flex flex-wrap items-center justify-between gap-2 border-t border-slate-100 bg-slate-50/50 px-5 py-4 sm:px-6">
            <p className="text-xs text-slate-500">
              Exibindo{' '}
              <strong className="text-slate-700">
                {filasFiltradas.length}
              </strong>{' '}
              de{' '}
              <strong className="text-slate-700">
                {filasDaAba.length}
              </strong>{' '}
              registros
            </p>

            <span className="inline-flex items-center gap-1.5 text-xs font-medium text-teal-700">
              <ShieldCheck size={14} />
              Gestão de filas
            </span>
          </div>
        )}
      </section>

      {/* Modal de confirmação */}
      {filaParaCancelar && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 p-4 backdrop-blur-sm"
          role="presentation"
        >
          <div
            role="dialog"
            aria-modal="true"
            aria-labelledby="titulo-cancelamento"
            className="w-full max-w-md overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-2xl"
          >
            <div className="border-b border-slate-100 p-6">
              <div className="flex items-start gap-4">
                <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-red-50 text-red-600">
                  <AlertCircle size={24} />
                </div>

                <div>
                  <span className="text-[10px] font-bold uppercase tracking-wider text-red-600">
                    Confirmação necessária
                  </span>

                  <h2
                    id="titulo-cancelamento"
                    className="mt-1 text-lg font-bold text-slate-900"
                  >
                    Cancelar entrada na fila?
                  </h2>
                </div>
              </div>
            </div>

            <div className="space-y-4 p-6">
              <p className="text-sm leading-relaxed text-slate-600">
                Deseja cancelar a entrada de{' '}
                <strong className="text-slate-900">
                  {filaParaCancelar.nomePaciente}
                </strong>
                ? Esta operação retirará o paciente do fluxo
                ativo desta entrada na fila.
              </p>

              <div className="flex flex-wrap items-center gap-2 rounded-xl border border-slate-100 bg-slate-50 p-3">
                <BadgeRisco
                  risco={filaParaCancelar.classificacaoRisco}
                />

                <BadgeStatus
                  status={filaParaCancelar.fila.status}
                />
              </div>

              <p className="text-xs leading-relaxed text-slate-500">
                Confirme somente após verificar que esta é a
                entrada correta.
              </p>
            </div>

            <div className="flex flex-col-reverse gap-3 border-t border-slate-100 bg-slate-50/70 p-5 sm:flex-row sm:justify-end">
              <button
                type="button"
                onClick={() => setFilaParaCancelar(null)}
                disabled={processando}
                className="inline-flex items-center justify-center rounded-xl border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 disabled:opacity-50"
              >
                Voltar
              </button>

              <button
                type="button"
                onClick={() => void handleCancelar()}
                disabled={processando}
                className="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-5 py-3 text-sm font-semibold text-white transition hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {processando ? (
                  <Loader2
                    size={17}
                    className="animate-spin"
                  />
                ) : (
                  <XCircle size={17} />
                )}

                {processando
                  ? 'Cancelando...'
                  : 'Confirmar cancelamento'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
