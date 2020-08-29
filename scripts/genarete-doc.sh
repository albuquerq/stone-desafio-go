#!/bin/bash

docker run -ti --rm -v $(pwd)/../docs:/docs humangeo/aglio aglio -i /docs/banking.apib -o ../docs/doc-api-v1.html
