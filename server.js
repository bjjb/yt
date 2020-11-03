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

const cors = (request, response) => {
  const CORS_METHODS = ['GET', 'OPTIONS']
  const CORS_HEADERS = ['*']
  const { headers: { origin } } = request
  if (origin) response.setHeader('Access-Control-Allow-Origin', origin)
  response.setHeader('Access-Control-Allow-Methods', CORS_METHODS.join(','))
  response.setHeader('Access-Control-Allow-Headers', CORS_HEADERS.join(','))
}

const handler = async (request, response) => {
  const { method, url, headers: { accept } } = request
  cors(request, response)
  if (method === 'OPTIONS') return response.sendStatus(200)

  const videoID = url.slice(1)
  response.setHeader('Content-Type', 'application/json')
  if (!videoID) return response.end(JSON.stringify({ status: 'ok' }))
  get(infoURL(videoID), (resp) => {
    const { statusCode, statusText } = resp
    const chunks = []
    resp.setEncoding('utf8')
    resp.on('data', chunk => chunks.push(chunk))
    resp.on('end', () => {
      const result = Object.assign(parse(chunks), { statusCode, statusText })
      response.end(JSON.stringify(result))
    })
  })
}

const server = createServer(handler)

server.listen(port, (err) => {
  if (err) {
    return console.error(err)
  }
  console.log(`server is listening on ${port}`)
})
