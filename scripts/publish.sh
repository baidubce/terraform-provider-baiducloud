# Only for the macOS
# Usage: cd scripts && sh publish.sh
echo "Start publishing"
echo "------------------"
echo

echo "cd $(pwd)"
cd ..
echo "------------------"
echo

echo "make fmt -f GNUmakefile"
make fmt -f GNUmakefile
echo "------------------"
echo

echo "go mod vendor"
go mod vendor
echo "------------------"
echo

echo "Generate documents..."
go run gendocs/main.go gendocs/template.go
echo "------------------"
echo

echo "Manual inspection:"
git status
if read -t 60 -n1 -p "Continue to pubilsh? Please enter: [Y/N]" isPublishing
then
	case $isPublishing in
		(Y | y)
			echo
			echo "Continue publishing";;
		(N | n)
			echo
			echo "Exit. Goodbye."
			exit;;
		(*)
			echo
			echo "Error choice. Please enter: [Y/N]. Exit. Please try again."
			exit;;
	esac
else
	echo
	echo "Timeout. Exit..."
	exit
fi
echo "------------------"
echo

echo "Update CHANGELOG.md"
let currentMajorVersion=$(cat CHANGELOG.md | grep Unreleased | awk -F'.' '{print $1}' | awk -F' ' '{print $2}')
let currentMinorVersion=$(cat CHANGELOG.md | grep Unreleased | awk -F'.' '{print $2}')
let currentPatchVersion=$(cat CHANGELOG.md | grep Unreleased | awk -F'.' '{print $3}' | awk -F' ' '{print $1}')
currentVersion="$currentMajorVersion"."$currentMinorVersion"."$currentPatchVersion"
nextVersion="$currentMajorVersion"."$(expr $currentMinorVersion + 1)"."0"
sed -i '' "s/Unreleased/$(LANG=en_US.UTF-8, date +%B\ %d,\ %Y)/" CHANGELOG.md
sed -i '' "1i\ 
	## $nextVersion\ (Unreleased)
	" CHANGELOG.md
echo "------------------"
echo

echo "git add CHANGELOG.md"
git add CHANGELOG.md
git commit -m "v"$currentVersion
git tag "v"$currentVersion
echo "------------------"
echo

echo 'The repository is ready to push. Please use "git push && git push --tag" to trigger the github actions to publish the provider.'
