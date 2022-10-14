# whats-this

That's a little go application to export the last tweets from your twitter dot com account over the api 
to your local hugo repository. 

## How to use it
* You need an API key for twitter
* setup an shell script with ENV variables, for example like this

```sh
#!/bin/sh

ACCESS_TOKEN="XXXX" \
ACCESS_TOKEN_SECRET="XXXX" \
CONSUMER_KEY="XXXX" \
CONSUMER_KEY_SECRET="XXXX" \
HUGO_POST_PATH="/path/to/your/hugo/content/posts" \
HUGO_PATH_STATIC="/path/to/your/hugo/static" \
go run export.go
```