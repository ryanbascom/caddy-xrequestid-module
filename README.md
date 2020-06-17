# caddy-xrequestid-module
A caddy module that adds an X-Request-Id http header if missing from origin request.

## Usage
By default, if the incomming request does not contain an ```X-Request-Id``` header, or the value of the header is empty or all whitespace, this module generates and adds the header in the format ```X-Request-Id:{UUID}```. 

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
