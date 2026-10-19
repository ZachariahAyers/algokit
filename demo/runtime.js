"use strict";

async function loadGoInstance(importObject) {
  try {
    const { instance } = await WebAssembly.instantiateStreaming(fetch("algokit.wasm"), importObject);
    return instance;
  } catch {
    const response = await fetch("algokit.wasm");
    if (!response.ok) throw new Error(`WebAssembly download failed (${response.status}).`);
    const { instance } = await WebAssembly.instantiate(await response.arrayBuffer(), importObject);
    return instance;
  }
}

async function startGoRuntime({ makeGo = () => new Go(), load = loadGoInstance } = {}) {
  const go = makeGo();
  const instance = await load(go.importObject);
  let announceReady;
  const ready = new Promise((resolve) => {
    announceReady = resolve;
    globalThis.algokitReady = announceReady;
  });
  const stopped = Promise.resolve()
    .then(() => go.run(instance))
    .then(() => {
      throw new Error("The Go runtime stopped before it was ready.");
    });
  try {
    await Promise.race([ready, stopped]);
  } finally {
    if (globalThis.algokitReady === announceReady) delete globalThis.algokitReady;
  }
}

globalThis.algokitRuntime = Object.freeze({ start: startGoRuntime });
