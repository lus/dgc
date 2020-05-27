package dgc

// Middleware defines how a middleware looks like
type Middleware func(following ExecutionHandler) ExecutionHandler
