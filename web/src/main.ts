import './style.css'

type Feature = {
  title: string
  description: string
  tone: string
  icon: string
}

type Install = {
  title: string
  note: string
  command: string
}

const screenshotUrl = `${import.meta.env.BASE_URL}images/screenshot-terminal.png`
const docsUrl = 'https://github.com/Equationzhao/g/blob/master/docs/man.md'
const releasesUrl = 'https://github.com/Equationzhao/g/releases'
const githubUrl = 'https://github.com/Equationzhao/g'

const features: Feature[] = [
  {
    title: 'Fast Performance',
    description: 'Parallel processing for large directory listings and quick terminal feedback.',
    tone: 'lime',
    icon: `
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="M4 14a8 8 0 1 1 14.2 4.98" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
        <path d="M13.5 6 16 3.5 18.5 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    `,
  },
  {
    title: 'Visual Output',
    description: 'File icons, vibrant colors, and multiple layouts tuned for modern terminals.',
    tone: 'amber',
    icon: `
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <rect x="4" y="5" width="16" height="14" rx="3" stroke="currentColor" stroke-width="1.8"/>
        <path d="m7.5 15 3.1-3.2 2.7 2.6 3.2-4.4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    `,
  },
  {
    title: 'Git Integration',
    description: 'Surface branch, repo status, and file status indicators directly in listings.',
    tone: 'coral',
    icon: `
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="M8 5.5 4.5 9 8 12.5M16 18.5 19.5 15 16 11.5M9.5 13.5l5-5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    `,
  },
  {
    title: 'Cross-Platform Support',
    description: 'Works across macOS, Linux, and Windows without changing how you list files.',
    tone: 'cyan',
    icon: `
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <rect x="4" y="5" width="16" height="10" rx="2.5" stroke="currentColor" stroke-width="1.8"/>
        <path d="M9 19h6M12 15v4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>
      </svg>
    `,
  },
  {
    title: 'Metadata Details',
    description: 'See permissions, owners, timestamps, and sizes at a glance in long listings.',
    tone: 'blue',
    icon: `
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <circle cx="12" cy="12" r="8" stroke="currentColor" stroke-width="1.8"/>
        <path d="M12 8v4l2.5 2.5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    `,
  },
  {
    title: 'Filtering & Sorting',
    description: 'Version-sort and flexible listing options help cut through noisy directories.',
    tone: 'violet',
    icon: `
      <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
        <path d="M7 6v12M7 18l-2.5-2.5M7 18l2.5-2.5M17 18V6M17 6l-2.5 2.5M17 6l2.5 2.5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    `,
  },
]

const installs: Install[] = [
  {
    title: 'Homebrew',
    note: 'Recommended for macOS',
    command: 'brew install g-ls',
  },
  {
    title: 'Install Script',
    note: 'One-line shell install',
    command: 'bash -c "$(curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh)"',
  },
  {
    title: 'Go Install',
    note: 'Latest release via Go toolchain',
    command: 'go install -ldflags="-s -w" github.com/Equationzhao/g@latest',
  },
  {
    title: 'Windows Scoop',
    note: 'Install from the repo manifest',
    command: 'scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json',
  },
]

const copyIcon = `
  <svg viewBox="0 0 24 24" fill="none" aria-hidden="true">
    <rect x="9" y="9" width="10" height="10" rx="2" stroke="currentColor" stroke-width="1.7"/>
    <path d="M7 15H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h7a2 2 0 0 1 2 2v1" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"/>
  </svg>
`

const escapeHtml = (value: string) =>
  value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')

const featureMarkup = features
  .map(
    (feature) => `
      <article class="feature-card tone-${feature.tone}">
        <span class="feature-card__icon">${feature.icon}</span>
        <h3>${feature.title}</h3>
        <p>${feature.description}</p>
      </article>
    `,
  )
  .join('')

const installMarkup = installs
  .map(
    (install) => `
      <article class="install-card">
        <div class="install-card__header">
          <div>
            <h3>${escapeHtml(install.title)}</h3>
            <p class="install-card__note">${escapeHtml(install.note)}</p>
          </div>
        </div>
        <div class="command-row">
          <code class="command">${escapeHtml(install.command)}</code>
          <button class="copy-button" type="button" data-copy="${escapeHtml(install.command)}" aria-label="Copy ${escapeHtml(install.title)} command">
            ${copyIcon}
            <span class="visually-hidden">Copy command</span>
          </button>
        </div>
      </article>
    `,
  )
  .join('')

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div class="site-shell" data-menu-open="false">
    <header class="terminal-chrome">
      <a class="brand" href="#top" aria-label="g home">g</a>
      <button class="menu-toggle" type="button" aria-expanded="false" aria-label="Toggle navigation">
        <span></span>
        <span></span>
        <span></span>
      </button>
      <nav class="site-nav" aria-label="Primary">
        <a href="${docsUrl}" target="_blank" rel="noreferrer">Documentation</a>
        <a href="${releasesUrl}" target="_blank" rel="noreferrer">Releases</a>
        <a href="${githubUrl}" target="_blank" rel="noreferrer">GitHub</a>
      </nav>
    </header>

    <main class="page" id="top">
      <section class="hero section">
        <h1>
          <span class="accent">g</span>: Cross-platform for <em>ls</em> replacement.
        </h1>
        <p class="hero-copy">
          A fast, visual tool for directory listings with powerful features. Built with Go.
        </p>
        <div class="hero-actions">
          <a class="button button--primary" href="#installs">Install g</a>
          <a class="button button--secondary" href="${githubUrl}" target="_blank" rel="noreferrer">View on GitHub</a>
        </div>
      </section>

      <section class="section" id="features">
        <div class="section-heading">
          <h2>Features</h2>
        </div>
        <div class="feature-grid">
          ${featureMarkup}
        </div>
      </section>

      <section class="section section--showcase" aria-labelledby="showcase-title">
        <div class="section-heading">
          <h2 id="showcase-title">Terminal Preview</h2>
        </div>
        <div class="showcase-grid">
          <article class="terminal-window terminal-window--primary">
            <div class="window-header">
              <div class="window-dots" aria-hidden="true">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <p class="window-command">g -A --table --table-style=unicode --file --long</p>
            </div>
            <img
              src="${screenshotUrl}"
              alt="g command showing colorful file listings with metadata and git indicators."
            />
          </article>
        </div>
      </section>

      <section class="section" id="installs">
        <div class="section-heading">
          <h2>Installs</h2>
        </div>
        <div class="install-grid">
          ${installMarkup}
        </div>
      </section>

      <section class="section section--cta">
        <h2>Ready to modernize your terminal?</h2>
        <p>Explore and support 'g'.</p>
        <div class="footer-links">
          <a href="${docsUrl}" target="_blank" rel="noreferrer">Documentation</a>
          <a href="${releasesUrl}" target="_blank" rel="noreferrer">Releases</a>
          <a href="${githubUrl}" target="_blank" rel="noreferrer">Give us a Star</a>
        </div>
      </section>
    </main>
  </div>
`

const shell = document.querySelector<HTMLElement>('.site-shell')
const menuToggle = document.querySelector<HTMLButtonElement>('.menu-toggle')
const navLinks = document.querySelectorAll<HTMLAnchorElement>('.site-nav a')

menuToggle?.addEventListener('click', () => {
  const isOpen = shell?.dataset.menuOpen === 'true'

  shell?.setAttribute('data-menu-open', String(!isOpen))
  menuToggle.setAttribute('aria-expanded', String(!isOpen))
})

navLinks.forEach((link) => {
  link.addEventListener('click', () => {
    shell?.setAttribute('data-menu-open', 'false')
    menuToggle?.setAttribute('aria-expanded', 'false')
  })
})

const fallbackCopy = async (value: string) => {
  const buffer = document.createElement('textarea')
  buffer.value = value
  buffer.setAttribute('readonly', '')
  buffer.style.position = 'absolute'
  buffer.style.opacity = '0'
  document.body.appendChild(buffer)
  buffer.select()
  document.execCommand('copy')
  document.body.removeChild(buffer)
}

document.querySelectorAll<HTMLButtonElement>('[data-copy]').forEach((button) => {
  button.addEventListener('click', async () => {
    const command = button.dataset.copy

    if (!command) {
      return
    }

    try {
      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(command)
      } else {
        await fallbackCopy(command)
      }

      button.dataset.state = 'copied'
      window.setTimeout(() => {
        button.dataset.state = 'idle'
      }, 1600)
    } catch {
      button.dataset.state = 'error'
      window.setTimeout(() => {
        button.dataset.state = 'idle'
      }, 1600)
    }
  })
})
