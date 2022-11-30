# Bitcoin Rate app
This application was created from [template](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template) and can be used with [Eliona](https://www.eliona.io/) enviroment.


## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables the app provides an own API access.


### Registration in Eliona ###

To start and initialize an app in an Eliona environment, the app have to registered in Eliona. For this, an entry in the database table `public.eliona_app` is necessary.


### Environment variables ###

#### API_ENDPOINT
The `APPNAME` MUST be set to `bitcoin`. Some resources use this name to identify the app inside an Eliona environment. For running as a Docker container inside an Eliona environment, the `APPNAME` have to set in the [Dockerfile](Dockerfile). If the app runs outside an Eliona environment the `APPNAME` must be set explicitly.

```bash
export APPNAME=bitcoin
```

#### CONNECTION_STRING

The `CONNECTION_STRING` variable configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started.

```bash
export CONNECTION_STRING=postgres://user:pass@localhost:5432/iot
```

#### API_ENDPOINT

The `API_ENDPOINT` variable configures the endpoint to access the [Eliona API](https://github.com/eliona-smart-building-assistant/eliona-api). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started.

```bash
export API_ENDPOINT=http://localhost:8082/v2
```

#### DEBUG_LEVEL (optional)

The `DEBUG_LEVEL` variable defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.

```bash
export LOG_LEVEL=debug # optionally, default is 'info'
```

#### Example - Environment variables 
![img.png](environment_variables.png)

### Database tables ###

The app requires some configuration data that remains in the database. To do this, the app creates its own database schema `bitcoin` during initialization. The data in this schema should be made editable by eliona frontend. This allows the app to be configured by the user without direct database access.

A good practice is to initialize the app configuration with [default values](sql/defaults.sql). This allows the user to see how what needs to be configured.

In detail, you need the following configuration data in table `bitcoin.configuration (name, value)`.

#### Example - settings in the DB
```sql
-- bitcoin.configuration (name, value)
('endpoint', 'https://api.coindesk.com/v1/bpi/currentprice.json') -- where is the API located
('polling_interval', '10') -- with interval in seconds is used to poll the API 
```

In order to define the currencies which rates are to be read, an entry in the table `bitcoin.currencies (code, description, proj_id)` is required. Each rate is later mapped with an asset in eliona.


## API Reference

The bitcoin-rate app grabs bitcoin rate from [coindesk](https://api.coindesk.com/v1/bpi/currentprice.json) web service and writes these data to eliona as heap data of assets. The heap data is separated in `bitcoin.Input`, `bitcoin.Info` and `bitcoin.Status` heaps. These structures are used to write the heap data.

```json
{"code": "USD", "rate": 16793.8472}
{"daytime": "Nov 30, 2022 08:58:00 UTC"}
{"comment": "United States Dollar"}
```

In eliona these heaps are handled as `bitcoin_rate` asset type with appropriate attributes created during the [initialization](init/init.sql).

