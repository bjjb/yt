import 'https://unpkg.com/mocha@8.2.0/mocha.js'
import 'https://unpkg.com/chai@4.2.0/chai.js'
import 'https://unpkg.com/sinon-chai@3.5.0/lib/sinon-chai.js'
import { fake } from 'https://unpkg.com/sinon@9.2.1/pkg/sinon-esm.js'
import youtube from './youtube.js'
import { url } from './youtube.js'

const { expect } = chai

mocha.setup('bdd')

describe(`youtube`, () => {
  it('is a function which takes a fetch function and returns a simple YouTube client', () => {
    expect(youtube).to.be.a('function')
    expect(youtube({ fetch: () => {} })).to.be.an('object')
  })
  describe('the returned object', () => {
    it('has seach, to search for videos and channels', () => {
      expect(youtube({}).search).to.be.a('function')
    })
    it('has getSubscriptions, to list your subscriptions', () => {
      expect(youtube({}).getSubscriptions).to.be.a('function')
    })
    it('has getChannel, to get the details of a channel', () => {
      expect(youtube({}).getChannel).to.be.a('function')
    })
    it('has getVideo, to get the details of video', () => {
      expect(youtube({}).getVideo).to.be.a('function')
    })
    describe('.search({ q })', () => {
      const q = 'Search Term'
      const results = { items: [ { 'foo': 'bar' } ] }
      const success = { ok: true, status: 200, json: fake.resolves(results) }
      it('calls fetch using the query', async () => {
        const fetch = fake.resolves(success)
        await youtube({ fetch }).search({ q })
        expect(fetch).to.have.been.calledWith(url('search', { q }))
      })
      xit('returns a search results generator', async () => {
        const fetch = fake.resolves(success)
        const items = await youtube({ fetch }).search({ q })
        expect(items.next).to.be.a('function')
        results.items.forEach((item) => {
          const { value, done } = items.next()
          expect(value).to.deep.eq(item)
          expect(done).to.eq(true)
        })
      })
    })
    describe('.getSubscriptions()', () => {
      it('calls fetch')
      it('returns a subscriptions generator')
    })
    describe('.getChannel({ id })', () => {
      it('calls fetch using the id')
      describe('when the channel exists', () => {
        it('returns a channels generator')
      })
      describe('when the channel does not exist', () => {
        it('throws an error')
      })
    })
    describe('.getVideo({ id })', () => {
      it('calls fetch using the id')
      describe('when the video exists', () => {
        it('resolves to a video')
      })
      describe('when the video does not exist', () => {
        it('throws an error')
      })
    })
  })
})
