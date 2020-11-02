'use strict'

import youtube from './youtube.js'

const client_id = '990154653919-271dvn0pa9v2nb8dtud4ecamqe43b8pa.apps.googleusercontent.com'
const endpoint = 'https://accounts.google.com/o/oauth2/v2/auth'
const request_type = 'token'
const redirect_uri = location
const scopes = [ '' ]
const token = oauth({ client_id, endpoint, request_type, redirect_uri, scope, sessionStorage, location })

console.log(`▶️ Welcome to the yt.js app, developer!
This app makes use of the yt.js script to provide a simple API to YouTube.

There's a test page (/test.html) which describes the behaviour of this
application when running on a sufficiently modern browser.
`)

addEventListener('load', (event) => {
})
