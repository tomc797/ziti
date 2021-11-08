[CmdletBinding()]
Param(
  [Parameter()]
  [Alias('F')]
  [string]$Foo
)

$latestFromGitHub=(irm https://api.github.com/repos/openziti/ziti/releases/latest)
$version=($latestFromGitHub.tag_name)
$zitidl=($latestFromGitHub).assets | where {$_.browser_download_url -Match "ziti-windows.*zip"}
$downloadUrl=($zitidl.browser_download_url)
$name=$zitidl.name
$defaultFolder="$env:USERPROFILE\.ziti\bin"
$toDir=$(Read-Host "Where folder should be used for ziti? [default: ${defaultfolder}]")
if($toDir.Trim() -eq "") {
    $toDir=("${defaultfolder}")
}

$zipFile="${toDir}\${name}"
if($(Test-Path -Path $zipFile -PathType Leaf)) {
    echo "The file has already been downloading. No need to download again"
} else {
    echo "The file does not exist"
    mkdir "${toDir}" -ErrorAction SilentlyContinue
    echo "Downloading file "
    echo "    from: ${downloadUrl} "
    echo "      to: ${zipFile}"
    iwr ${downloadUrl} -OutFile "$zipFile"
}

Expand-Archive -Path $zipFile -DestinationPath "${toDir}\${version}" -ErrorAction SilentlyContinue

#$addToPath=$(Read-Host "Add ziti to your path? [default: y]")
#if($toDir.Trim() -eq "") {
#    $addToPath="y"
#}
$env:Path+=";${toDir}\${version}\ziti\"

echo "You can now add ziti to your path by running this command:"
echo "`${env:Path}+=`";${toDir}\${version}\ziti\`""



