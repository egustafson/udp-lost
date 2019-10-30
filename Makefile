all:
	@echo "\n  make build | make run\n"

build:
	docker build -t udploss .

run:
	docker run -it --rm --name udploss udploss
