# Setting up your objectstore

Currently, only AWS S3 and Garage are supported for object stores. Garage is S3 compatible, which means the configuration is the same.

1. Create a new bucket in S3 (`GARAGE_BUCKET` env variable or `S3_BUCKET` env variable used by the server). Note that when using the docker compose file in deployments and Garage, a bucket is already created by using `GARAGE_DEFAULT_BUCKET`.

2. Create the file `cors-config.json` in the current directory with the following contents:

```json
{
  "CORSRules": [
    {
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "PUT", "POST"],
      "AllowedHeaders": ["*"],
      "MaxAgeSeconds": 3000
    }
  ]
}
```

3. Configure the CORS by performing the following command (replace `drops` with your bucket name):

```sh
aws s3api put-bucket-cors --bucket <your-bucket-name> --cors-configuration file://cors-config.json
```
