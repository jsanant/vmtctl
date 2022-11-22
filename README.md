## vmtctl

A CLI tool to generate [multi-tenant](https://docs.victoriametrics.com/Cluster-VictoriaMetrics.html#multitenancy) URLs in victoria-metrics and run the clustered version locally. The tool helps in understanding on how to use the [multi-tenant URLs](https://docs.victoriametrics.com/Cluster-VictoriaMetrics.html#url-format) in clustered mode.

## Usage

**Note:** Every time you run the binary new tenants are generated.

```
$ ./bin/vmtctl -h
      --config string   Path to the .toml config file (default "sample_config.toml")
      --csv             Create a CSV file with endpoints
      --gen-endpoints   Generate only endpoints

```

- Clone the repo
- Edit the `sample_config.toml` according to your requirements
- Run `make dev`, which will

  - Build the binary
  - Generate the multi-tenant endpoints
  - Populate the vmselect and vminsert endpoints in `datasource.yml` and `docker-compose.yml` files respectively and bring up victoria-metrics clustered version

- To generate a CSV file with the endpoints, you can run this command:

```
./bin/vmtctl --csv
```

This will create a file called `vmtctl.csv`.

- If you want to generate only the endpoints and not run victoria-metrics, you can run this command:

```
./bin/vmtctl --gen-endpoints

{"level":"info","ts":1658677393.758634,"caller":"vmtctl/main.go:32","msg":"Generating endpoints"}
+-----------+------------+-----------+------------------------------------------------------+------------------------------------------------------+
| COMPONENT | ACCOUNTID  | PROJECTID |                    PROMETHEUS URL                    |                     INFLUXDB URL                     |
+-----------+------------+-----------+------------------------------------------------------+------------------------------------------------------+
| vminsert  | 283739781  | 0         | http://vminsert:8480/insert/283739781:0/prometheus/  | http://0.0.0.0:8480/insert/283739781:0/influx/write  |
+           +------------+-----------+------------------------------------------------------+------------------------------------------------------+
|           | 2781260307 | 0         | http://vminsert:8480/insert/2781260307:0/prometheus/ | http://0.0.0.0:8480/insert/2781260307:0/influx/write |
+-----------+------------+-----------+------------------------------------------------------+------------------------------------------------------+
| vmselect  | 283739781  | 0         | http://vmselect:8481/select/283739781:0/prometheus   |                                                      |
+           +------------+-----------+------------------------------------------------------+------------------------------------------------------+
|           | 2781260307 | 0         | http://vmselect:8481/select/2781260307:0/prometheus  |                                                      |
+-----------+------------+-----------+------------------------------------------------------+------------------------------------------------------+
```

### Send metrics

- Once the cluster is up and running, you can use the tenant endpoints generated under the InfluxDB URL column to send data.

```
curl -d 'measurement,tag1=value1,tag2=value2 field1=123,field2=1.23' -X POST 'http://0.0.0.0:8480/insert/283739781:0/influx/write'
```

More info how victoria-metrics translates data from [InfluxDB endpoints](https://docs.victoriametrics.com/#how-to-send-data-from-influxdb-compatible-agents-such-as-telegraf).

- Sending metrics in Prometheus text exposition format

```
curl -d 'metric_name{foo="bar"} 123' -X POST http://0.0.0.0:8480/insert/283739781:0/prometheus/api/v1/import/prometheus
```

- After this is done, you can navigate to the grafana explore section using this [link](http://localhost:3000/explore). You will see that the metric is present only for the tenant to which you have sent the data.

Default creds:

- login - `admin`
- password - `admin`

# Clean up

**NOTE:** This will remove all volumes, please make sure you do not have any existing volumes before running this command.

Once you are done with your testing, you can run `make clean`. This is will stop all containers and remove all the volumes that were generated.

# LICENSE

[LICENSE](https://github.com/jsanant/vmtctl/blob/main/LICENSE)
