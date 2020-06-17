# caddy-xrequestid-module
A caddy module that adds an X-Request-Id http header if missing from origin request.

## Usage
By default, if the incomming request does not contain an ```X-Request-Id``` header, or the value of the header is empty or all whitespace, this module generates and adds the header in the format ```X-Request-Id:{UUID}```.

If the request contains an ```X-Request-Id``` header with a value that is not empty or blank, the header is not changed and passed along with the request.

## Example configuration

```json
...
"handle": [
{
  "handler": "x_request_id",
  "disable": false
},
...
```
