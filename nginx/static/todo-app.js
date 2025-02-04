let token = null;
let g_username = null;
const API_ROUTE = "/api/v1"

// Handle Signup
async function signup() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const response = await fetch(`${API_ROUTE}/user/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    const result = await response.json();
    if (response.ok) {
        token = result.token;
        g_username = username;
        document.getElementById('auth-status').textContent = 'Signup successful, logged in!';
        document.getElementById('auth-section').style.display = 'none';
        document.getElementById('todo-section').style.display = 'block';
        await loadTodos();
    } else {
        document.getElementById('auth-status').textContent = 
            `Signup failed (${response.statusText}): ${result}`;
    }
}

// Handle Login
async function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const response = await fetch(`${API_ROUTE}/user/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    const result = await response.json();
    if (response.ok) {
        token = result.token;
        g_username = username;
        document.getElementById('auth-status').textContent = 'Login successful!';
        document.getElementById('auth-section').style.display = 'none';
        document.getElementById('todo-section').style.display = 'block';
        await loadTodos();
    } else {
        document.getElementById('auth-status').textContent = 
            `Login failed (${response.statusText}): ${result}`;
    }
}

// Load Todos
async function loadTodos() {
    const response = await fetch(`${API_ROUTE}/todos`, {
        headers: { 'Authentication': `Bearer ${token}` }
    });
    const headerTitle = document.getElementById('header-title');
    headerTitle.textContent = `${g_username}s todo list`;

    const todos = await response.json();
    const todoList = document.getElementById('todo-list');
    todoList.innerHTML = '';

    if (response.ok) {
        todos.forEach(todo => {
            const todoItem = document.createElement('div');
            todoItem.className = `todo-item ${todo.isDone ? 'completed' : ''}`;

            const taskTitle = document.createElement('h3');
            taskTitle.textContent = todo.task;

            const createdAt = document.createElement('h8');
            createdAt.textContent = `Created at: ${todo.createdAt}`;

            const toggleButton = document.createElement('button');
            toggleButton.textContent = todo.isDone ? 'Undo' : 'Complete';
            toggleButton.addEventListener('click', () => toggleTodoStatus(todo.id, todo.isDone));

            const deleteButton = document.createElement('button');
            deleteButton.textContent = 'Delete';
            deleteButton.addEventListener('click', () => deleteTodo(todo.id));

            todoItem.appendChild(taskTitle);
            todoItem.appendChild(createdAt);
            todoItem.appendChild(toggleButton);
            todoItem.appendChild(deleteButton);

            todoList.appendChild(todoItem);});
    }
}

// Create Todo
async function createTodo() {
    const task = document.getElementById('new-todo').value;

    const response = await fetch(`${API_ROUTE}/todo`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authentication': `Bearer ${token}`
        },
        body: JSON.stringify({ task })
    });

    if (response.ok) {
        await loadTodos();
    }
}

// Toggle Todo Status
async function toggleTodoStatus(id, currentStatus) {
    const response = await fetch(`${API_ROUTE}/todo/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authentication': `Bearer ${token}`
        },
        body: JSON.stringify({ isDone: !currentStatus })
    });

    if (response.ok) {
        await loadTodos();
    }
}

// Delete Todo
async function deleteTodo(id) {
    const response = await fetch(`${API_ROUTE}/todo/${id}`, {
        method: 'DELETE',
        headers: {
            'Authentication': `Bearer ${token}`
        }
    });

    if (response.ok) {
        await loadTodos();
    }
}

// Logout
function logout() {
    token = null;
    document.getElementById('auth-section').style.display = 'block';
    document.getElementById('todo-section').style.display = 'none';
}

