build:
	protoc.exe -I. --go_out=plugins=micro:. proto/user/user.proto