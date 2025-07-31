// --- Stock Chart and Data ---
const ctx = document.getElementById('stock-graph').getContext('2d');
let stockChart;

async function fetchStockData(symbol = 'RELIANCE.NSE') {
  const res = await fetch(`/api/stock?symbol=${symbol}`);
  if (!res.ok) throw new Error('Failed to fetch stock data');
  return res.json();
}

function renderStockChart(data) {
  const labels = data.map(row => row.date).reverse();
  const prices = data.map(row => row.close).reverse();
  if (stockChart) stockChart.destroy();
  stockChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels,
      datasets: [{
        label: 'Close Price',
        data: prices,
        borderColor: '#ff9800',
        backgroundColor: 'rgba(255,152,0,0.12)',
        tension: 0.3,
        pointRadius: 0,
        fill: true,
      }]
    },
    options: {
      plugins: { legend: { display: false } },
      scales: {
        x: { display: false },
        y: { color: '#bfc2cf', beginAtZero: false }
      },
      responsive: true,
      maintainAspectRatio: false,
    }
  });
}

async function updateStockPanel(symbol = 'RELIANCE.NSE') {
  try {
    const data = await fetchStockData(symbol);
    renderStockChart(data.slice(0, 60));
    // Example: update equity, holdings, etc.
    document.getElementById('equity-value').textContent = '₹' + (data[0]?.close || '--');
    document.getElementById('holdings-current').textContent = '₹' + (data[0]?.close || '--');
    document.getElementById('holdings-investment').textContent = '₹' + (data[data.length-1]?.close || '--');
    const pnl = data[0]?.close && data[data.length-1]?.close ? (data[0].close - data[data.length-1].close).toFixed(2) : '--';
    document.getElementById('holdings-pnl').textContent = pnl > 0 ? '+' + pnl : pnl;
    // Holdings bar
    const percent = data[0]?.close && data[data.length-1]?.close ? Math.max(0, Math.min(100, ((data[0].close / data[data.length-1].close) * 100))) : 60;
    document.getElementById('holdings-bar').style.background = `linear-gradient(90deg, #ff9800 ${percent}%, #232733 ${percent}%)`;
  } catch (err) {
    document.getElementById('equity-value').textContent = '--';
    document.getElementById('holdings-current').textContent = '--';
    document.getElementById('holdings-investment').textContent = '--';
    document.getElementById('holdings-pnl').textContent = '--';
    document.getElementById('holdings-bar').style.background = `linear-gradient(90deg, #ff9800 60%, #232733 60%)`;
  }
}

// --- News Fetch ---
async function fetchNews() {
  const newsLoading = document.getElementById('news-loading');
  const newsList = document.getElementById('news-list');
  newsLoading.style.display = '';
  newsList.innerHTML = '';
  try {
    const res = await fetch('/api/news');
    if (!res.ok) throw new Error('Failed to fetch news');
    const data = await res.json();
    if (!data.length) {
      newsList.innerHTML = '<div class="news-loading">No news found.</div>';
      return;
    }
    data.forEach(item => {
      const div = document.createElement('div');
      div.className = 'news-item';
      div.innerHTML = `<a href="${item.url}" class="news-title" target="_blank">${item.title}</a><div class="news-summary">${item.summary || ''}</div>`;
      newsList.appendChild(div);
    });
  } catch (err) {
    newsList.innerHTML = `<div class="news-loading">${err.message}</div>`;
  } finally {
    newsLoading.style.display = 'none';
  }
}

window.addEventListener('DOMContentLoaded', () => {
  updateStockPanel();
  fetchNews();
});
