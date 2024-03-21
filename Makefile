GOIMPORTS = goimports

.PHONY: old-copy
old-copy:
	rm -rf internal/old
	mkdir internal/old
	cp -r ast field lex parse token internal/old/
	sed -r 's,"github.com/tsavola/dp/(ast|field|lex|parse|token)\b,"github.com/tsavola/dp/internal/old/\1,' -i internal/old/*/*.go
	sed -r 's,(old\w*) "github.com/tsavola/dp/(ast|field|lex|parse|token)\b,\1 "github.com/tsavola/dp/internal/old/\2,' -i internal/revise/*.go
	$(GOIMPORTS) -w internal/old/*/*.go internal/revise/*.go
