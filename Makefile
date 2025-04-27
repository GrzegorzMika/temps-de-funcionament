build:
	docker buildx build --platform linux/amd64 --build-arg ARCH=amd64 -t gregmika/temps_de_funcionament:$(version)-amd64 --push .
	docker buildx build --platform linux/arm64 --build-arg ARCH=arm64 -t gregmika/temps_de_funcionament:$(version)-arm64 --push .