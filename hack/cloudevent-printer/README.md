# cloudevent-printer

cloudevent-printer is a simple binary you can run to print out all cloudevents the server receives. This is useful to test what your clients e.g. ingestors, are sending

```
$ go run .
2024/01/01 23:19:29 Listening for events...

Context Attributes,
  specversion: 1.0
  type: ingenium.ingestor.data
  source: ingenium/ingestor/csv/
  id: da01e0fa-b20a-4994-bbce-797fca650821
  time: 2024-01-01T23:18:00.6933154Z
  datacontenttype: application/json
Data,
  {
    "Type": "data.type.ohlc",
    "Symbol": "PLTR",
    "Timestamp": "2021-06-18",
    "Data": {
      "Open": "25.590000",
      "High": "25.940001",
      "Low": "25.020000",
      "Close": "25.370001",
      "Volume": "65760500",
      "Period": ""
    }
  }

...
```
