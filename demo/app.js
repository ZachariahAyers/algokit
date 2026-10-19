"use strict";

const algorithms = {
  quick: ["Quick Sort", "Expected O(n log n); three-way partition; in place."],
  merge: ["Merge Sort", "O(n log n); stable; one scratch buffer."],
  heap: ["Heap Sort", "O(n log n); in place; no worst-case surprise."],
  shell: ["Shell Sort", "Ciura gap sequence; in place."],
  insertion: ["Insertion Sort", "O(n²), or O(n) when already ordered; stable."],
  bubble: ["Bubble Sort", "O(n²), with an early stop on ordered input; stable."],
  selection: ["Selection Sort", "O(n²); at most n−1 swaps."],
  counting: ["Counting Sort", "O(n + k); integer range is bounded by this lab."],
  "linear-search": ["Linear Search", "O(n); scans the sequence as entered."],
  "binary-search": ["Binary Search", "O(log n); runs on an ascending copy."],
  "jump-search": ["Jump Search", "O(√n); runs on an ascending copy."],
  "interpolation-search": ["Interpolation Search", "Expected O(log log n) for uniform keys; runs on an ascending copy."],
};

const samples = {
  reverse: "96, 84, 72, 60, 48, 36, 24, 12, 0, -12, -24, -36",
  repeats: "4, 2, 4, 1, 2, 4, 1, 1, 2, 4, 2, 1, 4, 2",
  nearly: "-8, -3, 0, 4, 9, 14, 12, 20, 27, 31, 35, 42",
};

const form = document.querySelector("#algorithm-form");
const algorithm = document.querySelector("#algorithm");
const algorithmNote = document.querySelector("#algorithm-note");
const valuesInput = document.querySelector("#values");
const targetField = document.querySelector("#target-field");
const targetInput = document.querySelector("#target");
const runButton = document.querySelector("#run-button");
const status = document.querySelector("#status");
const resultTitle = document.querySelector("#result-title");
const contract = document.querySelector("#contract");
const sequence = document.querySelector("#sequence");
const metrics = document.querySelector("#metrics");
const resultDetail = document.querySelector("#result-detail");

function isSearch() {
  return algorithm.value.endsWith("-search");
}

function updateAlgorithm() {
  algorithmNote.textContent = algorithms[algorithm.value][1];
  targetField.hidden = !isSearch();
}

function parseValues() {
  const tokens = valuesInput.value.trim().split(/[\s,]+/).filter(Boolean);
  if (tokens.length === 0) throw new Error("Enter at least one integer.");
  const values = tokens.map((token) => Number(token));
  if (values.some((value) => !Number.isSafeInteger(value))) {
    throw new Error("Every value must be a whole number.");
  }
  return values;
}

function formatDuration(nanoseconds) {
  if (nanoseconds === 0) return "< clock resolution";
  if (nanoseconds < 1_000) return `${nanoseconds} ns`;
  if (nanoseconds < 1_000_000) return `${(nanoseconds / 1_000).toFixed(2)} µs`;
  return `${(nanoseconds / 1_000_000).toFixed(2)} ms`;
}

function paintSequence(result) {
  sequence.replaceChildren();
  const minimum = Math.min(...result.values);
  const maximum = Math.max(...result.values);
  const spread = Math.max(1, maximum - minimum);
  result.values.forEach((value, index) => {
    const cell = document.createElement("span");
    cell.className = "value-cell";
    const strength = (value - minimum) / spread;
    cell.style.backgroundColor = `hsl(142 32% ${91 - strength * 32}%)`;
    cell.textContent = value;
    cell.title = `Index ${index}: ${value}`;
    if (result.search?.found && index === result.search.index) cell.classList.add("match");
    if (result.search && !result.search.found && index === result.search.index) cell.classList.add("insertion");
    sequence.append(cell);
  });
  if (result.search && !result.search.found && result.search.index === result.values.length) {
    const marker = document.createElement("span");
    marker.className = "value-cell insertion";
    marker.textContent = "end";
    marker.title = `Insertion point ${result.search.index}`;
    sequence.append(marker);
  }
}

function showResult(result) {
  const [name] = algorithms[result.algorithm];
  resultTitle.textContent = name;
  contract.textContent = result.contractVerified ? "Verified Result" : "Verification Failed";
  contract.classList.toggle("verified", result.contractVerified);
  paintSequence(result);

  const metricValues = metrics.querySelectorAll("dd");
  metricValues[0].textContent = result.values.length;
  metricValues[1].textContent = formatDuration(result.elapsedNanoseconds);
  metricValues[2].textContent = result.preparedAscending ? "Ascending Copy" : "As Entered";

  if (result.search) {
    resultDetail.textContent = result.search.found
      ? `Found ${result.search.target} at index ${result.search.index}. The outlined value is the first match.`
      : `${result.search.target} was not found. Its insertion point is index ${result.search.index}.`;
  } else {
    resultDetail.textContent = "The sequence above is the slice after the selected Go sort returned.";
  }
}

function completionMessage(result) {
  const [name] = algorithms[result.algorithm];
  if (!result.search) return `${name} complete. ${result.values.length} items sorted.`;
  if (result.search.found) {
    return `${name} complete. Found ${result.search.target} at index ${result.search.index}.`;
  }
  return `${name} complete. ${result.search.target} was not found; insertion point ${result.search.index}.`;
}

function submitRun(event) {
  event.preventDefault();
  status.className = "status";
  status.textContent = "";
  try {
    const request = { algorithm: algorithm.value, values: parseValues() };
    if (isSearch()) {
      request.target = Number(targetInput.value);
      if (!Number.isSafeInteger(request.target)) throw new Error("The target must be a whole number.");
    }
    const response = JSON.parse(globalThis.algokitRun(JSON.stringify(request)));
    if (response.error) throw new Error(response.error);
    showResult(response.result);
    status.textContent = completionMessage(response.result);
  } catch (error) {
    status.className = "status error";
    status.textContent = error instanceof Error ? error.message : "The run could not be completed.";
  }
}

algorithm.addEventListener("change", updateAlgorithm);
form.addEventListener("submit", submitRun);
document.querySelectorAll("[data-sample]").forEach((button) => {
  button.addEventListener("click", () => {
    valuesInput.value = samples[button.dataset.sample];
    valuesInput.focus();
  });
});

updateAlgorithm();
globalThis.algokitRuntime
  .start()
  .then(() => {
    runButton.disabled = false;
    runButton.textContent = "Run in Go";
    contract.textContent = "Go Runtime Ready";
    form.requestSubmit();
  })
  .catch((error) => {
    status.className = "status error";
    status.textContent = error instanceof Error ? error.message : "The Go runtime could not start.";
    runButton.textContent = "Go Runtime Unavailable";
    contract.textContent = "Runtime Unavailable";
  });
