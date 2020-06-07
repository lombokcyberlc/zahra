# build dan install
all:
	go build -o zahra main.o && sudo mv zahra /usr/local/zahra

# install semua depedencies yang dibutuhkan
go-get:
	go get ./...

# menjalankan main.go
run:
	go run main.go

# build main.go menjadi zahra
build: 
	go build -o zahra main.go

# hapus zahra exe
clean:
	sudo rm -rf zahra

# install zahra
install:
	sudo cp zahra /usr/local/zahra 

# install zahra service
install-service:
	sudo cp zahra.service /etc/systemd/system/

# uninstall zahra
uninstall:
	sudo rm -rf /usr/local/zahra

# melihat status service zahra
status:
	sudo systemctl status zahra

# start service zahra
start:
	sudo systemctl start zahra && sudo systemctl daemon-reload

# stop service zahra
stop:
	sudo systemctl stop zahra
