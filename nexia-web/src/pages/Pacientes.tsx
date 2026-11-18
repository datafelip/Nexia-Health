
import { useEffect, useMemo, useState } from 'react'
import {
  AlertCircle,
  CheckCircle2,
  ClipboardList,
  Loader2,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Trash2,
  UserRound,
  Users,
  X,
} from 'lucide-react'
import axios from 'axios'

import { pacienteService } from '../services/pacienteService'
import type { Paciente } from '../types/paciente'
import { filtrarPacientesVisiveis } from '../utils/pacienteUtils'
import PacienteForm from '../components/pacientes/PacienteForm'

function formatarCPF(cpf?: string | null): string {
  if (!cpf) return '—'

  const numeros = cpf.replace(/\D/g, '')

  if (numeros.length !== 11) return cpf

  return numeros.replace(
    /(\d{3})(\d{3})(\d{3})(\d{2})/,
    '$1.$2.$3-$4'
  )
}

function formatarDataNascimento(
  valor?: string | null
): string {
  if (!valor) return '—'

  const data = new Date(valor)

  if (Number.isNaN(data.getTime())) return '—'

  return data.toLocaleDateString('pt-BR', {
    timeZone: 'UTC',
  })
}

function obterIniciais(nome: string): string {
  const partes = nome.trim().split(/\s+/)

  return partes
    .slice(0, 2)
    .map((parte) => parte.charAt(0).toUpperCase())
    .join('')
}

export default function Pacientes() {
  const [pacientes, setPacientes] = useState<Paciente[]>([])
  const [pesquisa, setPesquisa] = useState('')
  const [carregando, setCarregando] = useState(true)
  const [erro, setErro] = useState('')
  const [mostrarFormulario, setMostrarFormulario] =
    useState(false)
  const [mensagemSucesso, setMensagemSucesso] = useState('')
  const [pacienteEmEdicao, setPacienteEmEdicao] =
    useState<Paciente | null>(null)
  const [pacienteExcluindo, setPacienteExcluindo] =
    useState<string | null>(null)

  async function carregarPacientes() {
    try {
      setCarregando(true)
      setErro('')

      const dados = await pacienteService.listar()

      const pacientesVisiveis = filtrarPacientesVisiveis(
        Array.isArray(dados) ? dados : []
      )

      setPacientes(pacientesVisiveis)
    } catch (error) {
      console.error('Erro ao carregar pacientes:', error)

      if (axios.isAxiosError(error)) {
        const mensagem = error.response?.data?.error

        setErro(
          typeof mensagem === 'string'
            ? mensagem
            : 'Não foi possível conectar à API.'
        )
      } else {
        setErro('Ocorreu um erro inesperado.')
      }
    } finally {
      setCarregando(false)
    }
  }

  useEffect(() => {
    void carregarPacientes()
  }, [])

  const pacientesFiltrados = useMemo(() => {
    const termo = pesquisa.toLowerCase().trim()
    const termoNumerico = pesquisa.replace(/\D/g, '')

    if (!termo) return pacientes

    return pacientes.filter((paciente) => {
      const nomeEncontrado = paciente.nome
        .toLowerCase()
        .includes(termo)

      const cpfEncontrado =
        termoNumerico.length > 0 &&
        paciente.cpf
          .replace(/\D/g, '')
          .includes(termoNumerico)

      return nomeEncontrado || cpfEncontrado
    })
  }, [pacientes, pesquisa])

  const totalComTelefone = useMemo(
    () =>
      pacientes.filter((paciente) =>
        Boolean(paciente.telefone?.trim())
      ).length,
    [pacientes]
  )

  function abrirFormulario() {
    setPacienteEmEdicao(null)
    setMensagemSucesso('')
    setMostrarFormulario(true)
  }

  function editarPaciente(paciente: Paciente) {
    setPacienteEmEdicao(paciente)
    setMensagemSucesso('')
    setMostrarFormulario(true)

    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
  }

  function fecharFormulario() {
    setMostrarFormulario(false)
    setPacienteEmEdicao(null)
  }

  async function cadastroRealizado() {
    const estavaEditando = Boolean(pacienteEmEdicao)

    setMostrarFormulario(false)
    setPacienteEmEdicao(null)
    setPesquisa('')

    setMensagemSucesso(
      estavaEditando
        ? 'Paciente atualizado com sucesso!'
        : 'Paciente cadastrado com sucesso!'
    )

    await carregarPacientes()
  }

  async function excluirPaciente(paciente: Paciente) {
    if (pacienteExcluindo !== null) return

    const confirmado = window.confirm(
      `Deseja realmente excluir o paciente "${paciente.nome}"?\n\n` +
        'Esta operação não poderá ser desfeita.'
    )

    if (!confirmado) return

    try {
      setPacienteExcluindo(paciente.id)
      setErro('')
      setMensagemSucesso('')

      await pacienteService.excluir(paciente.id)

      setPacientes((anteriores) =>
        anteriores.filter((item) => item.id !== paciente.id)
      )

      setMensagemSucesso(
        `Paciente "${paciente.nome}" excluído com sucesso!`
      )
    } catch (error) {
      console.error('Erro ao excluir paciente:', error)

      if (axios.isAxiosError(error)) {
        const mensagem = error.response?.data?.error

        setErro(
          typeof mensagem === 'string'
            ? mensagem
            : 'Não foi possível excluir o paciente. Verifique se existem registros vinculados.'
        )
      } else {
        setErro(
          'Ocorreu um erro inesperado ao excluir o paciente.'
        )
      }
    } finally {
      setPacienteExcluindo(null)
    }
  }

  return (
    <div className="space-y-7 pb-4">
      {/* Cabeçalho premium */}
      <section className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-6 shadow-sm sm:p-7">
        <div className="pointer-events-none absolute -right-16 -top-20 h-56 w-56 rounded-full bg-teal-50 blur-3xl" />

        <div className="relative flex flex-col justify-between gap-5 md:flex-row md:items-center">
          <div className="flex items-start gap-4">
            <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-teal-600 text-white shadow-lg shadow-teal-100">
              <Users size={25} />
            </div>

            <div>
              <div className="mb-2 flex flex-wrap items-center gap-2">
                <span className="rounded-full border border-teal-100 bg-teal-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-teal-700">
                  Gestão de pacientes
                </span>

                <span className="inline-flex items-center gap-1.5 text-xs font-medium text-slate-500">
                  <span className="h-2 w-2 rounded-full bg-teal-500" />
                  Cadastro integrado
                </span>
              </div>

              <h2 className="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
                Gerenciamento de Pacientes
              </h2>

              <p className="mt-2 max-w-xl text-sm leading-relaxed text-slate-500">
                Consulte, cadastre e atualize as informações
                dos pacientes registrados na unidade.
              </p>
            </div>
          </div>

          {!mostrarFormulario && (
            <button
              type="button"
              onClick={abrirFormulario}
              className="inline-flex shrink-0 items-center justify-center gap-2 self-start rounded-xl bg-teal-600 px-5 py-3 text-sm font-semibold text-white shadow-sm shadow-teal-100 transition hover:bg-teal-700"
            >
              <Plus size={18} />
              Novo paciente
            </button>
          )}
        </div>
      </section>

      {/* Mensagem de sucesso */}
      {mensagemSucesso && (
        <div
          role="status"
          className="flex items-center justify-between gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm font-medium text-emerald-700"
        >
          <div className="flex items-center gap-2">
            <CheckCircle2 size={19} className="shrink-0" />
            <span>{mensagemSucesso}</span>
          </div>

          <button
            type="button"
            onClick={() => setMensagemSucesso('')}
            aria-label="Fechar mensagem"
            className="rounded-lg p-1 hover:bg-emerald-100"
          >
            <X size={17} />
          </button>
        </div>
      )}

      {/* Formulário existente preservado */}
      {mostrarFormulario && (
        <PacienteForm
          key={pacienteEmEdicao?.id ?? 'novo'}
          paciente={pacienteEmEdicao}
          onSuccess={() => {
            void cadastroRealizado()
          }}
          onCancel={fecharFormulario}
        />
      )}

      {/* Indicadores */}
      <section className="grid gap-4 sm:grid-cols-3">
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="flex items-start justify-between">
            <div className="rounded-xl bg-blue-50 p-3 text-blue-600">
              <Users size={21} />
            </div>

            <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
              Total
            </span>
          </div>

          <p className="mt-4 text-sm font-medium text-slate-500">
            Pacientes cadastrados
          </p>

          <p className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            {carregando && pacientes.length === 0
              ? '—'
              : pacientes.length}
          </p>

          <p className="mt-1 text-xs text-slate-400">
            Registros disponíveis na listagem
          </p>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="flex items-start justify-between">
            <div className="rounded-xl bg-teal-50 p-3 text-teal-600">
              <ClipboardList size={21} />
            </div>

            <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
              Pesquisa
            </span>
          </div>

          <p className="mt-4 text-sm font-medium text-slate-500">
            Pacientes encontrados
          </p>

          <p className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            {carregando && pacientes.length === 0
              ? '—'
              : pacientesFiltrados.length}
          </p>

          <p className="mt-1 text-xs text-slate-400">
            Conforme os filtros aplicados
          </p>
        </div>

        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="flex items-start justify-between">
            <div className="rounded-xl bg-violet-50 p-3 text-violet-600">
              <ShieldCheck size={21} />
            </div>

            <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
              Contato
            </span>
          </div>

          <p className="mt-4 text-sm font-medium text-slate-500">
            Cadastros com telefone
          </p>

          <p className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            {carregando && pacientes.length === 0
              ? '—'
              : totalComTelefone}
          </p>

          <p className="mt-1 text-xs text-slate-400">
            Pacientes com contato informado
          </p>
        </div>
      </section>

      {/* Listagem */}
      <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
        <div className="flex flex-col justify-between gap-4 border-b border-slate-100 p-5 sm:p-6 lg:flex-row lg:items-center">
          <div>
            <h3 className="text-base font-bold text-slate-900">
              Lista de pacientes
            </h3>

            <p className="mt-1 text-xs text-slate-500">
              Consulte os dados e gerencie os cadastros.
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
                aria-label="Pesquisar pacientes por nome ou CPF"
                placeholder="Pesquisar por nome ou CPF..."
                value={pesquisa}
                onChange={(e) => setPesquisa(e.target.value)}
                className="w-full rounded-xl border border-slate-200 bg-slate-50 py-2.5 pl-10 pr-4 text-sm text-slate-700 outline-none transition placeholder:text-slate-400 focus:border-teal-500 focus:bg-white focus:ring-2 focus:ring-teal-500/10"
              />
            </div>

            <button
              type="button"
              onClick={() => void carregarPacientes()}
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
              size={30}
              className="animate-spin text-teal-600"
            />

            <p className="text-sm">
              Carregando pacientes...
            </p>
          </div>
        ) : pacientesFiltrados.length === 0 && !erro ? (
          <div className="flex flex-col items-center gap-3 px-5 py-16 text-center">
            <div className="rounded-2xl bg-slate-50 p-4 text-slate-400">
              <Users size={29} />
            </div>

            <p className="text-sm font-semibold text-slate-700">
              {pesquisa
                ? 'Nenhum paciente encontrado'
                : 'Nenhum paciente cadastrado'}
            </p>

            <p className="max-w-sm text-xs text-slate-400">
              {pesquisa
                ? 'Tente pesquisar por outro nome ou CPF.'
                : 'Os pacientes cadastrados aparecerão nesta lista.'}
            </p>
          </div>
        ) : pacientesFiltrados.length === 0 && erro ? (
          <div className="flex flex-col items-center gap-3 px-5 py-14 text-center">
            <div className="rounded-2xl bg-red-50 p-4 text-red-500">
              <AlertCircle size={28} />
            </div>

            <p className="text-sm font-medium text-slate-600">
              Não foi possível carregar a lista.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[850px] text-left">
              <thead className="bg-slate-50/80">
                <tr className="text-[11px] font-bold uppercase tracking-wider text-slate-500">
                  <th scope="col" className="px-6 py-4">
                    Paciente
                  </th>

                  <th scope="col" className="px-6 py-4">
                    CPF
                  </th>

                  <th scope="col" className="px-6 py-4">
                    Sexo
                  </th>

                  <th scope="col" className="px-6 py-4">
                    Telefone
                  </th>

                  <th scope="col" className="px-6 py-4">
                    Nascimento
                  </th>

                  <th scope="col" className="px-6 py-4 text-right">
                    Ações
                  </th>
                </tr>
              </thead>

              <tbody className="divide-y divide-slate-100">
                {pacientesFiltrados.map((paciente) => (
                  <tr
                    key={paciente.id}
                    className="transition-colors hover:bg-teal-50/30"
                  >
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-teal-50 text-xs font-bold text-teal-700">
                          {obterIniciais(paciente.nome)}
                        </div>

                        <div className="min-w-0">
                          <p className="font-semibold text-slate-800">
                            {paciente.nome}
                          </p>

                          <p className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
                            <UserRound size={12} />
                            Paciente
                          </p>
                        </div>
                      </div>
                    </td>

                    <td className="px-6 py-4 text-sm font-medium text-slate-600">
                      {formatarCPF(paciente.cpf)}
                    </td>

                    <td className="px-6 py-4">
                      <span className="inline-flex rounded-lg bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-600">
                        {paciente.sexo || '—'}
                      </span>
                    </td>

                    <td className="px-6 py-4 text-sm text-slate-600">
                      {paciente.telefone || '—'}
                    </td>

                    <td className="px-6 py-4 text-sm text-slate-600">
                      {formatarDataNascimento(
                        paciente.data_nascimento
                      )}
                    </td>

                    {/* Ações: editar e excluir */}
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <button
                          type="button"
                          onClick={() => editarPaciente(paciente)}
                          disabled={pacienteExcluindo !== null}
                          title={`Editar ${paciente.nome}`}
                          className="inline-flex items-center gap-2 rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-xs font-semibold text-teal-700 transition hover:border-teal-300 hover:bg-teal-50 disabled:cursor-not-allowed disabled:opacity-50"
                        >
                          <Pencil size={15} />
                          Editar
                        </button>

                        <button
                          type="button"
                          onClick={() => void excluirPaciente(paciente)}
                          disabled={pacienteExcluindo !== null}
                          title={`Excluir ${paciente.nome}`}
                          className="inline-flex items-center gap-2 rounded-xl border border-red-200 bg-white px-3.5 py-2 text-xs font-semibold text-red-600 transition hover:border-red-300 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-50"
                        >
                          {pacienteExcluindo === paciente.id ? (
                            <Loader2
                              size={15}
                              className="animate-spin"
                            />
                          ) : (
                            <Trash2 size={15} />
                          )}

                          {pacienteExcluindo === paciente.id
                            ? 'Excluindo...'
                            : 'Excluir'}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Rodapé da tabela */}
        {!carregando && (
          <div className="flex flex-wrap items-center justify-between gap-2 border-t border-slate-100 bg-slate-50/50 px-6 py-4">
            <p className="text-xs text-slate-500">
              Exibindo{' '}
              <strong className="text-slate-700">
                {pacientesFiltrados.length}
              </strong>{' '}
              de{' '}
              <strong className="text-slate-700">
                {pacientes.length}
              </strong>{' '}
              pacientes
            </p>

            <span className="inline-flex items-center gap-1.5 text-xs font-medium text-teal-700">
              <CheckCircle2 size={14} />
              Listagem atualizada
            </span>
          </div>
        )}
      </section>
    </div>
  )
}
