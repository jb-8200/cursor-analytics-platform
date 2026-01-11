package docs

import (
	"html/template"
	"net/http"
)

const (
	validUsername = "dox"
	validPassword = "dox-a3"
	sessionCookie = "docs_session"
)

// Handler handles documentation UI routes
type Handler struct {
	sessions      *SessionManager
	loginTemplate *template.Template
	docsTemplate  *template.Template
}

// NewHandler creates a new docs handler
func NewHandler() *Handler {
	h := &Handler{
		sessions: NewSessionManager(),
	}

	// Parse login template
	loginTmpl, err := template.New("login").Parse(loginHTML)
	if err != nil {
		panic(err)
	}
	h.loginTemplate = loginTmpl

	// Parse docs template
	docsTmpl, err := template.New("docs").Parse(docsHTML)
	if err != nil {
		panic(err)
	}
	h.docsTemplate = docsTmpl

	return h
}

// Index handles GET /docs - the main documentation page
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Check for valid session
	cookie, err := r.Cookie(sessionCookie)
	if err != nil || cookie == nil {
		http.Redirect(w, r, "/docs/login", http.StatusFound)
		return
	}

	username, valid := h.sessions.ValidateSession(cookie.Value)
	if !valid {
		http.Redirect(w, r, "/docs/login", http.StatusFound)
		return
	}

	// Render docs page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := map[string]string{
		"Username": username,
	}
	h.docsTemplate.Execute(w, data)
}

// Login handles GET /docs/login (show form) and POST /docs/login (authenticate)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Show login form
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h.loginTemplate.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		// Handle authentication
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username != validUsername || password != validPassword {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			data := map[string]string{
				"Error": "Invalid username or password",
			}
			h.loginTemplate.Execute(w, data)
			return
		}

		// Create session
		sessionID := h.sessions.CreateSession(username)

		// Set secure cookie
		cookie := &http.Cookie{
			Name:     sessionCookie,
			Value:    sessionID,
			Path:     "/docs",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   8 * 3600, // 8 hours
		}
		http.SetCookie(w, cookie)

		// Redirect to /docs
		http.Redirect(w, r, "/docs", http.StatusFound)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Logout handles GET /docs/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookie)
	if err == nil && cookie != nil {
		h.sessions.DeleteSession(cookie.Value)
	}

	// Clear the cookie
	clearCookie := &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/docs",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, clearCookie)

	// Redirect to login
	http.Redirect(w, r, "/docs/login", http.StatusFound)
}

// Static serves static files (OpenAPI specs, etc.)
func (h *Handler) Static(w http.ResponseWriter, r *http.Request) {
	// Remove /docs prefix from path
	path := r.URL.Path
	if len(path) > 5 && path[:5] == "/docs" {
		path = path[5:]
	}

	// Serve from embedded filesystem
	file, err := staticFiles.Open("static" + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	// Determine content type based on file extension
	if len(path) > 5 && path[len(path)-5:] == ".yaml" {
		w.Header().Set("Content-Type", "application/yaml")
	}

	// Read and serve file
	buf := make([]byte, 1024*1024) // 1MB max
	n, _ := file.Read(buf)
	w.Write(buf[:n])
}

const loginHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cursor Simulator - API Documentation Login</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            color: #e2e8f0;
        }

        .login-container {
            background: #1e293b;
            border: 1px solid #334155;
            border-radius: 8px;
            padding: 40px;
            width: 100%;
            max-width: 400px;
            box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.3);
        }

        .login-header {
            text-align: center;
            margin-bottom: 30px;
        }

        .login-header h1 {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 8px;
            color: #f1f5f9;
        }

        .login-header p {
            font-size: 14px;
            color: #94a3b8;
        }

        .form-group {
            margin-bottom: 20px;
        }

        label {
            display: block;
            font-size: 14px;
            font-weight: 500;
            margin-bottom: 8px;
            color: #cbd5e1;
        }

        input[type="text"],
        input[type="password"] {
            width: 100%;
            padding: 10px 12px;
            border: 1px solid #475569;
            border-radius: 6px;
            background: #0f172a;
            color: #e2e8f0;
            font-size: 14px;
            transition: border-color 0.2s;
        }

        input[type="text"]:focus,
        input[type="password"]:focus {
            outline: none;
            border-color: #3b82f6;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }

        .form-group button {
            width: 100%;
            padding: 10px 12px;
            background: #3b82f6;
            color: white;
            border: none;
            border-radius: 6px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: background-color 0.2s;
        }

        .form-group button:hover {
            background: #2563eb;
        }

        .form-group button:active {
            background: #1d4ed8;
        }

        .error-message {
            background: #7c2d12;
            border: 1px solid #92400e;
            color: #fed7aa;
            padding: 12px;
            border-radius: 6px;
            font-size: 14px;
            margin-bottom: 20px;
            display: none;
        }

        .error-message.show {
            display: block;
        }

        .credentials-hint {
            background: #1e3a8a;
            border: 1px solid #3730a3;
            color: #93c5fd;
            padding: 12px;
            border-radius: 6px;
            font-size: 13px;
            text-align: center;
            margin-top: 20px;
            line-height: 1.5;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <div class="login-header">
            <h1>API Documentation</h1>
            <p>Cursor Simulator Documentation Portal</p>
        </div>

        {{if .Error}}
        <div class="error-message show">
            {{.Error}}
        </div>
        {{end}}

        <form method="POST">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required autofocus>
            </div>

            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required>
            </div>

            <div class="form-group">
                <button type="submit">Sign In</button>
            </div>
        </form>

        <div class="credentials-hint">
            <strong>Demo Credentials</strong><br>
            Username: dox<br>
            Password: dox-a3
        </div>
    </div>
</body>
</html>
`

const docsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation - Cursor Simulator</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@scalar/themes@latest/style.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .header {
            background: #1e293b;
            border-bottom: 1px solid #334155;
            padding: 16px 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .header-title {
            color: #f1f5f9;
            font-size: 18px;
            font-weight: 600;
        }
        .header-user {
            color: #cbd5e1;
            font-size: 14px;
        }
        .header-logout {
            color: #3b82f6;
            text-decoration: none;
            font-size: 14px;
            margin-left: 20px;
            cursor: pointer;
        }
        .header-logout:hover {
            text-decoration: underline;
        }
        .spec-selector {
            margin: 20px;
            display: flex;
            gap: 10px;
            align-items: center;
        }
        .spec-selector label {
            color: #e2e8f0;
            font-weight: 500;
        }
        .spec-selector select {
            padding: 8px 12px;
            background: #0f172a;
            border: 1px solid #475569;
            border-radius: 6px;
            color: #e2e8f0;
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="header-title">Cursor API Documentation</div>
        <div>
            <span class="header-user">{{.Username}}</span>
            <a href="/docs/logout" class="header-logout">Logout</a>
        </div>
    </div>

    <div class="spec-selector">
        <label for="spec-select">API Specification:</label>
        <select id="spec-select" onchange="loadSpec(this.value)">
            <option value="cursor">Cursor API</option>
            <option value="github">GitHub Simulation API</option>
        </select>
    </div>

    <script
        id="api-reference"
        data-url="/docs/openapi/cursor-api.yaml"></script>

    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
        function loadSpec(type) {
            const specUrl = type === 'cursor'
                ? '/docs/openapi/cursor-api.yaml'
                : '/docs/openapi/github-sim-api.yaml';

            // Update the data-url attribute
            const scriptEl = document.getElementById('api-reference');
            scriptEl.setAttribute('data-url', specUrl);

            // Reload the page to apply new spec
            window.location.reload();
        }
    </script>
</body>
</html>
`
