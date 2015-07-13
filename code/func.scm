(define (append l . e)
  (define (ap x r)
    (if (null? (cdr r))
	(cons (car r) x)
	(cons (car r) (ap x (cdr r)))))
  (ap e l))

(define (map f l)
  (cond
   ((not (pair? l)) (f l)) ; improper
   ((null? l) '())
   (else (cons (f (car l)) (map f (cdr l))))
))

(define (filter f l))

(define (left-fold f l)) ; https://en.wikipedia.org/wiki/Fold_(higher-order_function)
(define (right-fold f l))

(define (not b)
  (if b #f #t))
