import 'https://unpkg.com/mocha@8.2.0/mocha.js'
import 'https://unpkg.com/chai@4.2.0/chai.js'
import { fake } from 'https://unpkg.com/sinon@9.2.1/pkg/sinon-esm.js'
import 'https://unpkg.com/sinon-chai@3.5.0/lib/sinon-chai.js'
import oauth, { url } from './oauth.js'
import youtube from './youtube.js'

mocha.setup('bdd')
addEventListener('load', () => mocha.run())

const { assert, expect } = chai

const endpoint = 'https://example.com/auth'
const request_type = 'token'
const client_id = 'CLIENTID';
const redirect_uri = 'https://localhost/callback'

describe(`oauth({ endpoint, request_type, client_id, redirect_uri })`, () => {
  const f = oauth({ endpoint, request_type, client_id, redirect_uri });

  ['endpoint', 'request_type', 'client_id'].forEach((missing) => {
    it(`throws an error if the ${missing} is missing`, () => {
      let options = { endpoint, request_type, client_id }
      delete options[missing]
      expect(() => oauth(options)).to.throw(`${missing} required`)
    })
  })

  it('returns a function, f', () => {
    expect(f).to.be.a('function')
  });

  describe('f({ sessionStorage: empty, location: empty })', () => {
    const sessionStorage = {}
    const hash = ''
    const replace = fake()
    const location = { hash, replace }
    const redirect = url({ endpoint, client_id, request_type, redirect_uri })
    it('redirects to the auth url', () => {
      f({ sessionStorage, location })()
      expect(replace).to.have.been.calledWith(redirect)
    })
  })
  describe('f({ sessionStorage: expired, location: empty })', () => {
    const sessionStorage = { accessToken: 'expired', accessTokenExpiresAt: Date.now() }
    const hash = ''
    const replace = fake()
    const location = { hash, replace }
    const redirect = url({ endpoint, client_id, request_type, redirect_uri })
    it('redirects to the auth url', () => {
      f({ sessionStorage, location })()
      expect(replace).to.have.been.calledWith(redirect)
    })
  })
  describe('f({ sessionStorage: valid, location: empty })', () => {
    const sessionStorage = { accessToken: 'valid', accessTokenExpiresAt: Date.now() + 1000 * 3600 }
    const hash = ''
    const replace = fake()
    const location = { hash, replace }
    it('returns the token', () => {
      expect(f({ sessionStorage, location })()).to.eq('valid')
      expect(replace).not.to.have.been.called
    })
  })
  describe('f({ sessionStorage: empty, location: error })', () => {
    const sessionStorage = {}
    const hash = '#error=foo'
    const replace = fake()
    const location = { hash, replace }
    it('throws an error (and does not redirect)', () => {
      expect(f({ sessionStorage, location })).to.throw(/foo/)
      expect(replace).not.to.have.been.called
    })
  })
  describe('f({ sessionStorage: empty, location: success })', () => {
    const sessionStorage = {}
    const hash = '#access_token=foo&expires_in=3600&state=STATE'
    const replace = fake()
    const location = { hash, replace }
    it('throws an error because of mismatched state', () => {
      expect(f({ sessionStorage, location })).to.throw(/state mismatch/)
      expect(replace).not.to.have.been.called
    })
  })
  describe('f({ sessionStorage: pending, location: success })', () => {
    const now = Date.now()
    const sessionStorage = { accessTokenState: 'STATE' }
    const hash = '#access_token=TOKEN&expires_in=3600&state=STATE'
    const replace = fake()
    const location = { hash, replace }
    it('sets sessionStorage, clears location, and returns the token', () => {
      expect(f({ sessionStorage, location })()).to.eq('TOKEN')
      expect(sessionStorage.accessToken).to.eq('TOKEN')
      expect(sessionStorage.accessTokenExpiresAt).to.be.at.least(now + 3600 * 1000)
      expect(sessionStorage.accessTokenState).to.be.undefined
      expect(location.hash).to.eq('')
      expect(replace).not.to.have.been.called
    })
  })
})

describe(`youtube({ token })`, () => {
  const token = () => 'TOKEN'
  const fetch = fake()
  const o = youtube({ token, fetch })
  describe('channels', () => {
    it('returns an iterable list of channels', () => {
      expect(o.channels()).to.be.instanceof(Array)
    })
  })
})
