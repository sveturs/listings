#!/bin/bash
curl -I http://localhost:8080/api/v1/cars/available
curl -I https://localhost:8443/api/v1/cars/available --insecure
