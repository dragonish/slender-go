package model

// WEB_DIR defines web dir.
const WEB_DIR = "web"

// FAVICON_FILE defines favicon file path.
const FAVICON_FILE = WEB_DIR + "/favicon/favicon.ico"

// JS_PATH defines js files path.
const JS_PATH = WEB_DIR + "/assets/js"

// CSS_PATH defines css files paht.
const CSS_PATH = WEB_DIR + "/assets/css"

// TEMPLATES_PATH defines template files path.
const TEMPLATES_PATH = WEB_DIR + "/templates"

// MANAGER_PATH defines manager files path.
const MANAGER_PATH = WEB_DIR + "/manager"

const UPLOAD_FILES_PATH = DATA_DIR + "/uploads"

// PAGE routing path definition.
const (
	PAGE_HOME    = "/"
	PAGE_ADMIN   = "/admin"
	PAGE_LOGIN   = "/login"
	PAGE_LOGOUT  = "/logout"
	PAGE_MANAGER = "/manager"
)

// CONTEXT_IDENTITY defines identity context key.
const CONTEXT_IDENTITY = "identity"
