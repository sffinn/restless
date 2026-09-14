# Go REST API Example for DigitalOcean App Platform

A minimal, production-ready REST API written in Go that runs perfectly on [DigitalOcean App Platform](https://www.digitalocean.com/products/app-platform).

## Features

- Full CRUD for a simple `Item` resource (in-memory store)
- Health endpoint for App Platform health checks
- Listens on `0.0.0.0:$PORT` (required by App Platform)
- Multi-stage Dockerfile for small, secure images
- Optional App Spec (`.do/app.yaml`) for reproducible deployments
- Zero external dependencies beyond `google/uuid` (uses Go 1.22+ `net/http` method routing)

## API Endpoints

| Method | Path              | Description              |
|--------|-------------------|--------------------------|
| GET    | `/`               | API info                 |
| GET    | `/health`         | Health check             |
| GET    | `/api/items`      | List all items           |
| POST   | `/api/items`      | Create an item           |
| GET    | `/api/items/{id}` | Get a single item        |
| PUT    | `/api/items/{id}` | Update an item           |
| DELETE | `/api/items/{id}` | Delete an item           |

### Example requests

```bash
# List items
curl https://your-app.ondigitalocean.app/api/items

# Create
curl -X POST https://your-app.ondigitalocean.app/api/items \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk"}'

# Update
curl -X PUT https://your-app.ondigitalocean.app/api/items/<id> \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy oat milk","completed":true}'

# Delete
curl -X DELETE https://your-app.ondigitalocean.app/api/items/<id>

