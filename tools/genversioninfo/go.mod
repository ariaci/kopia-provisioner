module kopia-provisioner-tools/genversioninfo

go 1.26.1

require (
	github.com/tc-hib/winres v0.2.1
	kopia-provisioner v0.0.0-00010101000000-000000000000
)

require (
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	golang.org/x/image v0.12.0 // indirect
)

replace kopia-provisioner => ../..
