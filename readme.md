# Simple, reasonably secure web portal written in Go
This is a simple web portal which features 2-Factor Authentication for securely logging in.

## Security
- User passwords are never stored as cleartext, instead they are salted and hashed.

- 2FA randomly generated secrets are stored on disc, encrypted with a master secret. 

- The master secret is not stored in environment variables or files, it lives only in memory.