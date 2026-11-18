
import { useCallback, useEffect, useMemo, useState } from 'react'
import axios from 'axios'

import {
  Activity,
  AlertCircle,
  CheckCircle2,
  ClipboardList,
  Clock3,
  HeartPulse,
  Loader2,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Thermometer,
  UserRound,
  Users,
  X,
  Ban,
} from 'lucide-react'

import TriagemForm from '../components/triagens/TriagemForm'

import { triagemService } from '../services/triagemService'
import { pacienteService } from '../services/pacienteService'
import { profissionalService } from '../services/profissionalService'
import { usuarioService } from '../services/usuarioService'

import { filtrarPacientesVisiveis } from '../utils/pacienteUtils'

import type {
  Triagem,
  ClassificacaoRisco,
} from '../types/triagem'

interface TriagemExibicao {
  triagem: Triagem
  pacienteNome: string
  profissionalNome: string
}

const estilosRisco: Record<ClassificacaoRisco, string> = {
  Vermelho: 'border-red-200 bg-red-50 text-red-700',
  Amarelo: 'border-amber-200 bg-amber-50 text-amber-700',
  Verde: 'border-emerald-200 bg-emerald-50 text-emerald-700',
}

const coresRisco: Record<ClassificacaoRisco, string> = {
  Vermelho: 'bg-red-500',
  Amarelo: 'bg-amber-400',
  Verde: 'bg-emerald-500',
}

function formatarData(data: string): string {
  if (!data) return 'Não informada'

  const valor = new Date(data)

  if (Number.isNaN(valor.getTime())) {
    return 'Data inválida'
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
    const mensagem = erro.response?.data?.error

    if (typeof mensagem === 'string') {
      return mensagem
    }

    const mensagemAlternativa = erro.response?.data?.message

    if (typeof mensagemAlternativa === 'string') {
      return mensagemAlternativa
    }

    if (!erro.response) {
      return 'Não foi possível conectar à API.'
    }
  }

  return 'Não foi possível concluir a operação.'
}

export default function Triagens() {
  const [triagens, setTriagens] = useState<TriagemExibicao[]>([])
  const [pesquisa, setPesquisa] = useState('')
  const [mostrarFormulario, setMostrarFormulario] = useState(false)

  const [carregando, setCarregando] = useState(true)
  const [erro, setErro] = useState<string | null>(null)
  const [sucesso, setSucesso] = useState<string | null>(null)

  // Estados do cancelamento.
  const [triagemParaCancelar, setTriagemParaCancelar] =
    useState<TriagemExibicao | null>(null)

  const [motivoCancelamento, setMotivoCancelamento] =
    useState('')

  const [cancelando, setCancelando] = useState(false)
  const [erroCancelamento, setErroCancelamento] =
    useState<string | null>(null)

  const carregarTriagens = useCallback(async () => {
    try {
      setCarregando(true)
      setErro(null)

      const [
        triagensResponse,
        pacientesResponse,
        profissionaisResponse,
        usuariosResponse,
      ] = await Promise.all([
        triagemService.listar(),
        pacienteService.listar(),
        profissionalService.listar(),
        usuarioService.listar(),
      ])

      const pacientesVisiveis = filtrarPacientesVisiveis(
        pacientesResponse ?? []
      )

      const pacientesPorId = new Map(
        pacientesVisiveis.map((paciente) => [
          paciente.id,
          paciente.nome,
        ])
      )

      const usuariosPorId = new Map(
        (usuariosResponse ?? []).map((usuario) => [
          usuario.id_usuario,
          usuario.nome,
        ])
      )

      const profissionaisPorId = new Map(
        (profissionaisResponse ?? []).map((profissional) => [
          profissional.id_profissional,
          usuariosPorId.get(profissional.id_usuario) ??
            'Profissional não identificado',
        ])
      )

      const lista: TriagemExibicao[] = (triagensResponse ?? [])
        .filter((triagem) =>
          pacientesPorId.has(triagem.id_paciente)
        )
        .map((triagem) => ({
          triagem,
          pacienteNome:
            pacientesPorId.get(triagem.id_paciente) ??
            'Paciente não identificado',
          profissionalNome:
            profissionaisPorId.get(triagem.id_profissional) ??
            'Profissional não identificado',
        }))
        .sort((a, b) => {
          const dataA = new Date(
            a.triagem.data_triagem
          ).getTime()

          const dataB = new Date(
            b.triagem.data_triagem
          ).getTime()

          return (
            (Number.isFinite(dataB) ? dataB : 0) -
            (Number.isFinite(dataA) ? dataA : 0)
          )
        })

      setTriagens(lista)
    } catch (erroCapturado) {
      setErro(obterMensagemErro(erroCapturado))
    } finally {
      setCarregando(false)
    }
  }, [])

  useEffect(() => {
    void carregarTriagens()
  }, [carregarTriagens])

  const triagensFiltradas = useMemo(() => {
    const termo = pesquisa.trim().toLocaleLowerCase('pt-BR')

    if (!termo) return triagens

    return triagens.filter(
      ({ pacienteNome, profissionalNome, triagem }) =>
        pacienteNome.toLocaleLowerCase('pt-BR').includes(termo) ||
        profissionalNome.toLocaleLowerCase('pt-BR').includes(termo) ||
        triagem.classificacao_risco
          .toLocaleLowerCase('pt-BR')
          .includes(termo) ||
        (triagem.cancelada_em ? 'cancelada' : 'ativa')
          .includes(termo) ||
        (triagem.motivo_cancelamento ?? '')
          .toLocaleLowerCase('pt-BR')
          .includes(termo)
    )
  }, [triagens, pesquisa])

  const indicadores = useMemo(
    () => [
      {
        titulo: 'Triagens realizadas',
        valor: triagens.length,
        descricao: 'Avaliações registradas',
        icone: ClipboardList,
        corIcone: 'text-blue-600',
        corFundo: 'bg-blue-50',
        corBarra: 'bg-blue-500',
      },
      {
        titulo: 'Risco vermelho',
        valor: triagens.filter(
          ({ triagem }) =>
            triagem.classificacao_risco === 'Vermelho'
        ).length,
        descricao: 'Classificação vermelha',
        icone: HeartPulse,
        corIcone: 'text-red-600',
        corFundo: 'bg-red-50',
        corBarra: 'bg-red-500',
      },
      {
        titulo: 'Risco amarelo',
        valor: triagens.filter(
          ({ triagem }) =>
            triagem.classificacao_risco === 'Amarelo'
        ).length,
        descricao: 'Classificação amarela',
        icone: Activity,
        corIcone: 'text-amber-600',
        corFundo: 'bg-amber-50',
        corBarra: 'bg-amber-400',
      },
      {
        titulo: 'Risco verde',
        valor: triagens.filter(
          ({ triagem }) =>
            triagem.classificacao_risco === 'Verde'
        ).length,
        descricao: 'Classificação verde',
        icone: ShieldCheck,
        corIcone: 'text-emerald-600',
        corFundo: 'bg-emerald-50',
        corBarra: 'bg-emerald-500',
      },
    ],
    [triagens]
  )

  async function handleSucesso() {
    setMostrarFormulario(false)
    setSucesso('Triagem registrada com sucesso!')
    setPesquisa('')
    await carregarTriagens()
  }

  function abrirCancelamento(item: TriagemExibicao) {
    if (item.triagem.cancelada_em) return

    setTriagemParaCancelar(item)
    setMotivoCancelamento('')
    setErroCancelamento(null)
    setErro(null)
    setSucesso(null)
  }

  function fecharCancelamento() {
    if (cancelando) return

    setTriagemParaCancelar(null)
    setMotivoCancelamento('')
    setErroCancelamento(null)
  }

  async function confirmarCancelamento() {
    if (!triagemParaCancelar || cancelando) return

    const motivo = motivoCancelamento.trim()

    if (!motivo) {
      setErroCancelamento(
        'Informe o motivo do cancelamento.'
      )
      return
    }

    if ([...motivo].length > 1000) {
      setErroCancelamento(
        'O motivo deve ter no máximo 1000 caracteres.'
      )
      return
    }

    try {
      setCancelando(true)
      setErroCancelamento(null)

      const triagemAtualizada =
        await triagemService.cancelar(
          triagemParaCancelar.triagem.id_triagem,
          { motivo }
        )

      // Atualiza imediatamente a triagem cancelada na tabela.
      setTriagens((lista) =>
        lista.map((item) =>
          item.triagem.id_triagem ===
          triagemAtualizada.id_triagem
            ? {
                ...item,
                triagem: triagemAtualizada,
              }
            : item
        )
      )

      setTriagemParaCancelar(null)
      setMotivoCancelamento('')

      setSucesso(
        'Triagem cancelada com sucesso. O histórico clínico foi preservado e a fila reorganizada.'
      )

      // Sincroniza a listagem com o backend.
      await carregarTriagens()
    } catch (erroCapturado) {
      setErroCancelamento(
        obterMensagemErro(erroCapturado)
      )
    } finally {
      setCancelando(false)
    }
  }

  return (
    <div className="space-y-7 pb-4">
      {/* Cabeçalho */}
      <section className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-6 shadow-sm sm:p-7">
        <div className="pointer-events-none absolute -right-16 -top-20 h-56 w-56 rounded-full bg-teal-50 blur-3xl" />

        <div className="relative flex flex-col justify-between gap-5 md:flex-row md:items-center">
          <div className="flex items-start gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-teal-600 text-white shadow-lg shadow-teal-100">
              <Activity size={25} />
            </div>

            <div>
              <div className="mb-2 flex flex-wrap items-center gap-2">
                <span className="rounded-full border border-teal-100 bg-teal-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-teal-700">
                  Avaliação clínica
                </span>

                <span className="inline-flex items-center gap-1.5 text-xs font-medium text-slate-500">
                  <span className="h-2 w-2 rounded-full bg-teal-500" />
                  Gestão de triagens
                </span>
              </div>

              <h2 className="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
                Triagens
              </h2>

              <p className="mt-2 max-w-xl text-sm leading-relaxed text-slate-500">
                Acompanhe as avaliações iniciais e as classificações
                de risco dos pacientes da unidade.
              </p>
            </div>
          </div>

          {!mostrarFormulario && (
            <button
              type="button"
              onClick={() => {
                setSucesso(null)
                setMostrarFormulario(true)
              }}
              className="inline-flex shrink-0 items-center justify-center gap-2 self-start rounded-xl bg-teal-600 px-5 py-3 text-sm font-semibold text-white shadow-sm shadow-teal-100 transition hover:bg-teal-700"
            >
              <Plus size={18} />
              Nova triagem
            </button>
          )}
        </div>
      </section>

      {/* Mensagem de sucesso */}
      {sucesso && (
        <div
          role="status"
          className="flex items-center justify-between gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm font-medium text-emerald-700"
        >
          <div className="flex items-center gap-2">
            <CheckCircle2 size={19} className="shrink-0" />
            <span>{sucesso}</span>
          </div>

          <button
            type="button"
            aria-label="Fechar mensagem"
            onClick={() => setSucesso(null)}
            className="rounded-lg p-1 transition hover:bg-emerald-100"
          >
            <X size={17} />
          </button>
        </div>
      )}

      {/* Formulário existente */}
      {mostrarFormulario && (
        <TriagemForm
          onSuccess={() => {
            void handleSucesso()
          }}
          onCancel={() => setMostrarFormulario(false)}
        />
      )}

      {/* Indicadores */}
      <section>
        <div className="mb-4">
          <h3 className="text-base font-bold text-slate-900">
            Indicadores de triagem
          </h3>

          <p className="mt-1 text-xs text-slate-500">
            Distribuição das classificações registradas
          </p>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
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
                    Triagem
                  </span>
                </div>

                <p className="mt-5 text-sm font-medium text-slate-500">
                  {indicador.titulo}
                </p>

                <p className="mt-2 text-4xl font-bold tracking-tight text-slate-900">
                  {carregando && triagens.length === 0
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

      {/* Listagem */}
      <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="flex flex-col justify-between gap-4 border-b border-slate-100 p-5 sm:p-6 lg:flex-row lg:items-center">
          <div>
            <div className="flex items-center gap-2">
              <ClipboardList
                size={19}
                className="text-teal-600"
              />

              <h3 className="text-base font-bold text-slate-900">
                Triagens registradas
              </h3>
            </div>

            <p className="mt-1 text-xs text-slate-500">
              {triagensFiltradas.length} registro(s) encontrado(s)
            </p>
          </div>

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative w-full sm:w-80">
              <Search
                size={18}
                className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400"
              />

              <input
                type="search"
                value={pesquisa}
                onChange={(event) =>
                  setPesquisa(event.target.value)
                }
                placeholder="Paciente, profissional ou risco..."
                aria-label="Buscar triagens"
                className="w-full rounded-xl border border-slate-200 bg-slate-50 py-2.5 pl-10 pr-4 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-teal-500 focus:bg-white focus:ring-2 focus:ring-teal-500/10"
              />
            </div>

            <button
              type="button"
              onClick={() => {
                void carregarTriagens()
              }}
              disabled={carregando}
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

        {/* Erro */}
        {erro && (
          <div
            role="alert"
            className="m-5 flex items-start gap-3 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
          >
            <AlertCircle
              size={19}
              className="shrink-0"
            />
            <span>{erro}</span>
          </div>
        )}

        {/* Carregamento */}
        {carregando ? (
          <div
            role="status"
            className="flex flex-col items-center justify-center gap-3 py-16 text-slate-500"
          >
            <Loader2
              size={29}
              className="animate-spin text-teal-600"
            />

            <p className="text-sm">
              Carregando triagens...
            </p>
          </div>
        ) : erro ? (
          <div className="flex flex-col items-center gap-3 px-5 py-14 text-center">
            <div className="rounded-2xl bg-red-50 p-4 text-red-500">
              <AlertCircle size={28} />
            </div>

            <p className="text-sm font-medium text-slate-600">
              Não foi possível carregar a lista.
            </p>
          </div>
        ) : triagensFiltradas.length === 0 ? (
          <div className="flex flex-col items-center gap-3 px-5 py-16 text-center">
            <div className="rounded-2xl bg-slate-50 p-4 text-slate-400">
              <ClipboardList size={29} />
            </div>

            <p className="text-sm font-semibold text-slate-700">
              Nenhuma triagem encontrada
            </p>

            <p className="text-xs text-slate-400">
              {pesquisa
                ? 'Tente pesquisar por outro termo.'
                : 'As triagens registradas aparecerão aqui.'}
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[1100px] text-left">
              <thead className="bg-slate-50/80">
                <tr className="text-[11px] font-bold uppercase tracking-wider text-slate-500">
                  <th scope="col" className="px-5 py-4">
                    Paciente
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Profissional
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Risco
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Pressão
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Temperatura
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Data da triagem
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Situação
                  </th>

                  <th scope="col" className="px-5 py-4">
                    Ações
                  </th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {triagensFiltradas.map(
                  ({
                    triagem,
                    pacienteNome,
                    profissionalNome,
                  }) => {
                    const cancelada = Boolean(
                      triagem.cancelada_em
                    )

                    return (
                      <tr
                        key={triagem.id_triagem}
                        className={`transition-colors ${
                          cancelada
                            ? 'bg-slate-50/70'
                            : 'hover:bg-teal-50/30'
                        }`}
                      >
                        <td className="px-5 py-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-teal-50 text-xs font-bold text-teal-700">
                              {obterIniciais(pacienteNome)}
                            </div>

                            <div className="min-w-0">
                              <p className="font-semibold text-slate-800">
                                {pacienteNome}
                              </p>

                              <p className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
                                <UserRound size={12} />
                                Paciente
                              </p>
                            </div>
                          </div>
                        </td>

                        <td className="px-5 py-4">
                          <div className="flex items-center gap-2 text-sm text-slate-600">
                            <Users
                              size={15}
                              className="shrink-0 text-slate-400"
                            />
                            <span>{profissionalNome}</span>
                          </div>
                        </td>

                        <td className="px-5 py-4">
                          <span
                            className={`inline-flex items-center gap-2 rounded-full border px-3 py-1.5 text-xs font-bold ${
                              estilosRisco[
                                triagem.classificacao_risco
                              ] ??
                              'border-slate-200 bg-slate-50 text-slate-700'
                            }`}
                          >
                            <span
                              className={`h-2 w-2 rounded-full ${
                                coresRisco[
                                  triagem.classificacao_risco
                                ] ?? 'bg-slate-400'
                              }`}
                            />

                            {triagem.classificacao_risco}
                          </span>
                        </td>

                        <td className="px-5 py-4 text-sm text-slate-600">
                          <div className="flex items-center gap-1.5">
                            <Activity
                              size={15}
                              className="text-slate-400"
                            />
                            {triagem.pressao_arterial || '—'}
                          </div>
                        </td>

                        <td className="px-5 py-4 text-sm text-slate-600">
                          <div className="flex items-center gap-1.5">
                            <Thermometer
                              size={15}
                              className="text-slate-400"
                            />
                            {triagem.temperatura != null
                              ? `${triagem.temperatura} °C`
                              : '—'}
                          </div>
                        </td>

                        <td className="px-5 py-4">
                          <div className="flex items-center gap-1.5 text-sm text-slate-600">
                            <Clock3
                              size={15}
                              className="shrink-0 text-slate-400"
                            />
                            {formatarData(triagem.data_triagem)}
                          </div>
                        </td>

                        <td className="px-5 py-4">
                          {cancelada ? (
                            <div className="space-y-1.5">
                              <span className="inline-flex items-center gap-1.5 rounded-full border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-bold text-red-700">
                                <Ban size={13} />
                                Cancelada
                              </span>

                              <p className="text-xs text-slate-500">
                                {formatarData(
                                  triagem.cancelada_em!
                                )}
                              </p>

                              {triagem.motivo_cancelamento && (
                                <p
                                  className="max-w-[220px] break-words text-xs text-slate-500"
                                  title={
                                    triagem.motivo_cancelamento
                                  }
                                >
                                  Motivo:{' '}
                                  {triagem.motivo_cancelamento}
                                </p>
                              )}
                            </div>
                          ) : (
                            <span className="inline-flex items-center gap-1.5 rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1.5 text-xs font-semibold text-emerald-700">
                              <CheckCircle2 size={13} />
                              Não cancelada
                            </span>
                          )}
                        </td>

                        <td className="px-5 py-4">
                          {cancelada ? (
                            <span className="text-xs text-slate-400">
                              Sem ações
                            </span>
                          ) : (
                            <button
                              type="button"
                              onClick={() =>
                                abrirCancelamento({
                                  triagem,
                                  pacienteNome,
                                  profissionalNome,
                                })
                              }
                              className="inline-flex items-center gap-1.5 rounded-lg border border-red-200 bg-white px-3 py-2 text-xs font-semibold text-red-700 transition hover:bg-red-50"
                              title="O cancelamento só é permitido enquanto o paciente aguarda na fila."
                            >
                              <Ban size={14} />
                              Cancelar
                            </button>
                          )}
                        </td>
                      </tr>
                    )
                  }
                )}
              </tbody>
            </table>
          </div>
        )}

        {!carregando && !erro && (
          <div className="flex flex-wrap items-center justify-between gap-2 border-t border-slate-100 bg-slate-50/50 px-5 py-4 sm:px-6">
            <p className="text-xs text-slate-500">
              Exibindo{' '}
              <strong className="text-slate-700">
                {triagensFiltradas.length}
              </strong>{' '}
              de{' '}
              <strong className="text-slate-700">
                {triagens.length}
              </strong>{' '}
              triagens
            </p>

            <span className="inline-flex items-center gap-1.5 text-xs font-medium text-teal-700">
              <CheckCircle2 size={14} />
              Registros carregados
            </span>
          </div>
        )}
      </section>

      {/* Modal de cancelamento */}
      {triagemParaCancelar && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 p-4"
          role="presentation"
          onMouseDown={(event) => {
            if (event.target === event.currentTarget) {
              fecharCancelamento()
            }
          }}
        >
          <div
            role="dialog"
            aria-modal="true"
            aria-labelledby="titulo-cancelamento"
            className="w-full max-w-lg rounded-2xl border border-slate-200 bg-white p-6 shadow-2xl"
          >
            <div className="flex items-start justify-between gap-4">
              <div className="flex items-center gap-3">
                <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-red-50 text-red-600">
                  <Ban size={24} />
                </div>

                <div>
                  <h3
                    id="titulo-cancelamento"
                    className="text-lg font-bold text-slate-900"
                  >
                    Cancelar triagem
                  </h3>

                  <p className="mt-1 text-sm text-slate-500">
                    Confirme o cancelamento do registro.
                  </p>
                </div>
              </div>

              <button
                type="button"
                onClick={fecharCancelamento}
                disabled={cancelando}
                aria-label="Fechar"
                className="rounded-lg p-2 text-slate-400 transition hover:bg-slate-100 disabled:opacity-50"
              >
                <X size={19} />
              </button>
            </div>

            <div className="mt-5 rounded-xl border border-slate-200 bg-slate-50 p-4">
              <p className="text-xs font-medium text-slate-500">
                Paciente
              </p>

              <p className="mt-1 font-bold text-slate-900">
                {triagemParaCancelar.pacienteNome}
              </p>

              <p className="mt-2 text-xs text-slate-500">
                Classificação:{' '}
                <strong>
                  {
                    triagemParaCancelar.triagem
                      .classificacao_risco
                  }
                </strong>
              </p>
            </div>

            <div className="mt-4 rounded-xl border border-amber-200 bg-amber-50 p-4 text-sm text-amber-800">
              <div className="flex items-start gap-2">
                <AlertCircle
                  size={18}
                  className="mt-0.5 shrink-0"
                />

                <p>
                  O cancelamento removerá o paciente da fila
                  de espera, mas preservará o registro clínico
                  para consulta e auditoria. A operação só é
                  permitida enquanto o paciente estiver
                  aguardando.
                </p>
              </div>
            </div>

            <form
              onSubmit={(event) => {
                event.preventDefault()
                void confirmarCancelamento()
              }}
            >
              <label
                htmlFor="motivo-cancelamento"
                className="mt-5 block text-sm font-semibold text-slate-800"
              >
                Motivo do cancelamento *
              </label>

              <textarea
                id="motivo-cancelamento"
                value={motivoCancelamento}
                onChange={(event) => {
                  setMotivoCancelamento(event.target.value)
                  setErroCancelamento(null)
                }}
                required
                maxLength={1000}
                rows={4}
                disabled={cancelando}
                placeholder="Descreva por que esta triagem deve ser cancelada..."
                className="mt-2 w-full resize-none rounded-xl border border-slate-200 bg-white p-3 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-teal-500 focus:ring-2 focus:ring-teal-500/10 disabled:bg-slate-50"
              />

              <p className="mt-1 text-right text-xs text-slate-400">
                {[...motivoCancelamento].length}/1000 caracteres
              </p>

              {erroCancelamento && (
                <div
                  role="alert"
                  className="mt-4 flex items-start gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700"
                >
                  <AlertCircle
                    size={17}
                    className="shrink-0"
                  />
                  <span>{erroCancelamento}</span>
                </div>
              )}

              <div className="mt-6 flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
                <button
                  type="button"
                  onClick={fecharCancelamento}
                  disabled={cancelando}
                  className="rounded-xl border border-slate-200 bg-white px-5 py-2.5 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 disabled:opacity-50"
                >
                  Voltar
                </button>

                <button
                  type="submit"
                  disabled={
                    cancelando ||
                    !motivoCancelamento.trim() ||
                    [...motivoCancelamento.trim()].length > 1000
                  }
                  className="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  {cancelando ? (
                    <Loader2
                      size={17}
                      className="animate-spin"
                    />
                  ) : (
                    <Ban size={17} />
                  )}

                  {cancelando
                    ? 'Cancelando...'
                    : 'Confirmar cancelamento'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
