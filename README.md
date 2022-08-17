# Go Rate Limiter

This project includes the following implementations to rate limiting:
- **Window**
    - Fixed — Number of requests are fixed over a time period. For example, setting a limit of 5 requests per minute. 5 request limit is only applied for a fixed minute interval (10:00 to 10:01). Ultimately, a user can make 10 request under that 1 minute interval.

    - Sliding — Number of requests restricted over time period however that is determined based on the time of the inital invocation. If we apply the same 5 req/min if a user starts innvocation at 10:00:30 they are permitted to invoke 4 additional times until 10:01:30.

- **Bucket**
    - Token Bucket — Set number of tokens are placed in a hypothetical bucket. Tokens are replensihed at fixed time period. Requests will only be allowed if request has a corresponding token.
    
    - Leaky Bucket — Packets are processed from a queue and allowed if space permits and processed at a fixed rate. Each individual client is to have their own queue.

## Custom Rates:
Both Window and Bucket categories can make use of the Basic and Premium Subscripition. 
- **Basic** 
    - Window — 5 requests per minute
    - Bucket — Burst of 5 requests and token replenished every 12 seconds
- **Premium**
    - Window — 20 requests per minute
    - Bucket — Burst of 20 request and token replensished every 3 seconds  

# How to run locally:
Within the folder run the following command:

``` go run main.go ```

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
