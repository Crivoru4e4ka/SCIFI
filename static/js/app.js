const { createApp } = Vue;

// 1. Создаем логику "Клик вне элемента"
const clickOutside = {
  beforeMount: (el, binding) => {
    el.clickOutsideEvent = event => {
      // Если клик был НЕ по самому элементу (el) и НЕ по его внутренностям
      if (!(el == event.target || el.contains(event.target))) {
        // Вызываем функцию, которую передали (в нашем случае closeSearch)
        binding.value();
      }
    };
    document.addEventListener("click", el.clickOutsideEvent);
  },
  unmounted: el => {
    document.removeEventListener("click", el.clickOutsideEvent);
  },
};

const app = createApp({
            data() {
                return {
                    currentUserProjectRole: '', 
                    profile: {
                        full_name: "",
                        email: "",
                        initials: ""
                    },
                    project: {
                        name: "",
                        description: "",
                        status: "active",
                        start_date: "",
                        research_goal: "",
                        main_hypothesis: "",
                        novelty: "",
                        expected_result: "",
                        execution_type: "manual",
                        team_id: null},
                    projectError: "",
                    projectSuccess: "",
                    isSubmitting: false,
                    createProjectModal: null,
                    createTaskModal: null,
newTask: {
    // --- БАЗОВЫЕ ПОЛЯ ---
    projectId: "",
    type: "research",
    status: "todo",
    title: "",
    description: "",
    assigneeId: null,
    priority: "Medium",
    dueDate: "",
    tags: "",

    // --- ОБЩИЙ НАУЧНЫЙ ПАСПОРТ ---
    researchContribution: "",
    researchMethod: "",
    expectedScientificResult: "",
    conclusion: "",
    doi: "",

    // --- ИССЛЕДОВАНИЕ ---
    resQuestion: "",
    resObject: "",
    resSubject: "",
    resSources: [{ title: '', authors: '', year: '', doi: '' }],
    resPubCount: "",
    resDatabases: "",
    resApproach: "",

    // --- ЭКСПЕРИМЕНТ ---
    expGoal: "",
    expVariables: [{ type: 'Independent', name: '', value: '' }],
    expControls: "",
    expInputData: "",
    expSampleSize: "",
    expMetrics: "",
    expResults: "",
    expStatSignificance: "",

    // --- СБОР ДАННЫХ ---
    dsSourceType: "",
    dsMethod: "",
    dsLink: "",
    dsFormat: "csv",
    dsCompleteness: 100,
    sampleSize: "",

    // --- АНАЛИЗ ---
    anTools: "",
    anVisuals: "",
    statIndicators: "",
    statPValue: 0.05,
    hypConfirmed: false,

    // --- РАЗРАБОТКА ---
    devComponent: "",
    devStack: "",
    devRepo: "",
    devSwagger: "",
    devTestCoverage: 0,
    devVersion: "",

    // --- ПУБЛИКАЦИЯ ---
    docType: "article",
    journal: "",
    docStatus: "draft",
    coAuthors: ""
},
                    taskError: "",
                    taskSuccess: "",
                    isTaskSubmitting: false,
                    projects: [],
                    sprints: [],
                    sprintForm: {
    id: null,
    name: '',
    duration: '2', // По умолчанию 2 недели
    startDate: '',
    endDate: '',
    goal: ''
},
newCommentText: '',
replyText: '',
replyToId: null,
activeComments: [],
sprintTaskCount: 0,
startSprintModal: null,
                    draggedTask: null,
                    projectTasks: [],
                    projectAttachments: [],
                    searchQuery: '',           // Текст поиска
                    filterMyTasks: false,      // Флаг "Только мои задачи"
                    filterPriority: 'all',
                    filterSubPage: 'list', // 'list' — список фильтров, 'search' — поиск задач
allTasks: [],          // Сюда загрузим задачи изо всех проектов
navFilters: {
    query: '',
    project: 'all',
    assignee: 'all',
    status: 'all',
    priority: 'all'
},
                    isSearchOpen: false,
                    displayedTasks: [],     // Фильтр по приоритету
                    expandedTaskId: null, 
                    selectedFile: null, 
                    users: [],
                    teams: [],
                    inviteEmail: '',
                    manageInviteEmail: '',
                    showManageSuggestions: false,
                    showSuggestions: false, // Список команд
                    newTeam: { name: '', description: '', member_emails: [] },
createTeamModal: null,
selectedTeam: { name: '', description: '' },
currentTeamMembers: [],
manageTeamModal: null,
recentItems: [],
flaggedItems: [],
                    currentPage: "sections",
                    activeSidebarTab: 'sections',
                    activities: [],
                    selectedTemplate: {
                        id: "kanban",
                        title: "Kanban",
                        description: "Работайте эффективно с задачами в формате доски."
                    },
                    kanbanWizard: {
                        step: 1,
                        name: "",
                        managementType: "manual",
                        teamId: null,
                        access: "closed",
                        key: "",
                        memberRows: [{ user_id: null, role: 'researcher', searchText: '', showSuggestions: false }],             // Список уже добавленных объектов {user_id, full_name, email, role}
                        research_goal: "", 
                        main_hypothesis: "", 
                        novelty: "", 
                        expected_result: "",
                    },
                    createdKanban: {
                        id: null,
                        name: "",
                        key: "",
                        description: "",
                        viewMode: "board",
                        boardTasks: {
                            todo: [],
                            inProgress: [],
                            done: []
                        }
                    },
                    currentProject: null,
                    currentProjectView: "board",
                    projectBoardTasks: {
                        todo: [],
                        inProgress: [],
                        review: [],
                        done: []
                    },
                    draggedFromStatus: null,
                    templates: [
                        {
                            id: "kanban",
                            title: "Kanban",
                            description: "Работайте эффективно с задачами в формате доски.",
                            previewClass: "template-preview-kanban"
                        },
                        {
                            id: "scrum",
                            title: "Scrum",
                            description: "Планируйте спринты и отслеживайте прогресс команды.",
                            previewClass: "template-preview-scrum"
                        },
                        {
                            id: "empty",
                            title: "Пустой проект",
                            description: "Начните с чистого листа и настройте всё под себя.",
                            previewClass: "template-preview-empty"
                        }
                    ],
                    grants: [],
                    currentGrant: null,
                    grantForm: {
                        title: '',
                        code: '',
                        funding_organization: '',
                        country: '',
                        description: '',
                        scientific_direction: '',
                        grant_type: 'state',
                        status: 'draft',
                        total_amount: 0,
                        currency: 'RUB',
                        start_date: '',
                        end_date: '',
                        application_deadline: '',
                        principal_investigator_id: null
                    },
                    grantFilters: {
                        status: '',
                        type: '',
                        search: ''
                    },
                    grantError: '',
                    grantSuccess: '',
                    isGrantSubmitting: false,
                    createGrantModal: null,
                    projectGrants: [],
                    projectPermissions: [],
                    projectRoles: [],
                    projectActivities: [],
                    manageMembersModal: null,
                    currentProjectMembers: [],
                    wizardTeamMembers: [],
                    newProjectMemberEmail: '',
                    newProjectMemberRole: 'researcher',
                    filteredProjectMemberUsers: [],
                    grantProjectForm: {
                        project_id: '',
                        allocated_amount: 0,
                        funding_purpose: '',
                        funding_start_date: '',
                        funding_end_date: ''
                    },
                    addProjectToGrantModal: null,
                    grantProjectError: ''
                };
            },
            computed: {
                submitButtonText() {
                    return this.isSubmitting ? "Создание..." : "Создать";
                },
                availableUsersForWizard() {
                    const currentUserId = Number(localStorage.getItem('currentUserId'));
                    return this.users.filter(u => u.id !== currentUserId);
                },
                taskTypeLabel() {
                    const labels = {
                        research: 'Исследование',
                        experiment: 'Эксперимент',
                        data_collection: 'Сбор данных',
                        analysis: 'Анализ результатов',
                        dev: 'Разработка',
                        doc: 'Публикация / Отчёт'
                    };
                    return labels[this.newTask?.type] || 'Задача';
                },
                grantStatusBadgeClass() {
                    const map = {
                        draft: 'bg-secondary',
                        submitted: 'bg-info',
                        under_review: 'bg-warning text-dark',
                        approved: 'bg-success',
                        rejected: 'bg-danger',
                        active: 'bg-primary',
                        completed: 'bg-success',
                        suspended: 'bg-dark'
                    };
                    return map[this.currentGrant?.status] || 'bg-secondary';
                },
                grantStatusLabel() {
                    const map = {
                        draft: 'Черновик',
                        submitted: 'Подан',
                        under_review: 'На рассмотрении',
                        approved: 'Одобрен',
                        rejected: 'Отклонен',
                        active: 'Активен',
                        completed: 'Завершен',
                        suspended: 'Приостановлен'
                    };
                    return map[this.currentGrant?.status] || this.currentGrant?.status;
                },
                taskSubmitButtonText() {
                    return this.isTaskSubmitting ? "Создание..." : "Создать";
                },
                canEditProject() {
                    return this.projectPermissions.includes('project.edit');
                },
                canManageMembers() {
                    return this.projectPermissions.includes('project.manage_members');
                },
                canCreateTask() {
                    return this.projectPermissions.includes('task.create');
                },
                canDeleteTask() {
                    return this.projectPermissions.includes('task.delete');
                },
                canApproveExperiment() {
                    return this.projectPermissions.includes('experiment.approve');
                },
                canValidateResults() {
                    return this.projectPermissions.includes('results.validate');
                },
                canViewAudit() {
                    return this.projectPermissions.includes('audit.view');
                },
                filteredManageUsers() {
    if (!this.manageInviteEmail || this.manageInviteEmail.length < 1) return [];
    const search = this.manageInviteEmail.toLowerCase();
    return this.users.filter(user => {
        const alreadyIn = this.currentTeamMembers.some(m => m.email === user.email);
        return !alreadyIn && (user.email.toLowerCase().includes(search) || user.full_name.toLowerCase().includes(search));
    }).slice(0, 5);
},
    // Поиск по проектам
    filteredSearchProjects() {
        if (!this.searchQuery || this.searchQuery.length < 1) return [];
        
        const q = this.searchQuery.toLowerCase();
        return this.projects.filter(p => {
            return p.name.toLowerCase().includes(q) || 
                   (p.key && p.key.toLowerCase().includes(q));
        }).slice(0, 5); // Ограничим до 5 результатов
    },
    // Глобальный поиск для выпадающего списка в шапке
globalSearchResults() {
    // Если в строке поиска пусто — ничего не показываем
    if (!this.searchQuery || this.searchQuery.trim().length === 0) return [];
    
    const q = this.searchQuery.toLowerCase();
    return this.allTasks.filter(t => {
        return t.title.toLowerCase().includes(q) || 
               (this.getProjectKey(t.project_id).toLowerCase() + '-' + t.task_num).includes(q);
    }).slice(0, 10);
},
    navigatorFilteredTasks() {
        return this.allTasks.filter(t => {
            const matchesSearch = !this.navFilters.query || t.title.toLowerCase().includes(this.navFilters.query.toLowerCase());
            const matchesProj = this.navFilters.project === 'all' || t.project_id === Number(this.navFilters.project);
            const matchesUser = this.navFilters.assignee === 'all' || t.assignee_id === Number(this.navFilters.assignee);
            const matchesStatus = this.navFilters.status === 'all' || t.status.toLowerCase() === this.navFilters.status.toLowerCase();
            return matchesSearch && matchesProj && matchesUser && matchesStatus;
        });
    },
filteredUsers() {
    if (!this.inviteEmail || this.inviteEmail.length < 1) return [];
    
    const search = this.inviteEmail.toLowerCase();
    const myEmail = this.profile.email ? this.profile.email.toLowerCase() : '';

    return this.users.filter(user => {
        const userEmail = user.email.toLowerCase();
        const userName = user.full_name.toLowerCase();
        
        // 1. Исключаем себя
        const isNotMe = userEmail !== myEmail;
        
        // 2. ИСПРАВЛЕННАЯ ЛОГИКА: проверяем, что ПОЧТА начинается на вводимые символы
        // ИЛИ ИМЯ начинается на эти символы
        // ИЛИ любое слово в имени начинается на эти символы (например, "М" для "Самир Мурзагулов")
        const matchesQuery = userEmail.startsWith(search) || 
                             userName.startsWith(search) ||
                             userName.split(' ').some(word => word.startsWith(search));
                             
        // 3. Убеждаемся, что не добавлен уже
        const notAlreadyInList = !this.newTeam.member_emails.includes(user.email);
        
        return isNotMe && matchesQuery && notAlreadyInList;
    }).slice(0, 5);
}
            },
            methods: {
                async openManageTeam(team) {
    this.selectedTeam = team;
    this.loadProjectRoles();
    try {
        // Запрос к бэкенду за участниками этой команды
        const response = await fetch(`/teams/${team.id}/members`);
        if (response.ok) {
            this.currentTeamMembers = await response.json();
            this.manageTeamModal.show();
        }
    } catch (e) { console.error(e); }
},

async deleteTask(taskId) {
    if (!confirm("Вы уверены, что хотите удалить эту задачу? Это действие необратимо.")) return;
    try {
        const response = await fetch(`/tasks/${taskId}`, { method: 'DELETE' });
        if (response.ok) {
            this.loadProjectTasks(this.currentProject.id); // Обновляем список
        } else {
            alert("Ошибка при удалении задачи");
        }
    } catch (e) { console.error(e); }
},

async deleteProject(projectId) {
    if (!confirm("ВНИМАНИЕ! Вы удаляете весь проект со всеми задачами, файлами и обсуждениями. Продолжить?")) return;
    try {
        const response = await fetch(`/projects/${projectId}`, { method: 'DELETE' });
        if (response.ok) {
            this.currentProject = null; // Закрываем вид проекта
            this.loadProjects();        // Обновляем список на главной
            this.activeSidebarTab = 'sections';
        } else {
            alert("Ошибка при удалении проекта");
        }
    } catch (e) { console.error(e); }
},

// Подготовка редактирования (открыть ту же модалку создания, но заполненную)
prepareEditTask(task) {
    this.newTask = {
        id: task.id,
        projectId: task.project_id,
        title: task.title,
        description: task.description,
        status: task.status,
        priority: task.priority,
        type: task.type,
        assigneeId: task.assignee_id,
        dueDate: task.due_date ? task.due_date.substring(0, 10) : ''
    };
    this.taskSubmitButtonText = "Сохранить изменения"; 
    this.createTaskModal.show();
},

openProjectActivity() {
    this.currentProjectView = 'audit';
    this.projectActivities = []; // Очищаем старое
    if (this.currentProject) {
        // Загружаем задачи, чтобы метод getTaskFormattedKey мог найти Ключи
        this.loadProjectTasks(this.currentProject.id);
        this.loadProjectActivities(this.currentProject.id);
    }
},

openProjectDiscussion() {
    this.currentProjectView = 'discussion'; // Переключаем вкладку
    this.loadComments(this.currentProject.id, 'project'); // Загружаем комменты проекта
},

// Универсальная загрузка
async loadComments(entityId, entityType) {
    try {
        const response = await fetch(`/comments/${entityType}/${entityId}`);
        if (response.ok) {
            this.activeComments = await response.json();
        }
    } catch (e) { console.error(e); }
},

// Универсальная отправка
async postComment(entityId, entityType, parentId = null) {
    const text = parentId ? this.replyText : this.newCommentText;
    if (!text.trim()) return;

    try {
        const response = await fetch('/comments', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                entity_id: entityId,
                entity_type: entityType,
                parent_id: parentId,
                content: text
            })
        });
        if (response.ok) {
            this.newCommentText = '';
            this.replyText = '';
            this.replyToId = null;
            this.loadComments(entityId, entityType);
            // Обновляем ленту активности проекта
            if (this.currentProject) {
                this.loadProjectActivities(this.currentProject.id);
            }
        }
    } catch (e) { alert("Ошибка отправки"); }
},

// Удаление комментария (Soft Delete)
async deleteComment(commentId) {
    if (!confirm("Вы действительно хотите удалить этот комментарий?")) return;

    try {
        const response = await fetch(`/comments/${commentId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            // Если мы внутри задачи — обновляем комменты задачи, если в проекте — проекта
            const entityType = this.currentProjectView === 'discussion' ? 'project' : 'task';
            const entityId = entityType === 'project' ? this.currentProject.id : this.expandedTaskId;
            // Обновляем список в текущей открытой задаче
            this.loadComments(this.expandedTaskId, 'task');
            // Обновляем ленту активности
            if (this.currentProject) {
                this.loadProjectActivities(this.currentProject.id);
            }
        }
    } catch (e) {
        console.error("Ошибка при удалении:", e);
    }
},

getTaskFormattedKey(entityId, entityType) {
    if (entityType !== 'task') return ''; // Если это не задача, не выводим ключ задачи
    
    // Ищем задачу среди загруженных задач проекта
    const task = this.projectTasks.find(t => t.id === entityId);
    
    if (task) {
        const pKey = this.getProjectKey(task.project_id);
        return `${pKey}-${task.task_num}`;
    }
    
    // Если задача не найдена в списке (например, архивирована), 
    // возвращаем заглушку, чтобы не показывать системный ID
    return "Задача"; 
},

async loadActivities() {
    try {
        const response = await fetch('/activities', { credentials: 'include' });
        if (response.ok) {
            this.activities = await response.json();
        }
    } catch (e) { console.error(e); }
},

formatRelativeTime(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString('ru-RU', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' });
},

// Фильтр пользователей для конкретной строки поиска
filterWizardUsers(text) {
    if (!text || text.length < 1) return [];
    const q = text.toLowerCase();
    const currentUserId = Number(localStorage.getItem('currentUserId'));

    return this.users.filter(u => {
        const email = (u.email || "").toLowerCase();
        const name = (u.full_name || "").toLowerCase();

        // 1. Проверяем, начинается ли ПОЧТА на введённые символы
        const emailMatch = email.startsWith(q);

        // 2. Проверяем, начинается ли ИМЯ или ФАМИЛИЯ на эти символы
        // (разбиваем на слова, чтобы "М" находило и "Максим", и "Иванов Максим")
        const nameMatch = name.split(' ').some(word => word.startsWith(q));

        // Исключаем себя
        const notMe = u.id !== currentUserId;

        return (emailMatch || nameMatch) && notMe;
    }).slice(0, 5);
},
// Срабатывает при клике на подсказку
selectUserForRow(row, user) {
    row.user_id = user.id;          // записываем ID для сервера
    row.searchText = user.full_name; // записываем имя в инпут
    row.showSuggestions = false;    // закрываем подсказки
},
addMemberRow() {
    this.kanbanWizard.memberRows.push({ 
        user_id: null, 
        role: 'researcher', 
        searchText: '', 
        showSuggestions: false 
    });
},
getProjectHypothesis(projectId) {
    if (!projectId) return "Выберите проект";
    const p = this.projects.find(proj => proj.id === Number(projectId));
    return p && p.main_hypothesis ? p.main_hypothesis : "Гипотеза исследования не задана для этого проекта";
},
async openTaskFromNavigator(task) {
    // 1. Находим объект проекта, к которому относится задача
    const project = this.projects.find(p => p.id === task.project_id);

    if (project) {
        // 2. Закрываем окно поиска
        this.isSearchOpen = false;
        this.searchQuery = ""; // Очищаем поиск, чтобы не мешал

        // 3. Вызываем стандартный метод открытия проекта (загрузит задачи, спринты и т.д.)
        this.openProject(project);

        // 4. Принудительно ставим вид "Список" (в нем удобнее всего смотреть детали)
        this.currentProjectView = 'list';

        // 5. Устанавливаем ID развернутой задачи
        this.expandedTaskId = task.id;

        // 6. Небольшая задержка, чтобы Vue успел отрендерить список, и прокрутка к задаче
        setTimeout(() => {
            const element = document.querySelector('.table-active');
            if (element) {
                element.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }, 600); 

    } else {
        alert("Проект, к которому относится задача, не найден или недоступен.");
    }
},
// 1. Загрузка фильтров из памяти браузера
storageKey(base) {
    const userId = localStorage.getItem('currentUserId') || 'anonymous';
    return `scifi_${base}_${userId}`;
},
loadSavedFilters() {
    try {
        const saved = localStorage.getItem(this.storageKey('saved_filters'));
        this.userFilters = saved ? JSON.parse(saved) : [];
    } catch (e) {
        console.error("Ошибка загрузки фильтров:", e);
        this.userFilters = [];
    }
},
// 2. Применение выбранного фильтра
applySavedFilter(f) {
    // Заполняем поля навигатора данными из сохраненного фильтра
    this.navFilters.query = f.query || '';
    this.navFilters.project = f.project || 'all';
    this.navFilters.assignee = f.assignee || 'all';
    this.navFilters.status = f.status || 'all';
    this.navFilters.priority = f.priority || 'all';
    
    // Переключаемся на страницу результатов поиска
    this.filterSubPage = 'search'; 
},
// 3. Удаление фильтра
deleteFilter(id) {
    if (!confirm("Удалить этот фильтр?")) return;
    this.userFilters = this.userFilters.filter(f => f.id !== id);
    localStorage.setItem(this.storageKey('saved_filters'), JSON.stringify(this.userFilters));
},
// Метод для кнопки "Просмотреть все задачи" в поиске
openIssueNavigator() {
    this.activeSidebarTab = 'filters';
    this.filterSubPage = 'search';
    this.currentProject = null;
    this.isSearchOpen = false;
    this.loadAllTasks();
},

// Загрузка вообще всех задач для навигатора
async loadAllTasks() {
    try {
        const response = await fetch('/user/tasks');
        if (response.ok) {
            this.allTasks = await response.json();
            console.log("Глобальные задачи загружены для поиска:", this.allTasks.length);
        }
    } catch (e) {
        console.error("Ошибка загрузки задач для поиска:", e);
    }
},

// Сброс фильтров навигатора
resetNavFilters() {
    this.navFilters = { query: '', project: 'all', assignee: 'all', status: 'all', priority: 'all' };
},

// Сохранение фильтра из навигатора
saveFilterFromNav() {
    const name = prompt("Назовите ваш фильтр:");
    if (!name) return;
    
    const newFilter = {
        id: Date.now(),
        name: name,
        ...this.navFilters
    };
    
    let saved = JSON.parse(localStorage.getItem(this.storageKey('saved_filters')) || '[]');
    saved.push(newFilter);
    localStorage.setItem(this.storageKey('saved_filters'), JSON.stringify(saved));
    this.userFilters = saved;
    alert("Фильтр сохранен!");
},
// Переключить статус (добавить/удалить из избранного)
toggleFlag(item, type) {
    let flagged = JSON.parse(localStorage.getItem(this.storageKey('flagged')) || '[]');
    
    const index = flagged.findIndex(f => f.id === item.id && f.type === type);
    
    if (index > -1) {
        // Если уже есть — удаляем
        flagged.splice(index, 1);
    } else {
        // Если нет — добавляем
        flagged.push({
            id: item.id,
            title: item.title || item.name,
            key: item.key || (this.currentProject ? this.currentProject.key : 'TASK'),
            task_num: item.task_num || null,
            type: type
        });
    }
    
    localStorage.setItem(this.storageKey('flagged'), JSON.stringify(flagged));
    this.flaggedItems = flagged;
},

// Проверить, отмечен ли элемент (для покраски звездочки)
isFlagged(id, type) {
    return this.flaggedItems.some(f => f.id === id && f.type === type);
},

// Загрузка из памяти
loadFlaggedFromStorage() {
    this.flaggedItems = JSON.parse(localStorage.getItem(this.storageKey('flagged')) || '[]');
},
openProjectById(id) {
    const project = this.projects.find(p => p.id === id);
    if (project) {
        this.openProject(project);
    } else {
        alert("Загрузка проекта...");
    }
},

async openTaskById(item) {
    // 1. Находим проект, к которому относится задача
    const project = this.projects.find(p => p.key === item.key);
    if (project) {
        await this.openProject(project); // Открываем проект
        this.currentProjectView = 'list'; // Переключаемся на список
        this.expandedTaskId = item.id;   // Разворачиваем задачу
        
        // Скроллим к задаче через небольшую паузу
        setTimeout(() => {
            const el = document.querySelector('.table-active');
            if (el) el.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }, 300);
    }
},
// Сохранить элемент в список недавних
addToRecent(item, type) {
    let recent = JSON.parse(localStorage.getItem(this.storageKey('recent')) || '[]');
    
    // Создаем компактный объект для хранения
    const entry = {
        id: item.id,
        title: item.title || item.name,
        key: item.key || (this.currentProject ? this.currentProject.key : 'TASK'),
        task_num: item.task_num || null,
        type: type, // 'project' или 'task'
        timestamp: new Date().getTime()
    };

    // Удаляем дубликат, если он уже был в списке (чтобы поднять его наверх)
    recent = recent.filter(r => !(r.id === entry.id && r.type === entry.type));
    
    // Добавляем в начало
    recent.unshift(entry);
    
    // Оставляем только последние 10
    recent = recent.slice(0, 10);
    
    localStorage.setItem(this.storageKey('recent'), JSON.stringify(recent));
    this.recentItems = recent;
},

// Загрузить список недавних из памяти
loadRecentFromStorage() {
    this.recentItems = JSON.parse(localStorage.getItem(this.storageKey('recent')) || '[]');
},
// Функция, которая срабатывает при клике на подсказку (исправляет вашу ошибку)
selectUserForExistingTeam(user) {
    this.manageInviteEmail = user.email;
    this.showManageSuggestions = false;
    this.addMemberToExistingTeam(); // Сразу пытаемся добавить
},

// Сама логика отправки запроса на сервер
async addMemberToExistingTeam() {
    const email = this.manageInviteEmail.trim();
    if (!email) return;

    try {
        const response = await fetch(`/teams/${this.selectedTeam.id}/members`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: email })
        });

        if (response.ok) {
            const newMember = await response.json();
            // Добавляем нового участника в локальный список, чтобы он сразу появился в окне
            this.currentTeamMembers.push(newMember);
            this.manageInviteEmail = '';
            this.showManageSuggestions = false;
            
            // Обновляем счетчик участников на главной карточке команды
            if (this.selectedTeam) {
                this.selectedTeam.members_count++;
            }
        } else {
            const err = await response.text();
            alert("Не удалось добавить: " + err);
        }
    } catch (e) {
        console.error("Ошибка добавления участника:", e);
    }
},

// Метод для изменения роли участника (если будете использовать)
async changeMemberRole(member) {
    try {
        await fetch(`/teams/${this.selectedTeam.id}/members/${member.user_id}/role`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ role: member.role })
        });
    } catch (e) {
        console.error("Ошибка смены роли:", e);
    }
},
// 1. Метод удаления участника (исправляет вашу ошибку)
async removeMember(member) {
    if (!confirm(`Удалить участника ${member.full_name} из команды?`)) return;

    try {
        const response = await fetch(`/teams/${this.selectedTeam.id}/members/${member.user_id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            // Удаляем участника из списка в открытом окне
            this.currentTeamMembers = this.currentTeamMembers.filter(m => m.user_id !== member.user_id);
            
            // Уменьшаем счетчик участников на карточке команды
            if (this.selectedTeam) {
                this.selectedTeam.members_count--;
            }
            console.log("Участник удален");
        } else {
            alert("Ошибка сервера при удалении участника");
        }
    } catch (e) {
        console.error("Ошибка сети:", e);
    }
},

// 2. Метод удаления всей команды (чтобы кнопка внизу тоже работала)
async deleteTeam() {
    if (!confirm(`ВНИМАНИЕ: Вы действительно хотите полностью удалить команду "${this.selectedTeam.name}"? Это действие необратимо.`)) return;

    try {
        const response = await fetch(`/teams/${this.selectedTeam.id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            // Удаляем команду из общего списка на странице
            this.teams = this.teams.filter(t => t.id !== this.selectedTeam.id);
            
            // Закрываем модальное окно
            this.manageTeamModal.hide();
            alert("Команда успешно удалена");
        } else {
            alert("Не удалось удалить команду");
        }
    } catch (e) {
        console.error("Ошибка при удалении команды:", e);
    }
},

// 3. Метод обновления инфо о команде (для кнопки "Сохранить изменения")
async updateTeamInfo() {
    try {
        const response = await fetch(`/teams/${this.selectedTeam.id}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                name: this.selectedTeam.name,
                description: this.selectedTeam.description
            })
        });

        if (response.ok) {
            alert("Информация обновлена!");
            this.loadUserTeams(); // Перезагружаем список, чтобы обновить имя на карточке
        }
    } catch (e) {
        console.error(e);
    }
},
selectUser(user) {
    // Добавляем почту пользователя, если её еще нет в списке
    if (!this.newTeam.member_emails.includes(user.email)) {
        this.newTeam.member_emails.push(user.email);
    }
    this.inviteEmail = ''; // Очищаем поле ввода
    this.showSuggestions = false; // Скрываем подсказки
},
addEmailToTeam() {
    const email = this.inviteEmail.trim().toLowerCase();
    // Простая валидация и проверка на дубликаты
    if (email && email.includes('@') && !this.newTeam.member_emails.includes(email)) {
        this.newTeam.member_emails.push(email);
        this.inviteEmail = ''; // Очищаем поле
    } else if (email && !email.includes('@')) {
        alert("Пожалуйста, введите корректный email.");
    }
},
removeEmailFromTeam(email) {
    this.newTeam.member_emails = this.newTeam.member_emails.filter(e => e !== email);
},
async createTeam() {
    if (!this.newTeam.name) {
        alert("Введите название команды");
        return;
    }

    try {
        const response = await fetch('/teams', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.newTeam)
        });

        if (response.ok) {
            const createdTeam = await response.json();
            // Добавляем создателя (тебя) в счетчик участников для красоты
            createdTeam.members_count = 1;
            this.teams.unshift(createdTeam);
            
            this.createTeamModal.hide();
            this.newTeam = { name: '', description: '', member_emails: [] };
            alert("Команда создана!");
        }
    } catch (e) {
        console.error(e);
        alert("Ошибка при создании команды");
    }
},
selectSidebarTab(tabName) {
    this.activeSidebarTab = tabName;
    this.currentProject = null;
    this.currentPage = 'sections';
    this.filterSubPage = 'list';
    this.currentGrant = null;

    if (tabName === 'activity') this.loadActivities();
    
    if (tabName === 'filters') {
        this.loadSavedFilters();
        this.loadAllTasks();
    }
    if (tabName === 'teams') this.loadUserTeams();
    if (tabName === 'grants') this.loadGrants();
},
async loadUserTeams() {
    const currentUserId = Number(localStorage.getItem('currentUserId'));
    try {
        const response = await fetch(`/user/${currentUserId}/teams`);
        if (response.ok) {
            this.teams = await response.json() || [];
        }
    } catch (e) {
        console.error("Ошибка загрузки команд:", e);
    }
},
closeSearch() {
    this.isSearchOpen = false;
    this.showSuggestions = false;
},
startProjectCreation() {
    // Находим основной шаблон (теперь он может называться "Базовый" или "Стандартный")
    const mainTemplate = this.templates.find(t => t.id === 'kanban');
    
    this.kanbanWizard.step = 1;
    this.kanbanWizard.name = '';
    this.kanbanWizard.key = '';
    this.kanbanWizard.memberRows = [{ user_id: null, role: 'researcher' }];
    this.loadUserTeams();
    
    this.selectedTemplate = mainTemplate;
    this.currentPage = "kanbanSetup";
},

                openTemplatesPage() {
                    this.currentPage = "templates";
                },
                goBackToSections() {
                    this.currentPage = "sections";
                },
                goBackToTemplates() {
                    this.currentPage = "templates";
                    this.kanbanWizard.step = 1;
                },
async openProject(project) {
    this.addToRecent(project, 'project');
    console.log("Данные открытого проекта:", project);
    this.currentProject = project;
    this.currentProjectView = "board";
    
    // Сбрасываем данные перед загрузкой новых
    this.projectTasks = [];
    this.projectAttachments = []; 
    this.sprints = [];
    this.projectGrants = [];

        // Получаем список участников проекта, чтобы узнать свою роль
    try {
        const response = await fetch(`/projects/${project.id}/members`);
        if (response.ok) {
            const members = await response.json();
            const currentUserId = Number(localStorage.getItem('currentUserId'));
            const myMemberEntry = members.find(m => m.user_id === currentUserId);
            
            // Сохраняем роль (если не нашли - значит viewer по умолчанию)
            this.currentUserProjectRole = myMemberEntry ? myMemberEntry.role_name || myMemberEntry.role : 'viewer';
            console.log("Моя роль в проекте:", this.currentUserProjectRole);
        }
    } catch (e) {
        console.error("Не удалось получить роль в проекте", e);
    }

    this.loadProjectTasks(project.id);
    this.loadProjectSprints(project.id);
    this.loadProjectAttachments(project.id);
    this.loadProjectAssignableUsers(project.id);
    this.loadProjectGrants(project.id);
    this.loadProjectPermissions(project.id);
    this.loadProjectActivities(project.id);
},
async approveResearchTask(task) {
    if (!confirm("Вы подтверждаете научную достоверность результатов этого этапа?")) return;
    
    try {
        const response = await fetch(`/tasks/${task.id}/status`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ status: 'done' }) // Переводим в "ГОТОВО"
        });
        
        if (response.ok) {
            alert("Результаты исследования утверждены.");
            this.loadProjectTasks(this.currentProject.id); // Обновляем список
        }
    } catch (e) {
        alert("Ошибка при сохранении утверждения");
    }
},
                getSprintTasks(sprintId) {
    // Если массив задач пуст, возвращаем пустой список
                    if (!this.projectTasks) return [];
                    return this.projectTasks.filter(t => t.sprint_id === sprintId);
                },
                // Метод открытия модалки
openStartSprintModal(sprint) {
    this.sprintForm.id = sprint.id;
    this.sprintForm.name = sprint.name;
    this.sprintForm.goal = '';
    
    // Устанавливаем текущую дату и время для поля "Дата начала"
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    this.sprintForm.startDate = now.toISOString().slice(0, 16);
    
    // Считаем количество задач в этом спринте
    this.sprintTaskCount = this.getSprintTasks(sprint.id).length;
    
    // Рассчитываем дату окончания
    this.calculateSprintEndDate();
    
    if (!this.startSprintModal) {
        this.startSprintModal = new bootstrap.Modal(document.getElementById('startSprintModal'));
    }
    this.startSprintModal.show();
},

// Метод расчета даты окончания (чистая логика)
calculateSprintEndDate() {
    if (this.sprintForm.duration === 'custom' || !this.sprintForm.startDate) return;
    
    let start = new Date(this.sprintForm.startDate);
    let weeks = parseInt(this.sprintForm.duration);
    start.setDate(start.getDate() + (weeks * 7));
    
    // Форматируем для input datetime-local
    this.sprintForm.endDate = start.toISOString().slice(0, 16);
},

// РЕАЛЬНЫЙ запрос на сервер для старта спринта
async confirmStartSprint() {
    try {
        const response = await fetch(`/sprints/${this.sprintForm.id}/start`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                name: this.sprintForm.name,
                start_date: this.sprintForm.startDate,
                end_date: this.sprintForm.endDate,
                goal: this.sprintForm.goal
            })
        });
        

        if (response.ok) {
            this.startSprintModal.hide();
            // Перезагружаем спринты, чтобы увидеть смену статуса на "Активен"
            this.loadProjectSprints(this.currentProject.id);
            alert("Спринт запущен!");
        } else {
            alert("Ошибка сервера при запуске спринта");
        }
    } catch (err) {
        console.error("Ошибка сети:", err);
    }
},
// 1. Красивое форматирование даты для шапки (например, "27 мая")
formatDateShort(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
},

// 2. Метод завершения спринта
async completeSprint(sprint) {
    const confirmMsg = `Вы уверены, что хотите завершить спринт "${sprint.name}"?\nНезавершенные задачи будут возвращены в бэклог.`;
    if (!confirm(confirmMsg)) return;

    try {
        const response = await fetch(`/sprints/${sprint.id}/complete`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' }
        });

        if (response.ok) {
            alert("Спринт успешно завершен!");
            // Перезагружаем данные проекта
            if (this.currentProject) {
                await this.loadProjectSprints(this.currentProject.id);
                await this.loadProjectTasks(this.currentProject.id);
            }
        } else {
            alert("Не удалось завершить спринт на сервере.");
        }
    } catch (err) {
        console.error("Ошибка сети:", err);
        alert("Ошибка связи с сервером.");
    }
},
                getBacklogTasks() {
                    if (!this.projectTasks) return [];
    // В бэклог идут задачи, у которых нет sprint_id (null или 0)
                    return this.projectTasks.filter(t => !t.sprint_id || t.sprint_id === 0);
                },

                getStatusBadgeClass(status) {
                    if (!status) return 'badge bg-secondary';
                    const s = status.toLowerCase();
                    if (s === 'todo' || s === 'к выполнению') return 'badge bg-primary';
                    if (s === 'in_progress' || s === 'в работе') return 'badge bg-warning text-dark';
                    if (s === 'review' || s === 'на ревью' || s === 'на проверке') return 'badge bg-info text-dark';
                    if (s === 'done' || s === 'готово') return 'badge bg-success';
                    return 'badge bg-secondary';
                },

                async loadProjectSprints(projectId) {
                    try {
                    const response = await fetch(`/projects/${projectId}/sprints`);
                    if (response.ok) {
                        this.sprints = await response.json();
                    } else {
                        this.sprints = []; // Если ошибка, обнуляем
                    }
                    } catch (err) {
                        console.error('Ошибка загрузки спринтов:', err);
                        this.sprints = [];
                    }
                },
async loadProjectAttachments(projectId) {
    try {
        const response = await fetch(`/projects/${projectId}/attachments`);
        if (response.ok) {
            const data = await response.json();
            // Если сервер прислал null, ставим пустой массив, чтобы .length не ломался
            this.projectAttachments = data || [];
        } else {
            this.projectAttachments = [];
        }
    } catch (err) {
        console.error('Ошибка загрузки файлов:', err);
        this.projectAttachments = [];
    }
},
isImage(fileName) {
    if (!fileName) return false;
    const ext = fileName.toLowerCase().split('.').pop();
    return ['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'].includes(ext);
},
handleTaskFileSelect(event) {
    this.selectedFile = event.target.files[0];
},
                goBackToProjects() {
                    this.currentProject = null;
                    this.projectTasks = [];
                    this.projectBoardTasks = {
                        todo: [],
                        inProgress: [],
                        review: [],
                        done: []
                    };
                },
                async loadProjectTasks(projectId) {
                    try {
                        const response = await fetch(`/projects/${projectId}/tasks`, {
                            credentials: 'include'
                        });
                        if (response.ok) {
                            const tasks = await response.json();
                            console.log("Задачи из базы:", tasks);
                            this.projectTasks = tasks || [];
                            this.organizeTasksByStatus();
                        }
                    } catch (err) {
                        console.error('Failed to load project tasks:', err);
                    }
                },
                async loadProjectAssignableUsers(projectId) {
                    try {
                        const response = await fetch(`/projects/${projectId}/assignable-users`, {
                            credentials: 'include'
                        });
                        if (response.ok) {
                            this.projectAssignableUsers = await response.json();
                        }
                    } catch (err) {
                        console.error('Failed to load assignable users:', err);
                    }
                },
                async createNewSprint() {
    console.log("Кнопка нажата, проект:", this.currentProject);
    if (!this.currentProject) {
        alert("Ошибка: Проект не выбран");
        return;
    }

    const sprintName = "Спринт " + (this.sprints.length + 1);
    const payload = {
        name: sprintName,
        status: "planned"
    };

    try {
        const response = await fetch(`/projects/${this.currentProject.id}/sprints`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            const newSprint = await response.json();
            this.sprints.push(newSprint);
            console.log("Спринт создан:", newSprint);
        } else {
            const errorText = await response.text();
            console.error("Ошибка сервера:", errorText);
            alert("Сервер вернул ошибку при создании спринта");
        }
    } catch (err) {
        console.error("Ошибка сети:", err);
        alert("Не удалось связаться с сервером");
    }
},
// 1. Когда начинаем тянуть задачу для спринта
startSprintDrag(event, task) {
    this.draggedTask = task;
    event.dataTransfer.effectAllowed = 'move';
},

// 2. Когда отпускаем задачу над спринтом или бэклогом
async onDropToSprint(event, sprintId) {
    if (!this.draggedTask) return;

    const taskId = this.draggedTask.id;
    
    // Оптимистичное обновление UI (сразу переносим задачу в списке)
    this.draggedTask.sprint_id = sprintId;

    try {
        const response = await fetch(`/tasks/${taskId}/sprint`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ sprint_id: sprintId })
        });

        if (!response.ok) {
            throw new Error("Ошибка сервера");
        }
        
        console.log(`Задача ${taskId} успешно перенесена в спринт ${sprintId}`);
    } catch (err) {
        console.error(err);
        alert("Не удалось сохранить перемещение. Обновите страницу.");
        // В случае ошибки можно перезагрузить задачи
        this.loadProjectTasks(this.currentProject.id);
    } finally {
        this.draggedTask = null;
    }
},
organizeTasksByStatus() {
    const taskSource = (this.currentProject ? this.projectTasks : this.allTasks) || [];
    if (!this.projectTasks) return;

    const currentUserId = Number(localStorage.getItem('currentUserId'));

    // Применяем фильтры
    const filtered = taskSource.filter(task => {
        const q = this.searchQuery.toLowerCase();
        const matchesSearch = !this.searchQuery || 
            task.title.toLowerCase().includes(q) || 
            (task.description && task.description.toLowerCase().includes(q)) ||
            (task.tags && task.tags.toLowerCase().includes(q));

        return matchesSearch; 
    });

    // 2. Логика для ДОСКИ (Board)
    const activeSprint = this.sprints.find(s => s.status === 'active');
    
    // Если есть активный спринт — на доске только его задачи. 
    // Если активного спринта НЕТ — показываем все задачи (чтобы они не пропадали!)
    const boardTasksSource = activeSprint 
        ? filtered.filter(t => t.sprint_id === activeSprint.id)
        : filtered;

    this.projectBoardTasks = {
        todo: boardTasksSource.filter(t => ['todo', 'к выполнению'].includes((t.status || '').toLowerCase())),
        inProgress: boardTasksSource.filter(t => ['in_progress', 'в работе'].includes((t.status || '').toLowerCase())),
        review: boardTasksSource.filter(t => ['review', 'на ревью', 'на проверке'].includes((t.status || '').toLowerCase())),
        done: boardTasksSource.filter(t => ['done', 'готово'].includes((t.status || '').toLowerCase()))
    };
    
    this.displayedTasks = filtered; 
},

// Добавь также метод сброса:
resetFilters() {
    this.searchQuery = '';
    this.filterMyTasks = false;
    this.filterPriority = 'all';
    this.organizeTasksByStatus();
},
                startDrag(event, task, status) {
                    this.draggedTask = task;
                    this.draggedFromStatus = status;
                    event.dataTransfer.effectAllowed = 'move';
                    event.target.style.opacity = '0.5';
                },
                endDrag(event) {
                    event.target.style.opacity = '1';
                },
                async moveTask(event, newStatus) {
                    event.preventDefault();
                    
                    if (!this.draggedTask) return;
                    
                    // Если задача перемещена в тот же статус, ничего не делаем
                    if (this.draggedFromStatus === newStatus) {
                        this.draggedTask = null;
                        this.draggedFromStatus = null;
                        return;
                    }
                    
                    // Оптимистичное обновление UI
                    const task = this.draggedTask;
                    const oldStatus = this.draggedFromStatus;
                    
                    // Удаляем задачу из старого статуса
                    if (oldStatus === 'todo') {
                        this.projectBoardTasks.todo = this.projectBoardTasks.todo.filter(t => t.id !== task.id);
                    } else if (oldStatus === 'in_progress') {
                        this.projectBoardTasks.inProgress = this.projectBoardTasks.inProgress.filter(t => t.id !== task.id);
                    } else if (oldStatus === 'review') { // Добавлено
                        this.projectBoardTasks.review = this.projectBoardTasks.review.filter(t => t.id !== task.id);
                    } else if (oldStatus === 'done') {
                        this.projectBoardTasks.done = this.projectBoardTasks.done.filter(t => t.id !== task.id);
                    }
                    
                    // Обновляем статус задачи
                    task.status = newStatus;
                    
                    // И добавления в новый:
                    if (newStatus === 'todo') {
                        this.projectBoardTasks.todo.push(task);
                    } else if (newStatus === 'in_progress') {
                        this.projectBoardTasks.inProgress.push(task);
                    } else if (newStatus === 'review') {
                        this.projectBoardTasks.review.push(task);
                    } else if (newStatus === 'done') {
                        this.projectBoardTasks.done.push(task);
                    }
                    
                    // Отправляем обновление на сервер
                    try {
                        const response = await fetch(`/tasks/${task.id}/status`, {
                            method: 'PATCH',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify({ status: newStatus }),
                            credentials: 'include'
                        });
                        
                        if (!response.ok) {
                            console.error('Ошибка при обновлении статуса задачи');
                            // Откатываем изменение UI при ошибке
                            this.loadProjectTasks(this.currentProject.id);
                        }
                    } catch (error) {
                        console.error('Ошибка сети:', error);
                        // Откатываем изменение UI при ошибке
                        this.loadProjectTasks(this.currentProject.id);
                    }
                    
                    this.draggedTask = null;
                    this.draggedFromStatus = null;
                },
                getUserName(userId) {
                    const user = this.users.find(u => u.id === userId);
                    return user ? user.full_name : 'Unknown';
                },
                getTeamName(teamId) {
                    const team = this.teams.find(t => t.id === teamId);
                    return team ? team.name : '';
                },
                selectTemplate(templateItem) {
                    this.selectedTemplate = templateItem;
                    this.kanbanWizard.step = 1;
                    this.kanbanWizard.name = '';
                    this.kanbanWizard.key = '';
                    this.kanbanWizard.memberRows = [{ user_id: null, role: 'researcher' }];
                    this.loadUserTeams();
                    this.loadProjectRoles();
                    this.currentPage = "kanbanSetup";
                },
updateKanbanKey() {
    // 1. Проверяем, что имя вообще введено
    if (!this.kanbanWizard.name) {
        this.kanbanWizard.key = "";
        return;
    }
    
    const name = this.kanbanWizard.name.trim();
    if (!name) {
        this.kanbanWizard.key = "";
        return;
    }
    
    // 2. Генерируем ключ из первых букв каждого слова
    const words = name.split(/\s+/).filter(w => w.length > 0);
    let key = '';
    
    for (const word of words) {
        if (key.length >= 10) break;
        key += word[0]; // Берем первую букву слова
    }
    
    // 3. ПРАВКА ТУТ: Оставляем заглавную латиницу (A-Z), кириллицу (А-Я) и цифры
    key = key.toUpperCase().replace(/[^A-ZА-Я0-9]/g, "").slice(0, 10);
    
    // 4. Проверяем, что ключ начинается с буквы (любой: русской или английской)
    if (!key || !/^[A-ZА-Я]/.test(key)) {
        this.kanbanWizard.key = ""; 
        return;
    }
    
    this.kanbanWizard.key = key;
},
nextKanbanStep() {
    const name = this.kanbanWizard.name ? this.kanbanWizard.name.trim() : "";
    if (!name || name.length < 2) {
        alert("Пожалуйста, введите название проекта (минимум 2 символа).");
        return;
    }
    
    if (!this.kanbanWizard.key || this.kanbanWizard.key.length < 1) {
        alert("Не удалось сгенерировать ключ. Попробуйте ввести название латиницей или введите ключ вручную.");
        return;
    }

    if (this.kanbanWizard.managementType === 'team' && (!this.kanbanWizard.teamId || this.kanbanWizard.teamId <= 0)) {
        alert("Пожалуйста, выберите команду для проекта.");
        return;
    }
    
    this.kanbanWizard.step = 2;
},
                onWizardTeamChange() {
                    if (this.kanbanWizard.teamId) {
                        this.loadTeamMembersForWizard(this.kanbanWizard.teamId);
                    }
                },
                addMemberRow() {
                    this.kanbanWizard.memberRows.push({ user_id: null, role: 'researcher' });
                },
                removeMemberRow(index) {
                    if (this.kanbanWizard.memberRows.length > 1) {
                        this.kanbanWizard.memberRows.splice(index, 1);
                    }
                },
                getTeamMembersPreview(teamId) {
                    return this.wizardTeamMembers || [];
                },
                async loadTeamMembersForWizard(teamId) {
                    if (!teamId) {
                        this.wizardTeamMembers = [];
                        return;
                    }
                    try {
                        const response = await fetch(`/teams/${teamId}/members`);
                        if (response.ok) {
                            this.wizardTeamMembers = await response.json();
                        }
                    } catch (err) {
                        console.error('Failed to load team members:', err);
                    }
                },
async finishKanbanSetup() {
    this.isSubmitting = true;
    
    // 1. Валидация ключа
    if (!this.kanbanWizard.key || this.kanbanWizard.key.length < 2 || this.kanbanWizard.key.length > 10) {
        alert("Ошибка: Ключ должен содержать от 2 до 10 символов.");
        this.isSubmitting = false;
        return;
    }

    const currentUserId = Number(localStorage.getItem('currentUserId') || 1);

    // 2. Подготовка данных
    const projectPayload = {
        name: this.kanbanWizard.name.trim(),
        key: this.kanbanWizard.key.toUpperCase(),
        description: this.kanbanWizard.main_hypothesis
            ? `Научный проект. Основная гипотеза: ${this.kanbanWizard.main_hypothesis}`
            : "Научный проект",
        status: "active",
        start_date: new Date().toISOString(),
        created_by: currentUserId,
        research_goal: this.kanbanWizard.research_goal,
        main_hypothesis: this.kanbanWizard.main_hypothesis,
        novelty: this.kanbanWizard.novelty,
        expected_result: this.kanbanWizard.expected_result,
        visibility: this.kanbanWizard.access,
        execution_type: this.kanbanWizard.managementType,
        team_id: this.kanbanWizard.managementType === 'team' ? this.kanbanWizard.teamId : null,
    };

    try {
        // 3. Создаем проект
        const response = await fetch('/projects', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(projectPayload)
        });

        if (!response.ok) {
            const errorMsg = await response.text();
            throw new Error(errorMsg || 'Ошибка при создании проекта');
        }

        const newProject = await response.json();

        // Создатель автоматически получает роль project_lead на бэкенде (CreateProject).
        // Дополнительный вызов не требуется.

        // 4. Если ручное управление — добавляем остальных приглашенных
        if (this.kanbanWizard.managementType === 'manual') {
            for (const row of this.kanbanWizard.memberRows) {
                if (!row.user_id || row.user_id === currentUserId) continue;
                try {
                    await fetch('/project-members', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ 
                            project_id: newProject.id, 
                            user_id: row.user_id, 
                            role: row.role 
                        })
                    });
                } catch (e) {
                    console.error('Не удалось добавить участника:', e);
                }
            }
        }

        // 5. Сброс и выход
        this.kanbanWizard.name = "";
        this.kanbanWizard.key = "";
        this.kanbanWizard.memberRows = [{ user_id: null, role: 'researcher' }];
        this.currentPage = "sections";
        this.loadProjects(); 

    } catch (error) {
        console.error("Ошибка в процессе создания:", error);
        alert("Ошибка: " + error.message);
    } finally {
        this.isSubmitting = false;
    }
},
                openKanbanPage() {
                    window.location.href = '/kanban.html';
                },
applyTheme(theme) {
    if (theme === 'dark') {
        document.body.classList.add('dark-theme');
    } else {
        document.body.classList.remove('dark-theme');
    }
    localStorage.setItem('appTheme', theme);

    // ДОБАВЬТЕ ЭТИ ПРОВЕРКИ (if):
    const lightBtn = document.getElementById('themeLight');
    const darkBtn = document.getElementById('themeDark');
    
    if (lightBtn) {
        lightBtn.classList.toggle('active', theme === 'light');
    }
    if (darkBtn) {
        darkBtn.classList.toggle('active', theme === 'dark');
    }
},
                formatDateForApi(dateValue) {
                    if (!dateValue) {
                        return null;
                    }
                    return new Date(dateValue + 'T00:00:00').toISOString();
                },
                async createProjectFromModal() {
                    const form = document.getElementById('createProjectForm');
                    this.projectError = "";
                    this.projectSuccess = "";

                    if (!form.checkValidity()) {
                        form.classList.add('was-validated');
                        return;
                    }

                    const payload = {
                        name: this.project.name.trim(),
                        description: this.project.description.trim(),
                        status: this.project.status,
                        start_date: this.formatDateForApi(this.project.start_date),
                        research_goal: this.project.research_goal.trim(),
                        main_hypothesis: this.project.main_hypothesis.trim(),
                        novelty: this.project.novelty.trim(),
                        expected_result: this.project.expected_result.trim(),
                        created_by: Number(localStorage.getItem('currentUserId') || 1),
                        visibility: this.project.visibility || 'closed',
                        execution_type: this.project.execution_type || 'manual',
                        team_id: this.project.execution_type === 'team' ? this.project.team_id : null
                    };

                    this.isSubmitting = true;
                    try {
                        const response = await fetch('/projects', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(payload)
                        });
                        if (!response.ok) {
                            const errorText = await response.text();
                            throw new Error(errorText || 'Не удалось создать проект');
                        }
                        this.projectSuccess = 'Проект успешно создан';
                        this.project = { name: "", description: "", status: "active", start_date: "", research_goal: "", main_hypothesis: "", novelty: "", expected_result: "", execution_type: "manual", team_id: null};
                        form.classList.remove('was-validated');
                        setTimeout(() => {
                            if (this.createProjectModal) {
                                this.createProjectModal.hide();
                            }
                        }, 700);
                    } catch (error) {
                        this.projectError = error.message || 'Ошибка при создании проекта';
                    } finally {
                        this.isSubmitting = false;
                    }
                },
                async logout() {
    try {
        await fetch('/logout'); // Сервер удалит куку
        localStorage.removeItem('currentUserId');
        // replace заменяет текущую запись в истории, "назад" не сработает
        window.location.replace('/login');
    } catch (e) {
        window.location.href = '/login';
    }
},
                async loadCurrentUser() {
                    try {
                        const response = await fetch('/me');
                        if (response.ok) {
                            const user = await response.json();
                            const names = (user.full_name || "").trim().split(' ').filter(Boolean);
                            let initials = '';
                            if (names.length >= 2) {
                                initials = (names[0][0] + names[1][0]).toUpperCase();
                            } else if (names.length === 1) {
                                initials = names[0].substring(0, 2).toUpperCase();
                            }
                            this.profile = {
                                full_name: user.full_name,
                                email: user.email,
                                initials
                            };
                            localStorage.setItem('currentUserId', String(user.id));
                            this.loadRecentFromStorage();
                            this.loadFlaggedFromStorage();
                            this.loadSavedFilters();
                        } else if (response.status === 401) {
                            window.location.replace('/login');
                        }
                    } catch (err) {
                        console.error('Failed to load user profile:', err);
                        window.location.replace('/login');
                    }
                },
                async loadProjects() {
                    try {
                        const response = await fetch('/user/projects', {
                            credentials: 'include'
                        });
                        if (response.ok) {
                            this.projects = await response.json();
                        } else {
                            console.error('Failed to load projects:', response.status, response.statusText);
                        }
                    } catch (err) {
                        console.error('Failed to load projects:', err);
                    }
                },
                async loadUsers() {
                    try {
                        const response = await fetch('/users');
                        if (response.ok) {
                            this.users = await response.json();
                        }
                    } catch (err) {
                        console.error('Failed to load users:', err);
                    }
                },

                formatDate(dateString) {
                    if (!dateString) return '';
                    const date = new Date(dateString);
                    return date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'long', day: 'numeric' });
                },
                assignToCurrentUser() {
                    const userId = Number(localStorage.getItem('currentUserId'));
                    if (userId) {
                        this.newTask.assigneeId = userId;
                    }
                },
resetTaskForm() {
    this.newTask = {
        // --- БАЗОВЫЕ ПОЛЯ ---
        projectId: this.currentProject ? this.currentProject.id : "", // сохраняем проект, если мы внутри него
        type: "research",
        status: "todo",
        title: "",
        description: "",
        assigneeId: null,
        priority: "Medium",
        dueDate: "",
        tags: "",

        // --- ОБЩИЙ НАУЧНЫЙ ПАСПОРТ ---
        researchContribution: "",
        researchMethod: "",
        expectedScientificResult: "",
        conclusion: "",
        doi: "",

        // --- ИССЛЕДОВАНИЕ ---
        resQuestion: "",
        resObject: "",
        resSubject: "",
        resSources: [{ title: '', authors: '', year: '', doi: '' }],
        resPubCount: "",
        resDatabases: "",
        resApproach: "",

        // --- ЭКСПЕРИМЕНТ ---
        expGoal: "",
        expVariables: [{ type: 'Independent', name: '', value: '' }],
        expControls: "",
        expInputData: "",
        expSampleSize: "",
        expMetrics: "",
        expResults: "",
        expStatSignificance: "",

        // --- СБОР ДАННЫХ ---
        dsSourceType: "",
        dsMethod: "",
        dsLink: "",
        dsFormat: "csv",
        dsCompleteness: 100,
        sampleSize: "",

        // --- АНАЛИЗ ---
        anTools: "",
        anVisuals: "",
        statIndicators: "",
        statPValue: 0.05,
        hypConfirmed: false,

        // --- РАЗРАБОТКА ---
        devComponent: "",
        devStack: "",
        devRepo: "",
        devSwagger: "",
        devTestCoverage: 0,
        devVersion: "",

        // --- ПУБЛИКАЦИЯ ---
        docType: "article",
        journal: "",
        docStatus: "draft",
        coAuthors: ""
    };

    // Очистка системных уведомлений
    this.taskError = "";
    this.taskSuccess = "";

    // Сброс поля выбора файла (физически в DOM)
    const fileInput = document.getElementById('taskFile');
    if (fileInput) fileInput.value = '';
    this.selectedFile = null;
},
                splitTags(tags) {
    if (!tags) return [];
    if (Array.isArray(tags)) return tags;
    return tags.split(',').map(t => t.trim()).filter(t => t !== "");
},

openTaskDetails(task) {
    console.log("Открываем детали задачи:", task);
    // Пока просто выводим в консоль
},
getProjectKey(projectId) { 
    if (!this.projects) return 'TASK';
    const p = this.projects.find(proj => proj.id === projectId);
    return p && p.key ? p.key.toUpperCase() : 'TASK';
},
toggleTask(taskId) {
    this.expandedTaskId = this.expandedTaskId === taskId ? null : taskId;
    if (this.expandedTaskId) {
        const task = this.projectTasks.find(t => t.id === taskId);
        if (task) {
            this.addToRecent(task, 'task');
            this.loadComments(taskId, 'task'); 
        }
    }
},
prepareCreateTask() {
    this.resetTaskForm(); // Сначала полностью очищаем форму
    
    // Если сейчас открыт какой-то проект (currentProject не null)
    if (this.currentProject) {
        this.newTask.projectId = this.currentProject.id;
    }
    
    // Теперь показываем модалку
    this.createTaskModal.show();
},
async createTaskFromModal() {
    const form = document.getElementById('createTaskForm');
    this.taskError = "";
    this.taskSuccess = "";

    if (!form.checkValidity()) {
        form.classList.add('was-validated');
        return;
    }

    const currentUserId = Number(localStorage.getItem('currentUserId') || 1);

    // 1. Упаковываем ВСЕ специфические поля в один объект (для колонки JSONB "parameters")
    const researchParams = {
        // Общее
        expected_result: this.newTask.expectedScientificResult,
        
        // Исследование
        res_question: this.newTask.resQuestion,
        res_object: this.newTask.resObject,
        res_subject: this.newTask.resSubject,
        res_sources: this.newTask.resSources,
        res_pub_count: this.newTask.resPubCount,
        res_databases: this.newTask.resDatabases,
        res_approach: this.newTask.resApproach,
        
        // Эксперимент
        exp_goal: this.newTask.expGoal,
        exp_variables: this.newTask.expVariables,
        exp_controls: this.newTask.expControls,
        exp_input: this.newTask.expInputData,
        exp_sample_size: this.newTask.expSampleSize,
        exp_metrics: this.newTask.expMetrics,
        exp_results: this.newTask.expResults,
        exp_stat_significance: this.newTask.expStatSignificance,
        
        // Сбор данных
        ds_source_type: this.newTask.dsSourceType,
        ds_method: this.newTask.dsMethod,
        ds_link: this.newTask.dsLink,
        ds_format: this.newTask.dsFormat,
        ds_completeness: this.newTask.dsCompleteness,
        sample_size: this.newTask.sampleSize,
        
        // Анализ
        an_tools: this.newTask.anTools,
        an_visuals: this.newTask.anVisuals,
        stat_indicators: this.newTask.statIndicators,
        stat_p_value: this.newTask.statPValue,
        hyp_confirmed: this.newTask.hypConfirmed,
        
        // Разработка
        dev_component: this.newTask.devComponent,
        dev_stack: this.newTask.devStack,
        dev_repo: this.newTask.devRepo,
        dev_swagger: this.newTask.devSwagger,
        dev_test_coverage: this.newTask.devTestCoverage,
        dev_version: this.newTask.devVersion,
        
        // Публикация
        doc_type: this.newTask.docType,
        doc_journal: this.newTask.journal,
        doc_status: this.newTask.docStatus,
        doc_coauthors: this.newTask.coAuthors
    };

    // 2. Формируем итоговый объект для отправки на Go-бэкенд
    const payload = {
        project_id: Number(this.newTask.projectId),
        title: this.newTask.title.trim(),
        description: this.newTask.description.trim(),
        status: this.newTask.status,
        priority: this.newTask.priority,
        type: this.newTask.type, 
        due_date: this.newTask.dueDate || null,
        assignee_id: this.newTask.assigneeId || null,
        created_by: currentUserId,
        tags: this.newTask.tags,
        research_contribution: this.newTask.researchContribution,
        research_method: this.newTask.researchMethod,
        doi: this.newTask.doi,
        conclusion: this.newTask.conclusion || "",
        // Передаем объект — Go-бэкенд сам запишет его в JSONB колонку
        parameters: researchParams, 
        metrics: { goal: this.newTask.expMetrics }
    };

    this.isTaskSubmitting = true;

    try {
        // ШАГ 1: Создаем задачу
        const response = await fetch('/tasks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload), // Используем наш payload
            credentials: 'include'
        });

        if (!response.ok) throw new Error('Не удалось создать задачу');
        
        const createdTask = await response.json();
        console.log("Задача создана успешно, ID:", createdTask.id);

        // ШАГ 2: Загрузка файла (если выбран)
        if (this.selectedFile) {
            const formData = new FormData();
            formData.append('file', this.selectedFile);
            await fetch(`/tasks/${createdTask.id}/attachments`, {
                method: 'POST',
                body: formData
            });
        }

        this.taskSuccess = 'Задача создана успешно';
        this.resetTaskForm();
        
        if (this.currentProject) {
            this.loadProjectTasks(this.currentProject.id);
            this.loadProjectAttachments(this.currentProject.id);
        }

        setTimeout(() => {
            this.createTaskModal.hide();
            form.classList.remove('was-validated');
        }, 1000);

    } catch (error) {
        this.taskError = error.message;
    } finally {
        this.isTaskSubmitting = false;
    }
},

// ==================== GRANTS ====================
async loadGrants() {
    try {
        const params = new URLSearchParams();
        if (this.grantFilters.status) params.append('status', this.grantFilters.status);
        if (this.grantFilters.type) params.append('type', this.grantFilters.type);
        if (this.grantFilters.search) params.append('search', this.grantFilters.search);
        const response = await fetch('/grants?' + params.toString(), { credentials: 'include' });
        if (response.ok) {
            const data = await response.json();
            this.grants = data.grants || [];
        }
    } catch (err) {
        console.error('Failed to load grants:', err);
    }
},
openGrant(grant) {
    this.currentGrant = grant;
    this.loadGrantProjects(grant.id);
    this.loadGrantBudget(grant.id);
},
goBackToGrants() {
    this.currentGrant = null;
},
async loadGrantProjects(grantId) {
    try {
        const response = await fetch(`/grants/${grantId}/projects`, { credentials: 'include' });
        if (response.ok) {
            this.currentGrant.projects = await response.json();
        }
    } catch (err) {
        console.error('Failed to load grant projects:', err);
    }
},
async loadGrantBudget(grantId) {
    try {
        const response = await fetch(`/grants/${grantId}/budget`, { credentials: 'include' });
        if (response.ok) {
            this.currentGrant.budget = await response.json();
        }
    } catch (err) {
        console.error('Failed to load grant budget:', err);
    }
},
prepareCreateGrant() {
    this.grantForm = {
        title: '',
        code: '',
        funding_organization: '',
        country: '',
        description: '',
        scientific_direction: '',
        grant_type: 'state',
        status: 'draft',
        total_amount: 0,
        currency: 'RUB',
        start_date: '',
        end_date: '',
        application_deadline: '',
        principal_investigator_id: null
    };
    this.grantError = '';
    this.grantSuccess = '';
    this.createGrantModal.show();
},
prepareEditGrant(grant) {
    this.grantForm = {
        id: grant.id,
        title: grant.title,
        code: grant.code || '',
        funding_organization: grant.funding_organization,
        country: grant.country || '',
        description: grant.description || '',
        scientific_direction: grant.scientific_direction || '',
        grant_type: grant.grant_type,
        status: grant.status,
        total_amount: grant.total_amount,
        currency: grant.currency || 'RUB',
        start_date: grant.start_date ? grant.start_date.substring(0, 10) : '',
        end_date: grant.end_date ? grant.end_date.substring(0, 10) : '',
        application_deadline: grant.application_deadline ? grant.application_deadline.substring(0, 10) : '',
        principal_investigator_id: grant.principal_investigator_id
    };
    this.grantError = '';
    this.grantSuccess = '';
    this.createGrantModal.show();
},
async createGrantFromModal() {
    const form = document.getElementById('createGrantForm');
    this.grantError = '';
    this.grantSuccess = '';
    if (!form.checkValidity()) {
        form.classList.add('was-validated');
        return;
    }
    const payload = {
        title: this.grantForm.title.trim(),
        code: this.grantForm.code.trim(),
        funding_organization: this.grantForm.funding_organization.trim(),
        country: this.grantForm.country.trim(),
        description: this.grantForm.description.trim(),
        scientific_direction: this.grantForm.scientific_direction.trim(),
        grant_type: this.grantForm.grant_type,
        status: this.grantForm.status,
        total_amount: Number(this.grantForm.total_amount) || 0,
        currency: this.grantForm.currency.trim() || 'RUB',
        start_date: this.grantForm.start_date || null,
        end_date: this.grantForm.end_date || null,
        application_deadline: this.grantForm.application_deadline || null,
        principal_investigator_id: this.grantForm.principal_investigator_id || null
    };
    this.isGrantSubmitting = true;
    try {
        const isEdit = !!this.grantForm.id;
        const url = isEdit ? `/grants/${this.grantForm.id}` : '/grants';
        const method = isEdit ? 'PATCH' : 'POST';
        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
            credentials: 'include'
        });
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Ошибка при сохранении гранта');
        }
        this.grantSuccess = isEdit ? 'Грант обновлен' : 'Грант создан';
        form.classList.remove('was-validated');
        setTimeout(() => {
            this.createGrantModal.hide();
            this.loadGrants();
        }, 800);
    } catch (error) {
        this.grantError = error.message || 'Ошибка';
    } finally {
        this.isGrantSubmitting = false;
    }
},
async deleteGrant(grantId) {
    if (!confirm('Удалить грант? Это также удалит все связанные распределения финансирования.')) return;
    try {
        const response = await fetch(`/grants/${grantId}`, { method: 'DELETE', credentials: 'include' });
        if (!response.ok) throw new Error('Не удалось удалить');
        this.loadGrants();
        if (this.currentGrant && this.currentGrant.id === grantId) {
            this.currentGrant = null;
        }
    } catch (err) {
        alert(err.message);
    }
},
getGrantStatusBadgeClass(status) {
    const map = {
        draft: 'bg-secondary',
        submitted: 'bg-info',
        under_review: 'bg-warning text-dark',
        approved: 'bg-success',
        rejected: 'bg-danger',
        active: 'bg-primary',
        completed: 'bg-success',
        suspended: 'bg-dark'
    };
    return map[status] || 'bg-secondary';
},
getGrantStatusLabel(status) {
    const map = {
        draft: 'Черновик', submitted: 'Подан', under_review: 'На рассмотрении',
        approved: 'Одобрен', rejected: 'Отклонен', active: 'Активен',
        completed: 'Завершен', suspended: 'Приостановлен'
    };
    return map[status] || status;
},
getGrantTypeLabel(type) {
    const map = {
        state: 'Государственный', university: 'Университетский',
        international: 'Международный', corporate: 'Корпоративный', internal: 'Внутренний'
    };
    return map[type] || type;
},
async loadProjectGrants(projectId) {
    try {
        const response = await fetch(`/projects/${projectId}/grants`, { credentials: 'include' });
        if (response.ok) {
            this.projectGrants = await response.json();
        }
    } catch (err) {
        console.error('Failed to load project grants:', err);
    }
},
async loadProjectPermissions(projectId) {
    try {
        const response = await fetch(`/projects/${projectId}/my-permissions`, { credentials: 'include' });
        if (response.ok) {
            const data = await response.json();
            this.projectPermissions = data.permissions || [];
        }
    } catch (err) {
        console.error('Failed to load permissions:', err);
        this.projectPermissions = [];
    }
},
async loadProjectRoles() {
    try {
        const response = await fetch('/roles', { credentials: 'include' });
        if (response.ok) {
            this.projectRoles = await response.json();
        }
    } catch (err) {
        console.error('Failed to load roles:', err);
    }
},
async loadProjectActivities(projectId) {
    try {
        const response = await fetch(`/projects/${projectId}/activities`, { credentials: 'include' });
        if (response.ok) {
            this.projectActivities = await response.json();
        }
    } catch (err) {
        console.error('Failed to load project activities:', err);
    }
},
filterUsersForMember() {
    const q = this.newProjectMemberEmail.toLowerCase().trim();
    if (!q) {
        this.filteredProjectMemberUsers = [];
        return;
    }
    this.filteredProjectMemberUsers = this.users.filter(u => {
        if (this.currentProjectMembers.some(m => m.user_id === u.id)) return false;
        return (u.full_name && u.full_name.toLowerCase().includes(q)) || (u.email && u.email.toLowerCase().includes(q));
    }).slice(0, 5);
},
selectProjectMemberUser(user) {
    this.newProjectMemberEmail = user.email;
    this.filteredProjectMemberUsers = [];
},
async addProjectMemberFromModal() {
    const user = this.users.find(u => u.email === this.newProjectMemberEmail.trim());
    if (!user) {
        alert('Пользователь не найден');
        return;
    }
    try {
        const response = await fetch('/project-members', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ project_id: this.currentProject.id, user_id: user.id, role: this.newProjectMemberRole })
        });
        if (!response.ok) throw new Error('Не удалось добавить участника');
        this.newProjectMemberEmail = '';
        this.newProjectMemberRole = 'researcher';
        this.loadProjectMembers(this.currentProject.id);
    } catch (err) {
        alert(err.message);
    }
},
async assignMemberRole(projectId, userId, roleId) {
    try {
        const response = await fetch(`/projects/${projectId}/members/${userId}/role`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ role_id: roleId }),
            credentials: 'include'
        });
        if (!response.ok) throw new Error('Не удалось назначить роль');
        this.loadProjectMembers(projectId);
    } catch (err) {
        alert(err.message);
    }
},
async removeProjectMember(userId) {
    if (!confirm('Удалить участника из проекта?')) return;
    try {
        const response = await fetch(`/projects/${this.currentProject.id}/members/${userId}`, {
            method: 'DELETE',
            credentials: 'include'
        });
        if (!response.ok) throw new Error('Не удалось удалить участника');
        this.loadProjectMembers(this.currentProject.id);
    } catch (err) {
        alert(err.message);
    }
},
openManageMembers() {
    this.newProjectMemberEmail = '';
    this.newProjectMemberRole = 'researcher';
    this.filteredProjectMemberUsers = [];
    this.loadProjectRoles();
    this.loadProjectMembers(this.currentProject.id);
    this.manageMembersModal.show();
},
getRoleLabel(roleName) {
    const map = {
        project_lead: 'Руководитель проекта',
        scientific_supervisor: 'Научный руководитель',
        researcher: 'Исследователь',
        analyst: 'Аналитик',
        developer: 'Разработчик',
        reviewer: 'Рецензент',
        viewer: 'Наблюдатель',
        admin: 'Администратор',
        user: 'Пользователь',
        guest: 'Гость',
        member: 'Участник'
    };
    return map[roleName] || roleName;
},
async loadProjectMembers(projectId) {
    try {
        const response = await fetch(`/projects/${projectId}/members`, { credentials: 'include' });
        if (response.ok) {
            this.currentProjectMembers = await response.json();
        }
    } catch (err) {
        console.error('Failed to load project members:', err);
    }
},
prepareAddProjectToGrant() {
    this.grantProjectForm = {
        project_id: '',
        allocated_amount: 0,
        funding_purpose: '',
        funding_start_date: '',
        funding_end_date: ''
    };
    this.grantProjectError = '';
    // Загружаем проекты если еще не загружены
    if (!this.projects || this.projects.length === 0) {
        this.loadProjects();
    }
    this.addProjectToGrantModal.show();
},
async createGrantProjectFunding() {
    const form = document.getElementById('addProjectToGrantForm');
    this.grantProjectError = '';
    if (!form.checkValidity()) {
        form.classList.add('was-validated');
        return;
    }
    if (!this.currentGrant || !this.currentGrant.id) return;
    try {
        const response = await fetch(`/grants/${this.currentGrant.id}/projects`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                project_id: Number(this.grantProjectForm.project_id),
                allocated_amount: Number(this.grantProjectForm.allocated_amount) || 0,
                funding_purpose: this.grantProjectForm.funding_purpose || 'Финансирование проекта',
                funding_start_date: this.grantProjectForm.funding_start_date || null,
                funding_end_date: this.grantProjectForm.funding_end_date || null
            }),
            credentials: 'include'
        });
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Ошибка при прикреплении проекта');
        }
        form.classList.remove('was-validated');
        this.addProjectToGrantModal.hide();
        this.loadGrantProjects(this.currentGrant.id);
        this.loadGrantBudget(this.currentGrant.id);
    } catch (err) {
        this.grantProjectError = err.message || 'Ошибка';
    }
},
async removeProjectFromGrant(fundingId) {
    if (!confirm('Удалить связь финансирования?')) return;
    try {
        const response = await fetch(`/grants/${this.currentGrant.id}/projects/${fundingId}`, {
            method: 'DELETE',
            credentials: 'include'
        });
        if (!response.ok) throw new Error('Не удалось удалить');
        this.loadGrantProjects(this.currentGrant.id);
        this.loadGrantBudget(this.currentGrant.id);
    } catch (err) {
        alert(err.message);
    }
},
            },
mounted() {
    this.loadRecentFromStorage();
    this.loadFlaggedFromStorage();
    this.loadSavedFilters();
    const savedTheme = localStorage.getItem('appTheme') || 'light';
    this.applyTheme(savedTheme);

    this.loadCurrentUser();
    this.loadProjects();
    this.loadUsers();
    this.loadAllTasks();

    // Существующие модалки
    const modalElement = document.getElementById('createProjectModal');
    this.createProjectModal = new bootstrap.Modal(modalElement);

    const taskModalElement = document.getElementById('createTaskModal');
    this.createTaskModal = new bootstrap.Modal(taskModalElement);

    const sprintModalElement = document.getElementById('startSprintModal');
    this.startSprintModal = new bootstrap.Modal(sprintModalElement);

    const teamModalElement = document.getElementById('createTeamModal');
    if (teamModalElement) {
        this.createTeamModal = new bootstrap.Modal(teamModalElement);
    }
    this.manageTeamModal = new bootstrap.Modal(document.getElementById('manageTeamModal'));

    const grantModalElement = document.getElementById('createGrantModal');
    if (grantModalElement) {
        this.createGrantModal = new bootstrap.Modal(grantModalElement);
    }

    const addProjectModalElement = document.getElementById('addProjectToGrantModal');
    if (addProjectModalElement) {
        this.addProjectToGrantModal = new bootstrap.Modal(addProjectModalElement);
    }

    const manageMembersModalElement = document.getElementById('manageMembersModal');
    if (manageMembersModalElement) {
        this.manageMembersModal = new bootstrap.Modal(manageMembersModalElement);
    }
}
        });
        app.directive("click-outside", clickOutside);
        app.mount('#app');
