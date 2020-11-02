// Extracts { error, access_token, expires_in and state } from a string. If
// error isn't blank, it will be thrown in an Error. Otherwise, it'll return
// an object containing the other values.
const parse = (hash) => {
  const params = Object.fromEntries(new URLSearchParams(hash))
  const { error, access_token, expires_in, state, token_type } = params
  if (error) throw new Error(error)
  return { access_token, expires_in, state, token_type }
}

// Returns sessionStorage.accessToken if it exists, and either
// sessionStorage.accessTokenExpiresAt doesn't exist, or represents some date
// in the future.
const read = ({ sessionStorage }) => {
  const { accessToken, accessTokenExpiresAt } = sessionStorage
  if (!accessTokenExpiresAt) return accessToken
  if (accessTokenExpiresAt > Date.now()) return accessToken
}

// Generates a URL from the given endpoint and additional parameters. Each
// parameter (except state) may be a 0-arity function, in which case it'll be
// called.  endpoint, client_id, request_type and redirect_uri are required.
// scope may be an array, in which case it'll be joined together with commas.
const url = ({ endpoint, client_id, request_type, redirect_uri, scope, state }) => {
  if (typeof(endpoint) === 'function') endpoint = endpoint()
  if (typeof(client_id) === 'function') client_id = client_id()
  if (typeof(request_type) === 'function') request_type = request_type()
  if (typeof(redirect_uri) === 'function') redirect_uri = redirect_uri()
  if (!endpoint) throw new Error('endpoint required')
  if (!client_id) throw new Error('client_id required')
  if (!request_type) throw new Error('request_type required')
  if (!redirect_uri) throw new Error('redirect_uri required')
  if (typeof(scope) === 'function') scope = scope()
  if (scope instanceof Array) scope = scope.join(',')

  const url = new URL(endpoint)
  url.searchParams.set('client_id', client_id)
  url.searchParams.set('request_type', request_type)
  url.searchParams.set('redirect_uri', new URL(redirect_uri))
  if (scope) url.searchParams.set('scope', scope)
  if (state) url.searchParams.set('state', state)
  return url
}

// Redirects to a url built from the given parameters, by first setting the
// accessTokenState in sessionStorage to state (which may be a function), and
// thencalling location.replace with a url generated from the parameters. It
// will also delete any accessToken (and accessTokenExpiresAt) from
// sessionStorage beforehand.
const redirect = ({ endpoint, client_id, request_type, redirect_uri, scope, state, sessionStorage, location }) => {
  if (typeof(state) === 'function') state = state()
  sessionStorage.accessTokenState = state
  const u = url({ endpoint, client_id, request_type, redirect_uri, scope, state })
  delete sessionStorage.accessTokenExpiresAt
  delete sessionStorage.accessToken
  location.replace(u)
}

// If passing location.hash contains an access_token, it'll be stored in
// sessionStorage.accessToken. If there's also an expires_in, it'll be used to
// calculate sessionStorage.accessTokenExpiresAt. If there's a state, it's
// checked to make sure it matches sessionStorage.accessTokenState (throwing
// an error if it doesn't match), which is then removed. If there's a
// token_type, then that's stored as sessionStorage.accessTokenType.
const set = ({ location, sessionStorage }) => {
  const { access_token, expires_in, state, token_type } = parse(location.hash.slice(1))
  if (state && state !== sessionStorage.accessTokenState)
    throw new Error('state mismatch')
  delete sessionStorage.accessTokenState
  delete sessionStorage.accessToken
  delete sessionStorage.accessTokenExpiresAt
  delete sessionStorage.accessTokenType
  if (!access_token) return
  sessionStorage.accessToken = access_token
  if (expires_in)
    sessionStorage.accessTokenExpiresAt = new Date(Date.now() + expires_in * 1000).valueOf()
  if (token_type)
    sessionStorage.accessTokenType = token_type
  location.hash = ''
  return access_token
}

const oauth = ({ endpoint, client_id, request_type, redirect_uri, scope, state }) => {
  // 4 attributes are required
  if (!endpoint)      throw new TypeError(`endpoint required`)
  if (!client_id)     throw new TypeError(`client_id required`)
  if (!request_type)  throw new TypeError(`request_type required`)
  if (request_type != 'token' && request_type != 'code')
    throw new TypeError(`invalid request_type ${request_type}`)

  if (!state) state = () => btoa(Math.random()).slice(6, 17)

  return ({ location, sessionStorage }) => {
    if (typeof(location) === 'function') location = location()
    if (typeof(sessionStorage) === 'function') sessionStorage = sessionStorage()

    if (!'hash' in location)
      throw new TypeError('invalid location.hash')
    if (typeof(location.replace) != 'function')
      throw new TypeError('invalid location.replace')

    const ctx = Object.assign(
      { endpoint, client_id, request_type, redirect_uri, scope, state },
      { location, sessionStorage },
    )
    return () => read(ctx) || set(ctx) || redirect(ctx)
  }
}

export { parse, read, url, redirect, set }

export default oauth
