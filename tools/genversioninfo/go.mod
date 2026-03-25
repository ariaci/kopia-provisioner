module kopia-provisioner-tools/genversioninfo

go 1.26.1

require github.com/josephspurrier/goversioninfo v1.5.0 // direct

require kopia-provisioner v0.0.0-00010101000000-000000000000

require github.com/akavel/rsrc v0.10.2 // indirect

replace kopia-provisioner => ../..
