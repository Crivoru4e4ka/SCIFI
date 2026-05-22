        // Как только страница загрузилась, проверяем, нет ли уже активной сессии
async function checkExistingSession() {
    try {
        const response = await fetch('/me');
        if (response.ok) {
            // Если сервер ответил 200 OK, значит мы уже залогинены
            window.location.replace('/'); // Используем replace, чтобы нельзя было нажать "назад"
        }
    } catch (e) {
        // Если ошибка — значит мы не авторизованы, остаемся на странице логина
    }
}
checkExistingSession();
        const { createApp } = Vue;

        createApp({
            data() {
                return {
                    login: {
                        email: "",
                        password: ""
                    },
                    register: {
                        email: "",
                        full_name: "",
                        role: "user",
                        password: ""
                    }
                };
            },
            methods: {
                async submitLogin() {
                    try {
                        const response = await fetch('/auth/login', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(this.login)
                        });
                        if (response.ok) {
                            window.location.href = '/';
                        } else {
                            const error = await response.text();
                            alert('Ошибка: ' + error);
                        }
                    } catch (err) {
                        alert('Ошибка сети: ' + err.message);
                    }
                },
                async submitRegister() {
                    try {
                        const response = await fetch('/auth/register', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json' },
                            body: JSON.stringify(this.register)
                        });
                        if (response.ok) {
                            alert('Регистрация прошла успешно. Войдите в систему.');
                            const loginTab = new bootstrap.Tab(document.getElementById('login-tab'));
                            loginTab.show();
                            this.register = { email: "", full_name: "", role: "user", password: "" };
                        } else {
                            const error = await response.text();
                            alert('Ошибка: ' + error);
                        }
                    } catch (err) {
                        alert('Ошибка сети: ' + err.message);
                    }
                }
            }
        }).mount('#authApp');
