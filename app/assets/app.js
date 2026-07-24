import 'htmx.org'
import Alpine from 'alpinejs'
import maplibregl from 'maplibre-gl'
import { animate, stagger } from 'motion'

const themeKey = 'skyvisor-theme'
const reducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches

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

const AIRPORT_COORDS = {
  LIS: [-9.1359, 38.7813],
  LHR: [-0.4543, 51.47],
  JFK: [-73.7781, 40.6413],
  LAX: [-118.4085, 33.9425],
  CDG: [2.5479, 49.0097],
  FRA: [8.5622, 50.0379],
  DXB: [55.3657, 25.2532],
  SIN: [103.9915, 1.3644],
  HND: [139.7798, 35.5494],
  SYD: [151.177, -33.9399],
  GRU: [-46.473, -23.4356],
  ORD: [-87.9048, 41.9786],
  ATL: [-84.4281, 33.6367],
  DFW: [-97.038, 32.8998],
  MIA: [-80.2906, 25.7959],
  AMS: [4.7639, 52.3105],
  MAD: [-3.5676, 40.4983],
  FCO: [12.2389, 41.8003],
  IST: [28.8146, 41.2753],
  DOH: [51.608, 25.2731],
}

function parseNumber(value) {
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : null
}

function greatCirclePoints(from, to, steps = 64) {
  const points = []
  for (let i = 0; i <= steps; i += 1) {
    const t = i / steps
    points.push([from[0] + (to[0] - from[0]) * t, from[1] + (to[1] - from[1]) * t])
  }
  return points
}

function routePointAtProgress(from, to, progress) {
  const points = greatCirclePoints(from, to, 48)
  const index = Math.min(points.length - 1, Math.max(0, Math.round((progress / 100) * (points.length - 1))))
  return points[index]
}

function flightLayers(map, element) {
  const depLon = parseNumber(element.dataset.depLon)
  const depLat = parseNumber(element.dataset.depLat)
  const arrLon = parseNumber(element.dataset.arrLon)
  const arrLat = parseNumber(element.dataset.arrLat)
  const liveLon = parseNumber(element.dataset.liveLon)
  const liveLat = parseNumber(element.dataset.liveLat)
  if (depLon == null || depLat == null || arrLon == null || arrLat == null) return

  const dep = [depLon, depLat]
  const arr = [arrLon, arrLat]
  const route = greatCirclePoints(dep, arr)
  const progress = parseNumber(element.dataset.progress) ?? 42
  const planePos = liveLon != null && liveLat != null ? [liveLon, liveLat] : routePointAtProgress(dep, arr, progress)

  if (!map.getSource('flight-route')) {
    map.addSource('flight-route', {
      type: 'geojson',
      data: { type: 'Feature', geometry: { type: 'LineString', coordinates: route }, properties: {} },
    })
    map.addLayer({
      id: 'flight-route-line',
      type: 'line',
      source: 'flight-route',
      paint: { 'line-color': '#60a5fa', 'line-width': 2.5, 'line-opacity': 0.85 },
    })
    map.addSource('flight-plane', {
      type: 'geojson',
      data: { type: 'Feature', geometry: { type: 'Point', coordinates: planePos }, properties: {} },
    })
    map.addLayer({
      id: 'flight-plane-dot',
      type: 'circle',
      source: 'flight-plane',
      paint: {
        'circle-radius': 7,
        'circle-color': '#38bdf8',
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
      },
    })
  } else {
    map.getSource('flight-route').setData({ type: 'Feature', geometry: { type: 'LineString', coordinates: route }, properties: {} })
    map.getSource('flight-plane').setData({ type: 'Feature', geometry: { type: 'Point', coordinates: planePos }, properties: {} })
  }

  const bounds = route.reduce((acc, coord) => acc.extend(coord), new maplibregl.LngLatBounds(route[0], route[0]))
  map.fitBounds(bounds, { padding: 72, maxZoom: 6.5, duration: reducedMotion() ? 0 : 900 })
}

function fleetLayers(map, element) {
  let markers = []
  try {
    markers = JSON.parse(element.dataset.markers || '[]')
  } catch {
    markers = []
  }
  const features = markers
    .map((marker) => {
      const coord = AIRPORT_COORDS[marker.iata]
      if (!coord) return null
      return {
        type: 'Feature',
        geometry: { type: 'Point', coordinates: coord },
        properties: { label: marker.flight || marker.iata },
      }
    })
    .filter(Boolean)

  const data = { type: 'FeatureCollection', features }
  if (!map.getSource('fleet-markers')) {
    map.addSource('fleet-markers', { type: 'geojson', data })
    map.addLayer({
      id: 'fleet-markers-dot',
      type: 'circle',
      source: 'fleet-markers',
      paint: {
        'circle-radius': 5,
        'circle-color': '#38bdf8',
        'circle-stroke-width': 1.5,
        'circle-stroke-color': '#ffffff',
      },
    })
  } else {
    map.getSource('fleet-markers').setData(data)
  }
}

// Great-circle routes between hub airports for the marketing globe.
const globeRoutes = [
  { from: [-9.13, 38.77], to: [-73.78, 40.64] }, // LIS -> JFK
  { from: [-0.45, 51.47], to: [55.36, 25.25] }, // LHR -> DXB
  { from: [103.99, 1.36], to: [151.18, -33.95] }, // SIN -> SYD
  { from: [-46.47, -23.43], to: [2.55, 49.01] }, // GRU -> CDG
  { from: [139.78, 35.55], to: [-122.38, 37.62] }, // HND -> SFO
  { from: [28.24, -26.14], to: [72.87, 19.09] }, // JNB -> BOM
]

function globeArcFeatures() {
  return {
    type: 'FeatureCollection',
    features: globeRoutes.map(({ from, to }) => {
      const generator = new window.arc.GreatCircle(
        { x: from[0], y: from[1] },
        { x: to[0], y: to[1] },
      )
      return generator.Arc(64, { offset: 10 }).json()
    }),
  }
}

function globeLayers(map) {
  map.setProjection({ type: 'globe' })
  map.addSource('globe-arcs', { type: 'geojson', data: globeArcFeatures() })
  map.addLayer({
    id: 'globe-arcs-line',
    type: 'line',
    source: 'globe-arcs',
    paint: {
      'line-color': '#60a5fa',
      'line-width': 1.6,
      'line-opacity': 0.85,
    },
  })
  map.addSource('globe-hubs', {
    type: 'geojson',
    data: {
      type: 'FeatureCollection',
      features: globeRoutes.flatMap(({ from, to }) => [from, to]).map((coordinates) => ({
        type: 'Feature',
        geometry: { type: 'Point', coordinates },
        properties: {},
      })),
    },
  })
  map.addLayer({
    id: 'globe-hubs-dot',
    type: 'circle',
    source: 'globe-hubs',
    paint: {
      'circle-radius': 2.5,
      'circle-color': '#93c5fd',
      'circle-opacity': 0.9,
    },
  })
}

// Slow idle rotation; pauses on interaction, disabled for reduced motion.
function spinGlobe(map, element) {
  if (reducedMotion()) return
  let userInteracting = false
  const spin = () => {
    if (userInteracting || map.getZoom() > 3) return
    const center = map.getCenter()
    center.lng += 0.4
    map.easeTo({ center, duration: 1000, easing: (n) => n })
  }
  map.on('mousedown', () => { userInteracting = true })
  map.on('dragstart', () => { userInteracting = true })
  map.on('mouseup', () => { userInteracting = false })
  map.on('touchend', () => { userInteracting = false })
  map.on('moveend', spin)
  const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) spin()
      else map.stop()
    })
  }, { threshold: 0.1 })
  observer.observe(element)
  spin()
}

function initMaps(root = document) {
  root.querySelectorAll('[data-skyvisor-map]').forEach((element) => {
    const mode = element.dataset.mapMode || 'world'
    const needsReinit = element.dataset.mapReady === 'true' && mode === 'flight'
    if (element.dataset.mapReady === 'true' && !needsReinit) return

    if (needsReinit && element._skyvisorMap) {
      element._skyvisorMap.remove()
      element._skyvisorMap = null
      element.dataset.mapReady = 'false'
    }
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
    element._skyvisorMap = map
    map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'bottom-right')
    map.addControl(new maplibregl.AttributionControl({ compact: true }))

    map.on('load', () => {
      if (mode === 'flight') flightLayers(map, element)
      if (mode === 'fleet') fleetLayers(map, element)
      if (mode === 'globe') {
        globeLayers(map)
        spinGlobe(map, element)
      }
    })
  })
}

function animateCounter(element) {
  const target = Number(element.dataset.motionCounter || element.textContent.trim())
  if (!Number.isFinite(target)) return
  if (reducedMotion()) {
    element.textContent = String(Math.round(target))
    return
  }
  animate(0, target, {
    duration: 0.75,
    ease: 'easeOut',
    onUpdate: (value) => {
      element.textContent = String(Math.round(value))
    },
  })
}

function animateProgressBars(root) {
  root.querySelectorAll('[data-motion-progress]').forEach((element) => {
    const target = Number(element.dataset.motionProgress || 0)
    const bar = element.querySelector('.flight-progress-bar')
    const plane = element.querySelector('.flight-progress-plane')
    if (!bar || !plane) return
    if (reducedMotion()) {
      bar.style.width = `${target}%`
      plane.style.left = `calc(${target}% - 5px)`
      return
    }
    animate(0, target, {
      duration: 0.9,
      ease: 'easeOut',
      onUpdate: (value) => {
        bar.style.width = `${value}%`
        plane.style.left = `calc(${value}% - 5px)`
      },
    })
  })
}

function animateEnter(root) {
  if (reducedMotion()) return
  const items = root.querySelectorAll('[data-motion-enter]')
  if (!items.length) return
  animate(
    items,
    { opacity: [0, 1], y: [14, 0] },
    { duration: 0.45, delay: stagger(0.06), ease: 'easeOut' },
  )
}

function pulseRiskMetrics(root) {
  root.querySelectorAll('[data-motion-pulse]').forEach((element) => {
    if (reducedMotion()) return
    animate(element, { scale: [1, 1.04, 1] }, { duration: 1.2, ease: 'easeInOut' })
  })
}

function initMotion(root = document) {
  root.querySelectorAll('[data-motion-counter]').forEach((element) => {
    if (element.dataset.motionCounterReady === 'true') return
    element.dataset.motionCounterReady = 'true'
    animateCounter(element)
  })
  animateEnter(root)
  animateProgressBars(root)
  pulseRiskMetrics(root)
}

function initFleetFilmstrip(root = document) {
  root.querySelectorAll('[data-fleet-filmstrip]').forEach((strip) => {
    if (strip.dataset.fleetFilmstripReady === 'true') return
    strip.dataset.fleetFilmstripReady = 'true'
    const mapEl = strip.parentElement?.querySelector('[data-map-mode="fleet"]')
    strip.querySelectorAll('[data-fleet-focus]').forEach((button) => {
      button.addEventListener('click', () => {
        strip.querySelectorAll('[data-fleet-focus]').forEach((item) => {
          item.classList.remove('border-primary/50', 'bg-primary/5')
        })
        button.classList.add('border-primary/50', 'bg-primary/5')
        const iata = button.dataset.fleetIata
        const coord = iata && AIRPORT_COORDS[iata]
        const map = mapEl?._skyvisorMap
        if (coord && map) {
          map.flyTo({ center: coord, zoom: 5.5, duration: reducedMotion() ? 0 : 900 })
        }
      })
    })
  })
}

function initLiveRefresh(root = document) {
  if (typeof EventSource === 'undefined') return
  root.querySelectorAll('[data-live-refresh-url]').forEach((element) => {
    if (element.dataset.liveRefreshReady === 'true') return
    element.dataset.liveRefreshReady = 'true'
    const source = new EventSource('/events')
    element._skyvisorEventSource = source
    let refreshTimer
    const refresh = () => {
      window.clearTimeout(refreshTimer)
      refreshTimer = window.setTimeout(() => {
        if (!document.documentElement.contains(element)) return
        htmx.ajax('GET', element.dataset.liveRefreshUrl, {
          target: element.dataset.liveRefreshTarget,
          swap: 'outerHTML',
        })
      }, 450)
    }
    ;['flight.updated', 'flight.delayed', 'flight.cancelled', 'gate.changed'].forEach((name) => {
      source.addEventListener(name, refresh)
    })
  })
}

window.Alpine = Alpine
window.maplibregl = maplibregl
window.SkyVisor = { applyTheme, initMaps, initLiveRefresh, initMotion, initFleetFilmstrip }

Alpine.start()

function boot() {
  initMaps()
  initLiveRefresh()
  initMotion()
  initFleetFilmstrip()
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', boot)
} else {
  boot()
}

document.body.addEventListener('htmx:afterSwap', (event) => {
  initMaps(event.detail.target)
  initLiveRefresh(event.detail.target)
  initMotion(event.detail.target)
  initFleetFilmstrip(event.detail.target)
})
document.body.addEventListener('htmx:beforeCleanupElement', (event) => {
  const source = event.detail.elt && event.detail.elt._skyvisorEventSource
  if (source) source.close()
})

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
  if ((localStorage.getItem(themeKey) || 'system') === 'system') applyTheme('system')
})
