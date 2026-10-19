"use strict";

const assert = require("node:assert/strict");
const { test } = require("node:test");

require("./runtime.js");

test("startup completes only after the Go adapter announces readiness", { timeout: 1_000 }, async () => {
  let ran = false;
  await globalThis.algokitRuntime.start({
    makeGo: () => ({
      importObject: {},
      run: () => {
        ran = true;
        queueMicrotask(() => globalThis.algokitReady());
        return new Promise(() => {});
      },
    }),
    load: async () => ({}),
  });
  assert.equal(ran, true);
  assert.equal(globalThis.algokitReady, undefined);
});

test("startup rejects when go.run fails before readiness", { timeout: 1_000 }, async () => {
  await assert.rejects(
    globalThis.algokitRuntime.start({
      makeGo: () => ({
        importObject: {},
        run: async () => {
          throw new Error("runtime failed");
        },
      }),
      load: async () => ({}),
    }),
    /runtime failed/,
  );
  assert.equal(globalThis.algokitReady, undefined);
});

test("startup rejects when the WebAssembly module cannot load", { timeout: 1_000 }, async () => {
  await assert.rejects(
    globalThis.algokitRuntime.start({
      makeGo: () => ({ importObject: {} }),
      load: async () => {
        throw new Error("download failed");
      },
    }),
    /download failed/,
  );
});
