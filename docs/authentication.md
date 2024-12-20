# Authentication

The platform uses JSON Web Tokens (JWT) for authentication.

## JWT

Registered claims:

### jti (JWT ID):

It's a unique identifier for each token
Perfect for token revocation by storing in a database
You can blacklist specific tokens by their jti
Useful for scenarios like user logout or security breaches
Common pattern: Store jti in Redis with token's expiry time

### aud (Audience):

Specifies which services/applications can use this token
Example: If you have multiple services (mobile app, web app, API), you can restrict tokens to specific platforms
Helps prevent token misuse across different parts of your system
Your services verify if they're the intended audience before accepting the token

## iss (Issuer):

Identifies which service generated the token
Useful in microservices architecture where multiple services issue tokens
Helps track token origin
Common value would be your service name or domain

## exp (Expiration Time):

Unix timestamp indicating when the token expires
Required to prevent tokens from being valid indefinitely
Best practice: Short-lived tokens (15-60 minutes)
Can be used with refresh tokens for better security

## sub (Subject):

Identifies the user/entity the token represents
Usually contains user ID or unique identifier
Essential for user-specific operations
Should be immutable for the user

## Security Best Practices:

1. Transport Security:
   - Always use HTTPS/TLS
   - Set secure and httpOnly cookie flags

2. Token Storage:
   - Store tokens securely (httpOnly cookies preferred)
   - Never store in localStorage due to XSS risks
   - Implement proper CSRF protection

3. Token Lifecycle:
   - Implement token rotation
   - Use refresh token patterns
   - Clear tokens on logout
   - Maintain token blacklist/revocation list
