# Other helpers and Parakeet methods

## Get Information about a model

```golang
llm.ShowModelInformation(url, model string) (llm.ModelInformation, int, error)
```

`ShowModelInformation` retrieves information about a model from the specified URL.

**Parameters**:
  - url: the base URL of the API.
  - model: the name of the model to retrieve information for.

**Returns**:
  - ModelInformation: the information about the model.
  - int: the HTTP status code of the response.
  - error: an error if the request fails.

**✋ Remark**: if the model does not exist, it will return an error with a status code of 404.

## Pull a model

```golang
llm.PullModel(url, model string) (llm.PullResult, int, error)
```

`PullModel` sends a POST request to the specified URL to pull a model with the given name.

**Parameters**:
  - url: The URL to send the request to.
  - model: The name of the model to pull.

**Returns**:
  - PullResult: The result of the pull operation.
  - int: The HTTP status code of the response.
  - error: An error if the request fails.
<!-- split -->