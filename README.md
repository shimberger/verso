# Verso

Simple reverse proxy for local development

## Usage

	verso --backend [target url] /assets:./build

This will proxy all requests to `http://localhost:8080` except all requests to `http://localhost:8080/assets/...` will be served from local build directory.