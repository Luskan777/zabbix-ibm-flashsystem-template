# zabbix-ibm-flashsystem-template

#### Zabbix template for IBM FlashSystem storage.

Adapted from original bash script created by [Yarli](https://github.com/Yarli) to Golang Script created by [me](https://github.com/Luskan777).

Supported IBM versions: IBM FS-3500, IBM FS-5000, IBM FS-7000 (needs testing)

#### Build the Golang script:

Build with Go binary (needs Go installed):

`cd src && env GOOS=linux GOARCH=amd64 go build -o IBMStorageFS5K  *.go`

Build with Docker (needs Docker and docker-compose installed):

`docker compose run --rm golang`

How to use Template:
 - Copy the script above from source path (`./src/IBMStorageFS5K.go`) to zabbix external script path (`/usr/lib/zabbix/externalscripts/`)
 - Import the template file into your templates section within the Zabbix ui.
 - Create a new user on your storage system which Zabbix can use.
 - Create host and add following user macros:
    - {$CABIP1} - this is the IP address of your storage system
    - {$CABPASS} - Password of your new user
    - {$CABUSER} - Username of your new user

Then execute the 5 discovery rules, which should then populate all the items you need for monitoring.


