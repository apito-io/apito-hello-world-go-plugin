package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	sdk "github.com/apito-io/go-apito-plugin-sdk"
)

// debugContextValues safely prints all known context values without panicking
func debugContextValues(ctx context.Context) {
	log.Printf("🔍 [hc-hello-world-plugin] === Context Debug Information ===")

	// List of known context keys to check
	knownKeys := []string{
		"project_id",
		"plugin_id",
		"cache",
		"user_id",
		"tenant_id",
		"request_id",
		"session_id",
		"application_id",
		"database",
		"config",
		"selectionSet",
		"variables",
	}

	for _, key := range knownKeys {
		val := ctx.Value(key)
		if val == nil {
			log.Printf("🔍 [hc-hello-world-plugin] %s: <nil>", key)
			continue
		}

		// Print the raw value first
		log.Printf("🔍 [hc-hello-world-plugin] %s (raw): %v (type: %T)", key, val, val)

		// Try safe type assertions for common types
		switch v := val.(type) {
		case string:
			log.Printf("🔍 [hc-hello-world-plugin] %s (string): %s", key, v)
		case int:
			log.Printf("🔍 [hc-hello-world-plugin] %s (int): %d", key, v)
		case int64:
			log.Printf("🔍 [hc-hello-world-plugin] %s (int64): %d", key, v)
		case bool:
			log.Printf("🔍 [hc-hello-world-plugin] %s (bool): %t", key, v)
		case map[string]interface{}:
			log.Printf("🔍 [hc-hello-world-plugin] %s (map): %+v", key, v)
		default:
			// For unknown types, just print the value and type
			log.Printf("🔍 [hc-hello-world-plugin] %s (unknown type %T): %v", key, v, v)
		}
	}

	log.Printf("🔍 [hc-hello-world-plugin] === End Context Debug ===")
}

func main() {
	log.Printf("🎯 [hc-hello-world-plugin] Starting plugin initialization...")

	// Start plugin normally - delve debugging is handled externally by the host
	startNormalPlugin()
}

// GraphQL Resolvers - Same business logic, much cleaner setup!

func helloWorldResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {

	log.Printf("🚀 [hc-hello-world-plugin] helloWorldResolver called with args: %+v", rawArgs)

	// Safe way to debug and print all context values without panicking
	debugContextValues(ctx)

	// NEW: Access context data passed from the host using SDK helpers
	pluginID := sdk.GetPluginID(rawArgs)
	projectID := sdk.GetProjectID(rawArgs)
	userID := sdk.GetUserID(rawArgs)
	tenantID := sdk.GetTenantID(rawArgs)

	// Also demonstrate direct context access
	pluginIDFromCtx := sdk.GetPluginIDFromContext(ctx)

	log.Printf("🔐 [hc-hello-world-plugin] Context Data from Host:")
	log.Printf("   - Plugin ID: %s", pluginID)
	log.Printf("   - Project ID: %s", projectID)
	log.Printf("   - User ID: %s", userID)
	log.Printf("   - Tenant ID: %s", tenantID)
	log.Printf("   - Plugin ID from Context: %s", pluginIDFromCtx)

	// Get all context data for debugging
	allContextData := sdk.GetAllContextData(rawArgs)
	log.Printf("🔍 [hc-hello-world-plugin] All Context Data: %+v", allContextData)

	// Use the SDK's automatic argument parsing based on field definition
	args := sdk.ParseArgsForResolver("helloWorldQuery", rawArgs)

	log.Printf("📝 [hc-hello-world-plugin] Parsed args: %+v", args)

	var result strings.Builder
	result.WriteString("Hello World Plugin Response (SDK Version with Auto-Parsing):\n")

	// Add context information to the response
	if pluginID != "" {
		result.WriteString(fmt.Sprintf("🔐 Plugin ID: %s\n", pluginID))
	}
	if projectID != "" {
		result.WriteString(fmt.Sprintf("🏗️  Project ID: %s\n", projectID))
	}
	if userID != "" {
		result.WriteString(fmt.Sprintf("👤 User ID: %s\n", userID))
	}
	if tenantID != "" {
		result.WriteString(fmt.Sprintf("🏢 Tenant ID: %s\n", tenantID))
	}

	// Handle name parameter - now type-safe!
	name := sdk.GetStringArg(args, "name", "World")
	log.Printf("👋 [hc-hello-world-plugin] Greeting name: %s", name)
	result.WriteString(fmt.Sprintf("Hello, %s!\n", name))

	// Handle object parameter - automatically parsed!
	if obj := sdk.GetObjectArg(args, "object"); len(obj) > 0 {
		log.Printf("📦 [hc-hello-world-plugin] Object parameter received: %+v", obj)
		result.WriteString("Object received: ")
		objName := sdk.GetStringArg(obj, "name")
		objAge := sdk.GetIntArg(obj, "age")
		result.WriteString(fmt.Sprintf("name=%s age=%d\n", objName, objAge))
	}

	// Handle arrayofObjects parameter - automatically parsed!
	if arrObjs := sdk.GetArrayArg(args, "arrayofObjects"); len(arrObjs) > 0 {
		log.Printf("📊 [hc-hello-world-plugin] Array of objects received: %d items", len(arrObjs))
		result.WriteString("Array of Objects received:\n")
		for i, obj := range arrObjs {
			if objMap, ok := obj.(map[string]interface{}); ok {
				objName := sdk.GetStringArg(objMap, "name")
				objAge := sdk.GetIntArg(objMap, "age")
				result.WriteString(fmt.Sprintf("  Object %d: name=%s age=%d\n", i+1, objName, objAge))
			}
		}
	}

	log.Printf("✅ [hc-hello-world-plugin] helloWorldResolver completed successfully")
	return result.String(), nil
}

func processComplexDataResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("processComplexData", rawArgs)

	var result strings.Builder
	result.WriteString("Processing complex data (SDK Version with Auto-Parsing):\n")

	// Process single user object - type-safe!
	if user := sdk.GetObjectArg(args, "user"); len(user) > 0 {
		result.WriteString("User: ")
		id := sdk.GetIntArg(user, "id")
		name := sdk.GetStringArg(user, "name")
		email := sdk.GetStringArg(user, "email")
		age := sdk.GetIntArg(user, "age")
		active := sdk.GetBoolArg(user, "active")
		result.WriteString(fmt.Sprintf("ID=%d Name=%s Email=%s Age=%d Active=%t\n", id, name, email, age, active))
	}

	// Process array of strings (tags) - automatically converted!
	if tagSlice, ok := args["tags"].([]string); ok {
		result.WriteString("Tags: ")
		for i, tag := range tagSlice {
			result.WriteString(tag)
			if i < len(tagSlice)-1 {
				result.WriteString(", ")
			}
		}
		result.WriteString("\n")
	}

	// Process array of integers (numbers) - automatically converted!
	if numberSlice, ok := args["numbers"].([]int); ok {
		result.WriteString("Numbers: ")
		for i, num := range numberSlice {
			result.WriteString(fmt.Sprintf("%d", num))
			if i < len(numberSlice)-1 {
				result.WriteString(", ")
			}
		}
		result.WriteString("\n")
	}

	// Process array of user objects (users) - automatically parsed!
	if users := sdk.GetArrayArg(args, "users"); len(users) > 0 {
		result.WriteString("Users:\n")
		for i, user := range users {
			if userMap, ok := user.(map[string]interface{}); ok {
				id := sdk.GetIntArg(userMap, "id")
				name := sdk.GetStringArg(userMap, "name")
				email := sdk.GetStringArg(userMap, "email")
				result.WriteString(fmt.Sprintf("  User %d: ID=%d Name=%s Email=%s\n", i+1, id, name, email))
			}
		}
	}

	// Process array of optional user objects (optionalUsers)
	if optionalUsers := sdk.GetArrayArg(args, "optionalUsers"); len(optionalUsers) > 0 {
		result.WriteString("Optional Users:\n")
		for i, user := range optionalUsers {
			if user != nil {
				if userMap, ok := user.(map[string]interface{}); ok {
					name := sdk.GetStringArg(userMap, "name")
					email := sdk.GetStringArg(userMap, "email")
					result.WriteString(fmt.Sprintf("  Optional User %d: Name=%s Email=%s\n", i+1, name, email))
				}
			} else {
				result.WriteString(fmt.Sprintf("  Optional User %d: null\n", i+1))
			}
		}
	}

	return result.String(), nil
}

func sayHelloResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("sayHelloMutation", rawArgs)

	// Type-safe argument extraction with default value
	message := sdk.GetStringArg(args, "message", "Hello!")

	return fmt.Sprintf("Plugin says: %s (from hc-hello-world-plugin using SDK with Auto-Parsing)", message), nil
}

// REST Handlers - Much simpler than managing protobuf structs!

func helloRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"message":   "Hello World from REST API (SDK Version)!",
		"timestamp": time.Now().Format(time.RFC3339),
		"plugin":    "hc-hello-world-plugin",
		"version":   "2.0.0-sdk",
	}, nil
}

func customHelloRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	name := "World"
	message := "Hello"

	if nameArg, ok := args["name"].(string); ok && nameArg != "" {
		name = nameArg
	}
	if msgArg, ok := args["message"].(string); ok && msgArg != "" {
		message = msgArg
	}

	return map[string]interface{}{
		"greeting": fmt.Sprintf("%s, %s! (SDK Version)", message, name),
		"plugin":   "hc-hello-world-plugin",
		"version":  "2.0.0-sdk",
	}, nil
}

func statusRESTHandler(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
		"status":  "running",
		"version": "2.0.0-sdk",
		"sdk":     "github.com/apito-io/go-apito-plugin-sdk",
		"features": []string{
			"GraphQL Queries",
			"GraphQL Mutations",
			"REST APIs",
			"Custom Functions",
		},
	}, nil
}

// Custom Functions

func customFunction(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return "Custom function executed successfully (SDK Version)", nil
}

// getUserProfileResolver demonstrates returning a complex User object
func getUserProfileResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	log.Printf("🚀 [hc-hello-world-plugin] getUserProfileResolver called with args: %+v", rawArgs)

	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("getUserProfile", rawArgs)
	userID := sdk.GetStringArg(args, "userId", "default-user")

	log.Printf("👤 [hc-hello-world-plugin] Fetching user profile for ID: %s", userID)

	// Return a complex User object structure with nested objects
	user := map[string]interface{}{
		"id":       userID,
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"username": "johndoe",
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
			"state":  "NY",
			"zip":    "10001",
		},
		"tags": []interface{}{
			map[string]interface{}{
				"key": "department",
				"val": "engineering",
			},
			map[string]interface{}{
				"key": "level",
				"val": "senior",
			},
			map[string]interface{}{
				"key": "team",
				"val": "backend",
			},
		},
		"active":    true,
		"createdAt": time.Now().Format(time.RFC3339),
	}

	log.Printf("[NESTED-OBJECT-DEBUG] [PLUGIN] getUserProfileResolver returning user: %+v", user)
	if address, exists := user["address"]; exists {
		log.Printf("[NESTED-OBJECT-DEBUG] [PLUGIN] User address: %+v (type: %T)", address, address)
	}
	return user, nil
}

// getUsersResolver demonstrates returning an array of User objects
func getUsersResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	log.Printf("🚀 [hc-hello-world-plugin] getUsersResolver called with args: %+v", rawArgs)

	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("getUsers", rawArgs)
	limit := sdk.GetIntArg(args, "limit", 10)
	offset := sdk.GetIntArg(args, "offset", 0)
	activeFilter := sdk.GetBoolArg(args, "active", true)

	log.Printf("📊 [hc-hello-world-plugin] Query params - limit: %d, offset: %d, active: %t", limit, offset, activeFilter)

	// Generate sample users array with nested objects
	users := []interface{}{
		map[string]interface{}{
			"id":       "1",
			"name":     "John Doe",
			"email":    "john.doe@example.com",
			"username": "johndoe",
			"address": map[string]interface{}{
				"street": "123 Main St",
				"city":   "New York",
				"state":  "NY",
				"zip":    "10001",
			},
			"tags": []interface{}{
				map[string]interface{}{"key": "department", "val": "engineering"},
				map[string]interface{}{"key": "level", "val": "senior"},
			},
			"active":    true,
			"createdAt": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		map[string]interface{}{
			"id":       "2",
			"name":     "Jane Smith",
			"email":    "jane.smith@example.com",
			"username": "janesmith",
			"address": map[string]interface{}{
				"street": "456 Oak Ave",
				"city":   "Los Angeles",
				"state":  "CA",
				"zip":    "90210",
			},
			"tags": []interface{}{
				map[string]interface{}{"key": "department", "val": "design"},
				map[string]interface{}{"key": "level", "val": "mid"},
			},
			"active":    false,
			"createdAt": time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
		},
		map[string]interface{}{
			"id":       "3",
			"name":     "Bob Johnson",
			"email":    "bob.johnson@example.com",
			"username": "bobjohnson",
			"address": map[string]interface{}{
				"street": "789 Pine Rd",
				"city":   "Chicago",
				"state":  "IL",
				"zip":    "60601",
			},
			"tags": []interface{}{
				map[string]interface{}{"key": "department", "val": "marketing"},
				map[string]interface{}{"key": "level", "val": "junior"},
			},
			"active":    true,
			"createdAt": time.Now().Add(-72 * time.Hour).Format(time.RFC3339),
		},
	}

	// Apply active filter
	var filteredUsers []interface{}
	for _, user := range users {
		userMap := user.(map[string]interface{})
		if userMap["active"].(bool) == activeFilter {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(filteredUsers) {
		start = len(filteredUsers)
	}
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	paginatedUsers := filteredUsers[start:end]

	log.Printf("[NESTED-OBJECT-DEBUG] [PLUGIN] getUsersResolver returning %d users", len(paginatedUsers))
	for i, user := range paginatedUsers {
		log.Printf("[NESTED-OBJECT-DEBUG] [PLUGIN] User %d: %+v", i, user)
		if userMap, ok := user.(map[string]interface{}); ok {
			if address, exists := userMap["address"]; exists {
				log.Printf("[NESTED-OBJECT-DEBUG] [PLUGIN] User %d address: %+v (type: %T)", i, address, address)
			}
		}
	}
	return paginatedUsers, nil
}

// getProductsPaginatedResolver demonstrates returning a paginated response
func getProductsPaginatedResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	log.Printf("🚀 [hc-hello-world-plugin] getProductsPaginatedResolver called with args: %+v", rawArgs)

	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("getProductsPaginated", rawArgs)
	page := sdk.GetIntArg(args, "page", 1)
	pageSize := sdk.GetIntArg(args, "pageSize", 5)
	category := sdk.GetStringArg(args, "category", "")

	log.Printf("📄 [hc-hello-world-plugin] Pagination params - page: %d, pageSize: %d, category: %s", page, pageSize, category)

	// Generate sample products
	allProducts := []interface{}{
		map[string]interface{}{
			"id":          "1",
			"name":        "Laptop",
			"description": "High-performance laptop",
			"price":       999.99,
			"stock":       10,
			"tags":        []string{"electronics", "computers"},
			"categories":  []string{"electronics", "office"},
		},
		map[string]interface{}{
			"id":          "2",
			"name":        "Coffee Mug",
			"description": "Ceramic coffee mug",
			"price":       12.99,
			"stock":       50,
			"tags":        []string{"kitchen", "drinkware"},
			"categories":  []string{"home", "kitchen"},
		},
		map[string]interface{}{
			"id":          "3",
			"name":        "Book",
			"description": "Programming book",
			"price":       29.99,
			"stock":       25,
			"tags":        []string{"education", "programming"},
			"categories":  []string{"books", "education"},
		},
	}

	// Filter by category if provided
	var filteredProducts []interface{}
	if category != "" {
		for _, product := range allProducts {
			productMap := product.(map[string]interface{})
			categories := productMap["categories"].([]string)
			for _, cat := range categories {
				if cat == category {
					filteredProducts = append(filteredProducts, product)
					break
				}
			}
		}
	} else {
		filteredProducts = allProducts
	}

	// Calculate pagination
	total := len(filteredProducts)
	totalPages := (total + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	// Get page items
	var pageItems []interface{}
	if offset < total {
		end := offset + pageSize
		if end > total {
			end = total
		}
		pageItems = filteredProducts[offset:end]
	}

	// Return paginated response structure (updated for v0.1.6 simplified structure)
	// Convert products to string array for simplified pagination
	var itemStrings []string
	for _, item := range pageItems {
		if productMap, ok := item.(map[string]interface{}); ok {
			itemStrings = append(itemStrings, fmt.Sprintf("%s - %s ($%.2f)",
				productMap["name"], productMap["description"], productMap["price"]))
		}
	}

	response := map[string]interface{}{
		"items":           itemStrings,
		"totalCount":      total,
		"pageSize":        pageSize,
		"currentPage":     page,
		"totalPages":      totalPages,
		"hasNextPage":     page < totalPages,
		"hasPreviousPage": page > 1,
		"success":         true,
		"message":         fmt.Sprintf("Retrieved %d products", len(pageItems)),
	}

	log.Printf("✅ [hc-hello-world-plugin] getProductsPaginatedResolver completed")
	return response, nil
}

// createUserResolver demonstrates returning a wrapped response for mutations
func createUserResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	log.Printf("🚀 [hc-hello-world-plugin] createUserResolver called with args: %+v", rawArgs)

	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("createUser", rawArgs)
	input := sdk.GetObjectArg(args, "input")

	name := sdk.GetStringArg(input, "name", "")
	email := sdk.GetStringArg(input, "email", "")
	username := sdk.GetStringArg(input, "username", "")

	log.Printf("👤 [hc-hello-world-plugin] Creating user - name: %s, email: %s, username: %s", name, email, username)

	// Validate input
	if name == "" || email == "" || username == "" {
		return map[string]interface{}{
			"success": false,
			"message": "Name, email, and username are required",
			"data":    nil,
			"errors": []interface{}{
				map[string]interface{}{
					"code":    "VALIDATION_ERROR",
					"message": "Missing required fields",
					"field":   "name,email,username",
					"details": []string{"All fields are required for user creation"},
				},
			},
		}, nil
	}

	// Create new user (simulated)
	newUser := map[string]interface{}{
		"id":        fmt.Sprintf("user_%d", time.Now().Unix()),
		"name":      name,
		"email":     email,
		"username":  username,
		"active":    true,
		"createdAt": time.Now().Format(time.RFC3339),
	}

	// Return success response
	response := map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"data":    newUser,
		"errors":  nil,
	}

	log.Printf("✅ [hc-hello-world-plugin] createUserResolver completed successfully")
	return response, nil
}

// getProductResolver demonstrates returning a single Product object
func getProductResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	log.Printf("🚀 [hc-hello-world-plugin] getProductResolver called with args: %+v", rawArgs)

	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("getProduct", rawArgs)
	productID := sdk.GetStringArg(args, "productId", "default-product")

	log.Printf("📦 [hc-hello-world-plugin] Fetching product for ID: %s", productID)

	// Return a complex Product object structure
	product := map[string]interface{}{
		"id":          productID,
		"name":        "Sample Product",
		"description": "This is a sample product from the plugin",
		"price":       29.99,
		"stock":       100,
		"tags":        []string{"sample", "plugin", "demo"},
		"categories":  []string{"electronics", "gadgets"},
	}

	log.Printf("✅ [hc-hello-world-plugin] getProductResolver completed")
	return product, nil
}

// processBulkTagsResolver demonstrates the new ArrayObjectArg functionality
func processBulkTagsResolver(ctx context.Context, rawArgs map[string]interface{}) (interface{}, error) {
	log.Printf("🚀 [hc-hello-world-plugin] processBulkTagsResolver called with args: %+v", rawArgs)

	// Use the SDK's automatic argument parsing
	args := sdk.ParseArgsForResolver("processBulkTags", rawArgs)
	userId := sdk.GetStringArg(args, "userId", "default-user")

	// ========================================
	// NEW: Demonstrate ArrayObjectArg with GetArrayObjectArg
	// ========================================
	tags := sdk.GetArrayObjectArg(args, "tags")

	log.Printf("👤 [hc-hello-world-plugin] Processing tags for user ID: %s", userId)
	log.Printf("🔍 [hc-hello-world-plugin] Received %d tag objects", len(tags))

	var result strings.Builder
	result.WriteString(fmt.Sprintf("✅ ArrayObjectArg Demo - Processing %d tags for user: %s\n\n", len(tags), userId))

	// Process each tag object using the new SDK helper functions
	for i, tagMap := range tags {
		result.WriteString(fmt.Sprintf("🔖 Tag %d:\n", i+1))

		// Use SDK helper functions for type-safe extraction
		tagID := sdk.GetStringArg(tagMap, "tag_id", "")
		name := sdk.GetStringArg(tagMap, "name", "")
		value := sdk.GetStringArg(tagMap, "value", "")
		weight := sdk.GetFloatArg(tagMap, "weight", 0.0)
		active := sdk.GetBoolArg(tagMap, "active", false)
		metadata := sdk.GetStringArg(tagMap, "metadata", "")

		result.WriteString(fmt.Sprintf("   📛 ID: %s\n", tagID))
		result.WriteString(fmt.Sprintf("   🏷️  Name: %s\n", name))
		result.WriteString(fmt.Sprintf("   💾 Value: %s\n", value))
		result.WriteString(fmt.Sprintf("   ⚖️  Weight: %.2f\n", weight))
		result.WriteString(fmt.Sprintf("   🟢 Active: %t\n", active))
		result.WriteString(fmt.Sprintf("   📋 Metadata: %s\n", metadata))
		result.WriteString("\n")

		log.Printf("📋 [hc-hello-world-plugin] Processed tag %d: ID=%s, Name=%s, Weight=%.2f, Active=%t",
			i+1, tagID, name, weight, active)
	}

	result.WriteString("🎉 ArrayObjectArg processing completed successfully!\n")
	result.WriteString("📊 This demonstrates:\n")
	result.WriteString("   ✅ sdk.ArrayObjectArg() for schema definition\n")
	result.WriteString("   ✅ sdk.GetArrayObjectArg() for typed extraction\n")
	result.WriteString("   ✅ sdk.GetFloatArg() for float type conversion\n")
	result.WriteString("   ✅ Complex object arrays with proper validation\n")

	log.Printf("✅ [hc-hello-world-plugin] processBulkTagsResolver completed successfully")
	return result.String(), nil
}

// startNormalPlugin starts the plugin normally
func startNormalPlugin() {
	log.Printf("🎯 [hc-hello-world-plugin] Starting normal plugin initialization...")

	// Check if debug mode is enabled via environment variable from engine
	debugMode := os.Getenv("PLUGIN_DEBUG_MODE")
	if debugMode == "true" {
		// ANSI color codes for colored output
		const (
			ColorReset  = "\033[0m"
			ColorRed    = "\033[31m"
			ColorGreen  = "\033[32m"
			ColorYellow = "\033[33m"
			ColorBlue   = "\033[34m"
			ColorPurple = "\033[35m"
			ColorCyan   = "\033[36m"
			ColorWhite  = "\033[37m"
			ColorBold   = "\033[1m"
		)

		pid := os.Getpid()

		// Display colorful debug information
		fmt.Printf("\n%s%s╔══════════════════════════════════════════════╗%s\n", ColorBold, ColorCyan, ColorReset)
		fmt.Printf("%s%s║              🐛 DEBUG MODE ENABLED           ║%s\n", ColorBold, ColorCyan, ColorReset)
		fmt.Printf("%s%s╠══════════════════════════════════════════════╣%s\n", ColorBold, ColorCyan, ColorReset)
		fmt.Printf("%s%s║ Plugin PID: %s%s%-4d%s%s                         ║%s\n", ColorBold, ColorCyan, ColorReset, ColorBold+ColorGreen, pid, ColorReset, ColorBold+ColorCyan, ColorReset)
		fmt.Printf("%s%s║ Ready for delve attachment! 🎯             ║%s\n", ColorBold, ColorCyan, ColorReset)
		fmt.Printf("%s%s║                                              ║%s\n", ColorBold, ColorCyan, ColorReset)
		fmt.Printf("%s%s║ %sTo attach delve:%s                          ║%s\n", ColorBold, ColorCyan, ColorYellow, ColorReset+ColorBold+ColorCyan, ColorReset)
		fmt.Printf("%s%s║ %s./runDebug.sh %d%s                     ║%s\n", ColorBold, ColorCyan, ColorGreen, pid, ColorReset+ColorBold+ColorCyan, ColorReset)
		fmt.Printf("%s%s║                                              ║%s\n", ColorBold, ColorCyan, ColorReset)
		fmt.Printf("%s%s║ %sThen in VSCode:%s                          ║%s\n", ColorBold, ColorCyan, ColorYellow, ColorReset+ColorBold+ColorCyan, ColorReset)
		fmt.Printf("%s%s║ %s'Attach to HashiCorp Plugin (Remote Debug)'%s ║%s\n", ColorBold, ColorCyan, ColorPurple, ColorReset+ColorBold+ColorCyan, ColorReset)
		fmt.Printf("%s%s╚══════════════════════════════════════════════╝%s\n\n", ColorBold, ColorCyan, ColorReset)

		log.Printf("🐛 [DEBUG] Plugin PID: %d - Ready for delve attachment!", pid)
	}

	// Initialize the plugin - replaces 50+ lines of handshake/gRPC boilerplate
	plugin := sdk.Init("hc-hello-world-plugin", "2.0.0-sdk", "apito-plugin-key")

	log.Printf("📋 [hc-hello-world-plugin] Registering GraphQL queries...")

	// ========================================
	// SIMPLE STRING EXAMPLE (Original)
	// ========================================

	// Register GraphQL queries - replaces 100+ lines of protobuf struct creation
	plugin.RegisterQuery("helloWorldQueryFahim",
		sdk.FieldWithArgs("String", "Hello World Plugin Query with Arguments", map[string]interface{}{
			"name": sdk.StringArg("Name to greet (optional)"),
			"object": sdk.ObjectArg("Object argument", map[string]interface{}{
				"name": sdk.StringProperty("Object name"),
				"age":  sdk.IntProperty("Object age"),
			}),
			"arrayofObjects": sdk.ListArg("Object", "Array of objects"),
		}),
		helloWorldResolver)

	// ========================================
	// COMPLEX OBJECT EXAMPLES (New)
	// ========================================

	// Define an Address object type (nested object)
	addressType := sdk.NewObjectType("Address", "A user's address").
		AddStringField("street", "Street address", false).
		AddStringField("city", "City", false).
		AddStringField("state", "State", false).
		AddStringField("zip", "Zip code", false).
		Build()

	// Define a Tag object type for the tags array
	tagType := sdk.NewObjectType("Tag", "A tag with key and value").
		AddStringField("key", "Tag key", false).
		AddStringField("val", "Tag value", false).
		Build()

	// Define a User object type with nested objects
	userType := sdk.NewObjectType("User", "A user in the system").
		AddStringField("id", "User ID", false).
		AddStringField("name", "User's full name", false).
		AddStringField("email", "User's email address", true).
		AddStringField("username", "User's username", true).
		AddObjectField("address", "User's address", addressType, true).
		AddObjectListField("tags", "User tags with key-value pairs", tagType, true, false).
		AddBooleanField("active", "Whether the user is active", false).
		AddStringField("createdAt", "When the user was created", true).
		Build()

	// Query that returns a single User object
	plugin.RegisterQuery("getUserProfile",
		sdk.ComplexObjectFieldWithArgs("Get user profile by ID", userType, map[string]interface{}{
			"userId": sdk.StringArg("User ID to fetch"),
		}),
		getUserProfileResolver)

	// Query that returns an array of User objects
	plugin.RegisterQuery("getUsers",
		sdk.ListOfObjectsFieldWithArgs("Get a list of users", userType, map[string]interface{}{
			"limit":  sdk.IntArg("Maximum number of users to return"),
			"offset": sdk.IntArg("Number of users to skip"),
			"active": sdk.BooleanArg("Filter by active status"),
		}),
		getUsersResolver)

	// Define a Product object type with nested structures
	productType := sdk.NewObjectType("Product", "A product in our catalog").
		AddStringField("id", "Product ID", false).
		AddStringField("name", "Product name", false).
		AddStringField("description", "Product description", true).
		AddFloatField("price", "Product price", false).
		AddIntField("stock", "Stock quantity", false).
		AddStringListField("tags", "Product tags", true, false).
		AddStringListField("categories", "Product categories", true, false).
		Build()

	// Query that returns a single product
	plugin.RegisterQuery("getProduct",
		sdk.ComplexObjectFieldWithArgs("Get product by ID", productType, map[string]interface{}{
			"productId": sdk.StringArg("Product ID to fetch"),
		}),
		getProductResolver)

	// Query that returns a paginated list of products
	paginatedProductType := sdk.PaginatedResponseType("Product")
	plugin.RegisterQuery("getProductsPaginated",
		sdk.ComplexObjectFieldWithArgs("Get paginated list of products", paginatedProductType, map[string]interface{}{
			"page":     sdk.IntArg("Page number (1-based)"),
			"pageSize": sdk.IntArg("Number of items per page"),
			"category": sdk.StringArg("Filter by category"),
		}),
		getProductsPaginatedResolver)

	// ========================================
	// REGISTER MUTATIONS
	// ========================================

	// Response wrapper type for mutations
	userResponseType := sdk.ResponseWrapperType("User")

	plugin.RegisterMutation("createUser",
		sdk.ComplexObjectFieldWithArgs("Create a new user", userResponseType, map[string]interface{}{
			"input": sdk.ObjectArg("User creation data", map[string]interface{}{
				"name":     sdk.StringProperty("User's full name"),
				"email":    sdk.StringProperty("User's email address"),
				"username": sdk.StringProperty("User's username"),
			}),
		}),
		createUserResolver)

	// ========================================
	// NEW: ARRAY OBJECT ARGUMENT EXAMPLE
	// ========================================

	// Demonstrates the new ArrayObjectArg functionality
	plugin.RegisterMutation("processBulkTags",
		sdk.FieldWithArgs("String", "Process multiple tag objects - demonstrates ArrayObjectArg", map[string]interface{}{
			"userId": sdk.StringArg("User ID to process tags for"),
			"tags": sdk.ArrayObjectArg("Array of tag objects with structured data", map[string]interface{}{
				"tag_id":   sdk.StringProperty("Tag identifier"),
				"name":     sdk.StringProperty("Tag name"),
				"value":    sdk.StringProperty("Tag value"),
				"weight":   sdk.FloatProperty("Tag weight/importance"),
				"active":   sdk.BooleanProperty("Whether tag is active"),
				"metadata": sdk.StringProperty("Additional metadata"),
			}),
		}),
		processBulkTagsResolver)

	// Register custom functions
	plugin.RegisterFunction("customFunction", customFunction)

	// ========================================
	// REGISTER REST APIS (examples)
	// ========================================

	plugin.RegisterRESTAPI(sdk.RESTEndpoint{
		Method:      "GET",
		Path:        "/hello",
		Description: "Simple hello endpoint",
		Schema:      map[string]interface{}{},
	}, helloRESTHandler)

	plugin.RegisterRESTAPI(sdk.RESTEndpoint{
		Method:      "POST",
		Path:        "/custom-hello",
		Description: "Custom hello endpoint with POST data",
		Schema: map[string]interface{}{
			"name": "string",
		},
	}, customHelloRESTHandler)

	plugin.RegisterRESTAPI(sdk.RESTEndpoint{
		Method:      "GET",
		Path:        "/status",
		Description: "Plugin status endpoint",
		Schema:      map[string]interface{}{},
	}, statusRESTHandler)

	log.Printf("🚀 [hc-hello-world-plugin] Plugin registration complete, starting server...")
	plugin.Serve()
}
