'use strict'

import Feadán from './feadán.js'

self.Feadán = Feadán

console.log(`▶️ Welcome to the Feadán app, developer!

I've added a new class (Feadán → %o) to the window, which you can use to poke
around. You can instantiate a new instance, and use it (for example) to log
into YouTube and search for videos, or to see which playlists you have set up
in your local database.

There's also a test page (/test.html) which describes the behaviour of this
application when running on a sufficiently modern browser.
`, Feadán)

addEventListener('load', (event) => {
  console.debug('Loaded: %o', event)
})

customElements.define('fadawn-playlist', Feadán.Playlist, { extends: 'ol' })
