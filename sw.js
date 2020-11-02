addEventListener('install', (event) => {
  console.debug('install event=%o, this=%o', event, this)
})

addEventListener('fetch', (event) => {
  console.debug('fetch: event=%o, this=%o', event, this)
})
