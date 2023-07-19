if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js', { scope: '/' }).then(() => {
    console.log('Service Worker registered successfully.');
  }).catch(error => {
    console.log('Service Worker registration failed:', error);
  });
}

const cacheName = 'librarian-v0.10.0-beta';
const files = [
  '/static/css/channel.css',
  '/static/css/frontpage.css',
  '/static/css/plyr.css',
  '/static/css/search.css',
  '/static/css/video.css',
  '/static/img/librarian.svg',
  '/static/img/plyr.svg',
  '/static/favicon/android-chrome-192x192.png',
  '/static/favicon/android-chrome-512x512.png',
  '/static/favicon/apple-touch-icon.png',
  '/static/favicon/favicon-16x16.png',
  '/static/favicon/favicon-32x32.png',
  '/static/favicon/favicon.ico',
  '/static/favicon/mstile-70x70.png',
  '/static/favicon/mstile-144x144.png',
  '/static/favicon/mstile-150x150.png',
  '/static/favicon/mstile-310x150.png',
  '/static/favicon/mstile-310x310.png',
  '/static/favicon/safari-pinned-tab.svg',
  '/static/fonts/Material-Icons-Outlined.css',
  '/static/fonts/Material-Icons-Outlined.woff2',
  '/static/js/plyr.js',
  '/static/blank.mp4'
];

self.addEventListener('install', e => {
  e.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll(files);
    })
  );
});

self.addEventListener('fetch', event => {
  if (event.request.method === 'GET') {
    let url = event.request.url.indexOf(self.location.origin) !== -1 ?
      event.request.url.split(`${self.location.origin}/`)[1] :
      event.request.url;
    let isFileCached = files.indexOf(url) !== -1;

    if (isFileCached) {
      event.respondWith(
        fetch(event.request).catch(err =>
          self.cache.open(cache_name).then(cache => cache.match("/offline.html"))
        )
      );
    }
  }
});