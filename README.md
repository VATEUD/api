# VATEUD API

VATEUD API that's used by most of our websites.

## Usage

Requirement - [Go](https://golang.org/)

1. Copy and edit the env file
```bash
cp .env.example .env
```
2. Migrate the tables (you can find them [here](https://github.com/VATEUD/api/blob/master/database/api.sql))
```bash
mysql -u username -p database < api/api.sql
```
3. Build the API
```bash
go build cmd/api/main.go
```
4. Start/Restart the daemon service (a simple example can be found [here](https://github.com/VATEUD/api/blob/master/scripts/api.service))

Start 
```bash
sudo service name_of_the_service enable && sudo service name_of_the_service start
```
Restart
```bash
sudo service name_of_the_service restart
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
