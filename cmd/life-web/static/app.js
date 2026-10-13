const canvas = document.querySelector('#canvas');
const context = canvas.getContext('2d');
const fileInput = document.querySelector('#file');
const playButton = document.querySelector('#play');
const speedInput = document.querySelector('#speed');
const speedValue = document.querySelector('#speed-value');
const wrapInput = document.querySelector('#wrap');
const status = document.querySelector('#status');

let grid = null;
let generation = 0;
let playing = false;
let timer = null;

function showError(message) {
  status.textContent = `Error: ${message}`;
  status.classList.add('error');
}

function showGeneration() {
  status.textContent = `Generation ${generation}`;
}

function draw() {
  if (!grid) return;
  const rows = grid.length;
  const columns = grid[0].length;
  const bounds = canvas.getBoundingClientRect();
  const cell = Math.max(1, Math.floor(Math.min(bounds.width / columns, 600 / rows)));
  canvas.width = cell * columns;
  canvas.height = cell * rows;
  context.fillStyle = '#030712';
  context.fillRect(0, 0, canvas.width, canvas.height);
  context.fillStyle = '#67e8f9';
  for (let row = 0; row < rows; row++) {
    for (let column = 0; column < columns; column++) {
      if (grid[row][column] === 1) context.fillRect(column * cell, row * cell, cell, cell);
    }
  }
  if (cell >= 4) {
    context.strokeStyle = '#1f2937';
    context.lineWidth = 1;
    for (let row = 0; row <= rows; row++) context.beginPath(), context.moveTo(0, row * cell + .5), context.lineTo(canvas.width, row * cell + .5), context.stroke();
    for (let column = 0; column <= columns; column++) context.beginPath(), context.moveTo(column * cell + .5, 0), context.lineTo(column * cell + .5, canvas.height), context.stroke();
  }
}

async function parse(text) {
  const response = await fetch('/api/parse', { method: 'POST', headers: { 'Content-Type': 'text/plain' }, body: text });
  const data = await response.json();
  if (!response.ok) throw new Error(data.error || 'invalid grid');
  grid = data.grid;
  generation = 0;
  showGeneration();
  draw();
}

async function step() {
  const response = await fetch('/api/step', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ grid, wrap: wrapInput.checked }) });
  const data = await response.json();
  if (!response.ok) throw new Error(data.error || 'step failed');
  grid = data.grid;
  generation++;
  showGeneration();
  draw();
}

function setPlaying(value) {
  playing = value;
  playButton.textContent = playing ? 'Pause' : 'Play';
  if (timer) clearInterval(timer);
  timer = playing ? setInterval(() => step().catch(error => { setPlaying(false); showError(error.message); }), Number(speedInput.value)) : null;
}

fileInput.addEventListener('change', async () => {
  const file = fileInput.files[0];
  if (!file) return;
  try { setPlaying(false); await parse(await file.text()); }
  catch (error) { showError(error.message); }
});
playButton.addEventListener('click', () => { if (grid) setPlaying(!playing); });
speedInput.addEventListener('input', () => { speedValue.textContent = `${speedInput.value} ms`; if (playing) setPlaying(true); });
window.addEventListener('resize', draw);

fetch('/sample.txt').then(response => response.text()).then(parse).catch(error => showError(error.message));
