const cacheName = "app-" + "0814491b2c7fe69babb1d86effba45689dbc8602";

self.addEventListener("install", event => {
  console.log("installing app worker 0814491b2c7fe69babb1d86effba45689dbc8602");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "/trendycalculator",
        "/trendycalculator/app.css",
        "/trendycalculator/app.js",
        "/trendycalculator/manifest.webmanifest",
        "/trendycalculator/wasm_exec.js",
        "/trendycalculator/web/app.css",
        "/trendycalculator/web/app.wasm",
        "https://storage.googleapis.com/murlok-github/icon-192.png",
        "https://storage.googleapis.com/murlok-github/icon-512.png",
        
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker 0814491b2c7fe69babb1d86effba45689dbc8602 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
