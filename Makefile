MAKEFLAGS	+=	--quiet

WHT	= \033[0;37m
BLK	= \033[0;30m
RED	= \033[0;31m
YEL	= \033[0;33m
BLU	= \033[0;34m
GRN	= \033[0;32m

NAME	=	bit

DIR_S	=	src

RM		=	rm -fdr

$(NAME):
			cd $(DIR_S) && go build -o ../bin/$(NAME)

install:
			cd $(DIR_S) && go install
			printf "$(WHT)[$(GRN)$(NAME) PROGRAM INSTALLED$(WHT)]\n"

build:		$(NAME)
			printf "$(WHT)[$(GRN)$(NAME) PROGRAM COMPILED$(WHT)]\n"

all:		test build install

clean:
			$(RM) bin $(DIR_S)/cover.out $(DIR_S)/cover.html
			printf "$(WHT)[$(YEL)$(NAME) BINARIES AND COVERAGE REMOVED$(WHT)]\n"

test:
			cd $(DIR_S) && go test ./... -v -coverprofile=cover.out
			cd $(DIR_S) && go tool cover -html=cover.out -o=cover.html

.PHONY:		all build install clean test