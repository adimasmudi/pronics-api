# 🛠️ Pronics Backend

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-v6.0-47A248?logo=mongodb&logoColor=white)

**Pronics Backend** is the core service for the Pronics On-Demand Technician Marketplace. It connects users with local technicians for various services (AC Repair, Electronic Fixes, Cleaning, etc.).

Unlike traditional e-commerce, technician services vary greatly in their requirements. This project utilizes **Go** for high-throughput concurrency and **MongoDB** (NoSQL) to handle **flexible data structures**, allowing different service categories to have unique attributes without rigid schema constraints.

## Key Features

- **Dynamic Service Schema:** leverages MongoDB's document model to store varying service attributes (e.g., _AC Repair_ vs _House Cleaning_) dynamically.
- **Real-time Booking:** Manages booking states from request to completion.
- **Technician Management:** Profiling and verification system for service providers.

## Tech Stack

- **Language:** Go (Golang)
- **Database:** MongoDB
- **Driver:** Mongo Go Driver
- **Architecture:** RESTful API

---

## Getting Started

Follow these steps to set up the project locally.

### 1. Prerequisites

- Go (v1.18+)
- MongoDB (Running locally or via Atlas)

### 2. Installation

```bash
# Clone the repository
git clone (https://github.com/adimasmudi/pronics-api.git)
cd pronics-api

# Install dependencies
go mod download
```
