package main

// To avoid key collisions with 3rd-party packages.
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
