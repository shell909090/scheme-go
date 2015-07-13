(define-syntax cond
  (syntax-rules (else)
    ((_ (rule1 value1) (else value2) body ...)
     (if rule1 value1 value2))
    ((_ (rule1 value1) body ...)
     (if rule1 value1 (cond body ...)))))
(define-syntax let)
(define-syntax when)
(define-syntax unless)
(define-syntax and)
(define-syntax or)
(define-syntax error
  (syntax-rules (error)
    ((_ p ...)
     (display p ...))))

(define-syntax caar
  (syntax-rules ()
    ((_ p)
     (car (car p)))))
(define-syntax cadr
  (syntax-rules ()
    ((_ p)
     (cdr (car p)))))
(define-syntax cdar
  (syntax-rules ()
    ((_ p)
     (car (cdr p)))))
(define-syntax cddr
  (syntax-rules ()
    ((_ p)
     (cdr (cdr p)))))
