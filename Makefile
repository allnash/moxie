# Build Moxie
default: linux

.PHONY: moxie_linux
linux:
	@echo "Building moxie binary to './builds/moxie'"
	@(cd cmd/; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-s -w" -o ../builds/moxie)

.PHONY: moxie_osx
osx:
	@echo "Building moxie(moxie_osx) binary to './builds/moxie_osx'"
	@(cd cmd/; CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build --ldflags "-s -w" -o ../builds/moxie_osx)

.PHONY: moxie_win
windows:
	@echo "Building moxie(moxie_windows) binary to './builds/moxie_win.exe'"
	@(cd cmd/; CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build --ldflags "-s -w" -o ../builds/moxie_win.exe)

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf moxie_win.exe moxie_osx moxie
 
install: install_linux

.PHONY: moxie_install_linux
install_linux:
	@echo "Installing Moxie Proxy to /usr/sbin/moxie directory"
	@cp builds/moxie /usr/sbin/moxie
	@mkdir -p /etc/moxie
	@mkdir -p /var/log/moxie
	@cp app.env /etc/moxie
	@echo "[Unit]\
          Description=Moxie the Reverse Proxy\
          \
          [Service]\
          Type=simple\
          Restart=always\
          RestartSec=5s\
          ExecStart=/usr/sbin/moxie\
		  \
          [Install]\
          WantedBy=multi-user.target" > /lib/systemd/system/moxie-proxy.service
	@echo "Start Moxie service using \
		   sudo service moxie-proxy start"

