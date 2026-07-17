(function () {
  const stored = localStorage.getItem('skyvisor-theme') || 'system'
  const dark = stored === 'dark' || (stored === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.dataset.theme = stored
  document.documentElement.classList.toggle('dark', dark)
})()
