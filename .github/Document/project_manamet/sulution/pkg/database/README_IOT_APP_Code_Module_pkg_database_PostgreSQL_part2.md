## PostgreSQL and GORM: A Comprehensive Guide to Basic Queries, Conditions, Subqueries, and Dynamic Control

---

### หลักการ (Concept)

#### คืออะไร?

**PostgreSQL** is a powerful, open-source object-relational database system that uses and extends the SQL language. It is known for its reliability, feature robustness, and performance. PostgreSQL is highly extensible, supporting various data types, indexing methods, and procedural languages, making it suitable for a wide range of applications, from small projects to large-scale enterprise systems[reference:0].

**GORM** is an Object-Relational Mapping (ORM) library for the Go programming language. It provides a developer-friendly API that maps Go structs to database tables. This abstraction allows developers to interact with databases using Go code instead of raw SQL, significantly simplifying data access logic, improving development speed, and reducing boilerplate code. GORM supports various databases, including PostgreSQL, MySQL, SQLite, and SQL Server, through dedicated drivers[reference:1].

#### มีกี่แบบ?

**PostgreSQL Query Types:**
*   **Basic `SELECT`**: The fundamental operation for retrieving data from one or more tables.
*   **Conditional Queries (`WHERE`)**: Filters rows based on specified conditions using operators like `=`, `>`, `<`, `LIKE`, `IN`, `BETWEEN`, and logical operators (`AND`, `OR`, `NOT`).
*   **Subqueries (Nested Queries)**: Queries nested within another query, used with `IN`, `EXISTS`, `ANY`, `ALL`, or in `FROM` and `SELECT` clauses.
*   **Conditional Logic (`CASE`)**: An expression that provides if-then-else logic within a query, allowing for value transformation based on conditions.
*   **Set Operations (`UNION`, `INTERSECT`, `EXCEPT`)**: Combine results from multiple `SELECT` statements.
*   **Grouping and Aggregation (`GROUP BY`, `HAVING`)**: Groups rows that share a property and applies aggregate functions (e.g., `COUNT`, `SUM`, `AVG`).
*   **Sorting (`ORDER BY`)**: Sorts the result set based on one or more columns[reference:2].
*   **Limiting (`LIMIT` and `OFFSET`)**: Restricts the number of rows returned and skips a specified number of rows for pagination.

**GORM Query Patterns:**
*   **Basic CRUD**: `Create`, `First`, `Find`, `Update`, `Delete`.
*   **Conditional Queries**: Using `Where`, `Or`, `Not`, and `Select` to filter and shape data.
*   **Subqueries**: Implementing nested queries within `Where` or `From` clauses.
*   **Dynamic Queries**: Building queries conditionally using Go's control structures (`if`, `else`, `switch`), `Scopes`, or `Clauses`.
*   **Preloading (`Preload`)**: Eagerly loading associated data to avoid N+1 query problems.
*   **Raw SQL and `Scan`**: Executing custom SQL queries and mapping results to structs.

**ข้อห้ามสำคัญ:** ใช้ `*` เสมอใน Production เพราะอาจดึงข้อมูลเกินจำเป็น กดดัน Network I/O ไม่แนะนำให้ใช้ GORM กับ Query ที่ซับซ้อนมาก (Complex Reporting) ควรใช้ Raw SQL จะดีกว่า

### การออกแบบ Workflow และ Dataflow

The following diagram illustrates a typical workflow for building a dynamic query using GORM with conditional logic in a Go application:

```mermaid
graph TD
    A[Start: Receive Request with Filters] --> B[Initialize Base Query: db.Model(&User{})];
    B --> C{Apply Filters Conditionally};
    C --> D[if name filter present: db = db.Where("name LIKE ?", name)];
    C --> E[if age range filter: db = db.Where("age BETWEEN ? AND ?", min, max)];
    C --> F[if status filter: db = db.Where("status = ?", status)];
    D --> G[Execute Final Query: db.Find(&users)];
    E --> G;
    F --> G;
    G --> H[Return Results];
```

Data flows from the application layer, where Go structs represent data models, through the GORM API. GORM translates method calls into SQL statements, executes them on the PostgreSQL database, and scans the returned rows back into Go structs.

---

### ตัวอย่างโค้ดที่รันได้จริง

This section provides a complete, runnable example that covers installation, configuration, basic queries, conditional queries, subqueries, and dynamic queries using control structures (`if`, `else`, `switch`).

#### Project Setup and Installation

1.  **Initialize a Go Module:**
    ```bash
    mkdir gorm-postgres-demo && cd gorm-postgres-demo
    go mod init gorm-postgres-demo
    ```

2.  **Install Required Packages:**
    ```bash
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/postgres
    ```
    These commands install the core GORM library and the PostgreSQL driver[reference:3].

#### Configuration and Database Connection

Create a `main.go` file with the following structure:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User Model
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:255;not null"`
	Email     string         `gorm:"size:255;uniqueIndex;not null"`
	Age       int
	Status    string         `gorm:"default:active"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Order Model for subquery example
type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint
	Product   string
	Amount    float64
	CreatedAt time.Time
}

func main() {
	// Database connection configuration
	dsn := "host=localhost user=postgres password=yourpassword dbname=gorm_demo port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable logging
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Auto migrate schemas
	db.AutoMigrate(&User{}, &Order{})
	
	// ... rest of the code examples
}
```

#### Basic Queries

```go
// Create a new user
user := User{Name: "John Doe", Email: "john@example.com", Age: 30, Status: "active"}
result := db.Create(&user)
if result.Error != nil {
	log.Printf("Error creating user: %v", result.Error)
}
fmt.Printf("Created user with ID: %d\n", user.ID)

// Read user by ID
var fetchedUser User
if err := db.First(&fetchedUser, user.ID).Error; err != nil {
	log.Printf("Error fetching user: %v", err)
} else {
	fmt.Printf("Fetched user: %+v\n", fetchedUser)
}

// Update user
db.Model(&fetchedUser).Update("Age", 31)

// Delete user (soft delete)
db.Delete(&fetchedUser)
```

#### Conditional Queries with `Where`

```go
var users []User

// Simple condition
db.Where("age > ?", 25).Find(&users)
fmt.Printf("Users older than 25: %d\n", len(users))

// Multiple conditions with AND
db.Where("age > ? AND status = ?", 20, "active").Find(&users)

// Using struct (non-zero values only)
db.Where(&User{Status: "active", Age: 30}).Find(&users)

// Using map (includes zero values)
db.Where(map[string]interface{}{"status": "inactive", "age": 0}).Find(&users)

// IN clause
db.Where("id IN ?", []uint{1, 2, 3}).Find(&users)

// LIKE clause
db.Where("name LIKE ?", "%John%").Find(&users)
```

#### Advanced Conditional Logic with `Or` and `Not`

```go
// OR condition - note the nested Where for correct parentheses
db.Where(db.Where("age > ?", 30).Or("status = ?", "inactive")).Find(&users)

// NOT condition
db.Not("status = ?", "inactive").Find(&users)

// Complex combination
db.Where("age BETWEEN ? AND ?", 18, 35).
   Where(db.Where("status = ?", "active").Or("name LIKE ?", "%Admin%")).
   Find(&users)
```

#### Subqueries

```go
// Subquery in WHERE clause
subQuery := db.Model(&Order{}).Select("user_id").Where("amount > ?", 100)
db.Where("id IN (?)", subQuery).Find(&users)

// Subquery using IN with EXISTS
db.Where("EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)", 500).Find(&users)
```

#### Conditional Expressions with `CASE`

```go
type UserStatus struct {
	Name   string
	Status string
	Level  string
}

var userStatuses []UserStatus

db.Model(&User{}).
	Select("name, status, CASE WHEN age < 18 THEN 'Minor' WHEN age BETWEEN 18 AND 65 THEN 'Adult' ELSE 'Senior' END as level").
	Find(&userStatuses)
```

#### Dynamic Queries with `if`, `else`, `switch`

```go
func searchUsers(db *gorm.DB, name string, minAge, maxAge int, status string, sortBy string) ([]User, error) {
	var users []User
	query := db.Model(&User{})
	
	// Dynamic conditions using if
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	
	if minAge > 0 && maxAge > 0 {
		query = query.Where("age BETWEEN ? AND ?", minAge, maxAge)
	} else if minAge > 0 {
		query = query.Where("age >= ?", minAge)
	} else if maxAge > 0 {
		query = query.Where("age <= ?", maxAge)
	}
	
	// Using switch for sorting
	switch sortBy {
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	case "age_asc":
		query = query.Order("age ASC")
	case "age_desc":
		query = query.Order("age DESC")
	case "created_desc":
		fallthrough
	default:
		query = query.Order("created_at DESC")
	}
	
	// Handle status with if-else
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Find(&users).Error
	return users, err
}
```

#### Using Scopes for Reusable Dynamic Queries

```go
// Scope for active users
func ActiveUsers(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", "active")
}

// Scope for age range
func AgeRange(min, max int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age BETWEEN ? AND ?", min, max)
	}
}

// Usage
db.Scopes(ActiveUsers, AgeRange(18, 65)).Find(&users)
```

---

### วิธีใช้งาน module นี้

To use GORM with PostgreSQL in your project:

1.  **Import the required packages** in your Go files:
    ```go
    import (
        "gorm.io/gorm"
        "gorm.io/driver/postgres"
    )
    ```

2.  **Define your models** as Go structs with GORM tags for column names, constraints, and relationships.

3.  **Establish a connection** using `gorm.Open()` with a PostgreSQL DSN (Data Source Name).

4.  **Use the `db` instance** to perform CRUD and query operations.

5.  **Always close the database connection** when your application terminates (optional, as GORM manages connection pooling).

---

### การติดตั้ง

1.  **Install Go 1.16+** from [golang.org](https://golang.org/dl/).
2.  **Install PostgreSQL 12+** from [postgresql.org](https://www.postgresql.org/download/).
3.  **Create a database** for your application:
    ```sql
    CREATE DATABASE gorm_demo;
    ```
4.  **Install Go packages**:
    ```bash
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/postgres
    ```
5.  **(Optional) Use Docker** for PostgreSQL:
    ```yaml
    # docker-compose.yml
    version: '3.8'
    services:
      postgres:
        image: postgres:16
        environment:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: yourpassword
          POSTGRES_DB: gorm_demo
        ports:
          - "5432:5432"
        volumes:
          - postgres_data:/var/lib/postgresql/data
    volumes:
      postgres_data:
    ```
    Run `docker-compose up -d` to start PostgreSQL in a container[reference:4].

---

### การตั้งค่า configuration

GORM provides several configuration options when opening a database connection:

```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    // Skip default transaction for performance
    SkipDefaultTransaction: true,
    
    // Naming strategy for tables and columns
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "tbl_",     // Table name prefix
        SingularTable: true,       // Use singular table names
        NoLowerCase:   false,      // Skip snake_casing
    },
    
    // Logger configuration
    Logger: logger.Default.LogMode(logger.Info),
    
    // Disable foreign key constraint when migrating
    DisableForeignKeyConstraintWhenMigrating: true,
})

// Configure connection pool
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

### การรวมกับ GROM

**Note:** *GROM* might be a typo; assuming you meant **GORM**.

GORM is designed to be easily integrated with any Go application. The integration involves:

1.  **Creating a database package** that initializes and provides the GORM instance.
2.  **Injecting the DB instance** into your service or handler structs.
3.  **Using dependency injection** for better testability.

Example integration structure:

```go
// database/db.go
package database

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

func InitDB(dsn string) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// services/user_service.go
package services

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) GetActiveUsers() ([]User, error) {
    var users []User
    err := s.db.Where("status = ?", "active").Find(&users).Error
    return users, err
}
```

---

### การใช้งานจริง

**Scenario:** Building a REST API endpoint for searching users with multiple optional filters.

```go
// GET /users?name=john&min_age=18&max_age=65&status=active&sort=age_asc
func GetUsersHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parse query parameters
        name := r.URL.Query().Get("name")
        minAge := r.URL.Query().Get("min_age")
        maxAge := r.URL.Query().Get("max_age")
        status := r.URL.Query().Get("status")
        sort := r.URL.Query().Get("sort")
        
        // Build dynamic query
        query := db.Model(&User{})
        
        if name != "" {
            query = query.Where("name LIKE ?", "%"+name+"%")
        }
        
        if minAge != "" {
            query = query.Where("age >= ?", minAge)
        }
        
        if maxAge != "" {
            query = query.Where("age <= ?", maxAge)
        }
        
        if status != "" {
            query = query.Where("status = ?", status)
        }
        
        // Apply sorting
        switch sort {
        case "name_asc":
            query = query.Order("name ASC")
        case "age_asc":
            query = query.Order("age ASC")
        case "age_desc":
            query = query.Order("age DESC")
        default:
            query = query.Order("created_at DESC")
        }
        
        var users []User
        if err := query.Find(&users).Error; err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        json.NewEncoder(w).Encode(users)
    }
}
```

---

### ตารางสรุป Components

| Component | Description | Key Features |
|-----------|-------------|--------------|
| **PostgreSQL Driver** | `gorm.io/driver/postgres` | Enables GORM to communicate with PostgreSQL databases |
| **GORM Core** | `gorm.io/gorm` | Main ORM library with chainable API, callbacks, and migrations |
| **Model Structs** | User-defined Go types | Map to database tables using struct tags |
| **Chainable API** | Methods like `Where`, `Select`, `Order` | Build queries incrementally with method chaining[reference:5] |
| **Scopes** | Reusable query functions | Encapsulate common query logic for reuse[reference:6] |
| **Subquery Support** | Nested queries using `*gorm.DB` | Build complex queries programmatically |
| **Raw SQL** | `db.Raw()` and `db.Exec()` | Execute custom SQL when ORM is insufficient |

---

### แบบฝึกหัดท้าย module (5 ข้อ)

1.  **Basic Query with Conditions:** Write a GORM query to find all users with age between 20 and 30 (inclusive) who are active, ordered by name ascending.

2.  **Dynamic Query with `switch`:** Create a function that accepts a `sortBy` parameter (`"name"`, `"age"`, `"created_at"`) and a `desc` boolean, returning users sorted accordingly.

3.  **Subquery Practice:** Write a query using a subquery to find all users who have placed orders with a total amount exceeding 500. Use the `User` and `Order` models from the examples.

4.  **Conditional Logic with `CASE`:** Use a `CASE` expression to categorize users into `"Young"` (age < 25), `"Middle"` (age 25-50), and `"Senior"` (age > 50), and count how many users fall into each category.

5.  **Scope Implementation:** Create a GORM scope called `RecentlyActive` that filters users who have been active within the last 7 days. Then use this scope in combination with an age filter to fetch relevant users.

---

### แหล่งอ้างอิง

1.  PostgreSQL SELECT Documentation [8†L5-L46]
2.  PostgreSQL Subquery Expressions (EXISTS, IN) [9†L2-L28]
3.  PostgreSQL Conditional Expressions (CASE) [10†L8-L30]
4.  GORM Chainable API Documentation [11†L4-L30]
5.  GORM Where Conditions and SQL Injection Prevention [12†L4-L27]
6.  GORM Conditional Query Best Practices [13†L4-L27]
7.  GORM Dynamic Queries with Scopes and Clauses [14†L4-L32]
8.  GORM Nested Queries with Preload [15†L5-L19]
9.  GORM PostgreSQL Connection Example [16†L3-L23]
10. GORM CRUD with PostgreSQL and Docker [17†L3-L27]