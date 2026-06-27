// ============================================
// TaskFlow Dashboard — JavaScript
// ============================================

let currentFilter = { category: '', priority: '', status: 'all' };
let tasks = [];

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    loadTasks();
    setupNavigation();
    setupKeyboardShortcuts();
    createToastContainer();
});

// ============================================
// Navigation & Filters
// ============================================

function setupNavigation() {
    document.querySelectorAll('.nav-item').forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const filterType = item.dataset.filter;
            const value = item.dataset.value;

            // Reset other filters when switching filter type
            if (filterType === 'status') {
                currentFilter = { category: '', priority: '', status: value };
            } else if (filterType === 'category') {
                currentFilter.category = currentFilter.category === value ? '' : value;
                currentFilter.status = 'all';
            } else if (filterType === 'priority') {
                currentFilter.priority = currentFilter.priority === value ? '' : value;
                currentFilter.status = 'all';
            }

            // Update active state
            document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
            item.classList.add('active');

            updateFilterIndicator();
            loadTasks();
        });
    });
}

function updateFilterIndicator() {
    const indicator = document.getElementById('filter-indicator');
    let text = 'Mostrando ';

    if (currentFilter.status === 'pending') text += 'tareas pendientes';
    else if (currentFilter.status === 'completed') text += 'tareas completadas';
    else text += 'todas las tareas';

    if (currentFilter.category) text += ` — ${currentFilter.category}`;
    if (currentFilter.priority) text += ` — prioridad ${currentFilter.priority}`;

    indicator.textContent = text;
}

// ============================================
// Task CRUD
// ============================================

async function loadTasks() {
    try {
        const params = new URLSearchParams();
        if (currentFilter.category) params.set('category', currentFilter.category);
        if (currentFilter.priority) params.set('priority', currentFilter.priority);
        if (currentFilter.status && currentFilter.status !== 'all') params.set('status', currentFilter.status);

        const res = await fetch(`/api/tasks?${params.toString()}`);
        if (!res.ok) {
            if (res.status === 401) {
                window.location.href = '/';
                return;
            }
            throw new Error('Error loading tasks');
        }

        tasks = await res.json();
        renderTasks();
        updateStats();
    } catch (err) {
        console.error('Error:', err);
        showToast('Error al cargar tareas', 'error');
    }
}

function renderTasks() {
    const list = document.getElementById('task-list');
    const empty = document.getElementById('empty-state');

    if (!tasks || tasks.length === 0) {
        list.innerHTML = '';
        empty.style.display = 'block';
        return;
    }

    empty.style.display = 'none';
    list.innerHTML = tasks.map((task, index) => `
        <div class="task-item ${task.completed ? 'completed' : ''}" style="animation-delay: ${index * 0.05}s">
            <div class="task-checkbox ${task.completed ? 'checked' : ''}" onclick="toggleTask(${task.id})"></div>
            <div class="task-content">
                <div class="task-title">${escapeHtml(task.title)}</div>
                <div class="task-meta">
                    <span class="task-category category-${task.category}">${getCategoryLabel(task.category)}</span>
                    <span class="task-priority priority-${task.priority}">${getPriorityLabel(task.priority)}</span>
                    <span class="task-date">${formatDate(task.created_at)}</span>
                </div>
            </div>
            <div class="task-actions">
                <button class="task-action-btn edit" onclick="editTask(${task.id})" title="Editar">✏️</button>
                <button class="task-action-btn delete" onclick="deleteTask(${task.id})" title="Eliminar">🗑️</button>
            </div>
        </div>
    `).join('');
}

function updateStats() {
    const total = tasks ? tasks.length : 0;
    const completed = tasks ? tasks.filter(t => t.completed).length : 0;
    const pending = total - completed;
    const progress = total > 0 ? Math.round((completed / total) * 100) : 0;

    animateNumber('stat-total', total);
    animateNumber('stat-pending', pending);
    animateNumber('stat-done', completed);
    document.getElementById('stat-progress').textContent = progress + '%';
    document.getElementById('progress-fill').style.width = progress + '%';

    // Update badges
    document.getElementById('badge-all').textContent = total;
    document.getElementById('badge-pending').textContent = pending;
    document.getElementById('badge-completed').textContent = completed;
}

function animateNumber(elementId, target) {
    const el = document.getElementById(elementId);
    const current = parseInt(el.textContent) || 0;
    if (current === target) return;

    const duration = 400;
    const start = performance.now();

    function update(now) {
        const elapsed = now - start;
        const progress = Math.min(elapsed / duration, 1);
        const eased = 1 - Math.pow(1 - progress, 3);
        el.textContent = Math.round(current + (target - current) * eased);
        if (progress < 1) requestAnimationFrame(update);
    }
    requestAnimationFrame(update);
}

async function toggleTask(id) {
    try {
        const res = await fetch(`/api/tasks/${id}/toggle`, { method: 'PATCH' });
        if (!res.ok) throw new Error();
        
        // Update local state for immediate feedback
        const task = tasks.find(t => t.id === id);
        if (task) {
            task.completed = !task.completed;
            renderTasks();
            updateStats();
            showToast(task.completed ? 'Tarea completada ✓' : 'Tarea reactivada', 'success');
        }
    } catch {
        showToast('Error al actualizar tarea', 'error');
    }
}

async function deleteTask(id) {
    if (!confirm('¿Eliminar esta tarea?')) return;

    try {
        const res = await fetch(`/api/tasks/${id}`, { method: 'DELETE' });
        if (!res.ok) throw new Error();
        
        tasks = tasks.filter(t => t.id !== id);
        renderTasks();
        updateStats();
        showToast('Tarea eliminada', 'info');
    } catch {
        showToast('Error al eliminar tarea', 'error');
    }
}

function editTask(id) {
    const task = tasks.find(t => t.id === id);
    if (!task) return;

    document.getElementById('task-id').value = task.id;
    document.getElementById('task-title').value = task.title;
    document.getElementById('task-description').value = task.description || '';
    document.getElementById('task-category').value = task.category;
    document.getElementById('task-priority').value = task.priority;
    document.getElementById('modal-title').textContent = 'Editar Tarea';

    openModal();
}

async function handleTaskSubmit(e) {
    e.preventDefault();

    const id = document.getElementById('task-id').value;
    const data = {
        title: document.getElementById('task-title').value,
        description: document.getElementById('task-description').value,
        category: document.getElementById('task-category').value,
        priority: document.getElementById('task-priority').value,
    };

    try {
        let res;
        if (id) {
            res = await fetch(`/api/tasks/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
        } else {
            res = await fetch('/api/tasks', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
        }

        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.error || 'Error');
        }

        closeModal();
        loadTasks();
        showToast(id ? 'Tarea actualizada' : 'Tarea creada ✓', 'success');
    } catch (err) {
        showToast(err.message, 'error');
    }
}

// ============================================
// Modal
// ============================================

function openModal() {
    document.getElementById('modal-overlay').classList.add('active');
    document.body.style.overflow = 'hidden';
    setTimeout(() => document.getElementById('task-title').focus(), 200);
}

function closeModal(e) {
    if (e && e.target !== e.currentTarget) return;
    document.getElementById('modal-overlay').classList.remove('active');
    document.body.style.overflow = '';
    resetForm();
}

function resetForm() {
    document.getElementById('task-form').reset();
    document.getElementById('task-id').value = '';
    document.getElementById('modal-title').textContent = 'Nueva Tarea';
}

// ============================================
// Sidebar Toggle (Mobile)
// ============================================

function toggleSidebar() {
    document.getElementById('sidebar').classList.toggle('open');
}

// ============================================
// Keyboard Shortcuts
// ============================================

function setupKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') closeModal();
        if (e.key === 'n' && e.ctrlKey) {
            e.preventDefault();
            openModal();
        }
    });
}

// ============================================
// Helpers
// ============================================

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function getCategoryLabel(cat) {
    const labels = { trabajo: '💼 Trabajo', personal: '🏠 Personal', estudio: '📚 Estudio' };
    return labels[cat] || cat;
}

function getPriorityLabel(pri) {
    const labels = { alta: '🔴 Alta', media: '🟡 Media', baja: '🟢 Baja' };
    return labels[pri] || pri;
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    const now = new Date();
    const diff = now - date;
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) return 'Hoy';
    if (days === 1) return 'Ayer';
    if (days < 7) return `Hace ${days} días`;
    return date.toLocaleDateString('es-ES', { day: 'numeric', month: 'short' });
}

// ============================================
// Toasts
// ============================================

function createToastContainer() {
    if (!document.querySelector('.toast-container')) {
        const container = document.createElement('div');
        container.className = 'toast-container';
        document.body.appendChild(container);
    }
}

function showToast(message, type = 'info') {
    const container = document.querySelector('.toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    const icons = { success: '✓', error: '✕', info: 'ℹ' };
    toast.innerHTML = `<span>${icons[type] || ''}</span> ${message}`;

    container.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'toast-out 0.3s ease forwards';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}
