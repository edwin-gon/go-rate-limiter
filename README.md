# Go Rate Limiter

This project includes the following approaches to rate limiting:
- Fixed Window
- Sliding Window
- Token Bucket
- Leaky Bucket (Coming Soon)  

Each implementation is made available via API endpoints. The API construction handles centralize response mapping to provide appropriate status codes and messages. 

# How to run locally:
Within the folder run the following command:

``` go run . ```

This will create a HTTP server running on port 5050.

## Endpoints:
- **/sliding/client**
- **/fixed/client** 
### Future Features:
- Testing
- Logging
- Dependency Injection (Wire)
- Redis/DB Reads

## Resources:
- 