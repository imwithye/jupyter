build:
	GOOS=linux GOARCH=amd64 go build -o jupyterlab_linux_amd64.exe
	GOOS=windows GOARCH=amd64 go build -o jupyterlab_windows_amd64.exe
	GOOS=darwin GOARCH=arm64 go build -o jupyterlab_darwin_arm64.exe

clean:
	rm -rf *.exe

.PHONY: build clean
