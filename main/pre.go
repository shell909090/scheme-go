package main

const (
	PreDefineMacro = `
(define-syntax cond
  (syntax-rules (else)
    ((_ (rule1 value1) (else value2) body ...)
     (if rule1 value1 value2))
    ((_ (rule1 value1) body ...)
     (if rule1 value1 (cond body ...)))))
(define-syntax error
  (syntax-rules (error)
  ((_ p ...)
   (display p ...))))
`
)
