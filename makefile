.PHONY: build install

build: # build the binary with the proper access rights and file name
	@echo "Running packr2..."
	GO111MODULE=on packr2
	@echo "Done."
	@echo "Building binary..."
	go build -o csvset
	@echo "Done."
	@echo "Modifying access rights to binary (755)..."
	chmod 755 csvset
	@echo "Done."
	@echo "Cleaning up..."
	packr2 clean
	@echo "Done."

install: build
	@echo "Placing binary in installation folder"
	mv csvset /usr/local/bin/
	@echo "Done."

clean:
	@echo "Removing installation contents"
	rm /usr/local/bin/csvset
	@echo "Done."