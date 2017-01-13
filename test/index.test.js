/* global describe, it, beforeEach, afterEach */
'use strict'

const request = require('supertest')

describe('loading express', () => {
  var server
  beforeEach(() => {
    server = require('../')
  })
  afterEach(() => server.close())

  it('GET / is a 404', () =>
    request(server)
      .get('/')
      .expect(404)
  )

  it('404 everything else', () =>
    request(server)
      .get('/foo/bar')
      .expect(404)
  )

  it('POST gets the request body straight back as the response body', () =>
    request(server)
      .post('/')
      .send({hello: 'world'})
      .expect(200, {hello: 'world'})
  )
})
