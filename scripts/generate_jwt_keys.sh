#!/bin/bash
# generate_jwt_keys.sh

echo "Generating JWT ECDSA keys for production..."

# Generate encrypted private key
echo "Enter password for private key encryption:"
openssl ecparam -genkey -name prime256v1 -noout | openssl ec -aes256 -out jwt_private_key.pem

# Extract public key
openssl ec -in jwt_private_key.pem -pubout -out jwt_public_key.pem

# Set secure permissions
chmod 600 jwt_private_key.pem
chmod 644 jwt_public_key.pem

echo "Keys generated successfully!"
echo "Private key: jwt_private_key.pem (encrypted)"
echo "Public key: jwt_public_key.pem"

# Show environment variable format
echo ""
echo "To use in your application:"
echo "export JWT_PRIVATE_KEY=\"\$(cat jwt_private_key.pem)\""
echo "export JWT_PUBLIC_KEY=\"\$(cat jwt_public_key.pem)\""