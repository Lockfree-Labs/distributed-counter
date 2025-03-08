# 

# Distributed Counter API


The Distributed Counter Service is built to allow tenants to maintain their own sets of
counters within isolated namespaces. This ensures that counter operations (increment, 
decrement, read, etc.) occur in a secure, conflict-free environment. The service is 
distributed and designed to handle a high volume of concurrent operations across 
multiple data centers.

### Key Features
- **Multi-Tenancy**: Each tenant is provided a dedicated namespace ensuring isolation between tenants.
- **Scalability:** Built on a distributed architecture, it can handle thousands of concurrent counter updates.
- **High Availability**: Designed for fault tolerance with load balancing and redundancy.
- **Flexible Counter Management**: Tenants can create and manage multiple counters within their namespace.
- **Consistent Performance**: Optimized for low-latency operations in a distributed environment.
- **Eventually Consistent**: All counters are eventually consistent.
- **Auto Expiring Counters**: Counters by default expire after 24 hours of inactivity.

## Endpoints
- `POST /counter/{counter_id}/increment` increments the counter.
- `GET /counter/{counter_id}` retrieves the current counter.

## Usage Instructions
The API is hosted on rapidapi for ease of use.
https://rapidapi.com/malikanshul29/api/distributed-counter

