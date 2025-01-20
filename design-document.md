| authors                           | state |
| --------------------------------- | ----- |
| Joseph Eid (josephceid@gmail.com) | draft |

# RFD - File Browser Application

## What

This document outlines the design for a secure file browser web application that allows authenticated users to browse directory contents on a remote server.

## Why

The purpose for the implementation of this application is to showcase that the author is capable of implementing a full-stack web application that demonstrates the following:

-   The ability to use Go to create an API which serves a React frontend
-   The ability to use React to implement a UI
-   The ability to implement authentication and authorization
-   The ability to take encryption and web security best practices into account

## Details

## Proposed User Experience

### User Stories

-   Display the filename, type (file or directory), and human-readable size for files.
-   Add support for filtering the directory contents based on filename. Filtering should be performed client-side, and a simple substring match is sufficient.
-   Add support for sorting directory contents based on filename, type, and size.
-   Include breadcrumbs that show the current location in the directory. The breadcrumbs should be clickable for easy navigation to parent directories.
-   Implement URL navigation. The state of the app should be encoded in the URL. No state should be lost upon a page refresh.

### Wireframes

https://www.figma.com/design/3Y8OlARrN0USs8pE7h4h0z/StoutExplorer?node-id=0-1&p=f&t=1oXnKeu9mxSpI5x9-0

### User Flows

1. **Authentication Flow**

    - Unauthenticated users who try to visit a path are redirected to /login?redirect={path}
        - The redirect path should be validated in the frontend to check that it starts with / (a valid path) and that it does not contain a 'javascript:' or 'vbscript:' protocol to prevent XSS attacks
        - In the backend once the redirect has occurred, the path will be validated as part of the request validation process using `filepath.Clean` to prevent directory traversal
    - After successful login, users return to their originally requested path
    - Session status persists across page refreshes
    - Logout option available in user menu

2. **Browsing Flow**
    - Users can navigate directories via:
        - Clicking directory names
        - Using breadcrumb navigation
        - Direct URL entry
    - File listing shows name, type, and size
    - Client-side filtering via search box
    - Sorting by name, type, or size

## Proposed API structure

### API Endpoints

```
GET /api/v1/browse?path={path} - Returns the contents of the directory at the given path. Making the path a query parameter means only a single parameter needs to be encoded in the frontend and decoded in the backend.
Response:
{
  "name": string,
  "type": "dir",
  "size": number,
  "contents": [
    {
      "name": string,
      "type": "file"|"dir",
      "size": number
    }
  ]
}

POST /api/v1/auth/login - Logs in the user and returns a session token
Request:
{
  "username": string,
  "password": string
}
Response:
{
  "username": string
}

POST /api/v1/auth/logout - Logs out the user and invalidates the session token
Request: None
Response: None
```

## URL Structure

-   `/login` - Authentication page
-   `/browse` - Root directory view
-   `/browse/{path}` - Specific directory view - path will need to be url-encoded as it is likely to contain special characters
-   URL query parameters:
    -   `sort`: Column to sort by (name|type|size)
    -   `direction`: Sort direction (asc|desc)
    -   `filter`: Search/filter string - should be url-encoded as it may contain special characters

## Security Implementation

1. **TLS Configuration**

    - Self-signed certificates for development
    - Strong cipher suites only
    - TLS 1.2+ only: TLS 1.0 and 1.1 have been deprecated by various major browsers (Chrome, Firefox, Edge and Safari) and both versions do not support modern cryptographic algorithms and ciphers. We have chosen to use TLS 1.2 as the minimum version as it is the most widely supported secure version of TLS and some older versions of major browsers lack the ability to use TLS 1.3.
    - HTTP/2 enabled

2. **Authentication & Session Management**
    - Session-based authentication using secure cookies
    - Sessions stored in memory with the following structure:

```go
type Session struct {
    ID        string    // 256 bit base64 encoded string
    UserID    string
    CreatedAt time.Time
    ExpiresAt time.Time
}
```

-   Secure cookie configuration:

```go
cookie := &http.Cookie{
    Name:     "session_id",
    Value:    sessionID,   // 256 bit base64 encoded string
    Path:     "/",
    HttpOnly: true,        // Prevents JavaScript access
    Secure:   true,        // Only sent over HTTPS
    SameSite: http.SameSiteStrictMode, // For now there are no plans to support cross-site requests or third-party integrations and all requests will be coming from the same origin, so we can use Strict mode
    MaxAge:   86400,       // 24 hours
}
```

-   Session flow:
    1.  User logs in with credentials
    2.  Server validates credentials and creates new session
    3.  Session ID stored in secure cookie - this will be a 256 bit base64 encoded string, this level of entropy makes it infeasible to brute force the session ID and future-proofs against quantum computing threats
    4.  Subsequent requests include cookie automatically
    5.  Server validates session ID on each request
-   Security measures:
    -   Password hashing using bcrypt - decided to use bcrypt because it is widely used, has a simple API and has built in salt to prevent rainbow table attacks
    -   Automatic session expiration after 24 hours
    -   Session invalidation on logout
    -   Timing attacks will be mitigated by always taking the same amount of time to perform a login attempt, this will be done by having a fall-back password comparison for non-existant users, which will return an invalid login attempt, the same way an incorrect password for a valid user would.

3. **Additional Security Measures**
    - HTTP Security Headers:
        - Content-Security-Policy - set to default-src 'self' to prevent XSS attacks
        - X-Frame-Options - set to DENY to prevent clickjacking
        - X-Content-Type-Options - set to nosniff to prevent MIME type sniffing
        - Strict-Transport-Security - set to max-age=31536000; includeSubDomains to prevent SSL/TLS attacks
    - Input validation on all API endpoints
    - Path traversal protection (see Implementation Details)
    - Secure cookie attributes (HttpOnly, Secure, SameSite - see http.Cookie struct above)

## Implementation Details

1. **File System Access**

    - Root directory will be configurable via environment variable
    - Path sanitization to prevent directory traversal, this can be achieved by using the `filepath` package in Go, specifically the `filepath.Clean` function.
    - Error handling for permission issues and non-existent paths which should redirect to the login or the 404 page depending on whether or not the user is authenticated

2. **Linting and Formatting**

    - Go code will be linted using `golangci-lint`
    - Go code will be formatted using `gofmt`
    - React code will be linted using `eslint`
    - React code will be formatted using `prettier`

3. **Mocked In-Memory Authentication Table**
    - The authentication table will be mocked in-memory in the Go backend
    - Usernames and hashed passwords will be stored in this mocked authentication table
    - The authentication table will also store the session data

## Proposed Testing Strategy

### Backend

-   Unit tests surrounding the API endpoints
-   Integration tests surrounding the API endpoints
-   Security tests surrounding the API endpoints

### Frontend

-   Unit tests surrounding the UI using React Testing Library and Jest
-   End-to-end tests based on defined user flows in the Proposed User Experience section

## Future Improvements

-   Rate limiting on login attempts (5 attempts per 15 minutes) - prevents brute force attacks
-   If we introduce more POST endpoints in the future, we should implement CSRF protection using anti-CSRF tokens - this can be achieved by using the `gorilla/csrf` package in Go. Both the login and the logout endpoints will require a CSRF token to be sent in the request body as they are both state changing operations. The GET request used for file browsing will not require a CSRF token as it is a read-only operation.
