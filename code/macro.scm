; let let*
(define-syntax let
  (syntax-rules ()
    ((_ ((name1 value1)) body ...)
     ((lambda (name1) body ...) value1)
    ((_ ((name1 value1) others ...) body ...)
     (let (others ...)
       ((lambda (name1) body ...) value1))))))

(define-syntax cond
  (syntax-rules (else)
    ((_ (rule1 value1))
     (if rule1 value1))
    ((_ (rule1 value1) (else value2) body ...)
     (if rule1 value1 value2))
    ((_ (rule1 value1) body ...)
     (if rule1 value1 (cond body ...)))))
(define-syntax when
  (syntax-rules ()
    ((_ test body ...)
     (if test
	 (begin body ...)))))
(define-syntax unless
  (syntax-rules ()
    ((_ test body ...)
     (if (not test)
	 (begin body ...)))))

(define-syntax and
  (syntax-rules ()
    ((_ a b)
     (let (t a)
       (if (not t) t b)))
    ((_ a b c ...)
     (and (and a b) c ...))))
(define-syntax or
  (syntax-rules ()
    ((_ a b)
     (let (t a)
       (if t t b)))
    ((_ a b c ...)
     (or (or a b) c ...))))

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
