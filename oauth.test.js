import 'https://unpkg.com/mocha@8.2.0/mocha.js'
import 'https://unpkg.com/chai@4.2.0/chai.js'
import 'https://unpkg.com/sinon-chai@3.5.0/lib/sinon-chai.js'
import { fake } from 'https://unpkg.com/sinon@9.2.1/pkg/sinon-esm.js'
import { url } from './oauth.js'
import oauth from './oauth.js'

const { expect } = chai

mocha.setup('bdd')

describe(`oauth`, () => {
  const options = {
    endpoint:     () => 'https://example.com/auth',
    response_type: () => 'token',
    client_id:    () => 'CLIENTID',
    redirect_uri: () => 'https://localhost/callback',
  }

  it('is a function, which expects a sessionStorage, a location, and some extra options for building auth URLs', () => {
    expect(oauth).to.be.a('function')
  })

  describe('when the sessionStorage is empty', () => {
    const sessionStorage = {}
    const location = { hash: '', assign: fake() }
    it('calls location.assign with an auth url which is configured with the options', () => {
      const state = 'blubber'
      oauth({ sessionStorage, location, state, ...options })
      expect(location.assign).to.have.been.calledWith(url(options))
      expect(sessionStorage).to.deep.eq({ accessTokenState: state })
    })
  })
  describe('when the sessionStorage.accessTokenExpiresAt is past', () => {
    const sessionStorage = { accessToken: 'expired', accessTokenExpiresAt: Date.now() }
    const location = { hash: '', assign: fake() }
    const state = () => 'blubber'
    it('redirects to the auth url', () => {
      oauth({ sessionStorage, location, state, ...options })
      expect(location.assign).to.have.been.calledWith(url({ state, ...options }))
      expect(sessionStorage).to.deep.eq({ accessTokenState: state() })
    })
  })
  describe('when the session is valid', () => {
    const sessionStorage = { accessToken: 'tOkEn', accessTokenExpiresAt: Date.now() + 9999 }
    const location = { hash: '', assign: fake() }
    it('does not redirect, but returns the token', () => {
      expect(oauth({ sessionStorage, location, ...options })).to.eq('tOkEn')
      expect(location.assign).not.to.have.been.called
    })
  })
  describe('when the location hash indicates an error', () => {
    const sessionStorage = {}
    const location = { hash: '#error=ERROR', assign: fake() }
    it('throws the error (and does not redirect)', () => {
      expect(() => oauth({ sessionStorage, location, ...options })).to.throw(/ERROR/)
      expect(location.assign).not.to.have.been.called
    })
  })
  describe(`when the location state doesn't match the session state`, () => {
    const sessionStorage = { accessTokenState: 'A' }
    const location = { hash: '#access_token=foo&state=B', assign: fake() }
    it('throws an error because of mismatched state', () => {
      expect(() => oauth({ sessionStorage, location, ...options })).to.throw(/state mismatch/)
      expect(location.assign).not.to.have.been.called
    })
  })
  describe('when the location indicates a new token (and the state matches)', () => {
    const now = Date.now()
    const sessionStorage = { accessTokenState: 'X' }
    const location = { hash: '#access_token=TOKEN&expires_in=1&state=X', assign: fake() }
    it('returns the token, having stored it and cleared the location', () => {
      expect(oauth({ sessionStorage, location, ...options })).to.eq('TOKEN')
      expect(sessionStorage.accessToken).to.eq('TOKEN')
      expect(sessionStorage.accessTokenExpiresAt).to.be.at.least(now + 1000)
      expect(sessionStorage.accessTokenState).to.be.undefined
      expect(location.hash).to.eq('')
      expect(location.assign).not.to.have.been.called
    })
  })
})
