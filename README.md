# Pagination Implementation Guide

This repository demonstrates the implementation of two pagination techniques: **Limit/Offset Pagination** and **Cursor-Based Pagination**. Below, you'll find a detailed explanation of each technique, including its definition, use cases, pros and cons, and performance considerations. Additionally, the project uses the following tech stack:

---

## **Tech Stack**
- **Programming Language**: Go (Golang 1.23) >=
- **API Server**: Built using Go's standard library.
- **Database**: PostgreSQL (running in a Docker container).
- **Testing**: `testcontainers` for integration testing with PostgreSQL.
- **Data Generation**: `github.com/icrowley/fake` for generating fake data for testing.

---

## **1. Limit/Offset Pagination**

### **Definition**
Limit/Offset pagination divides a dataset into smaller chunks (pages) by specifying:
- **Limit**: The number of items to return per page.
- **Offset**: The number of items to skip before fetching the data.

Example Query:
```sql
SELECT * FROM items LIMIT 10 OFFSET 20;
```
This fetches 10 items after skipping the first 20.

---

### **Use Cases**
- **Static Datasets**: When the dataset doesn’t change frequently.
- **Random Access**: When users need to jump to specific pages (e.g., page 5 of 10).
- **Simple Implementations**: When quick implementation is needed without complex logic.

---

### **Pros**
- **Easy to Implement**: Simple to understand and implement in most databases and APIs.
- **Flexible Navigation**: Allows users to jump to any page directly.
- **Predictable**: Works well for small to medium datasets.

---

### **Cons**
- **Performance Issues with Large Offsets**: As the offset increases, the database must scan and skip more rows, which can be slow and resource-intensive.
  - Example: `OFFSET 100000 LIMIT 10` requires the database to scan 100,010 rows, even though only 10 are returned.
- **Inconsistent Results**: If data changes (e.g., items are added or removed), the same offset may return different results, leading to duplicates or skipped items.
- **High Memory and CPU Usage**: Scanning and skipping rows can consume significant memory and CPU, especially for large datasets.

---

### **Performance**
- **Memory**: High memory usage for large offsets due to scanning and skipping rows.
- **CPU**: High CPU usage for computing offsets in large datasets.
- **Network**: Minimal impact, as only the requested data is sent over the network.

---

## **2. Cursor-Based Pagination**

### **Definition**
Cursor-based pagination uses a unique identifier (cursor) to fetch items after a specific point. The cursor is typically a sortable column like an ID or timestamp.

Example Query:
```sql
SELECT * FROM items WHERE id > 20 LIMIT 10;
```
This fetches 10 items after the item with ID 20.

---

### **Use Cases**
- **Dynamic Datasets**: When the dataset changes frequently (e.g., real-time data like social media feeds).
- **Sequential Access**: When users navigate through pages sequentially (e.g., "Next" button).
- **Large Datasets**: When performance is critical for large datasets.

---

### **Pros**
- **Efficient for Large Datasets**: No need to scan and skip rows, making it faster and more resource-efficient.
- **Consistent Results**: Immune to data changes (e.g., additions or deletions) because it relies on a unique identifier.
- **Scalable**: Performs well even with millions of records.

---

### **Cons**
- **Complex Implementation**: Requires a unique, sortable column (e.g., ID, timestamp) and additional logic to handle cursors.
- **No Random Access**: Users cannot jump to a specific page directly (e.g., page 5 of 10).
- **Cursor Management**: Cursors must be stored and managed correctly to avoid errors.

---

### **Performance**
- **Memory**: Low memory usage, as the database only fetches the required rows.
- **CPU**: Low CPU usage, as there’s no need to compute offsets.
- **Network**: Minimal impact, as only the requested data is sent over the network.

---

## **Comparison Table**

| Feature                  | Limit/Offset Pagination          | Cursor-Based Pagination          |
|--------------------------|----------------------------------|----------------------------------|
| **Ease of Implementation** | Easy                            | More complex                     |
| **Performance**           | Poor for large offsets          | Excellent for large datasets     |
| **Random Access**         | Supported                       | Not supported                    |
| **Consistency**           | Inconsistent with data changes  | Consistent                       |
| **Memory Usage**          | High for large offsets          | Low                              |
| **CPU Usage**             | High for large offsets          | Low                              |
| **Network Usage**         | Minimal                         | Minimal                          |
| **Use Case**              | Static datasets, random access  | Dynamic datasets, sequential access |

---

## **When to Use Each Technique**

### **Limit/Offset Pagination**
- Use when:
  - The dataset is small to medium-sized.
  - Users need random access to pages (e.g., jumping to page 5).
  - Implementation simplicity is a priority.
- Avoid when:
  - The dataset is large and performance is critical.
  - The dataset changes frequently.

---

### **Cursor-Based Pagination**
- Use when:
  - The dataset is large and performance is critical.
  - The dataset changes frequently (e.g., real-time data).
  - Users navigate sequentially (e.g., "Next" button).
- Avoid when:
  - Users need random access to pages.
  - Implementation complexity is a concern.

---

## **Hybrid Approach**
In some cases, you can combine both techniques:
- Use **cursor-based pagination** for sequential navigation (e.g., "Next" button).
- Use **limit/offset pagination** for random access (e.g., jumping to a specific page).

---

## **Tech Stack Implementation Details**

### **1. Go (Golang 1.23)**
- The API server is built using Go's standard library.
- The server exposes endpoints for both pagination techniques:
  - `/items?limit=10&offset=20` for limit/offset pagination.
  - `/items?cursor=20&limit=10` for cursor-based pagination.

### **2. Dockerized PostgreSQL**
- A PostgreSQL database runs in a Docker container for local development and testing.
- Docker Compose is used to manage the containerized environment.

### **3. Testcontainers**
- Integration tests are written using `testcontainers` to spin up a temporary PostgreSQL instance for testing.
- Ensures that pagination logic works correctly with a real database.

### **4. Fake Data Generation**
- The `github.com/icrowley/fake` library is used to generate fake data for testing pagination.
- Example: Generate fake names, emails, and timestamps to populate the database.

---

## **Getting Started**
To implement these pagination techniques in your project:
1. Clone this repository.
2. Set up the environment:
   - Install Docker and Docker Compose.
   - Run `docker-compose up` to start the PostgreSQL container and go-server.
3. Test the implementation:
   - Use tools like `curl` or Postman to test the API endpoints.
   - Run integration tests using `testcontainers`:
     ```bash
     go test -v ./...
     ```

---

## **Contributing**
Feel free to contribute to this project by submitting issues or pull requests. Your feedback and improvements are welcome!

---

## **License**
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Happy coding! 🚀
