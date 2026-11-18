
import { useEffect, useRef, useState } from 'react'
import { Outlet, useLocation } from 'react-router-dom'
import {
  Activity,
  ChevronRight,
  HeartPulse,
  Menu,
  X,
} from 'lucide-react'

import Sidebar from './Sidebar'

interface PageInfo {
  title: string
  description: string
}

const pageInfo: Record<string, PageInfo> = {
  '/': {
    title: 'Dashboard',
    description: 'Visão geral da operação hospitalar',
  },
  '/pacientes': {
    title: 'Pacientes',
    description: 'Cadastro e gerenciamento de pacientes',
  },
  '/triagens': {
    title: 'Triagens',
    description: 'Avaliação clínica e classificação de risco',
  },
  '/filas': {
    title: 'Fila de Atendimento',
    description: 'Organização e acompanhamento das chamadas',
  },
  '/atendimentos': {
    title: 'Atendimentos',
    description: 'Gestão dos atendimentos médicos',
  },
}

export default function MainLayout() {
  const location = useLocation()

  const [sidebarAberta, setSidebarAberta] = useState(false)

  const botaoMenuRef = useRef<HTMLButtonElement>(null)
  const botaoFecharRef = useRef<HTMLButtonElement>(null)
  const menuMobileRef = useRef<HTMLDivElement>(null)

  const pagina = pageInfo[location.pathname] ?? {
    title: 'Nexia Health',
    description: 'Sistema de gestão hospitalar',
  }

  function fecharSidebar() {
    setSidebarAberta(false)
  }

  // Fecha o menu quando o usuário navega para outra rota.
  useEffect(() => {
    setSidebarAberta(false)
  }, [location.pathname])

  // Bloqueia a rolagem do fundo e gerencia o teclado no mobile.
  useEffect(() => {
    if (!sidebarAberta) {
      return
    }

    const overflowAnterior = document.body.style.overflow

    document.body.style.overflow = 'hidden'

    botaoFecharRef.current?.focus()

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setSidebarAberta(false)
        botaoMenuRef.current?.focus()
        return
      }

      if (event.key !== 'Tab') {
        return
      }

      const elementos = menuMobileRef.current?.querySelectorAll<HTMLElement>(
        'a[href], button:not([disabled]), [tabindex]:not([tabindex="-1"])'
      )

      if (!elementos || elementos.length === 0) {
        return
      }

      const primeiro = elementos[0]
      const ultimo = elementos[elementos.length - 1]

      if (
        event.shiftKey &&
        document.activeElement === primeiro
      ) {
        event.preventDefault()
        ultimo.focus()
      } else if (
        !event.shiftKey &&
        document.activeElement === ultimo
      ) {
        event.preventDefault()
        primeiro.focus()
      }
    }

    function handleResize() {
      if (window.innerWidth >= 1024) {
        setSidebarAberta(false)
      }
    }

    document.addEventListener('keydown', handleKeyDown)
    window.addEventListener('resize', handleResize)

    return () => {
      document.body.style.overflow = overflowAnterior
      document.removeEventListener('keydown', handleKeyDown)
      window.removeEventListener('resize', handleResize)
    }
  }, [sidebarAberta])

  return (
    <div className="min-h-screen bg-[#F8FAFC] text-slate-900">
      {/* Link de acessibilidade */}
      <a
        href="#conteudo-principal"
        className="
          sr-only
          fixed left-4 top-4 z-[100]
          rounded-xl bg-teal-700 px-4 py-3
          font-semibold text-white shadow-lg
          focus:not-sr-only
        "
      >
        Pular para o conteúdo principal
      </a>

      {/* Sidebar desktop */}
      <aside
        aria-label="Navegação principal"
        className="
          fixed inset-y-0 left-0 z-40
          hidden w-64 lg:block
        "
      >
        <Sidebar />
      </aside>

      {/* Sidebar mobile */}
      {sidebarAberta && (
        <div className="fixed inset-0 z-50 lg:hidden">
          {/* Fundo escurecido */}
          <button
            type="button"
            aria-label="Fechar menu de navegação"
            onClick={fecharSidebar}
            className="
              absolute inset-0
              bg-slate-950/60
              backdrop-blur-sm
              transition-opacity
            "
          />

          {/* Painel lateral */}
          <div
            ref={menuMobileRef}
            role="dialog"
            aria-modal="true"
            aria-label="Menu de navegação"
            className="
              relative z-10
              h-full w-[min(18rem,85vw)]
              shadow-2xl shadow-slate-950/30
            "
          >
            <Sidebar onNavigate={fecharSidebar} />

            <button
              ref={botaoFecharRef}
              type="button"
              aria-label="Fechar menu"
              onClick={() => {
                fecharSidebar()
                botaoMenuRef.current?.focus()
              }}
              className="
                absolute right-3 top-4 z-20
                flex h-9 w-9 items-center justify-center
                rounded-xl bg-white/10 text-white
                transition
                hover:bg-white/20
                focus-visible:outline
                focus-visible:outline-2
                focus-visible:outline-offset-2
                focus-visible:outline-teal-300
              "
            >
              <X size={18} />
            </button>
          </div>
        </div>
      )}

      {/* Estrutura principal */}
      <div className="flex min-h-screen min-w-0 flex-col lg:ml-64">
        {/* Cabeçalho */}
        <header
          className="
            sticky top-0 z-30
            border-b border-slate-200/80
            bg-white/95 backdrop-blur-xl
          "
        >
          <div
            className="
              flex min-h-[72px] items-center
              justify-between gap-3
              px-4 py-3
              sm:min-h-20 sm:gap-4 sm:px-6
              lg:px-8
            "
          >
            {/* Título e navegação */}
            <div className="flex min-w-0 flex-1 items-center gap-3 sm:gap-4">
              {/* Botão mobile */}
              <button
                ref={botaoMenuRef}
                type="button"
                onClick={() => setSidebarAberta(true)}
                aria-label="Abrir menu de navegação"
                aria-expanded={sidebarAberta}
                aria-controls="menu-mobile-nexia"
                className="
                  flex h-10 w-10 shrink-0
                  items-center justify-center
                  rounded-xl border border-slate-200
                  bg-white text-slate-600
                  transition
                  hover:border-teal-200
                  hover:bg-teal-50
                  hover:text-teal-700
                  focus-visible:outline
                  focus-visible:outline-2
                  focus-visible:outline-offset-2
                  focus-visible:outline-teal-600
                  lg:hidden
                "
              >
                <Menu size={20} />
              </button>

              <div className="min-w-0 flex-1">
                {/* Breadcrumb */}
                <nav
                  aria-label="Localização atual"
                  className="
                    mb-1 hidden items-center gap-2
                    text-xs font-medium text-slate-400
                    sm:flex
                  "
                >
                  <span>Workspace</span>

                  <ChevronRight
                    size={13}
                    aria-hidden="true"
                  />

                  <span
                    aria-current="page"
                    className="truncate text-teal-700"
                  >
                    {pagina.title}
                  </span>
                </nav>

                {/* Título */}
                <h1
                  className="
                    truncate text-base font-bold
                    tracking-tight text-slate-900
                    sm:text-xl
                  "
                >
                  {pagina.title}
                </h1>

                {/* Descrição */}
                <p
                  className="
                    mt-0.5 hidden truncate
                    text-xs text-slate-500
                    sm:block
                  "
                >
                  {pagina.description}
                </p>
              </div>
            </div>

            {/* Identidade da plataforma */}
            <div className="flex shrink-0 items-center gap-3">
              {/* Indicador da plataforma */}
              <div
                className="
                  hidden items-center gap-2
                  rounded-full border border-teal-100
                  bg-teal-50 px-3 py-1.5
                  xl:flex
                "
              >
                <Activity
                  size={14}
                  className="text-teal-600"
                  aria-hidden="true"
                />

                <span className="text-xs font-semibold text-teal-700">
                  Nexia Health
                </span>
              </div>

              {/* Divisor */}
              <div className="hidden h-8 w-px bg-slate-200 sm:block" />

              {/* Identidade visual */}
              <div className="flex items-center gap-2.5">
                <div
                  className="
                    flex h-9 w-9 shrink-0
                    items-center justify-center
                    rounded-xl
                    bg-gradient-to-br
                    from-teal-500 to-teal-700
                    text-white
                    shadow-sm shadow-teal-200
                    sm:h-10 sm:w-10
                  "
                >
                  <HeartPulse
                    size={21}
                    aria-hidden="true"
                  />
                </div>

                <div className="hidden 2xl:block">
                  <p className="text-sm font-semibold text-slate-800">
                    Nexia Health
                  </p>

                  <p className="text-xs text-slate-400">
                    Medical Platform
                  </p>
                </div>
              </div>
            </div>
          </div>
        </header>

        {/* Conteúdo das páginas */}
        <main
          id="conteudo-principal"
          tabIndex={-1}
          className="
            mx-auto w-full min-w-0
            max-w-[1600px] flex-1
            px-4 py-5
            sm:px-6 sm:py-6
            lg:px-8 lg:py-8
          "
        >
          <Outlet />
        </main>

        {/* Rodapé */}
        <footer
          className="
            mx-auto flex w-full max-w-[1600px]
            flex-col gap-2
            border-t border-slate-200/70
            px-4 py-5
            text-xs text-slate-400
            sm:flex-row sm:items-center
            sm:justify-between sm:px-6
            lg:px-8
          "
        >
          <span>
            © {new Date().getFullYear()} Nexia Health
          </span>

          <span>
            Plataforma de gestão hospitalar
          </span>
        </footer>
      </div>
    </div>
  )
}
