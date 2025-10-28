const API_ROOT = 'http://localhost:8080/api/v1';

const qs = (sel, root=document) => root.querySelector(sel);
const qsa = (sel, root=document) => [...root.querySelectorAll(sel)];

const state = {
  token: localStorage.getItem('token') || '',
  me: null,
  teamId: null,
  teamName: null,
};

function show(id) { qs('#'+id).classList.remove('hidden'); }
function hide(id) { qs('#'+id).classList.add('hidden'); }
function setTab(tab) {
  qsa('.tab').forEach(t=>t.classList.toggle('active', t.dataset.tab===tab));
  qsa('.tab-content').forEach(c=>c.classList.add('hidden'));
  qs('#tab-'+tab).classList.remove('hidden');
}

function headers(json=true) {
  const h = { 'Authorization': 'Bearer ' + state.token };
  if (json) h['Content-Type'] = 'application/json';
  return h;
}

async function api(path, opts={}) {
  const headersIn = opts.headers || {};
  const h = { ...headersIn };
  if (state.token && !('Authorization' in h)) {
    h['Authorization'] = 'Bearer ' + state.token;
  }
  // If body is present and no content-type set, default to JSON
  if (opts.body && !('Content-Type' in h)) {
    h['Content-Type'] = 'application/json';
  }
  const res = await fetch(API_ROOT + path, { ...opts, headers: h });
  if (!res.ok) {
    if (res.status === 401) {
      // token missing/invalid -> force logout UI-side
      logout();
    }
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  const ct = res.headers.get('content-type') || '';
  if (ct.includes('application/json')) return res.json();
  return res.text();
}

function notifyError(el, msg) {
  el.textContent = '';
  try {
    const j = JSON.parse(msg);
    el.textContent = j.error || msg;
  } catch { el.textContent = msg; }
}

function logout() {
  state.token = '';
  localStorage.removeItem('token');
  hide('main-screen');
  show('auth-screen');
}

async function refreshMyTeams() {
  const list = await api('/me/teams', { headers: headers(false) });
  const el = qs('#teams-list');
  el.innerHTML = '';
  if (!Array.isArray(list) || list.length === 0) {
    el.innerHTML = '<div class="list-item">Команды не найдены</div>';
    return;
  }
  list.forEach(t => {
    const div = document.createElement('div');
    div.className = 'list-item';
    div.innerHTML = `<h3>${t.name}</h3><p>${t.description||''}</p>`;
    div.onclick = ()=> openTeam(t.id, t.name);
    el.appendChild(div);
  });
}

async function openTeam(id, name) {
  state.teamId = id; state.teamName = name;
  qs('#team-name').textContent = name;
  hide('teams-section');
  show('team-view');
  setTab('artifacts');
  try { await loadArtifacts(); } catch(e){ console.warn('artifacts load failed', e); }
  try { await loadContacts(); } catch(e){ console.warn('contacts load failed', e); }
  try { await loadRequestsSafe(); } catch(e){ console.warn('requests load failed', e); }
}

async function loadArtifacts() {
  let arr = await api(`/teams/${state.teamId}/artifacts`, { headers: headers(false) });
  if (!Array.isArray(arr)) arr = [];
  const el = qs('#artifacts-list');
  el.innerHTML='';
  arr.forEach(a=>{
    const item = document.createElement('div');
    item.className='list-item';
    item.innerHTML=`<h3>${a.name} <span class=\"badge\">${a.type}</span></h3>
      <p>${a.description||''}</p>
      <div class=\"meta\">Проект: ${a.project_name} • ID: ${a.id}</div>
      <div class=\"actions\">\
        <button class=\"btn btn-small\" data-act=\"show-fields\">Поля</button>\
        <button class=\"btn btn-small\" data-act=\"add-field\">+ Поле</button>\
        <button class=\"btn btn-small btn-secondary\" data-act=\"delete\">Удалить</button>\
      </div>
      <div class=\"fields\" style=\"margin-top:10px; display:none;\"></div>`;
    item.onclick = async (e)=>{
      const act = e.target?.dataset?.act;
      if (act==='delete') {
        if (confirm('Удалить артефакт?')) {
          await api(`/teams/${state.teamId}/artifacts/${a.id}`, { method:'DELETE', headers: headers(false) });
          await loadArtifacts();
        }
        e.stopPropagation();
      }
      if (act==='add-field') {
        openModalField(a.id, async ()=>{ // callback to refresh fields
          await renderFields(item, a.id);
        });
        e.stopPropagation();
      }
      if (act==='show-fields') {
        await renderFields(item, a.id);
        e.stopPropagation();
      }
    };
    el.appendChild(item);
  });
}

async function loadContacts() {
  let arr = await api(`/teams/${state.teamId}/contacts`, { headers: headers(false) });
  if (!Array.isArray(arr)) arr = [];
  const el = qs('#contacts-list');
  el.innerHTML='';
  arr.forEach(c=>{
    const item = document.createElement('div');
    item.className='list-item';
    item.innerHTML=`<h3>${c.name}</h3><div class="meta">TG: ${c.telegram_contact||''} • ID: ${c.id}</div>
      <div class="actions">
        <button class="btn btn-small btn-secondary" data-act="delete">Удалить</button>
      </div>`;
    item.onclick = async (e)=>{
      if (e.target?.dataset?.act==='delete') {
        if (confirm('Удалить контакт?')) {
          await api(`/teams/${state.teamId}/contacts/${c.id}`, { method:'DELETE', headers: headers(false) });
          await loadContacts();
        }
        e.stopPropagation();
      }
    };
    el.appendChild(item);
  });
}

async function loadRequestsSafe() {
  // может вернуть 403 для member/viewer — просто скрываем
  try {
    const arr = await api(`/teams/${state.teamId}/requests`, { headers: headers(false) });
    const el = qs('#requests-list');
    el.innerHTML='';
    arr.forEach(r=>{
      const item = document.createElement('div');
      item.className='list-item';
      const badge = r.status==='pending'?'badge-pending':(r.status==='approved'?'badge-approved':'badge-rejected');
      item.innerHTML=`<h3>Заявка #${r.id} <span class="badge ${badge}">${r.status}</span></h3>
        <div class="meta">От пользователя: ${r.user_id} • ${new Date(r.created_at).toLocaleString()}</div>
        ${r.status==='pending'?'<div class="actions">\
          <button class="btn btn-small" data-act="approve">Одобрить</button>\
          <button class="btn btn-small btn-secondary" data-act="reject">Отклонить</button>\
        </div>':''}`;
      item.onclick = async (e)=>{
        const act = e.target?.dataset?.act;
        if (act==='approve' || act==='reject') {
          await api(`/teams/${state.teamId}/requests/${r.id}/${act}`, { method:'POST', headers: headers(false) });
          await loadRequestsSafe();
        }
      };
      el.appendChild(item);
    });
    qs('[data-tab="requests"]').classList.remove('hidden');
  } catch (e) {
    qs('[data-tab="requests"]').classList.add('hidden');
  }
}

function openModal(html) {
  qs('#modal-body').innerHTML = html;
  qs('#modal').classList.remove('hidden');
}
function closeModal(){ qs('#modal').classList.add('hidden'); }

function openModalCreateTeam() {
  openModal(`
    <h3>Создать команду</h3>
    <input type="text" id="m-team-name" placeholder="Название"/>
    <textarea id="m-team-desc" placeholder="Описание"></textarea>
    <div class="btn-group">
      <button id="m-team-save" class="btn">Создать</button>
    </div>
  `);
  qs('#m-team-save').onclick = async ()=>{
    const name = qs('#m-team-name').value.trim();
    const description = qs('#m-team-desc').value.trim();
    if (!name) return;
    await api('/teams', { method:'POST', headers: headers(), body: JSON.stringify({name, description}) });
    closeModal();
    await refreshMyTeams();
  };
}

function openModalArtifact() {
  openModal(`
    <h3>Создать артефакт</h3>
    <input id="m-art-name" placeholder="Имя"/>
    <select id="m-art-type">
      <option value="table">table</option>
      <option value="view">view</option>
      <option value="procedure">procedure</option>
      <option value="function">function</option>
      <option value="index">index</option>
      <option value="dataset">dataset</option>
      <option value="api">api</option>
      <option value="file">file</option>
    </select>
    <input id="m-art-project" placeholder="Проект"/>
    <textarea id="m-art-desc" placeholder="Описание"></textarea>
    <div class="btn-group">
      <button id="m-art-save" class="btn">Создать</button>
    </div>
  `);
  qs('#m-art-save').onclick = async ()=>{
    const body = {
      name: qs('#m-art-name').value.trim(),
      type: qs('#m-art-type').value,
      description: qs('#m-art-desc').value.trim(),
      project_name: qs('#m-art-project').value.trim(),
      developer_id: 1
    };
    if (!body.name || !body.project_name) return;
    await api(`/teams/${state.teamId}/artifacts`, { method:'POST', headers: headers(), body: JSON.stringify(body)});
    closeModal();
    await loadArtifacts();
  };
}

async function renderFields(container, artifactId){
  const box = container.querySelector('.fields');
  if (!box) return;
  const arr = await api(`/teams/${state.teamId}/artifacts/${artifactId}/fields`, { headers: headers(false) });
  box.style.display = 'block';
  if (!Array.isArray(arr) || arr.length===0){ box.innerHTML = '<div style="color:#777;">Нет полей</div>'; return; }
  box.innerHTML = arr.map(f=>`<div style=\"font-size:13px; padding:6px 0; border-top:1px dashed #ddd;\">`+
    `<b>${f.field_name}</b>: ${f.data_type} ${f.is_pk?'(PK)':''}</div>`).join('');
}

function openModalField(artifactId, onCreated) {
  openModal(`
    <h3>Добавить поле</h3>
    <input id="m-field-name" placeholder="Имя поля"/>
    <input id="m-field-type" placeholder="Тип данных"/>
    <textarea id="m-field-desc" placeholder="Описание"></textarea>
    <label><input type="checkbox" id="m-field-pk"/> Первичный ключ</label>
    <div class="btn-group"><button id="m-field-save" class="btn">Добавить</button></div>
  `);
qs('#m-field-save').onclick = async ()=>{
    const body = {
      field_name: qs('#m-field-name').value.trim(),
      data_type: qs('#m-field-type').value.trim(),
      description: qs('#m-field-desc').value.trim(),
      is_pk: qs('#m-field-pk').checked
    };
    if (!body.field_name || !body.data_type) return;
    try {
      await api(`/teams/${state.teamId}/artifacts/${artifactId}/fields`, { method:'POST', headers: headers(), body: JSON.stringify(body)});
      closeModal();
      if (typeof onCreated === 'function') { await onCreated(); }
      alert('Поле создано');
    } catch(e){ alert('Ошибка: ' + e.message); }
  };
}

function openModalContact() {
  openModal(`
    <h3>Создать контакт</h3>
    <input id="m-contact-name" placeholder="Имя"/>
    <input id="m-contact-tg" placeholder="Telegram @username"/>
    <div class="btn-group"><button id="m-contact-save" class="btn">Создать</button></div>
  `);
  qs('#m-contact-save').onclick = async ()=>{
    const body = { name: qs('#m-contact-name').value.trim(), telegram_contact: qs('#m-contact-tg').value.trim() };
    if (!body.name) return;
    await api(`/teams/${state.teamId}/contacts`, { method:'POST', headers: headers(), body: JSON.stringify(body)});
    closeModal();
    await loadContacts();
  };
}

function bindUI() {
  // Auth
  qs('#btn-login').onclick = async ()=>{
    const email = qs('#email').value.trim();
    const password = qs('#password').value;
    const err = qs('#auth-error'); err.textContent='';
    try {
      const data = await api('/auth/login', { method:'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({email,password})});
      state.token = data.token; localStorage.setItem('token', state.token);
      qs('#user-email').textContent = email;
      hide('auth-screen'); show('main-screen');
      await refreshMyTeams();
    } catch (e) { notifyError(err, e.message); }
  };

  qs('#btn-register').onclick = async ()=>{
    const email = qs('#email').value.trim();
    const password = qs('#password').value;
    const name = qs('#name').value.trim();
    const err = qs('#auth-error'); err.textContent='';
    try {
      const data = await api('/auth/register', { method:'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({email,password,name})});
      state.token = data.token; localStorage.setItem('token', state.token);
      qs('#user-email').textContent = email;
      hide('auth-screen'); show('main-screen');
      await refreshMyTeams();
    } catch (e) { notifyError(err, e.message); }
  };

  qs('#btn-logout').onclick = logout;

  // Teams
  qs('#btn-create-team').onclick = openModalCreateTeam;
  qs('#btn-search-teams').onclick = async ()=>{
    const q = qs('#team-search').value.trim();
    if (!q) return;
    const res = await api(`/teams?search=${encodeURIComponent(q)}`, { headers: headers(false) });
    const el = qs('#search-results'); el.classList.remove('hidden'); el.innerHTML='';
    res.forEach(t=>{
      const div = document.createElement('div');
      div.className = 'list-item';
      div.innerHTML = `<h3>${t.name}</h3><p>${t.description||''}</p><div class="actions"><button class="btn btn-small">Запросить доступ</button></div>`;
      div.onclick = async (e)=>{
        if (e.target.tagName==='BUTTON') {
          await api(`/teams/${t.id}/join`, { method:'POST', headers: headers(false) });
          alert('Запрос отправлен');
          e.stopPropagation();
        } else {
          openTeam(t.id, t.name);
        }
      };
      el.appendChild(div);
    });
  };

  qs('#btn-back-to-teams').onclick = ()=>{
    hide('team-view'); show('teams-section');
  };

  // Tabs
  qsa('.tab').forEach(btn=>{
    btn.onclick = ()=> setTab(btn.dataset.tab);
  });

  // Creates
  qs('#btn-create-artifact').onclick = openModalArtifact;
  qs('#btn-create-contact').onclick = openModalContact;

  // Modal close
  qs('.modal .close').onclick = closeModal;
  qs('#modal').addEventListener('click', (e)=>{ if (e.target.id==='modal') closeModal(); });
}

window.addEventListener('DOMContentLoaded', async ()=>{
  bindUI();
  if (state.token) {
    // try to load
    show('main-screen'); hide('auth-screen');
    await refreshMyTeams().catch(()=>logout());
  }
});
