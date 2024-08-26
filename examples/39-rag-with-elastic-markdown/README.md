# How to run this example

Ref: https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html

## Elastic
To store embeddings in Elasticsearch using the Golang client, you'll need to follow these steps:

**Start Elasticsearch and Kibana**
```bash
docker compose up -d
```

Wait for some seconds until the services are up and running.
Then, a `certs` directory will be created with the certificates needed to connect to Elasticsearch.

> To stop the Docker Compose stack, run: `docker compose down`

### Test the connection with the certificate

```bash
export ELASTIC_PASSWORD=iloveparakeets
curl --cacert ./certs/ca/ca.crt -u elastic:$ELASTIC_PASSWORD https://localhost:9200
```

You should get something like this:
```bash
{
  "name" : "602f990409e9",
  "cluster_name" : "docker-cluster",
  "cluster_uuid" : "IaF2jmr0SwS8RrsiHMrDcw",
  "version" : {
    "number" : "8.15.0",
    "build_flavor" : "default",
    "build_type" : "docker",
    "build_hash" : "1a77947f34deddb41af25e6f0ddb8e830159c179",
    "build_date" : "2024-08-05T10:05:34.233336849Z",
    "build_snapshot" : false,
    "lucene_version" : "9.11.1",
    "minimum_wire_compatibility_version" : "7.17.0",
    "minimum_index_compatibility_version" : "7.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

## Run the samples

This example demonstrates how to store embeddings in Elasticsearch and perform a vector similarity search using the Go client.

### Create the embeddings

```bash
(cd create-embeddings; go run main.go)
```

### Check if the embeddings are stored

- Go to Kibana: http://0.0.0.0:5701/app/management/data/index_management/indices
- Log in with the user "elastic" and password "iloveparakeets"
- Open the console and run the following query:
    ```bash
    GET /chronicles-index/_search
    ```
You should see the embeddings stored in the index `chronicles-index`:

![Kibana](./imgs/kibana.png)

### Perform a vector similarity search

```bash
(cd use-embeddings; go run main.go)
```

[This program](use-embeddings/main.go) completes the prompt: "What's new with TLS client?". If everything works, you'll see it answered from embeddings derived from [go1.23.md](create-embeddings/go1.23.md).

While your results may vary, here's an example output:
```
🔎 searching for similarity...
📝 doc: 19 score: 1.3576632
📝 doc: 27 score: 1.3363576
📝 doc: 20 score: 1.2898484
📝 doc: 1 score: 1.2662655
📝 doc: 23 score: 1.2095306

🤖 answer:
According to the crypto/tls section, the TLS client now supports the Encrypted Client Hello (ECH) draft specification. This feature can be enabled by setting the Config.EncryptedClientHelloConfigList field to an encoded ECHConfigList for the host that is being connected to.

Additionally, the QUICConn type used by QUIC implementations includes new events reporting on the state of session resumption, and provides a way for the QUIC layer to add data to session tickets and session cache entries.

Also, 3DES cipher suites were removed from the default list used when Config.CipherSuites is nil, and the experimental post-quantum key exchange mechanism X25519Kyber768Draft00 is now enabled by default when Config.CurvePreferences is nil.
```
