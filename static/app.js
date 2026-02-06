const POLL_INTERVAL = 1000;

let rangeSeconds = 60;
function getRangeMs() {
  return rangeSeconds * 1000; // миллисекунды
}
let history = [];

function now() {
  return Date.now();
}

async function fetchStatus() {
  try {
    const res = await fetch("/api/v1/stats");
    if (!res.ok) return;

    const s = await res.json();

    history.push({
      t: now(),
      inV: s.InputVoltage,
      outV: s.OutputVoltage,
      load: s.OutputCurrentPct,
      bat: s.BatteryVoltage,
      temp: s.Temperature,
    });

    trimHistory();
    updateCharts();
  } catch {}
}

function trimHistory() {
  const cutoff = now() - getRangeMs();
  history = history.filter((p) => p.t >= cutoff);
}

function makeChart(ctx, label, color, getValue) {
  const chart = new Chart(ctx, {
    type: "line",
    data: {
      datasets: [
        {
          label,
          data: [],
          borderColor: color,
          tension: 0.2,
          pointRadius: 0,
        },
      ],
    },
    options: {
      animation: false,
      responsive: true,
      scales: {
        x: {
          type: "time",
          time: {
            unit: "second",
            displayFormats: { second: "HH:mm:ss" },
            tooltipFormat: "HH:mm:ss",
          },
          ticks: {
            autoSkip: true,
            maxTicksLimit: 6,
          },
          grid: {
            color: "rgba(255,255,255,0.1)", // лёгкая сетка
          },
        },
        y: {
          beginAtZero: false,
          grid: {
            color: "rgba(255,255,255,0.1)", // лёгкая сетка
          },
        },
      },
    },
  });

  chart.getValue = getValue;
  return chart;
}

const charts = {
  voltage: makeChart(
    document.getElementById("voltageChart"),
    "Voltage(V)",
    "#38bdf8",
    (p) => p.outV,
  ),

  load: makeChart(
    document.getElementById("loadChart"),
    "Load (%)",
    "#facc15",
    (p) => p.load,
  ),
  battery: makeChart(
    document.getElementById("batteryChart"),
    "Battery (V)",
    "#4ade80",
    (p) => p.bat,
  ),
  temp: makeChart(
    document.getElementById("tempChart"),
    "Temperature (°C)",
    "#fb7185",
    (p) => p.temp,
  ),
};

function updateCharts() {
  const minX = now() - getRangeMs();
  const maxX = now();

  Object.values(charts).forEach((chart) => {
    const data = history.map((p) => ({
      x: p.t,
      y: chart.getValue(p),
    }));

    chart.data.datasets[0].data = data;

    // Явный диапазон по X
    chart.options.scales.x.min = minX;
    chart.options.scales.x.max = maxX;

    // Авто Y: min/max ±5% для плавного вида
    if (data.length > 0) {
      const ys = data.map((p) => p.y);
      const minY = Math.min(...ys);
      const maxY = Math.max(...ys);
      const padding = (maxY - minY) * 0.05 || 1; // если разница 0
      chart.options.scales.y.min = minY - padding;
      chart.options.scales.y.max = maxY + padding;
    }

    chart.update();
  });
}

// Range buttons
document.querySelectorAll("button[data-range]").forEach((btn) => {
  btn.onclick = () => {
    document
      .querySelectorAll("button")
      .forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");

    rangeSeconds = Number(btn.dataset.range);
    trimHistory();
    updateCharts();
  };
});

// Polling
setInterval(fetchStatus, POLL_INTERVAL);
fetchStatus();
