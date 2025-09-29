#!/bin/bash

echo "Fixing Swagger annotations in all Go files..."

# Исправляем utils.SuccessResponseSwag на backend_pkg_utils.SuccessResponseSwag
find /data/hostel-booking-system/backend -name "*.go" -type f -exec sed -i 's/utils\.SuccessResponseSwag/backend_pkg_utils.SuccessResponseSwag/g' {} \;

# Исправляем utils.ErrorResponseSwag на backend_pkg_utils.ErrorResponseSwag
find /data/hostel-booking-system/backend -name "*.go" -type f -exec sed -i 's/utils\.ErrorResponseSwag/backend_pkg_utils.ErrorResponseSwag/g' {} \;

# Исправляем example для массивов - заменяем неправильный формат [1 на [1]
find /data/hostel-booking-system/backend -name "*.go" -type f -exec sed -i 's/example:"\[1"/example:"[1]"/g' {} \;
find /data/hostel-booking-system/backend -name "*.go" -type f -exec sed -i 's/example:"\[1,/example:"[1,/g' {} \;

echo "Swagger annotations fixed!"
