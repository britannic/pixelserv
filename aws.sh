#!/bin/bash
# Updates the AWS repository and pushes it to github

VERSION="1.0"
host=${1}
pkg=${2}
tag=${3}

# echo ${*}
ssh -tt ${host} <<-EOF
	cd repositories/blacklist/
	reprepro includedeb wheezy /tmp/${pkg}*.deb
	cd ..
	git add --all
	git commit -am"Package repository release ${pkg}package"
	git tag "${pkg}package"
	git push origin master
	git push --tags
	exit
EOF

