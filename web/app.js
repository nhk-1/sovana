const fileInput = document.getElementById('file');
const analyzeBtn = document.getElementById('analyze');
const statusEl = document.getElementById('status');
const listEl = document.getElementById('subscriptions');
const countEl = document.getElementById('count');

analyzeBtn.addEventListener('click', async () => {
  if (!fileInput.files.length) {
    setStatus('Ajoutez un fichier CSV.', 'error');
    return;
  }

  const formData = new FormData();
  formData.append('file', fileInput.files[0]);

  setStatus('Analyse en cours...', 'info');
  analyzeBtn.disabled = true;

  try {
    const res = await fetch('/api/upload', {
      method: 'POST',
      body: formData,
    });
    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || 'Erreur lors de l’analyse');
    }
    const data = await res.json();
    renderSubscriptions(data.subscriptions || []);
    setStatus(`Analyse terminée : ${data.count || 0} abonnement(s) détecté(s).`, 'success');
  } catch (err) {
    console.error(err);
    setStatus(err.message, 'error');
  } finally {
    analyzeBtn.disabled = false;
  }
});

function renderSubscriptions(subscriptions) {
  countEl.textContent = subscriptions.length;
  listEl.innerHTML = '';

  if (!subscriptions.length) {
    listEl.classList.add('empty');
    listEl.innerHTML = '<p class="empty__text">Aucun abonnement détecté. Essayez un autre relevé.</p>';
    return;
  }

  listEl.classList.remove('empty');
  subscriptions.forEach((sub) => {
    const el = document.createElement('article');
    el.className = 'subscription';
    el.innerHTML = `
      <div>
        <h3 class="subscription__title">${sub.merchant}</h3>
        <p class="subscription__meta">${sub.amount.toFixed(2)} € / mois • Vue le ${formatDate(sub.first_seen)} → ${formatDate(sub.last_seen)}</p>
        <p class="subscription__meta">Coût annuel estimé : <strong>${sub.annual_cost.toFixed(2)} €</strong></p>
      </div>
      <div class="subscription__cta">${renderCancelLink(sub.cancel_url)}</div>
    `;
    listEl.appendChild(el);
  });
}

function renderCancelLink(url) {
  if (!url) return '<span class="subscription__meta">Lien de résiliation non trouvé</span>';
  return `<a class="link" href="${url}" target="_blank" rel="noopener">Comment résilier</a>`;
}

function formatDate(value) {
  if (!value) return '—';
  const date = new Date(value);
  return date.toLocaleDateString('fr-FR');
}

function setStatus(text, variant) {
  statusEl.textContent = text;
  statusEl.style.color = variant === 'success' ? '#22c55e' : variant === 'error' ? '#ef4444' : '#e2e8f0';
}
