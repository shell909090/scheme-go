(define (append l . e))
(define (map f l))
(define (filter f l))
(define (left-fold f l)) ; https://en.wikipedia.org/wiki/Fold_(higher-order_function)
(define (right-fold f l))

(define (not b)
  (if b #f #t))
