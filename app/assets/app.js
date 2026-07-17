import 'htmx.org'
import Alpine from 'alpinejs'
import maplibregl from 'maplibre-gl'

const themeKey = 'skyvisor-theme'

function applyTheme(mode) {
  const dark = mode === 'dark' || (mode === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.dataset.theme = mode
  document.documentElement.classList.toggle('dark', dark)
  localStorage.setItem(themeKey, mode)
  window.dispatchEvent(new CustomEvent('skyvisor:theme', { detail: { mode, dark } }))
}

Alpine.data('appShell', () => ({
  menuOpen: false,
  theme: localStorage.getItem(themeKey) || 'system',
  cycleTheme() {
    const modes = ['system', 'light', 'dark']
    this.theme = modes[(modes.indexOf(this.theme) + 1) % modes.length]
    applyTheme(this.theme)
  },
}))

Alpine.data('flightSearch', () => ({
  flight: '',
  submit() {
    this.flight = this.flight.trim().toUpperCase().replace(/\s+/g, '')
  },
}))

function initMaps(root = document) {
  root.querySelectorAll('[data-skyvisor-map]').forEach((element) => {
    if (element.dataset.mapReady === 'true') return
    element.dataset.mapReady = 'true'

    const center = [Number(element.dataset.longitude || 0), Number(element.dataset.latitude || 28)]
    const map = new maplibregl.Map({
      container: element,
      style: element.dataset.styleUrl || 'https://demotiles.maplibre.org/style.json',
      center,
      zoom: Number(element.dataset.zoom || 1.25),
      attributionControl: false,
      cooperativeGestures: true,
    })
    map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'bottom-right')
    map.addControl(new maplibregl.AttributionControl({ compact: true }))
  })
}

window.Alpine = Alpine
window.maplibregl = maplibregl
window.SkyVisor = { applyTheme, initMaps }

Alpine.start()

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => initMaps())
} else {
  initMaps()
}

document.body.addEventListener('htmx:afterSwap', (event) => initMaps(event.detail.target))

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
  if ((localStorage.getItem(themeKey) || 'system') === 'system') applyTheme('system')
})
