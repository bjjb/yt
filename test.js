import 'https://unpkg.com/mocha@8.2.0/mocha.js'
import 'https://unpkg.com/chai@4.2.0/chai.js'
import 'https://unpkg.com/sinon-chai@3.5.0/lib/sinon-chai.js'
import { fake } from 'https://unpkg.com/sinon@9.2.1/pkg/sinon-esm.js'
import oauth, { url } from './oauth.js'
import youtube from './youtube.js'

mocha.setup('bdd')
addEventListener('load', () => mocha.run())

const { assert, expect } = chai

const options = {
  endpoint:     () => 'https://example.com/auth',
  request_type: () => 'token',
  client_id:    () => 'CLIENTID',
  redirect_uri: () => 'https://localhost/callback',
}

describe(`oauth`, () => {
  it('is a function', () => {
    expect(oauth).to.be.a('function')
  })

  describe('when the session is empty', () => {
    const sessionStorage = {}
    const location = { hash: '', replace: fake() }
    it('redirects to the auth url', () => {
      const state = 'blubber'
      oauth({ sessionStorage, location, state, ...options })
      expect(location.replace).to.have.been.calledWith(url(options))
      expect(sessionStorage).to.deep.eq({ accessTokenState: state })
    })
  })
  describe('when the session is expired', () => {
    const sessionStorage = { accessToken: 'expired', accessTokenExpiresAt: Date.now() }
    const location = { hash: '', replace: fake() }
    const state = () => 'blubber'
    it('redirects to the auth url', () => {
      oauth({ sessionStorage, location, state, ...options })
      expect(location.replace).to.have.been.calledWith(url({ state, ...options }))
      expect(sessionStorage).to.deep.eq({ accessTokenState: state() })
    })
  })
  describe('when the session is valid', () => {
    const sessionStorage = { accessToken: 'valid', accessTokenExpiresAt: Date.now() + 9999 }
    const location = { hash: '', replace: fake() }
    it('returns the token', () => {
      expect(oauth({ sessionStorage, location, ...options })).to.eq('valid')
      expect(location.replace).not.to.have.been.called
    })
  })
  describe('when the location indicates an error', () => {
    const sessionStorage = {}
    const location = { hash: '#error=ERROR', replace: fake() }
    it('throws an error (and does not redirect)', () => {
      expect(() => oauth({ sessionStorage, location, ...options })).to.throw(/ERROR/)
      expect(location.replace).not.to.have.been.called
    })
  })
  describe(`when the location state doesn't match the session state`, () => {
    const sessionStorage = { accessTokenState: 'A' }
    const location = { hash: '#access_token=foo&state=B', replace: fake() }
    it('throws an error because of mismatched state', () => {
      expect(() => oauth({ sessionStorage, location, ...options })).to.throw(/state mismatch/)
      expect(location.replace).not.to.have.been.called
    })
  })
  describe('when the location indicates a new token', () => {
    const now = Date.now()
    const sessionStorage = { accessTokenState: 'X' }
    const location = { hash: '#access_token=TOKEN&expires_in=1&state=X', replace: fake() }
    it('returns the token, having stored it and cleared the location', () => {
      expect(oauth({ sessionStorage, location, ...options })).to.eq('TOKEN')
      expect(sessionStorage.accessToken).to.eq('TOKEN')
      expect(sessionStorage.accessTokenExpiresAt).to.be.at.least(now + 1000)
      expect(sessionStorage.accessTokenState).to.be.undefined
      expect(location.hash).to.eq('')
      expect(location.replace).not.to.have.been.called
    })
  })
})

describe(`youtube`, () => {
  it('is a function', () => {
    expect(youtube).to.be.a('function')
  })
  describe('youtube({ token, fetch })', () => {
    const token = fake.returns('TOKEN')
    const results = [{}]
    const fetch = fake.resolves({ json: fake.resolves(JSON.stringify(results)) })
    const yt = youtube({ token, fetch })
    describe('.search({ q })', () => {
      it('calls token', () => {
        youtube({ token, fetch })
        expect(token).to.have.been.called
      })
    })
  })
})
