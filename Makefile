NAME	:= gotta-go-fast-api.bin

all: $(NAME)

$(NAME):
	go build -o $(NAME)

build:
	go build -o $(NAME)

run: $(NAME)
	./$(NAME)

debug:
	go build -ldflags="-X 'main.buildMode=debug'" -o $(NAME)

clean:
	rm $(NAME)

.PHONY: all build run debug clean