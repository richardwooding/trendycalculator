const cacheName = "app-" + "8a9e281b7d0fde851818ff43d7b71637f91b3b20";

self.addEventListener("install", event => {
  console.log("installing app worker 8a9e281b7d0fde851818ff43d7b71637f91b3b20");
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
  console.log("app worker 8a9e281b7d0fde851818ff43d7b71637f91b3b20 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
