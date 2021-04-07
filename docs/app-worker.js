const cacheName = "app-" + "abb49f8ab33cb52d498a4f4ec994208c6eac842f";

self.addEventListener("install", event => {
  console.log("installing app worker abb49f8ab33cb52d498a4f4ec994208c6eac842f");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "/trendycalculator",
        "/trendycalculator/app.css",
        "/trendycalculator/app.js",
        "/trendycalculator/manifest.webmanifest",
        "/trendycalculator/wasm_exec.js",
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
  console.log("app worker abb49f8ab33cb52d498a4f4ec994208c6eac842f is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
