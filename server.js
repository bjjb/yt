// A very simple HTTP service which accepts a video ID in the request path and
// queries the YouTube get_video_info endpoint, returning the result as JSON.
const { createServer } = require('http')
const { get } = require('https')
const port = process.env.PORT || 3000

const parse = (chunks) => {
  return Object.fromEntries(new URLSearchParams(chunks.join('')))
}

const infoURL = (videoID) => {
  const url = new URL('https://www.youtube.com/get_video_info')
  url.searchParams.set('video_id', videoID)
  return url
}

const handler = async (request, response) => {
  const id = request.url.slice(1)
  if (!id) return response.end('What are you looking for?')
  const headers = { Accept: 'x-www-url-form-encoded' }
  get(infoURL(id), { headers }, (resp) => {
    const { statusCode, statusText } = resp
    const chunks = []
    resp.setEncoding('utf8')
    resp.on('data', chunk => chunks.push(chunk))
    resp.on('end', () => {
      response.setHeader('Content-Type', 'application/json')
      const result = Object.assign(parse(chunks), { statusCode, statusText })
      response.end(JSON.stringify(result))
    })
  })
}

const server = createServer(handler)

server.listen(3000, (err) => {
  if (err) {
    return console.error(err)
  }
  console.log(`server is listening on ${port}`)
})
