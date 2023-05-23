package hyperv_winrm

import (
	"context"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type createDvdArgs struct {
	Path string
	Ip   string
}

var createDvdTemplate = template.Must(template.New("CreateDvd").Parse(`
$ErrorActionPreference = 'Stop'
$path='{{.Path}}'
$ip='{{.Ip}}'

$yamlContent = @{
    "network"=@{
        "ethernets"=@{
            "eth0"=@{
                "dhcp4"="no"
                "gateway4"="172.16.1.254"
                "addresses" = @("$ip/16")
                "nameservers"=@{
                    "addresses"=@("172.16.14.27")
                }
            }
        }
    }
}

$folderPath = Split-Path -Path $path -Parent

if (-not (Test-Path -Path $folderPath -PathType Container)){
    New-Item -ItemType Directory -Path $folderPath | Out-Null
}

$tmpPath = Split-Path -Path $folderPath -Parent
$tmpPath += "\tmp"

if (-not (Test-Path -Path $tmpPath -PathType Container)){
    New-Item -ItemType Directory -Path $tmpPath | Out-Null
}

$yamlContent | ConvertTo-Yaml | Out-File -FilePath "$tmpPath\network_settings.yaml" -Encoding UTF8 
oscdimg -n -d -m $tmpPath $path
Remove-Item -LiteralPath $tmpPath -Force -Recurse

`))

func (c *ClientConfig) CreateDvd(ctx context.Context, path string, ip string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, createDvdTemplate, createDvdArgs{
		Path: path,
		Ip:   ip,
	})

	return err
}

type getDvdArgs struct {
	Path string
	Ip   string
}

var getDvdTemplate = template.Must(template.New("GetDvd").Parse(`
$ErrorActionPreference = 'Stop'
$path='{{.Path}}'
$ip='{{.Ip}}'

if (Test-Path $path) {
	$dvd = @{
        Path=$path
        Ip=$ip
    }
    $dvd = ConvertTo-Json -InputObject $dvd
    $dvd
} else {
	"{}"
}
`))

func (c *ClientConfig) GetDvd(ctx context.Context, path string, ip string) (result api.Dvd, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getDvdTemplate, getDvdArgs{
		Path: path,
		Ip:   ip,
	}, &result)

	return result, err
}

type deleteDvdArgs struct {
	Path string
}

var deleteDvdTemplate = template.Must(template.New("DeleteDvd").Parse(`
$ErrorActionPreference = 'Stop'

$targetDirectory = (split-path '{{.Path}}' -Parent)
$targetName = (split-path '{{.Path}}' -Leaf)
$targetName = $targetName.Substring(0,$targetName.LastIndexOf('.')).split('\')[-1]

Get-ChildItem -Path $targetDirectory |?{$_.BaseName.StartsWith($targetName)} | %{
	Remove-Item $_.FullName -Force
}

`))

func (c *ClientConfig) DeleteDvd(ctx context.Context, path string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, deleteDvdTemplate, deleteDvdArgs{
		Path: path,
	})

	return err
}
